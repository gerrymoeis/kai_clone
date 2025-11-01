package cmd

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var genEdgeCmd = &cobra.Command{
	Use:   "gen-edge",
	Short: "Generate Cloudflare Pages Functions from annotated Go handlers",
	Long: `Scans app/routes/ for Go handlers marked with //gforge:edge comments
and generates equivalent JavaScript functions in functions/ directory.

Example annotation:
  //gforge:edge path=/api/users method=POST ttl=60s
  func CreateUser(w http.ResponseWriter, r *http.Request) { ... }

This generates: functions/api/users.js with onRequestPost handler`,
	RunE: func(cmd *cobra.Command, args []string) error {
		banner()
		return generateEdgeFunctions()
	},
}

func init() {
	rootCmd.AddCommand(genEdgeCmd)
}

type edgeAnnotation struct {
	Path       string
	Method     string
	TTL        string
	Auth       bool
	CORS       bool
	SourceFile string
	FuncName   string
	Comment    string
}

func generateEdgeFunctions() error {
	fmt.Println("🔍 Scanning for edge function annotations...")
	fmt.Println()

	routesDir := "app/routes"
	annotations, err := scanForEdgeAnnotations(routesDir)
	if err != nil {
		return fmt.Errorf("failed to scan routes: %w", err)
	}

	if len(annotations) == 0 {
		fmt.Println("❌ No edge function annotations found")
		fmt.Println()
		fmt.Println("Add annotations to your route handlers:")
		fmt.Println("  //gforge:edge path=/api/hello method=GET")
		fmt.Println("  func HandleHello(w http.ResponseWriter, r *http.Request) {")
		fmt.Println("    // your code")
		fmt.Println("  }")
		fmt.Println()
		return nil
	}

	fmt.Printf("✅ Found %d edge function annotation(s)\n", len(annotations))
	fmt.Println()

	functionsDir := "functions"
	if err := os.MkdirAll(functionsDir, 0755); err != nil {
		return fmt.Errorf("failed to create functions directory: %w", err)
	}

	generated := 0
	for _, ann := range annotations {
		if err := generateJSFunction(ann, functionsDir); err != nil {
			fmt.Printf("⚠️  Warning: failed to generate %s: %v\n", ann.Path, err)
			continue
		}
		generated++
	}

	fmt.Println()
	fmt.Printf("🎉 Generated %d edge function(s) in %s/\n", generated, functionsDir)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Review generated functions in functions/")
	fmt.Println("  2. Test locally: gforge dev")
	fmt.Println("  3. Deploy: gforge export && gforge deploy pages --run")
	fmt.Println()

	return nil
}

func scanForEdgeAnnotations(dir string) ([]edgeAnnotation, error) {
	var annotations []edgeAnnotation

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		fileAnnotations, err := parseFileForAnnotations(path)
		if err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
		annotations = append(annotations, fileAnnotations...)
		return nil
	})

	return annotations, err
}

func parseFileForAnnotations(filename string) ([]edgeAnnotation, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var annotations []edgeAnnotation

	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// Check if function has edge annotation in doc comments
		if funcDecl.Doc == nil {
			continue
		}

		for _, comment := range funcDecl.Doc.List {
			text := strings.TrimSpace(comment.Text)
			if !strings.HasPrefix(text, "//gforge:edge") {
				continue
			}

			ann := parseEdgeComment(text)
			if ann.Path == "" {
				continue
			}

			ann.SourceFile = filename
			ann.FuncName = funcDecl.Name.Name
			ann.Comment = text

			annotations = append(annotations, ann)
		}
	}

	return annotations, nil
}

func parseEdgeComment(comment string) edgeAnnotation {
	// Parse: //gforge:edge path=/api/users method=POST ttl=60s auth=true cors=true
	ann := edgeAnnotation{
		Method: "GET", // default
		Auth:   false,
		CORS:   true, // default enabled
	}

	// Remove //gforge:edge prefix
	parts := strings.Fields(comment)
	for _, part := range parts[1:] { // Skip "//gforge:edge"
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(kv[0]))
		val := strings.TrimSpace(kv[1])

		switch key {
		case "path":
			ann.Path = val
		case "method":
			ann.Method = strings.ToUpper(val)
		case "ttl":
			ann.TTL = val
		case "auth":
			ann.Auth = val == "true"
		case "cors":
			ann.CORS = val == "true"
		}
	}

	return ann
}

func generateJSFunction(ann edgeAnnotation, baseDir string) error {
	// Convert path to file path
	// /api/users/[id] -> api/users/[id].js
	// /api/hello -> api/hello.js

	path := strings.TrimPrefix(ann.Path, "/")
	if path == "" {
		return fmt.Errorf("empty path")
	}

	// Create directory structure
	parts := strings.Split(path, "/")
	var filePath string

	if len(parts) == 1 {
		// Root level: /hello -> hello.js
		filePath = filepath.Join(baseDir, parts[0]+".js")
	} else {
		// Nested: /api/users -> api/users.js
		dir := filepath.Join(baseDir, filepath.Join(parts[:len(parts)-1]...))
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		filePath = filepath.Join(dir, parts[len(parts)-1]+".js")
	}

	// Generate JavaScript function
	js := generateJavaScript(ann)

	// Check if file exists
	if _, err := os.Stat(filePath); err == nil {
		// File exists, backup
		backupPath := filePath + ".bak"
		if err := os.Rename(filePath, backupPath); err != nil {
			fmt.Printf("⚠️  Could not backup %s\n", filePath)
		}
	}

	if err := os.WriteFile(filePath, []byte(js), 0644); err != nil {
		return err
	}

	fmt.Printf("✅ Generated: %s (%s %s)\n", filePath, ann.Method, ann.Path)
	return nil
}

