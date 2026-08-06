package scaffold

// The seed catalogues.
//
// Three locales ship, and two of them are real translations rather than copies
// of the English with the keys renamed. A stub catalogue is worse than none: it
// makes the feature look finished, and the first person to switch locale finds
// out otherwise.
//
// English, French and Swahili. Swahili because the scaffold is built in Uganda
// and East Africa is not a rounding error, and because a framework whose only
// second language is French has quietly decided who it is for.
//
// Keys are dotted and grouped by area. Flat JSON rather than nested, because a
// flat map needs no traversal code and a translator editing it cannot break the
// shape by mis-indenting a brace.

func i18nLocaleEN() string {
	return `{
  "errors.not_found": "Not found",
  "errors.resource_not_found": "%s not found",
  "errors.unauthorized": "You need to sign in to do that",
  "errors.forbidden": "You do not have permission to do that",
  "errors.validation": "Some of the details are not right",
  "errors.internal": "Something went wrong on our end",
  "errors.invalid_request": "That request could not be understood",
  "errors.rate_limited": "Too many requests. Try again shortly",
  "errors.conflict": "That conflicts with something that already exists",
  "errors.payload_too_large": "That upload is too large",

  "auth.invalid_credentials": "That email and password do not match",
  "auth.account_locked": "This account is locked. Try again later",
  "auth.email_taken": "An account with that email already exists",
  "auth.email_not_verified": "Please verify your email address first",
  "auth.token_expired": "Your session has expired. Please sign in again",
  "auth.token_invalid": "That sign-in link is not valid",
  "auth.registered": "Account created",
  "auth.signed_in": "Signed in",
  "auth.signed_out": "Signed out",
  "auth.password_reset_sent": "If that email exists, a reset link is on its way",
  "auth.password_changed": "Password changed. Other devices have been signed out",
  "auth.totp_required": "Enter the code from your authenticator app",
  "auth.totp_invalid": "That code is not right",

  "validation.required": "%s is required",
  "validation.email": "%s must be a valid email address",
  "validation.min_length": "%s must be at least %d characters",
  "validation.max_length": "%s must be at most %d characters",
  "validation.unique": "That %s is already taken",
  "validation.numeric": "%s must be a number",
  "validation.positive": "%s must be greater than zero",

  "resource.created": "%s created",
  "resource.updated": "%s updated",
  "resource.deleted": "%s deleted",
  "resource.restored": "%s restored",

  "common.yes": "Yes",
  "common.no": "No",
  "common.saved": "Saved",
  "common.cancelled": "Cancelled"
}
`
}

func i18nLocaleFR() string {
	return `{
  "errors.not_found": "Introuvable",
  "errors.resource_not_found": "%s introuvable",
  "errors.unauthorized": "Vous devez vous connecter pour faire cela",
  "errors.forbidden": "Vous n'avez pas la permission de faire cela",
  "errors.validation": "Certaines informations sont incorrectes",
  "errors.internal": "Une erreur est survenue de notre côté",
  "errors.invalid_request": "Cette requête n'a pas pu être comprise",
  "errors.rate_limited": "Trop de requêtes. Réessayez dans un instant",
  "errors.conflict": "Cela entre en conflit avec un élément existant",
  "errors.payload_too_large": "Ce fichier est trop volumineux",

  "auth.invalid_credentials": "Cet e-mail et ce mot de passe ne correspondent pas",
  "auth.account_locked": "Ce compte est verrouillé. Réessayez plus tard",
  "auth.email_taken": "Un compte existe déjà avec cet e-mail",
  "auth.email_not_verified": "Veuillez d'abord vérifier votre adresse e-mail",
  "auth.token_expired": "Votre session a expiré. Veuillez vous reconnecter",
  "auth.token_invalid": "Ce lien de connexion n'est pas valide",
  "auth.registered": "Compte créé",
  "auth.signed_in": "Connecté",
  "auth.signed_out": "Déconnecté",
  "auth.password_reset_sent": "Si cet e-mail existe, un lien de réinitialisation est en route",
  "auth.password_changed": "Mot de passe modifié. Les autres appareils ont été déconnectés",
  "auth.totp_required": "Saisissez le code de votre application d'authentification",
  "auth.totp_invalid": "Ce code est incorrect",

  "validation.required": "%s est obligatoire",
  "validation.email": "%s doit être une adresse e-mail valide",
  "validation.min_length": "%s doit contenir au moins %d caractères",
  "validation.max_length": "%s doit contenir au plus %d caractères",
  "validation.unique": "Ce %s est déjà utilisé",
  "validation.numeric": "%s doit être un nombre",
  "validation.positive": "%s doit être supérieur à zéro",

  "resource.created": "%s créé",
  "resource.updated": "%s mis à jour",
  "resource.deleted": "%s supprimé",
  "resource.restored": "%s restauré",

  "common.yes": "Oui",
  "common.no": "Non",
  "common.saved": "Enregistré",
  "common.cancelled": "Annulé"
}
`
}

