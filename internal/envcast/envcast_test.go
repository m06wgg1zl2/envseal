package envcast_test

import (
	"testing"

	"github.com/yourorg/envseal/internal/envcast"
)

func TestToString_ReturnsValue(t *testing.T) {
	env := map[string]string{"APP_NAME": "myapp"}
	if got := envcast.ToString(env, "APP_NAME", "default"); got != "myapp" {
		t.Errorf("expected myapp, got %q", got)
	}
}

func TestToString_ReturnsDefault(t *testing.T) {
	env := map[string]string{}
	if got := envcast.ToString(env, "MISSING", "fallback"); got != "fallback" {
		t.Errorf("expected fallback, got %q", got)
	}
}

func TestToInt_Valid(t *testing.T) {
	env := map[string]string{"PORT": "8080"}
	n, err := envcast.ToInt(env, "PORT", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 8080 {
		t.Errorf("expected 8080, got %d", n)
	}
}

func TestToInt_Invalid(t *testing.T) {
	env := map[string]string{"PORT": "abc"}
	_, err := envcast.ToInt(env, "PORT", 0)
	if err == nil {
		t.Error("expected error for invalid int")
	}
}

func TestToInt_Missing_ReturnsDefault(t *testing.T) {
	env := map[string]string{}
	n, err := envcast.ToInt(env, "PORT", 3000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 3000 {
		t.Errorf("expected 3000, got %d", n)
	}
}

func TestToBool_TrueVariants(t *testing.T) {
	for _, val := range []string{"true", "1", "yes", "on", "TRUE"} {
		env := map[string]string{"FLAG": val}
		b, err := envcast.ToBool(env, "FLAG", false)
		if err != nil {
			t.Errorf("unexpected error for %q: %v", val, err)
		}
		if !b {
			t.Errorf("expected true for %q", val)
		}
	}
}

func TestToBool_FalseVariants(t *testing.T) {
	for _, val := range []string{"false", "0", "no", "off"} {
		env := map[string]string{"FLAG": val}
		b, err := envcast.ToBool(env, "FLAG", true)
		if err != nil {
			t.Errorf("unexpected error for %q: %v", val, err)
		}
		if b {
			t.Errorf("expected false for %q", val)
		}
	}
}

func TestToBool_Invalid(t *testing.T) {
	env := map[string]string{"FLAG": "maybe"}
	_, err := envcast.ToBool(env, "FLAG", false)
	if err == nil {
		t.Error("expected error for invalid bool")
	}
}

func TestToFloat64_Valid(t *testing.T) {
	env := map[string]string{"RATIO": "3.14"}
	f, err := envcast.ToFloat64(env, "RATIO", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f != 3.14 {
		t.Errorf("expected 3.14, got %f", f)
	}
}

func TestToFloat64_Invalid(t *testing.T) {
	env := map[string]string{"RATIO": "notanumber"}
	_, err := envcast.ToFloat64(env, "RATIO", 0)
	if err == nil {
		t.Error("expected error for invalid float")
	}
}

func TestToStringSlice_CommaSeparated(t *testing.T) {
	env := map[string]string{"HOSTS": "a.com, b.com, c.com"}
	slice := envcast.ToStringSlice(env, "HOSTS", ",")
	if len(slice) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(slice))
	}
	if slice[1] != "b.com" {
		t.Errorf("expected b.com, got %q", slice[1])
	}
}

func TestToStringSlice_Missing(t *testing.T) {
	env := map[string]string{}
	if s := envcast.ToStringSlice(env, "HOSTS", ","); s != nil {
		t.Errorf("expected nil slice, got %v", s)
	}
}

func TestValidate_CollectsErrors(t *testing.T) {
	env := map[string]string{"PORT": "bad", "FLAG": "maybe"}
	_, portErr := envcast.ToInt(env, "PORT", 0)
	_, flagErr := envcast.ToBool(env, "FLAG", false)
	results := []envcast.Result{
		{Key: "PORT", Raw: "bad", Error: portErr},
		{Key: "FLAG", Raw: "maybe", Error: flagErr},
	}
	errs := envcast.Validate(results)
	if len(errs) != 2 {
		t.Errorf("expected 2 errors, got %d", len(errs))
	}
}
