package cancom

import (
	"testing"
)

func TestProvider(t *testing.T) {
	p := Provider()

	if err := p.InternalValidate(); err != nil {
		t.Fatalf("provider internal validation failed: %v", err)
	}

	tokenSchema, ok := p.Schema["token"]
	if !ok {
		t.Fatal("expected token schema attribute")
	}
	if !tokenSchema.Sensitive {
		t.Error("expected provider token attribute to be marked Sensitive")
	}
}
