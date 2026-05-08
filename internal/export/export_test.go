package export

import (
	"strings"
	"testing"
)

const sampleEnv = `# Database config
DB_HOST=localhost
DB_PORT=5432
DB_NAME="myapp"

# App
APP_SECRET=supersecret
APP_ENV=production
`

func TestWrite_ShellFormat(t *testing.T) {
	var buf strings.Builder
	err := Write(&buf, sampleEnv, Options{Format: FormatShell})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `export DB_HOST="localhost"`) {
		t.Errorf("expected shell export for DB_HOST, got:\n%s", out)
	}
	if !strings.Contains(out, `export APP_SECRET="supersecret"`) {
		t.Errorf("expected shell export for APP_SECRET, got:\n%s", out)
	}
}

func TestWrite_DotenvFormat(t *testing.T) {
	var buf strings.Builder
	err := Write(&buf, sampleEnv, Options{Format: FormatDotenv})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_HOST=localhost") {
		t.Errorf("expected dotenv line for DB_HOST, got:\n%s", out)
	}
	if strings.Contains(out, "# Database") {
		t.Errorf("comments should be stripped, got:\n%s", out)
	}
}

func TestWrite_JSONFormat(t *testing.T) {
	var buf strings.Builder
	err := Write(&buf, sampleEnv, Options{Format: FormatJSON})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"DB_HOST": "localhost"`) {
		t.Errorf("expected JSON entry for DB_HOST, got:\n%s", out)
	}
	if !strings.HasPrefix(out, "{") || !strings.HasSuffix(strings.TrimSpace(out), "}") {
		t.Errorf("expected valid JSON object braces, got:\n%s", out)
	}
}

func TestWrite_FilterKeys(t *testing.T) {
	var buf strings.Builder
	err := Write(&buf, sampleEnv, Options{
		Format: FormatDotenv,
		Keys:   []string{"DB_HOST", "APP_ENV"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output")
	}
	if strings.Contains(out, "APP_SECRET") {
		t.Errorf("APP_SECRET should be filtered out")
	}
	if strings.Contains(out, "DB_PORT") {
		t.Errorf("DB_PORT should be filtered out")
	}
}

func TestWrite_UnknownFormat(t *testing.T) {
	var buf strings.Builder
	err := Write(&buf, sampleEnv, Options{Format: Format("xml")})
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestWrite_EmptyEnv(t *testing.T) {
	var buf strings.Builder
	err := Write(&buf, "", Options{Format: FormatShell})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.String() != "" {
		t.Errorf("expected empty output for empty env, got: %q", buf.String())
	}
}
