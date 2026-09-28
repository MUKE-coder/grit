package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// grit scale: which scaling stage is this deployment at, and what is the one
// thing to do next.
//
// The measuring happens in the API, because these numbers are only true where
// the load is. A laptop's p99 and a laptop's connection count describe a
// laptop. This command points at a running deployment, asks it, and renders the
// answer.
//
// It deliberately gives one recommendation. The guide this follows has a rule,
// "add exactly one thing", and a tool that lists six possible improvements has
// quietly handed the hardest part of the job, choosing, back to the reader.

type scaleReport struct {
	Latency struct {
		Samples int     `json:"samples"`
		P50     float64 `json:"p50_ms"`
		P95     float64 `json:"p95_ms"`
		P99     float64 `json:"p99_ms"`
		Max     float64 `json:"max_ms"`
		RPS     float64 `json:"requests_per_second"`
	} `json:"latency"`
	DB struct {
		Driver          string  `json:"driver"`
		OpenConnections int     `json:"open_connections"`
		MaxConnections  int     `json:"max_connections"`
		PoolMax         int     `json:"pool_max_per_instance"`
		Replicas        int     `json:"replicas"`
		ReplicationLag  float64 `json:"replication_lag_seconds"`
		SlowQueries     []struct {
			Query   string  `json:"query"`
			Calls   int64   `json:"calls"`
			MeanMS  float64 `json:"mean_ms"`
			TotalMS float64 `json:"total_ms"`
		} `json:"slowest_queries"`
		SeqScanTables []struct {
			Table    string `json:"table"`
			SeqScans int64  `json:"seq_scans"`
			LiveRows int64  `json:"live_rows"`
		} `json:"tables_scanned_end_to_end"`
		QueryStats bool `json:"query_stats_available"`
	} `json:"database"`
	Cache struct {
		Configured bool    `json:"configured"`
		Hits       int64   `json:"hits"`
		Misses     int64   `json:"misses"`
		HitRate    float64 `json:"hit_rate"`
	} `json:"cache"`
	Stage string   `json:"stage"`
	Next  string   `json:"next"`
	Why   string   `json:"why"`
	Notes []string `json:"notes"`
}

func scaleCmd() *cobra.Command {
	var apiURL, token string
	var asJSON bool

	cmd := &cobra.Command{
		Use:   "scale",
		Short: "Measure this deployment and name the one thing to do next",
		Long: "Asks a running Grit API for its request percentiles, connection use, slowest queries\n" +
			"and cache hit rate, then names the scaling stage it is at and the single next step.\n\n" +
			"Point it at production. These numbers describe wherever they are measured, and a\n" +
			"laptop under no load is always healthy.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if apiURL == "" {
				apiURL = firstNonEmpty(os.Getenv("GRIT_API_URL"), os.Getenv("API_URL"), "http://localhost:8080")
			}
			if token == "" {
				token = os.Getenv("GRIT_ADMIN_TOKEN")
			}

			report, err := fetchScale(apiURL, token)
			if err != nil {
				return err
			}

			if asJSON {
				out, _ := json.MarshalIndent(report, "", "  ")
				fmt.Println(string(out))
				return nil
			}

			printLogo()
			renderScale(report, apiURL)
			return nil
		},
	}

	cmd.Flags().StringVar(&apiURL, "api", "", "Base URL of the running API (default: API_URL, or localhost:8080)")
	cmd.Flags().StringVar(&token, "token", "", "Admin bearer token (default: GRIT_ADMIN_TOKEN)")
	cmd.Flags().BoolVar(&asJSON, "json", false, "Print the raw measurement")

	return cmd
}

