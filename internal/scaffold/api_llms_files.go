package scaffold

// internal/llms, and the two files an agent reads before it calls a generated
// API.
//
// The OpenAPI spec at /docs/openapi.json is the contract and is the better
// artifact for generating a client. It is also silent about everything the
// endpoints share: the response envelope, the stable error codes, the version
// rewrite, the page ceiling, which header carries the token. An agent that
// reads the spec learns the shape of each endpoint and none of that.
//
// So llms.txt is the orientation at the path models look for it, and
// llms-full.txt is every route the router holds, which is the one thing the
// spec cannot answer: the spec documents the routes somebody wrote an override
// for, and this is the router's own table.

// apiLLMSGo emits internal/llms/llms.go.
func apiLLMSGo() string { return tmpl("api/llms/llms.go") }

// apiLLMSTestGo emits internal/llms/llms_test.go.
func apiLLMSTestGo() string { return tmpl("api/llms/llms_test.go") }

// llmsRoutesBlock mounts the two files, beside the API reference and under the
// same condition: they describe the whole surface, admin routes included.
const llmsRoutesBlock = `
	// llms.txt and llms-full.txt, for an agent pointed at this service.
	//
	// Beside the reference and gated with it: a route list is not a secret and
	// is also not something to hand out by default. See internal/llms.
	if cfg.AppEnv != "production" || cfg.APIDocsPublic {
		registerLLMSFiles(r, cfg)
	}
`

// llmsRegisterFunc is the function that block calls, written beside
// registerAPIDocs for the same reason: the wiring is one line in routes.go and
// the detail is somewhere a reader can find it.
const llmsRegisterFunc = `
// registerLLMSFiles serves /llms.txt and /llms-full.txt.
//
// The route table is read on each request rather than cached. It is a few
// hundred entries and these two paths are not on a hot path, and a cache here
// would be a copy that goes stale the first time something registers a route
// after boot.
func registerLLMSFiles(r *gin.Engine, cfg *config.Config) {
	conf := llms.Config{
		AppName:     cfg.AppName,
		BaseURL:     cfg.AppURL,
		Version:     APIVersion,
		MaxPageSize: paginate.MaxPageSize,
	}
	// Collected once per request from the engine, so a route a plugin or your
	// own code adds is in the file without anybody maintaining a list.
	table := func() []llms.Route {
		infos := r.Routes()
		out := make([]llms.Route, 0, len(infos))
		for _, info := range infos {
			out = append(out, llms.Route{
				Method:  info.Method,
				Path:    info.Path,
				Handler: shortHandlerName(info.Handler),
			})
		}
		return out
	}

	r.GET("/llms.txt", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(llms.Index(conf)))
	})
	r.GET("/llms-full.txt", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(llms.Full(conf, table())))
	})
}

// shortHandlerName turns what Gin records, which is a fully qualified function
// name with a closure suffix, into the part worth reading:
// "acme/apps/api/internal/handlers.(*UserHandler).List-fm" becomes
// "UserHandler.List".
func shortHandlerName(full string) string {
	name := full
	if i := strings.LastIndex(name, "/"); i >= 0 {
		name = name[i+1:]
	}
	name = strings.TrimSuffix(name, "-fm")
	// Drop the package qualifier and the pointer receiver's punctuation.
	if i := strings.Index(name, "."); i >= 0 {
		name = name[i+1:]
	}
	name = strings.ReplaceAll(name, "(*", "")
	name = strings.ReplaceAll(name, ")", "")
	// An anonymous handler ends in func1, func2 and so on, which names nothing:
	// a closure declared inside Setup arrives here as "Setup.func1".
	if last := name[strings.LastIndex(name, ".")+1:]; name == "" || isAnonymous(last) {
		return ""
	}
	return name
}

// isAnonymous reports whether a name segment is Go's funcN for a closure.
func isAnonymous(segment string) bool {
	if !strings.HasPrefix(segment, "func") || len(segment) == len("func") {
		return false
	}
	for _, r := range segment[len("func"):] {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
`
