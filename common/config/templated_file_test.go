package config

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	uconfig "go.uber.org/config"
)

func TestTemplatedFileRender(t *testing.T) {
	tests := []struct {
		name     string
		contents string
		env      map[string]string
		want     string
		wantErr  string
	}{
		{
			name:     "no template directives",
			contents: "log:\n  level: info\n",
			want:     "log:\n  level: info\n",
		},
		{
			name:     "default with unset variable",
			contents: `level: {{ default .Env.LOG_LEVEL "info" }}`,
			want:     "level: info",
		},
		{
			name:     "default with empty variable",
			contents: `level: {{ default .Env.LOG_LEVEL "info" }}`,
			env:      map[string]string{"LOG_LEVEL": ""},
			want:     "level: info",
		},
		{
			name:     "default with set variable",
			contents: `level: {{ default .Env.LOG_LEVEL "info" }}`,
			env:      map[string]string{"LOG_LEVEL": "debug"},
			want:     "level: debug",
		},
		{
			name:     "lower",
			contents: `{{- $db := default .Env.DB "Cassandra" | lower -}}` + "\ndb: {{ $db }}",
			env:      map[string]string{"DB": "MySQL"},
			want:     "db: mysql",
		},
		{
			name:     "split",
			contents: "hosts:\n{{- range $seed := (split .Env.SEEDS \",\") }}\n  - {{ . }}\n{{- end }}",
			env:      map[string]string{"SEEDS": "host1:7933,host2:7933"},
			want:     "hosts:\n  - host1:7933\n  - host2:7933",
		},
		{
			name:     "default with no default value",
			contents: `level: {{ default .Env.LOG_LEVEL }}`,
			wantErr:  "default called with only nil values",
		},
		{
			name:     "default with non string default value",
			contents: `port: {{ default .Env.PORT 7933 }}`,
			wantErr:  "default called with non-string value",
		},
		{
			name:     "unknown function",
			contents: `level: {{ upper .Env.LOG_LEVEL }}`,
			wantErr:  "as a template",
		},
		{
			name:     "invalid template syntax",
			contents: `level: {{ .Env.LOG_LEVEL`,
			wantErr:  "as a template",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.env {
				t.Setenv(k, v)
			}

			file := writeTempFile(t, tc.contents)
			got, err := (&templatedFile{path: file}).render()

			if tc.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got.String())
		})
	}
}

func TestTemplatedFileRead(t *testing.T) {
	t.Setenv("LOG_LEVEL", "debug")
	contents := `level: {{ default .Env.LOG_LEVEL "info" }}`
	want := "level: debug"

	file := &templatedFile{path: writeTempFile(t, contents)}

	// reads are streamed in whatever size the caller asks for
	var got []byte
	buf := make([]byte, 5)
	for {
		n, err := file.Read(buf)
		got = append(got, buf[:n]...)
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
	}
	assert.Equal(t, want, string(got))

	// once exhausted it keeps reporting EOF
	n, err := file.Read(buf)
	assert.Zero(t, n)
	assert.Equal(t, io.EOF, err)
}

func TestTemplatedFileReadAll(t *testing.T) {
	t.Setenv("LOG_LEVEL", "debug")
	file := &templatedFile{path: writeTempFile(t, `level: {{ default .Env.LOG_LEVEL "info" }}`)}

	got, err := io.ReadAll(file)
	require.NoError(t, err)
	assert.Equal(t, "level: debug", string(got))
}

func TestTemplatedFileReadMissingFile(t *testing.T) {
	file := &templatedFile{path: filepath.Join(t.TempDir(), "does-not-exist.yaml")}

	n, err := file.Read(make([]byte, 16))
	assert.Zero(t, n)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read config file")
}

func TestTemplatedFileWithYAMLProvider(t *testing.T) {
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("SEEDS", "host1,host2")

	contents := "log:\n" +
		`  level: {{ default .Env.LOG_LEVEL "info" }}` + "\n" +
		`  stdout: {{ default .Env.LOG_STDOUT "true" }}` + "\n" +
		"hosts:\n{{- range $seed := (split .Env.SEEDS \",\") }}\n  - {{ . }}\n{{- end }}\n"

	provider, err := uconfig.NewYAML(TemplatedFile(writeTempFile(t, contents)))
	require.NoError(t, err)

	var cfg struct {
		Log struct {
			Level  string `yaml:"level"`
			Stdout bool   `yaml:"stdout"`
		} `yaml:"log"`
		Hosts []string `yaml:"hosts"`
	}
	require.NoError(t, provider.Get(uconfig.Root).Populate(&cfg))

	assert.Equal(t, "debug", cfg.Log.Level)
	assert.True(t, cfg.Log.Stdout)
	assert.Equal(t, []string{"host1", "host2"}, cfg.Hosts)
}

func TestDefaultValue(t *testing.T) {
	tests := []struct {
		name    string
		args    []interface{}
		want    string
		wantErr string
	}{
		{name: "no values", wantErr: "default called with no values"},
		{name: "value present", args: []interface{}{"a", "b"}, want: "a"},
		{name: "empty value", args: []interface{}{"", "b"}, want: "b"},
		{name: "nil value", args: []interface{}{nil, "b"}, want: "b"},
		{name: "non string value", args: []interface{}{7933, "b"}, wantErr: "default called with non-string value"},
		{name: "missing default", args: []interface{}{""}, wantErr: "default called with only nil values"},
		{name: "nil default", args: []interface{}{"", nil}, wantErr: "default called with only nil values"},
		{name: "non string default", args: []interface{}{"", 7933}, wantErr: "default called with non-string value"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := defaultValue(tc.args...)
			if tc.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestEnvironment(t *testing.T) {
	t.Setenv("CADENCE_TEST_ENV_VAR", "value")

	env := environment()

	assert.Equal(t, "value", env["CADENCE_TEST_ENV_VAR"])
	assert.Empty(t, env["CADENCE_TEST_ENV_VAR_UNSET"])
	assert.Len(t, env, len(os.Environ()))
}

func writeTempFile(t *testing.T, contents string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yaml")
	require.NoError(t, os.WriteFile(path, []byte(contents), 0600))
	return path
}