func fetchScale(base, token string) (*scaleReport, error) {
	url := strings.TrimSuffix(base, "/") + "/api/v1/scale"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := (&http.Client{Timeout: 20 * time.Second}).Do(req)
	if err != nil {
		return nil, fmt.Errorf(
			"could not reach the API at %s: %w\n"+
				"    Start it, or pass --api https://your-deployment", base, err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusUnauthorized, http.StatusForbidden:
		return nil, fmt.Errorf(
			"the API refused: this endpoint is admin only.\n" +
				"    Pass --token, or set GRIT_ADMIN_TOKEN to an admin's access token")
	case http.StatusNotFound:
		return nil, fmt.Errorf(
			"that API has no /admin/scale endpoint.\n" +
				"    It predates grit scale: run grit upgrade in the project first")
	default:
		return nil, fmt.Errorf("the API answered %d: %s", resp.StatusCode, truncate(string(body), 200))
	}

	var envelope struct {
		Data scaleReport `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("could not read the measurement: %w", err)
	}
	return &envelope.Data, nil
}

func renderScale(r *scaleReport, apiURL string) {
	bold := color.New(color.Bold)
	muted := color.New(color.FgHiBlack)
	green := color.New(color.FgHiGreen, color.Bold)
	yellow := color.New(color.FgHiYellow, color.Bold)

	muted.Printf("  Measured at %s\n\n", apiURL)

	bold.Println("  Requests")
	if r.Latency.Samples == 0 {
		muted.Println("    Nothing measured yet. This API has served no requests since it started.")
	} else {
		fmt.Printf("    p50 %6.0fms   p95 %6.0fms   p99 %6.0fms   max %6.0fms\n",
			r.Latency.P50, r.Latency.P95, r.Latency.P99, r.Latency.Max)
		muted.Printf("    %d requests in the window, %.1f/s since start\n", r.Latency.Samples, r.Latency.RPS)
	}

	fmt.Println()
	bold.Println("  Database")
	if r.DB.MaxConnections > 0 {
		pct := float64(r.DB.OpenConnections) / float64(r.DB.MaxConnections) * 100
		fmt.Printf("    %d of %d connections in use (%.0f%%), pool max %d per instance\n",
			r.DB.OpenConnections, r.DB.MaxConnections, pct, r.DB.PoolMax)
	} else {
		muted.Printf("    %s\n", r.DB.Driver)
	}
	if r.DB.Replicas > 0 {
		fmt.Printf("    %d read replica(s), %.1fs behind\n", r.DB.Replicas, r.DB.ReplicationLag)
	} else {
		muted.Println("    No read replicas")
	}

	if len(r.DB.SeqScanTables) > 0 {
		fmt.Println()
		yellow.Println("  Tables being read end to end")
		for _, t := range r.DB.SeqScanTables {
			fmt.Printf("    %-28s %8d scans of %d rows\n", t.Table, t.SeqScans, t.LiveRows)
		}
		muted.Println("    An index is free. A bigger database is not.")
	}

	if len(r.DB.SlowQueries) > 0 {
		fmt.Println()
		bold.Println("  Costliest queries, by total time")
		for _, q := range r.DB.SlowQueries {
			fmt.Printf("    %7.0fms total  %7d calls  %6.1fms each\n", q.TotalMS, q.Calls, q.MeanMS)
			muted.Printf("      %s\n", truncate(collapse(q.Query), 92))
		}
	}

	fmt.Println()
	bold.Println("  Cache")
	if !r.Cache.Configured {
		muted.Println("    Not configured")
	} else if r.Cache.Hits+r.Cache.Misses == 0 {
		muted.Println("    Configured, nothing cached yet")
	} else {
		fmt.Printf("    %.0f%% hit rate (%d hits, %d misses)\n",
			r.Cache.HitRate*100, r.Cache.Hits, r.Cache.Misses)
	}

	fmt.Println()
	if strings.HasPrefix(r.Stage, "Healthy") {
		green.Printf("  %s\n", r.Stage)
	} else {
		yellow.Printf("  %s\n", r.Stage)
	}
	muted.Printf("    %s\n\n", r.Why)
	bold.Printf("  Do this next: ")
	fmt.Printf("%s\n", r.Next)

	if len(r.Notes) > 0 {
		fmt.Println()
		muted.Println("  Worth knowing, but not yet:")
		for _, n := range r.Notes {
			muted.Printf("    %s\n", n)
		}
	}
	fmt.Println()
	muted.Println("  The whole guide: https://gritframework.dev/docs/scaling")
	fmt.Println()
}

// collapse puts a SQL statement on one line so a table of them lines up.
func collapse(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
