package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// blogListSanitize and blogSlugSanitize clean post HTML as the public blog API
// serves it. The web app renders posts on the server now, where DOMPurify has no
// DOM to work with, so the API's own sanitiser is the layer that covers posts
// stored before it sanitised on write.
const (
	blogListSanitize = `	for i := range blogs {
		blogs[i].Content = sanitize.HTML(blogs[i].Content)
	}

`
	blogSlugSanitize = `	blog.Content = sanitize.HTML(blog.Content)

`
	blogListCacheLine = "\tc.Header(\"Cache-Control\", \"public, max-age=300\")\n"
	blogSlugCacheLine = "\tc.Header(\"Cache-Control\", \"public, max-age=3600\")\n"
)

// composeWebAPIEnv gives the web container the API's address on the Docker
// network, for the server-rendered pages.
const composeWebAPIEnv = `    environment:
      # Server-rendered pages call the API over the Docker network.
      API_INTERNAL_URL: http://api:8080
`

// repairPublicPages brings a project up to the fix for H22 in the contact-app
// review. The Next.js pages themselves arrive whole with the web files. The two
// changes they rely on are the developer's files: the blog handler sanitises the
// HTML it serves, and the production web container knows where the API is.
func repairPublicPages(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	handler := filepath.Join(apiRoot, "internal", "handlers", "blog_handler.go")
	if fileExists(handler) && fileContains(filepath.Join(apiRoot, "internal", "sanitize", "html.go"), "func HTML(") {
		module := opts.Module()
		if err := repairSourceFile(root, m, handler, func(src string) (string, []string, []string) {
			return repairBlogSanitizeSource(src, module)
		}); err != nil {
			return err
		}
	}
	web := filepath.Join(root, "apps", "web")
	next := fileExists(filepath.Join(web, "next.config.ts")) || fileExists(filepath.Join(web, "next.config.js")) || fileExists(filepath.Join(web, "next.config.mjs"))
	if prod := filepath.Join(root, "docker-compose.prod.yml"); next && fileExists(prod) {
		if err := repairTextFile(root, m, prod, repairComposeWebAPIEnvSource); err != nil {
			return err
		}
	}
	return nil
}

func repairBlogSanitizeSource(src, module string) (string, []string, []string) {
	if !strings.Contains(src, "func (h *BlogHandler) GetBySlug(") || strings.Contains(src, "sanitize.HTML(") {
		return src, nil, nil
	}
	if strings.Count(src, blogListCacheLine) != 1 || strings.Count(src, blogSlugCacheLine) != 1 {
		return src, nil, []string{"blog_handler.go is not the file Grit wrote: sanitise blog.Content with sanitize.HTML before the public list and slug endpoints return it, since the web app now renders posts on the server"}
	}
	out := strings.Replace(src, blogListCacheLine, blogListSanitize+blogListCacheLine, 1)
	out = strings.Replace(out, blogSlugCacheLine, blogSlugSanitize+blogSlugCacheLine, 1)
	var ok bool
	if out, ok = addImportGroup(out, module+"/internal/sanitize"); !ok {
		return src, nil, []string{"could not add the sanitize import to blog_handler.go"}
	}
	return out, []string{"the public blog endpoints serve sanitised HTML, posts stored before sanitising on write included"}, nil
}

func repairComposeWebAPIEnvSource(src string) (string, []string, []string) {
	start, end, ok := serviceBlock(src, "web")
	if !ok || strings.Contains(src[start:end], "API_INTERNAL_URL") {
		return src, nil, nil
	}
	block := src[start:end]
	anchor := "        NEXT_PUBLIC_API_URL: ${API_URL:-http://localhost:8080}\n"
	if strings.Count(block, anchor) != 1 {
		return src, nil, []string{"the web service in docker-compose.prod.yml is not the one Grit wrote: give it API_INTERNAL_URL: http://api:8080, or its server-rendered pages cannot reach the API"}
	}
	block = strings.Replace(block, anchor, anchor+composeWebAPIEnv, 1)
	return src[:start] + block + src[end:],
		[]string{"the web container reaches the API over the Docker network for its server-rendered pages"}, nil
}
