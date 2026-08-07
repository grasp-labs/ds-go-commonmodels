package errors_test

import (
	"strings"
	"testing"

	e "github.com/grasp-labs/ds-go-commonmodels/v3/commonmodels/enum/errors"
)

// Every language tag Accept-Language negotiation can produce for these two
// languages must reach the right catalog. A service defaulting to "no" once fell
// through to English here while its custom messages came back in Norwegian.
func TestHumanMessageLocale_AcceptsLanguageTags(t *testing.T) {
	norwegian := e.HumanMessageLocale("nb", e.Required, "e-post")
	for _, locale := range []string{"no", "nb", "nn", "nb-NO", "no_NO", "NB", " nb "} {
		if got := e.HumanMessageLocale(locale, e.Required, "e-post"); got != norwegian {
			t.Errorf("locale %q: got %q, want the Norwegian message %q", locale, got, norwegian)
		}
	}

	english := e.HumanMessageLocale("en", e.Required, "email")
	for _, locale := range []string{"en", "en-US", "EN", "", "pl", "de-DE"} {
		if got := e.HumanMessageLocale(locale, e.Required, "email"); got != english {
			t.Errorf("locale %q: got %q, want the English message %q", locale, got, english)
		}
	}
}

func TestCustomHumanMessageLocale_AgreesWithCatalogLocales(t *testing.T) {
	c := e.CustomMessage{En: "Hello", No: "Hei"}
	for _, locale := range []string{"no", "nb", "nb-NO", "nn"} {
		if got := e.CustomHumanMessageLocale(locale, c); got != c.No {
			t.Errorf("locale %q: got %q, want %q", locale, got, c.No)
		}
	}
	for _, locale := range []string{"en", "en-GB", "pl", ""} {
		if got := e.CustomHumanMessageLocale(locale, c); got != c.En {
			t.Errorf("locale %q: got %q, want %q", locale, got, c.En)
		}
	}
	if got := e.CustomHumanMessageLocale("nb", e.CustomMessage{En: "Hello"}); got != "Hello" {
		t.Errorf("a missing translation should fall back to English, got %q", got)
	}
}

// A template and its call sites drift. Whatever the mismatch, the client must
// never see fmt's complaints about it.
func TestHumanMessageLocale_NeverLeaksFormatNoise(t *testing.T) {
	cases := map[string]string{
		"missing arg":          e.HumanMessageLocale("en", e.Required),
		"missing arg nb":       e.HumanMessageLocale("nb", e.Required),
		"extra arg":            e.HumanMessageLocale("en", e.ValidationFailed, "memory_mb"),
		"extra arg nb":         e.HumanMessageLocale("nb", e.ValidationFailed, "memory_mb"),
		"too many args":        e.HumanMessageLocale("en", e.Required, "id", "name"),
		"status with arg":      e.HumanMessageLocale("nb", e.InvalidStatus, "bogus"),
		"unknown code":         e.HumanMessageLocale("en", "no_such_code"),
		"unknown code untrans": e.HumanMessageLocale("nb", "no_such_code"),
	}
	for name, msg := range cases {
		for _, noise := range []string{"%s", "%!", "%d", "%v", "EXTRA", "MISSING", "NOVERB"} {
			if strings.Contains(msg, noise) {
				t.Errorf("%s: message %q contains %q", name, msg, noise)
			}
		}
		if msg == "" {
			t.Errorf("%s: message is empty", name)
		}
	}
}

func TestHumanMessageLocale_FillsMissingArgWithASubject(t *testing.T) {
	if got := e.HumanMessageLocale("en", e.Required); got != "This field is required." {
		t.Errorf("got %q", got)
	}
	if got := e.HumanMessageLocale("nb", e.Required); got != "Dette feltet er påkrevd." {
		t.Errorf("got %q", got)
	}
}

func TestCustomMessageLocale(t *testing.T) {
	t.Run("test hit", func(t *testing.T) {
		c := e.CustomMessage{
			En: "Hello",
			No: "Hei",
		}
		msg := e.CustomHumanMessageLocale("no", c)
		if msg != c.No {
			t.Fatalf("expected %s, got %s", c.No, msg)
		}
	})
}

func TestCustomMessageLocale_InvalidRef(t *testing.T) {
	t.Run("test hit", func(t *testing.T) {
		c := e.CustomMessage{
			En: "Hello",
			No: "Hei",
		}
		msg := e.CustomHumanMessageLocale("pl", c)
		if msg != c.En {
			t.Fatalf("expected %s, got %s", c.En, msg)
		}
	})
}
