package heartbeat

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsXCodePlayground(t *testing.T) {
	tests := map[string]struct {
		Dir      string
		Expected bool
	}{
		"playground directory": {
			Dir:      setupTestXCodePlayground(t, "zeedas.playground"),
			Expected: true,
		},
		"xcplayground directory": {
			Dir:      setupTestXCodePlayground(t, "zeedas.xcplayground"),
			Expected: true,
		},
		"xcplaygroundpage directory": {
			Dir:      setupTestXCodePlayground(t, "zeedas.xcplaygroundpage"),
			Expected: true,
		},
		"not playground": {
			Dir:      setupTestXCodePlayground(t, "zeedas"),
			Expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ret := isXCodePlayground(test.Dir)

			assert.Equal(t, test.Expected, ret)
		})
	}
}

func setupTestXCodePlayground(t *testing.T, dir string) string {
	tmpDir := t.TempDir()

	err := os.Mkdir(filepath.Join(tmpDir, dir), os.FileMode(int(0700)))
	require.NoError(t, err)

	return filepath.Join(tmpDir, dir)
}
