package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"go.uber.org/config"
)

// templateFuncs are the functions made available to config templates.
// This is a subset of what dockerize provides
var templateFuncs = template.FuncMap{
	"default": defaultValue,
	"lower":   strings.ToLower,
	"split":   strings.Split,
}

// templateData is the data made available to config templates.
type templateData struct {
	// Env contains all environment variables of the current process.
	Env map[string]string
}

type templatedFile struct {
	path   string
	result *bytes.Buffer
}

func (t *templatedFile) Read(p []byte) (int, error) {
	// We need to cache this result because Read is copying the output to the specified slice.
	// If their buffer is too small (it is) then they get an error, resize it, and call again. We want to avoid
	// repeatedly loading and executing the template during this "right-sizing" process of their buffer.
	if t.result == nil {
		// No point in caching the error, they're not going to call us again
		buffer, err := t.render()
		if err != nil {
			return 0, err
		}
		t.result = buffer
	}
	return t.result.Read(p)
}

func (t *templatedFile) render() (*bytes.Buffer, error) {
	contents, err := os.ReadFile(t.path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %q: %w", t.path, err)
	}

	tmpl, err := template.New(filepath.Base(t.path)).Funcs(templateFuncs).Parse(string(contents))
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file %q as a template: %w", t.path, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, templateData{Env: environment()}); err != nil {
		return nil, fmt.Errorf("failed to render config file %q: %w", t.path, err)
	}

	return &buf, nil
}

// TemplatedFile returns a config.YAMLOption which renders the file at the given
// path as a Go text/template before parsing it as YAML.
func TemplatedFile(path string) config.YAMLOption {
	return config.Source(&templatedFile{
		path: path,
	})
}

func environment() map[string]string {
	environ := os.Environ()
	env := make(map[string]string, len(environ))
	for _, kv := range environ {
		key, value, found := strings.Cut(kv, "=")
		if !found {
			continue
		}
		env[key] = value
	}
	return env
}

// defaultValue returns the first non-empty string argument. Expected usage is to have a literal as the final parameter
// so that it falls back from potentially nil values to that value as a default.
func defaultValue(args ...interface{}) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("default called with no values")
	}
	for _, arg := range args {
		if arg == nil {
			continue
		}
		asString, ok := arg.(string)
		if !ok {
			return "", fmt.Errorf("default called with non-string value")
		}
		if asString != "" {
			return asString, nil
		}

	}

	return "", fmt.Errorf("default called with only nil values")
}
