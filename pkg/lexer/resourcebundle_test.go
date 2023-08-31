package lexer_test

import (
	"os"
	"testing"

	"github.com/zeedas/zeedas-cli/pkg/lexer"

	"github.com/stretchr/testify/assert"
)

func TestResourceBundle_AnalyseText(t *testing.T) {
	data, err := os.ReadFile("testdata/resource.txt")
	assert.NoError(t, err)

	l := lexer.ResourceBundle{}.Lexer()

	assert.Equal(t, float32(1.0), l.AnalyseText(string(data)))
}
