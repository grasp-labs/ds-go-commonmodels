package errors

import (
	"strings"
	"testing"
)

// The catalogs must agree on placeholders. A translation that drops a %s makes
// every caller passing an argument emit fmt noise, but only in that locale, which
// is exactly the kind of bug that survives for months unnoticed.
func TestCatalogs_PlaceholderParity(t *testing.T) {
	for code, en := range messagesEN {
		nb, ok := messagesNB[code]
		if !ok {
			t.Errorf("code %q has no Norwegian translation", code)
			continue
		}
		if got, want := countVerbs(nb), countVerbs(en); got != want {
			t.Errorf("code %q: Norwegian has %d placeholders, English has %d\n  en=%q\n  nb=%q", code, got, want, en, nb)
		}
	}
	for code := range messagesNB {
		if _, ok := messagesEN[code]; !ok {
			t.Errorf("code %q is translated but missing from the English catalog", code)
		}
	}
}

// Every catalog entry must have a status mapping, or StatusFor silently answers
// 500 for a code whose message says something else entirely.
func TestCatalogs_EveryCodeHasAStatus(t *testing.T) {
	for code := range messagesEN {
		if _, ok := statusByCode[code]; !ok {
			t.Errorf("code %q has a message but no HTTP status mapping", code)
		}
	}
}

func TestCountVerbs(t *testing.T) {
	cases := map[string]int{
		"":                       0,
		"no verbs here.":         0,
		"%s is required.":        1,
		"%s and %s":              2,
		"100%% done":             0,
		"100%% done, %s":         1,
		"%d items":               1,
		"The provided %s is bad": 1,
	}
	for msg, want := range cases {
		if got := countVerbs(msg); got != want {
			t.Errorf("countVerbs(%q) = %d, want %d", msg, got, want)
		}
	}
}

func TestNormalizeLocale(t *testing.T) {
	cases := map[string]string{
		"en": LocaleEN, "en-US": LocaleEN, "EN": LocaleEN, "eng": LocaleEN,
		"no": LocaleNB, "nb": LocaleNB, "nn": LocaleNB, "nb-NO": LocaleNB, "no_NO": LocaleNB, " NB ": LocaleNB,
		"": "", "pl": "", "de-DE": "", "x": "",
	}
	for locale, want := range cases {
		if got := normalizeLocale(locale); got != want {
			t.Errorf("normalizeLocale(%q) = %q, want %q", locale, got, want)
		}
	}
}

// A locale key must have a generic subject to fall back on, or fill() would put
// an English word into a Norwegian sentence.
func TestGenericSubject_CoversEveryCatalog(t *testing.T) {
	for key := range catalogs {
		if strings.TrimSpace(genericSubject[key]) == "" {
			t.Errorf("catalog %q has no generic subject for a missing argument", key)
		}
	}
}
