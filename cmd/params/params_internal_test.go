package params

import (
	"regexp"
	"testing"

	"github.com/zeedas/zeedas-cli/pkg/regex"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseEditorFromPlugin(t *testing.T) {
	tests := map[string]struct {
		Plugin   string
		Expected string
	}{
		"editor/version plugin/version": {
			Plugin:   "vscode/1.51.1 vscode-zeedas/4.0.9",
			Expected: "vscode",
		},
		"plugin/version (no dash)": {
			Plugin:   "emacs-zeedas/1.0.2",
			Expected: "emacs",
		},
		"plugin/version (multiple dashes)": {
			Plugin:   "camunda-modeler-zeedas/0.4.3",
			Expected: "camunda-modeler",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			editor, err := parseEditorFromPlugin(test.Plugin)
			require.NoError(t, err)

			assert.Equal(t, test.Expected, editor)
		})
	}
}

func TestParseEditorFromPluginErr(t *testing.T) {
	editor, err := parseEditorFromPlugin("editor-zeedas")
	require.Error(t, err)

	assert.Empty(t, editor)
	assert.Equal(t, "plugin malformed: editor-zeedas", err.Error())
}

func TestParseBoolOrRegexList(t *testing.T) {
	tests := map[string]struct {
		Input    string
		Expected []regex.Regex
	}{
		"string empty": {
			Input:    " ",
			Expected: nil,
		},
		"false string": {
			Input:    "false",
			Expected: []regex.Regex{regexp.MustCompile("a^")},
		},
		"true string": {
			Input:    "true",
			Expected: []regex.Regex{regexp.MustCompile(".*")},
		},
		"valid regex": {
			Input: "\t.?\n\t\n \n\t\tzeedas.? \t\n",
			Expected: []regex.Regex{
				regexp.MustCompile(".?"),
				regexp.MustCompile("zeedas.?"),
			},
		},
		"valid regex with windows style": {
			Input: "\t.?\r\n\t\t\tzeedas.? \t\r\n",
			Expected: []regex.Regex{
				regexp.MustCompile(".?"),
				regexp.MustCompile("zeedas.?"),
			},
		},
		"valid regex with old mac style": {
			Input: "\t.?\r\t\t\tzeedas.? \t\r",
			Expected: []regex.Regex{
				regexp.MustCompile(".?"),
				regexp.MustCompile("zeedas.?"),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			regex, err := parseBoolOrRegexList(test.Input)
			require.NoError(t, err)

			assert.Equal(t, test.Expected, regex)
		})
	}
}
