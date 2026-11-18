package scaffold

// The seed catalogues.
//
// Four locales ship, and three of them are real translations rather than copies
// of the English with the keys renamed. A stub catalogue is worse than none: it
// makes the feature look finished, and the first person to switch locale finds
// out otherwise.
//
// English, Spanish, French and Swahili. Swahili because the scaffold is built in
// Uganda and East Africa is not a rounding error, and because a framework whose
// only second language is French has quietly decided who it is for. Spanish for
// the same reason at a larger scale: it is the second most spoken first
// language in the world, and a project that has to add it by hand is one where
// the feature is present but the audience it names is not.
//
// Keys are dotted and grouped by area. Flat JSON rather than nested, because a
// flat map needs no traversal code and a translator editing it cannot break the
// shape by mis-indenting a brace.
//
// A message that takes a resource or field name as %s is worded so the sentence
// stays grammatical whatever that name is. Spanish and French inflect the
// participle and the article, so "%s creado" renders "Factura creado" and
// "Ese %s ya está en uso" renders "Ese Factura". The fix is in the wording, not
// in a lookup table: put the agreement on a word the catalogue owns ("Se creó
// %s", whose impersonal "se" leaves no noun to agree with) or on a fixed noun
// ("El campo %s es obligatorio", "La valeur de %s est déjà utilisée"). A
// gendered catalogue would need the grammatical gender of every resource, and
// no resource definition carries one.
//
// Swahili agrees by noun class rather than gender. The class-9 forms here cover
// most resource names and are wrong for a class-1 one (a person), which is the
// same missing information in a different grammar.

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

  "errors.account_disabled": "This account has been deactivated",
  "errors.invalid_backup_code": "That recovery code is not correct",
  "errors.sign_in_expired": "That sign-in has expired. Start again",
  "errors.sign_in_failed": "Your sign-in could not be completed. Try again",
  "errors.request_not_verified": "That request could not be verified. Reload the page and try again",
  "errors.totp_already_enabled": "Two-factor authentication is already on",
  "errors.totp_failed": "Two-factor setup could not be completed. Try again",
  "errors.mail_failed": "That email could not be sent. Try again",
  "errors.invalid_file_type": "That kind of file is not accepted here",
  "errors.invalid_file": "That file could not be read",
  "errors.upload_failed": "That upload did not finish. Try again",
  "errors.upload_not_found": "That file could not be found",
  "errors.storage_unavailable": "File storage is not available right now",
  "errors.service_unavailable": "That is not available right now. Try again shortly",
  "errors.ai_unavailable": "The assistant is not available right now",

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

