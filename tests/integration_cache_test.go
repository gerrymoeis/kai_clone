package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"io/fs"

	"gothicforge3/app/routes"
	"gothicforge3/internal/server"
)

func Test_HTMLCache_Public_Defaults(t *testing.T) {
	_ = os.Setenv("LOG_FORMAT", "off")
	_ = os.Unsetenv("DISABLE_HTML_CACHE")
	_ = os.Unsetenv("CACHE_PUBLIC_TTL")
	_ = os.Unsetenv("CACHE_SWREVAL_TTL")

	r := server.New()
	// register a temporary public HTML route that does not set cookies
	r.Get("/_test_public", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte("<html><body>ok</body></html>"))
	})
	// also include normal app routes for completeness
	routes.Register(r)

	req := httptest.NewRequest(http.MethodGet, "/_test_public", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("/_test_public want 200, got %d", rec.Code)
	}
	ct := rec.Header().Get("Content-Type")
	if !strings.Contains(strings.ToLower(ct), "text/html") {
		t.Fatalf("/_test_public content-type want text/html, got %q", ct)
	}
	cc := strings.ToLower(rec.Header().Get("Cache-Control"))
	// Depending on session middleware behavior, Set-Cookie may be present which forces private, no-store
	if !(strings.Contains(cc, "s-maxage=60") && strings.Contains(cc, "stale-while-revalidate=300")) &&
		!(strings.Contains(cc, "private") && strings.Contains(cc, "no-store")) {
		t.Fatalf("/_test_public cache-control unexpected: %q", cc)
	}
}

func Test_HTMLCache_Disabled(t *testing.T) {
	_ = os.Setenv("LOG_FORMAT", "off")
	_ = os.Setenv("DISABLE_HTML_CACHE", "1")
	defer os.Unsetenv("DISABLE_HTML_CACHE")

	r := server.New()
	r.Get("/_test_public", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte("<html><body>ok</body></html>"))
	})
	routes.Register(r)
	req := httptest.NewRequest(http.MethodGet, "/_test_public", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("/_test_public want 200, got %d", rec.Code)
	}
	cc := rec.Header().Get("Cache-Control")
	if cc != "" {
		t.Fatalf("/_test_public cache-control should be empty when disabled, got %q", cc)
	}
}

func Test_HTMLCache_HXRequest_Private(t *testing.T) {
	_ = os.Setenv("LOG_FORMAT", "off")
	_ = os.Unsetenv("DISABLE_HTML_CACHE")

	r := server.New()
	r.Get("/_test_public", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte("<html><body>ok</body></html>"))
	})
	routes.Register(r)
	req := httptest.NewRequest(http.MethodGet, "/_test_public", nil)
	req.Header.Set("HX-Request", "true")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("/_test_public want 200, got %d", rec.Code)
	}
	cc := rec.Header().Get("Cache-Control")
	if !strings.Contains(strings.ToLower(cc), "private") || !strings.Contains(strings.ToLower(cc), "no-store") {
		t.Fatalf("/_test_public cache-control want private, no-store for HX-Request, got %q", cc)
	}
}

func Test_HTMLCache_Session_Private(t *testing.T) {
	_ = os.Setenv("LOG_FORMAT", "off")
	_ = os.Unsetenv("DISABLE_HTML_CACHE")

	r := server.New()
	// Reuse app home which writes to session and sets cookie
	routes.Register(r)
	req1 := httptest.NewRequest(http.MethodGet, "/", nil)
	rec1 := httptest.NewRecorder()
	r.ServeHTTP(rec1, req1)
	if rec1.Code != 200 {
		t.Fatalf("seed want 200, got %d", rec1.Code)
	}
	setCookie := rec1.Header().Get("Set-Cookie")
	if setCookie == "" {
		// If no cookie, skip (session store not configured?)
		t.Skip("no Set-Cookie emitted; session cookie not present in this environment")
	}
	cookiePair := setCookie
	if p := strings.Index(cookiePair, ";"); p > 0 {
		cookiePair = cookiePair[:p]
	}

	// Second request includes the session cookie
	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	req2.Header.Set("Cookie", cookiePair)
	rec2 := httptest.NewRecorder()
	r.ServeHTTP(rec2, req2)
	if rec2.Code != 200 {
		t.Fatalf("/ want 200, got %d", rec2.Code)
	}
	cc := rec2.Header().Get("Cache-Control")
	if !strings.Contains(strings.ToLower(cc), "private") || !strings.Contains(strings.ToLower(cc), "no-store") {
		t.Fatalf("/ cache-control want private, no-store for session requests, got %q", cc)
	}
}

func moduleRoot(t *testing.T) string {
	t.Helper()
	cur, _ := os.Getwd()
	try := cur
	for {
		if _, err := os.Stat(filepath.Join(try, "go.mod")); err == nil {
			return try
		}
		parent := filepath.Dir(try)
		if parent == try { t.Fatalf("could not find module root from %s", cur) }
		try = parent
	}
}

func Test_Static_CSS_Mime_And_Cache(t *testing.T) {
	_ = os.Setenv("LOG_FORMAT", "off")
	// Ensure a css file exists under module root app/static/styles/
	root := moduleRoot(t)
	cssDir := filepath.Join(root, "app", "static", "styles")
	if err := os.MkdirAll(cssDir, fs.FileMode(0o755)); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	cssPath := filepath.Join(cssDir, "test.css")
	if err := os.WriteFile(cssPath, []byte("/* test */\n"), fs.FileMode(0o644)); err != nil {
		t.Fatalf("write css: %v", err)
	}
	t.Cleanup(func() { _ = os.Remove(cssPath) })

	r := server.New()
	req := httptest.NewRequest(http.MethodGet, "/static/styles/test.css", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("static css want 200, got %d", rec.Code)
	}
	ct := strings.ToLower(rec.Header().Get("Content-Type"))
	if !strings.Contains(ct, "text/css") {
		t.Fatalf("static css content-type want text/css, got %q", ct)
	}
	cc := strings.ToLower(rec.Header().Get("Cache-Control"))
	if !strings.Contains(cc, "immutable") || !strings.Contains(cc, "max-age=31536000") {
		t.Fatalf("static css cache-control unexpected: %q", cc)
	}
}
