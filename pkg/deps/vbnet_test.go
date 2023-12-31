package deps_test

import (
	"testing"

	"github.com/zeedas/zeedas-cli/pkg/deps"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParserVbNet_Parse(t *testing.T) {
	parser := deps.ParserVbNet{}

	dependencies, err := parser.Parse("testdata/vbnet.vb")
	require.NoError(t, err)

	assert.Equal(t, []string{
		"zeedas",
		"Math",
		"Proper",
	}, dependencies)
}