func i18nLocaleES() string {
	return `{
  "errors.not_found": "No encontrado",
  "errors.resource_not_found": "No se encontró %s",
  "errors.unauthorized": "Tienes que iniciar sesión para hacer eso",
  "errors.forbidden": "No tienes permiso para hacer eso",
  "errors.validation": "Algunos datos no son correctos",
  "errors.internal": "Algo salió mal por nuestra parte",
  "errors.invalid_request": "No se pudo entender esa solicitud",
  "errors.rate_limited": "Demasiadas solicitudes. Inténtalo de nuevo en un momento",
  "errors.conflict": "Eso entra en conflicto con algo que ya existe",
  "errors.payload_too_large": "Ese archivo es demasiado grande",

  "errors.account_disabled": "Esta cuenta ha sido desactivada",
  "errors.invalid_backup_code": "Ese código de recuperación no es correcto",
  "errors.sign_in_expired": "Ese inicio de sesión ha caducado. Empieza otra vez",
  "errors.sign_in_failed": "No se pudo completar el inicio de sesión. Inténtalo otra vez",
  "errors.request_not_verified": "No se pudo verificar la solicitud. Vuelve a cargar la página e inténtalo otra vez",
  "errors.totp_already_enabled": "La verificación en dos pasos ya está activada",
  "errors.totp_failed": "No se pudo completar la configuración en dos pasos. Inténtalo otra vez",
  "errors.mail_failed": "No se pudo enviar el correo. Inténtalo otra vez",
  "errors.invalid_file_type": "Ese tipo de archivo no se acepta aquí",
  "errors.invalid_file": "No se pudo leer ese archivo",
  "errors.upload_failed": "La subida no se completó. Inténtalo otra vez",
  "errors.upload_not_found": "No se encontró ese archivo",
  "errors.storage_unavailable": "El almacenamiento de archivos no está disponible ahora mismo",
  "errors.service_unavailable": "Eso no está disponible ahora mismo. Inténtalo en un momento",
  "errors.ai_unavailable": "El asistente no está disponible ahora mismo",

  "auth.invalid_credentials": "Ese correo y esa contraseña no coinciden",
  "auth.account_locked": "Esta cuenta está bloqueada. Inténtalo más tarde",
  "auth.email_taken": "Ya existe una cuenta con ese correo",
  "auth.email_not_verified": "Primero verifica tu correo electrónico",
  "auth.token_expired": "Tu sesión ha caducado. Vuelve a iniciar sesión",
  "auth.token_invalid": "Ese enlace de inicio de sesión no es válido",
  "auth.registered": "Cuenta creada",
  "auth.signed_in": "Sesión iniciada",
  "auth.signed_out": "Sesión cerrada",
  "auth.password_reset_sent": "Si ese correo existe, un enlace de restablecimiento está en camino",
  "auth.password_changed": "Contraseña cambiada. Se ha cerrado la sesión en los demás dispositivos",
  "auth.totp_required": "Introduce el código de tu aplicación de autenticación",
  "auth.totp_invalid": "Ese código no es correcto",

  "validation.required": "El campo %s es obligatorio",
  "validation.email": "%s debe ser una dirección de correo válida",
  "validation.min_length": "%s debe tener al menos %d caracteres",
  "validation.max_length": "%s debe tener como máximo %d caracteres",
  "validation.unique": "El valor de %s ya está en uso",
  "validation.numeric": "%s debe ser un número",
  "validation.positive": "%s debe ser mayor que cero",

  "resource.created": "Se creó %s",
  "resource.updated": "Se actualizó %s",
  "resource.deleted": "Se eliminó %s",
  "resource.restored": "Se restauró %s",

  "common.yes": "Sí",
  "common.no": "No",
  "common.saved": "Guardado",
  "common.cancelled": "Cancelado"
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

  "errors.account_disabled": "Ce compte a été désactivé",
  "errors.invalid_backup_code": "Ce code de secours n'est pas correct",
  "errors.sign_in_expired": "Cette connexion a expiré. Recommencez",
  "errors.sign_in_failed": "La connexion n'a pas pu aboutir. Réessayez",
  "errors.request_not_verified": "Cette requête n'a pas pu être vérifiée. Rechargez la page et réessayez",
  "errors.totp_already_enabled": "La double authentification est déjà activée",
  "errors.totp_failed": "La configuration de la double authentification n'a pas abouti. Réessayez",
  "errors.mail_failed": "Cet e-mail n'a pas pu être envoyé. Réessayez",
  "errors.invalid_file_type": "Ce type de fichier n'est pas accepté ici",
  "errors.invalid_file": "Ce fichier n'a pas pu être lu",
  "errors.upload_failed": "L'envoi n'est pas allé au bout. Réessayez",
  "errors.upload_not_found": "Ce fichier est introuvable",
  "errors.storage_unavailable": "Le stockage de fichiers n'est pas disponible pour le moment",
  "errors.service_unavailable": "Ce service n'est pas disponible pour le moment. Réessayez dans un instant",
  "errors.ai_unavailable": "L'assistant n'est pas disponible pour le moment",

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
  "validation.unique": "La valeur de %s est déjà utilisée",
  "validation.numeric": "%s doit être un nombre",
  "validation.positive": "%s doit être supérieur à zéro",

  "resource.created": "Création de %s",
  "resource.updated": "Mise à jour de %s",
  "resource.deleted": "Suppression de %s",
  "resource.restored": "Restauration de %s",

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

  "errors.account_disabled": "Akaunti hii imezimwa",
  "errors.invalid_backup_code": "Msimbo huo wa dharura si sahihi",
  "errors.sign_in_expired": "Kuingia huko kumeisha muda. Anza tena",
  "errors.sign_in_failed": "Kuingia hakukukamilika. Jaribu tena",
  "errors.request_not_verified": "Ombi hilo halikuthibitishwa. Pakia ukurasa upya na ujaribu tena",
  "errors.totp_already_enabled": "Uthibitishaji wa hatua mbili umewashwa tayari",
  "errors.totp_failed": "Uwekaji wa hatua mbili haukukamilika. Jaribu tena",
  "errors.mail_failed": "Barua pepe hiyo haikutumwa. Jaribu tena",
  "errors.invalid_file_type": "Aina hiyo ya faili haikubaliki hapa",
  "errors.invalid_file": "Faili hilo halikusomeka",
  "errors.upload_failed": "Upakiaji haukukamilika. Jaribu tena",
  "errors.upload_not_found": "Faili hilo halipatikani",
  "errors.storage_unavailable": "Hifadhi ya mafaili haipatikani kwa sasa",
  "errors.service_unavailable": "Hiyo haipatikani kwa sasa. Jaribu tena baadaye kidogo",
  "errors.ai_unavailable": "Msaidizi hapatikani kwa sasa",

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

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"{{MODULE}}/internal/respond"
)

func TestCataloguesResolve(t *testing.T) {
	if err := Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	for _, loc := range []string{"en", "es", "fr", "sw"} {
		if !Has(loc) {
			t.Fatalf("locale %q not compiled in; Available=%v", loc, Available())
		}
	}
	if got := T("fr", "errors.not_found"); got != "Introuvable" {
		t.Fatalf("fr errors.not_found = %q", got)
	}
	if got := T("es", "errors.not_found"); got != "No encontrado" {
		t.Fatalf("es errors.not_found = %q", got)
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

// TestCatalogueMarkersMatch holds every catalogue to the printf verbs its
// English message has. Keys alone are not enough: a translation that drops a
// %s keeps its key, compiles, and prints %!s(MISSING) to whoever reads it.
func TestCatalogueMarkersMatch(t *testing.T) {
	if err := Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	base := messages[Default]
	for locale, m := range messages {
		if locale == Default {
			continue
		}
		for key, translated := range m {
			got, want := verbs(translated), verbs(base[key])
			if got != want {
				t.Errorf("%s %q renders with %q, %s has %q", locale, key, got, Default, want)
			}
		}
	}
}

// verbs lists the printf verbs in a message, sorted, so two translations of one
// key compare equal however they order the sentence a %s sits in.
func verbs(message string) string {
	var found []string
	for _, v := range printfVerb.FindAllString(message, -1) {
		if v != "%%" {
			found = append(found, v)
		}
	}
	sort.Strings(found)
	return strings.Join(found, " ")
}

var printfVerb = regexp.MustCompile("%[-+# 0]*[0-9]*(\\.[0-9]+)?[A-Za-z%]")

// The framework's own error messages come out translated.
//
// This is the end of the thing i18n was half of until v3.545.0: the
// catalogues were compiled in and the frontend read them, and every message
// the API itself refused a request with was English, because all 441
// respond.Fail calls pass an English string and nothing translated it.
//
// respond now asks a translator for the general message of the code it was
// given, and internal/i18n/respond.go registers one. So this test is about
// the registration as much as the catalogue: an init that does not run, in a
// package nothing imports, is the way this silently would not work.
func TestTheFrameworksOwnErrorsAreTranslated(t *testing.T) {
	if respond.Translator == nil {
		t.Fatal("nothing registered a translator, so every API error is English " +
			"however many catalogues are compiled in")
	}

	gin.SetMode(gin.TestMode)

	cases := []struct {
		locale string
		code   respond.Code
		want   string
	}{
		{"es", respond.CodeForbidden, "No tienes permiso para hacer eso"},
		{"fr", respond.CodeForbidden, "Vous n'avez pas la permission de faire cela"},
		{"sw", respond.CodeNotFound, "Haipatikani"},
		{"es", respond.CodeValidationError, "Algunos datos no son correctos"},
	}
	for _, tc := range cases {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
		c.Set("locale", tc.locale)

		respond.Fail(c, tc.code, "the handler's own English message")

		var body struct {
			Error struct {
				Code    string ` + "`json:\"code\"`" + `
				Message string ` + "`json:\"message\"`" + `
			} ` + "`json:\"error\"`" + `
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("%s %s: %v\n%s", tc.locale, tc.code, err, rec.Body.String())
		}
		if body.Error.Code != string(tc.code) {
			t.Errorf("%s: code came back %q", tc.locale, body.Error.Code)
		}
		if body.Error.Message != tc.want {
			t.Errorf("%s %s: message %q, want %q",
				tc.locale, tc.code, body.Error.Message, tc.want)
		}
	}
}

// A code with no translation keeps the handler's message, in every locale.
//
// This is what makes adding a locale additive. There are 124 codes and a
// translation for the ones a person reads; the rest have to come out as
// English rather than as a key, or every locale would owe a translation for
// an API key rejection nobody will ever read.
func TestAnUntranslatedCodeKeepsItsEnglishMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	if respond.MessageKey(respond.CodeInvalidSignature) != "" {
		t.Skip("INVALID_SIGNATURE has a translation now; pick another untranslated code")
	}

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Set("locale", "es")

	respond.Fail(c, respond.CodeInvalidSignature, "That signature does not match")

	if !strings.Contains(rec.Body.String(), "That signature does not match") {
		t.Errorf("an untranslated code lost its message: %s", rec.Body.String())
	}
}

// English is untouched, which is the line that made this safe to ship.
//
// A handler's own message often says more than the code's general one: "Only
// an owner can archive an invoice" against "You do not have permission to do
// that". Translating the default locale as well would have discarded the
// better sentence at all 441 call sites at once, silently.
func TestEnglishStillGetsTheHandlersOwnMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Set("locale", "en")

	respond.Fail(c, respond.CodeForbidden, "Only an owner can archive an invoice")

	if !strings.Contains(rec.Body.String(), "Only an owner can archive an invoice") {
		t.Errorf("English lost the handler's specific message: %s", rec.Body.String())
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
