package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/generate"
	"github.com/MUKE-coder/grit/v3/internal/routeparser"
)

// asJSON renders a value as indented JSON. Tool output is consumed by a model,
// and JSON survives truncation and quoting better than an ASCII table.
func asJSON(v interface{}) (string, error) {
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("encoding result: %w", err)
	}
	return string(out), nil
}

func (s *Server) projectInfo() (string, error) {
	info := map[string]interface{}{
		"root": s.Root,
	}

	data, err := os.ReadFile(filepath.Join(s.Root, "grit.json"))
	if err != nil {
		return "", fmt.Errorf("reading grit.json: %w", err)
	}
	var cfg map[string]interface{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return "", fmt.Errorf("parsing grit.json: %w", err)
	}
	for k, v := range cfg {
		info[k] = v
	}

	if mod, err := goModulePath(s.Root); err == nil {
		info["go_module"] = mod
	}
	info["api_dir"] = relOrEmpty(s.Root, apiRoot(s.Root))

	return asJSON(info)
}

func (s *Server) listRoutes(args json.RawMessage) (string, error) {
	var a struct {
		Method   string `json:"method"`
		Contains string `json:"contains"`
	}
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil {
			return "", fmt.Errorf("invalid arguments: %w", err)
		}
	}

	routesFile, err := routeparser.FindRoutesFile(s.Root)
	if err != nil {
		return "", err
	}
	routes, err := routeparser.ParseProject(routesFile)
	if err != nil {
		return "", err
	}

	method := strings.ToUpper(strings.TrimSpace(a.Method))
	contains := strings.ToLower(strings.TrimSpace(a.Contains))

	out := make([]map[string]string, 0, len(routes))
	for _, r := range routes {
		if method != "" && r.Method != method {
			continue
		}
		if contains != "" && !strings.Contains(strings.ToLower(r.Path), contains) {
			continue
		}
		out = append(out, map[string]string{
			"method":  r.Method,
			"path":    r.Path,
			"handler": r.Handler,
			// Group is the coarse bucket, access is what a caller must present.
			// These were one field reporting the group under the name "access",
			// which told a model that /api/v1/auth/me was public.
			"group":  r.Group,
			"access": r.Access.Summary(),
		})
	}

	return asJSON(map[string]interface{}{
		"count":  len(out),
		"source": relOrEmpty(s.Root, routesFile),
		"routes": out,
	})
}

func (s *Server) describeModels(args json.RawMessage) (string, error) {
	var a struct {
		Model string `json:"model"`
	}
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil {
			return "", fmt.Errorf("invalid arguments: %w", err)
		}
	}

	modelsDir := filepath.Join(apiRoot(s.Root), "internal", "models")
	entries, err := os.ReadDir(modelsDir)
	if err != nil {
		return "", fmt.Errorf("reading models directory: %w", err)
	}

	want := strings.ToLower(strings.TrimSpace(a.Model))
	models := []map[string]interface{}{}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		structs, err := generate.ParseGoStructs(filepath.Join(modelsDir, entry.Name()))
		if err != nil {
			// One unparseable file should not blank the whole answer.
			continue
		}
		for _, st := range structs {
			if want != "" && strings.ToLower(st.Name) != want {
				continue
			}
			fields := make([]map[string]string, 0, len(st.Fields))
			for _, f := range st.Fields {
				fields = append(fields, map[string]string{
					"name":      f.Name,
					"go_type":   f.GoType,
					"json_name": f.JSONName,
					"gorm":      f.GORMTag,
				})
			}
			models = append(models, map[string]interface{}{
				"name":   st.Name,
				"file":   entry.Name(),
				"fields": fields,
			})
		}
	}

	sort.Slice(models, func(i, j int) bool {
		return models[i]["name"].(string) < models[j]["name"].(string)
	})

	if want != "" && len(models) == 0 {
		return "", fmt.Errorf("no model named %q in %s", a.Model, relOrEmpty(s.Root, modelsDir))
	}

	return asJSON(map[string]interface{}{
		"count":  len(models),
		"source": relOrEmpty(s.Root, modelsDir),
		"models": models,
	})
}

// apiRoot returns the directory holding the Go API — the project root for
// single/api-in-place layouts, apps/api for a monorepo.
func apiRoot(root string) string {
	mono := filepath.Join(root, "apps", "api")
	if _, err := os.Stat(filepath.Join(mono, "go.mod")); err == nil {
		return mono
	}
	return root
}

func goModulePath(root string) (string, error) {
	data, err := os.ReadFile(filepath.Join(apiRoot(root), "go.mod"))
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}
	return "", fmt.Errorf("no module line in go.mod")
}

// relOrEmpty renders a path relative to the project root so tool output does
// not leak the absolute layout of the machine it ran on.
func relOrEmpty(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return path
	}
	return filepath.ToSlash(rel)
}
