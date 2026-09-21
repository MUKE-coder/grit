package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// The Expo API client refreshed once per failing request. A screen that loads
// three things at once after the access token expires got three 401s and sent
// the same refresh token three times; the server rotates refresh tokens and
// treats a spent one as stolen, so the second use revoked the session and the
// user was signed out 15 minutes into using the app. The admin and desktop
// clients already queued requests behind one refresh. Found building the
// WhatsApp blueprint's mobile app.
const (
	expoRefreshBlockOld = "    // Try refresh if unauthorized\n" +
		"    if (res.status === 401 && !isAuthEndpoint) {\n" +
		"      const refreshToken = await this.getRefreshToken();\n" +
		"      if (refreshToken) {\n" +
		"        const refreshRes = await fetchWithTimeout(`${API_URL}/auth/refresh`, {\n" +
		"          method: \"POST\",\n" +
		"          headers: { \"Content-Type\": \"application/json\" },\n" +
		"          body: JSON.stringify({ refresh_token: refreshToken }),\n" +
		"        });\n" +
		"\n" +
		"        if (refreshRes.ok) {\n" +
		"          const data = await refreshRes.json();\n" +
		"          await this.setTokens(data.data.tokens.access_token, data.data.tokens.refresh_token);\n" +
		"\n" +
		"          headers[\"Authorization\"] = `Bearer ${data.data.tokens.access_token}`;\n" +
		"          res = await fetchWithTimeout(`${API_URL}${endpoint}`, {\n" +
		"            method: options.method || \"GET\",\n" +
		"            headers,\n" +
		"            body: options.body ? JSON.stringify(options.body) : undefined,\n" +
		"          });\n" +
		"        } else {\n" +
		"          await this.clearTokens();\n" +
		"          throw new Error(\"Session expired\");\n" +
		"        }\n" +
		"      }\n" +
		"    }\n"

	expoRefreshBlockNew = "    // Refresh once and retry. Requests that fail together share one refresh:\n" +
		"    // the server treats a refresh token used twice as stolen and ends the\n" +
		"    // session, so a refresh per request signed the user out.\n" +
		"    if (res.status === 401 && !isAuthEndpoint) {\n" +
		"      const fresh = await this.refreshAccessToken();\n" +
		"      if (fresh) {\n" +
		"        headers[\"Authorization\"] = `Bearer ${fresh}`;\n" +
		"        res = await fetchWithTimeout(`${API_URL}${endpoint}`, {\n" +
		"          method: options.method || \"GET\",\n" +
		"          headers,\n" +
		"          body: options.body ? JSON.stringify(options.body) : undefined,\n" +
		"        });\n" +
		"      }\n" +
		"    }\n"

	expoRequestSignature = "  private async request(endpoint: string, options: RequestOptions = {}) {\n"

	expoRefreshMethod = "  private refreshing: Promise<string | null> | null = null;\n" +
		"\n" +
		"  /**\n" +
		"   * Trades the refresh token for a new pair, once however many requests are\n" +
		"   * waiting on it. Resolves to the new access token, or null when there is no\n" +
		"   * refresh token; throws \"Session expired\" when the server refuses it.\n" +
		"   */\n" +
		"  private refreshAccessToken(): Promise<string | null> {\n" +
		"    this.refreshing ??= (async () => {\n" +
		"      const refreshToken = await this.getRefreshToken();\n" +
		"      if (!refreshToken) return null;\n" +
		"      const refreshRes = await fetchWithTimeout(`${API_URL}/auth/refresh`, {\n" +
		"        method: \"POST\",\n" +
		"        headers: { \"Content-Type\": \"application/json\" },\n" +
		"        body: JSON.stringify({ refresh_token: refreshToken }),\n" +
		"      });\n" +
		"      if (!refreshRes.ok) {\n" +
		"        await this.clearTokens();\n" +
		"        throw new Error(\"Session expired\");\n" +
		"      }\n" +
		"      const data = await refreshRes.json();\n" +
		"      await this.setTokens(data.data.tokens.access_token, data.data.tokens.refresh_token);\n" +
		"      return data.data.tokens.access_token as string;\n" +
		"    })().finally(() => {\n" +
		"      this.refreshing = null;\n" +
		"    });\n" +
		"    return this.refreshing;\n" +
		"  }\n" +
		"\n"
)

// repairExpoRefreshSource moves an Expo API client Grit wrote to one shared refresh.
func repairExpoRefreshSource(src string) (string, []string, []string) {
	if strings.Contains(src, "refreshAccessToken()") {
		return src, nil, nil
	}
	if strings.Count(src, expoRefreshBlockOld) != 1 || strings.Count(src, expoRequestSignature) != 1 {
		if strings.Contains(src, "/auth/refresh") {
			return src, nil, []string{"the API client is not the one Grit wrote: make requests that fail together share one refresh, or a screen that loads two things after the token expires signs the user out"}
		}
		return src, nil, nil
	}
	out := strings.Replace(src, expoRefreshBlockOld, expoRefreshBlockNew, 1)
	out = strings.Replace(out, expoRequestSignature, expoRefreshMethod+expoRequestSignature, 1)
	return out, []string{"requests that fail together share one refresh, instead of signing the user out"}, nil
}

// repairExpoRefresh applies repairExpoRefreshSource to apps/expo/lib/api.ts.
func repairExpoRefresh(root string) error {
	path := filepath.Join(root, "apps", "expo", "lib", "api.ts")
	if !fileExists(path) {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	return repairTextFile(root, m, path, repairExpoRefreshSource)
}
