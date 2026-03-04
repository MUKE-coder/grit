package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

func writeDesktopStudioFiles(root string, opts DesktopOptions) error {
	files := map[string]string{
		filepath.Join(root, "cmd", "studio", "main.go"): desktopStudioMainGo(),
	}

	for path, content := range files {
		content = strings.ReplaceAll(content, "<MODULE>", opts.ProjectName)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

func desktopStudioMainGo() string {
	return `package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"<MODULE>/internal/config"
	"<MODULE>/internal/db"
	"<MODULE>/internal/models"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()
	database := db.Connect(cfg)

	sqlDB, err := database.DB()
	if err != nil {
		log.Fatalf("failed to get underlying database: %v", err)
	}

	// Simple database browser endpoint
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, studioHTML())
	})

	mux.HandleFunc("/api/tables", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		tables := []string{}
		rows, err := sqlDB.Query("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' ORDER BY name")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var name string
			rows.Scan(&name)
			tables = append(tables, name)
		}
		fmt.Fprintf(w, "[")
		for i, t := range tables {
			if i > 0 {
				fmt.Fprintf(w, ",")
			}
			fmt.Fprintf(w, "\"" + t + "\"")
		}
		fmt.Fprintf(w, "]")
	})

	mux.HandleFunc("/api/query", func(w http.ResponseWriter, r *http.Request) {
		table := r.URL.Query().Get("table")
		if table == "" {
			http.Error(w, "table parameter required", 400)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		rows, err := sqlDB.Query("SELECT * FROM " + table + " LIMIT 100")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		cols, _ := rows.Columns()
		fmt.Fprintf(w, "{\"columns\":[")
		for i, c := range cols {
			if i > 0 {
				fmt.Fprintf(w, ",")
			}
			fmt.Fprintf(w, "\"" + c + "\"")
		}
		fmt.Fprintf(w, "],\"rows\":[")

		values := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))
		for i := range values {
			ptrs[i] = &values[i]
		}

		rowIdx := 0
		for rows.Next() {
			rows.Scan(ptrs...)
			if rowIdx > 0 {
				fmt.Fprintf(w, ",")
			}
			fmt.Fprintf(w, "[")
			for i, v := range values {
				if i > 0 {
					fmt.Fprintf(w, ",")
				}
				switch val := v.(type) {
				case nil:
					fmt.Fprintf(w, "null")
				case []byte:
					fmt.Fprintf(w, "\"" + string(val) + "\"")
				case string:
					fmt.Fprintf(w, "\"" + val + "\"")
				default:
					fmt.Fprintf(w, "%v", val)
				}
			}
			fmt.Fprintf(w, "]")
			rowIdx++
		}
		fmt.Fprintf(w, "]}")
	})

	// Register models for reference
	_ = []interface{}{
		&models.User{},
		&models.Blog{},
		&models.Contact{},
		// grit:studio-models
	}

	fmt.Println("GORM Studio running at http://localhost:4000")
	fmt.Println("Press Ctrl+C to stop")
	log.Fatal(http.ListenAndServe(":4000", mux))
}

func studioHTML() string {
	return ` + "`" + `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Grit Studio</title>
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body { font-family: system-ui, sans-serif; background: #0a0a0f; color: #e8e8f0; }
    .header { padding: 16px 24px; border-bottom: 1px solid #2a2a3a; display: flex; align-items: center; gap: 12px; }
    .header h1 { font-size: 18px; font-weight: 600; }
    .header span { color: #6c5ce7; }
    .container { display: flex; height: calc(100vh - 57px); }
    .sidebar { width: 220px; border-right: 1px solid #2a2a3a; padding: 12px; overflow-y: auto; }
    .sidebar button { display: block; width: 100%; text-align: left; padding: 8px 12px; background: none; border: none; color: #9090a8; cursor: pointer; border-radius: 6px; font-size: 14px; margin-bottom: 2px; }
    .sidebar button:hover { background: #22222e; color: #e8e8f0; }
    .sidebar button.active { background: #6c5ce7; color: white; }
    .main { flex: 1; overflow: auto; padding: 16px; }
    table { width: 100%; border-collapse: collapse; font-size: 13px; }
    th { text-align: left; padding: 8px 12px; background: #22222e; color: #9090a8; font-weight: 500; border-bottom: 1px solid #2a2a3a; position: sticky; top: 0; }
    td { padding: 8px 12px; border-bottom: 1px solid #1a1a2a; max-width: 300px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
    tr:hover td { background: #15151f; }
    .empty { color: #606078; text-align: center; padding: 40px; }
  </style>
</head>
<body>
  <div class="header"><span>&#9632;</span> <h1>Grit Studio</h1></div>
  <div class="container">
    <div class="sidebar" id="sidebar"></div>
    <div class="main" id="main"><div class="empty">Select a table</div></div>
  </div>
  <script>
    let activeTable = "";
    fetch("/api/tables").then(r=>r.json()).then(tables => {
      const sb = document.getElementById("sidebar");
      tables.forEach(t => {
        const btn = document.createElement("button");
        btn.textContent = t;
        btn.onclick = () => loadTable(t);
        sb.appendChild(btn);
      });
    });
    function loadTable(name) {
      activeTable = name;
      document.querySelectorAll(".sidebar button").forEach(b => b.classList.toggle("active", b.textContent === name));
      fetch("/api/query?table=" + name).then(r=>r.json()).then(data => {
        const main = document.getElementById("main");
        if (!data.rows || data.rows.length === 0) { main.innerHTML = '<div class="empty">No rows</div>'; return; }
        let html = "<table><thead><tr>";
        data.columns.forEach(c => html += "<th>" + c + "</th>");
        html += "</tr></thead><tbody>";
        data.rows.forEach(row => { html += "<tr>"; row.forEach(v => html += "<td>" + (v === null ? '<span style="color:#606078">NULL</span>' : v) + "</td>"); html += "</tr>"; });
        html += "</tbody></table>";
        main.innerHTML = html;
      });
    }
  </script>
</body>
</html>` + "`" + `
}
`
}