func generateJavaScript(ann edgeAnnotation) string {
	method := strings.ToLower(ann.Method)
	handlerName := fmt.Sprintf("onRequest%s", strings.Title(method))

	var js strings.Builder

	// Header comment
	js.WriteString(fmt.Sprintf(`/**
 * Cloudflare Pages Function
 * 
 * Auto-generated from: %s
 * Function: %s
 * Path: %s
 * Method: %s
 * 
 * Generated by: gforge gen-edge
 * 
 * IMPORTANT: This is a template. You need to implement the actual logic.
 * The Go function serves as a reference for the expected behavior.
 */

`, ann.SourceFile, ann.FuncName, ann.Path, ann.Method))

	// Main handler
	js.WriteString(fmt.Sprintf(`export async function %s(context) {
`, handlerName))

	// Auth check if needed
	if ann.Auth {
		js.WriteString(`  // TODO: Implement authentication
  // Check JWT token from cookies or headers
  // const token = context.request.headers.get('Authorization');
  // if (!token) {
  //   return new Response('Unauthorized', { status: 401 });
  // }

`)
	}

	// Request body parsing for POST/PUT/PATCH
	if method == "post" || method == "put" || method == "patch" {
		js.WriteString(fmt.Sprintf(`  try {
    // Parse request body
    const contentType = context.request.headers.get('Content-Type') || '';
    
    let data;
    if (contentType.includes('application/json')) {
      data = await context.request.json();
    } else if (contentType.includes('application/x-www-form-urlencoded')) {
      const formData = await context.request.formData();
      data = Object.fromEntries(formData);
    } else {
      data = await context.request.text();
    }

    // TODO: Implement your logic here
    // This should match the behavior of: %s
    console.log('Received data:', data);

    // Example response
    return new Response(
      JSON.stringify({ 
        success: true,
        message: 'Request processed',
        data: data
      }),
      {
        status: 200,
        headers: {
          'Content-Type': 'application/json',
`, ann.FuncName))
	} else {
		// GET/DELETE
		js.WriteString(fmt.Sprintf(`  try {
    // Extract URL parameters
    const url = new URL(context.request.url);
    const params = Object.fromEntries(url.searchParams);
    
    // Extract route parameters (if any)
    // const id = context.params.id; // For routes like /api/users/[id]

    // TODO: Implement your logic here
    // This should match the behavior of: %s
    console.log('Request params:', params);

    // Example response
    return new Response(
      JSON.stringify({ 
        success: true,
        message: 'Request processed',
        params: params
      }),
      {
        status: 200,
        headers: {
          'Content-Type': 'application/json',
`, ann.FuncName))
	}

	// CORS headers
	if ann.CORS {
		js.WriteString(fmt.Sprintf(`          'Access-Control-Allow-Origin': '*',
          'Access-Control-Allow-Methods': '%s, OPTIONS',
          'Access-Control-Allow-Headers': 'Content-Type, Authorization',
`, ann.Method))
	}

	// TTL cache header
	if ann.TTL != "" {
		ttlSeconds := parseTTL(ann.TTL)
		js.WriteString(fmt.Sprintf(`          'Cache-Control': 'public, max-age=%d',
`, ttlSeconds))
	}

	js.WriteString(`        },
      }
    );
  } catch (error) {
    console.error('Error:', error);
    return new Response(
      JSON.stringify({ 
        success: false, 
        error: error.message 
      }),
      {
        status: 500,
        headers: { 'Content-Type': 'application/json' },
      }
    );
  }
}
`)

	// OPTIONS handler for CORS
	if ann.CORS {
		js.WriteString(fmt.Sprintf(`
// Handle CORS preflight
export async function onRequestOptions(context) {
  return new Response(null, {
    status: 204,
    headers: {
      'Access-Control-Allow-Origin': '*',
      'Access-Control-Allow-Methods': '%s, OPTIONS',
      'Access-Control-Allow-Headers': 'Content-Type, Authorization',
      'Access-Control-Max-Age': '86400',
    },
  });
}
`, ann.Method))
	}

	return js.String()
}

func parseTTL(ttl string) int {
	// Parse ttl like "60s", "5m", "1h" to seconds
	re := regexp.MustCompile(`^(\d+)([smh])$`)
	matches := re.FindStringSubmatch(ttl)
	if matches == nil {
		return 0
	}

	num := 0
	fmt.Sscanf(matches[1], "%d", &num)

	switch matches[2] {
	case "s":
		return num
	case "m":
		return num * 60
	case "h":
		return num * 3600
	}

	return 0
}