func i18nLocaleSW() string {
	return `{
  "errors.not_found": "Haipatikani",
  "errors.resource_not_found": "%s haipatikani",
  "errors.unauthorized": "Unahitaji kuingia ili kufanya hivyo",
  "errors.forbidden": "Huna ruhusa ya kufanya hivyo",
  "errors.validation": "Baadhi ya taarifa si sahihi",
  "errors.internal": "Hitilafu imetokea kwa upande wetu",
  "errors.invalid_request": "Ombi hilo halikueleweka",
  "errors.rate_limited": "Maombi mengi mno. Jaribu tena baadaye kidogo",
  "errors.conflict": "Hiyo inagongana na kitu kilichopo tayari",
  "errors.payload_too_large": "Faili hilo ni kubwa mno",

  "auth.invalid_credentials": "Barua pepe na nenosiri hazilingani",
  "auth.account_locked": "Akaunti hii imefungwa. Jaribu tena baadaye",
  "auth.email_taken": "Akaunti yenye barua pepe hiyo tayari ipo",
  "auth.email_not_verified": "Tafadhali thibitisha barua pepe yako kwanza",
  "auth.token_expired": "Kipindi chako kimeisha. Tafadhali ingia tena",
  "auth.token_invalid": "Kiungo hicho cha kuingia si halali",
  "auth.registered": "Akaunti imetengenezwa",
  "auth.signed_in": "Umeingia",
  "auth.signed_out": "Umetoka",
  "auth.password_reset_sent": "Kama barua pepe hiyo ipo, kiungo cha kubadilisha kinakuja",
  "auth.password_changed": "Nenosiri limebadilishwa. Vifaa vingine vimetolewa",
  "auth.totp_required": "Weka msimbo kutoka programu yako ya uthibitishaji",
  "auth.totp_invalid": "Msimbo huo si sahihi",

  "validation.required": "%s inahitajika",
  "validation.email": "%s lazima iwe barua pepe halali",
  "validation.min_length": "%s lazima iwe na angalau herufi %d",
  "validation.max_length": "%s lazima isizidi herufi %d",
  "validation.unique": "%s hiyo tayari imetumika",
  "validation.numeric": "%s lazima iwe namba",
  "validation.positive": "%s lazima iwe kubwa kuliko sifuri",

  "resource.created": "%s imetengenezwa",
  "resource.updated": "%s imesasishwa",
  "resource.deleted": "%s imefutwa",
  "resource.restored": "%s imerejeshwa",

  "common.yes": "Ndiyo",
  "common.no": "Hapana",
  "common.saved": "Imehifadhiwa",
  "common.cancelled": "Imeghairiwa"
}
`
}

// i18nHandlerGo emits internal/handlers/i18n.go — the endpoint the frontend
// calls to discover which locales this build actually has, and to fetch one.
// Without it a language switcher has to hardcode the list and will offer
// languages the binary cannot serve.
func i18nHandlerGo() string {
	return `package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"{{MODULE}}/internal/i18n"
	"{{MODULE}}/internal/middleware"
)

// I18nHandler serves the catalogues to browser clients.
type I18nHandler struct{}

func NewI18nHandler() *I18nHandler { return &I18nHandler{} }

// LocalesResponse is the list of languages this build was compiled with.
type LocalesResponse struct {
	Locales []string ` + "`" + `json:"locales"` + "`" + `
	Current string   ` + "`" + `json:"current"` + "`" + `
	Default string   ` + "`" + `json:"default"` + "`" + `
}

// List reports the available locales and which one this request resolved to.
func (h *I18nHandler) List(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": LocalesResponse{
			Locales: i18n.Available(),
			Current: middleware.LocaleOf(c),
			Default: i18n.Default,
		},
	})
}
`
}

// i18nTestGo emits internal/i18n/i18n_test.go.
//
// Compiling proves the package parses. It does not prove a catalogue loads, a
// locale resolves, or that a translator's trailing comma has not silently
// emptied the French. This does.
func i18nTestGo() string {
	return `package i18n

import "testing"

func TestCataloguesResolve(t *testing.T) {
	if err := Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	for _, loc := range []string{"en", "fr", "sw"} {
		if !Has(loc) {
			t.Fatalf("locale %q not compiled in; Available=%v", loc, Available())
		}
	}
	if got := T("fr", "errors.not_found"); got != "Introuvable" {
		t.Fatalf("fr errors.not_found = %q", got)
	}
	// An unknown locale falls back to English rather than to nothing.
	if got := T("de", "errors.not_found"); got != "Not found" {
		t.Fatalf("de fallback = %q", got)
	}
	// An unknown key returns the key, so it shows up in the UI and gets fixed
	// instead of rendering as a blank space nobody notices.
	if got := T("en", "nope.missing"); got != "nope.missing" {
		t.Fatalf("missing key = %q, want the key back", got)
	}
	if got := T("fr", "errors.resource_not_found", "Produit"); got != "Produit introuvable" {
		t.Fatalf("fr with argument = %q", got)
	}
}

// Every catalogue must carry every key. A missing one silently falls back to
// English, which looks like a half-finished translation to the person reading
// it and like nothing at all to the person who shipped it.
func TestCataloguesAgreeOnKeys(t *testing.T) {
	if err := Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	base := messages[Default]
	for locale, m := range messages {
		if locale == Default {
			continue
		}
		for key := range base {
			if _, ok := m[key]; !ok {
				t.Errorf("%s is missing %q", locale, key)
			}
		}
		for key := range m {
			if _, ok := base[key]; !ok {
				t.Errorf("%s has %q, which %s does not", locale, key, Default)
			}
		}
	}
}

func TestNegotiate(t *testing.T) {
	cases := map[string]string{
		"fr":                     "fr",
		"fr-CA":                  "fr",
		"de":                     "en",
		"":                       "en",
		"de, fr;q=0.9, en;q=0.8": "fr",
		"en;q=0.3, sw;q=0.9":     "sw",
		"*":                      "en",
	}
	for header, want := range cases {
		if got := Negotiate(header); got != want {
			t.Fatalf("Negotiate(%q) = %q, want %q", header, got, want)
		}
	}
}
`
}
