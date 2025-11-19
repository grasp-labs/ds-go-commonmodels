package errors_test

import (
	"testing"

	e "github.com/grasp-labs/ds-go-commonmodels/v3/commonmodels/enum/errors"
)

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
