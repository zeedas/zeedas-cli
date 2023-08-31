package lexer

import (
	"github.com/zeedas/zeedas-cli/pkg/heartbeat"

	"github.com/alecthomas/chroma/v2"
)

// Boogie lexer.
type Boogie struct{}

// Lexer returns the lexer.
func (l Boogie) Lexer() chroma.Lexer {
	return chroma.MustNewLexer(
		&chroma.Config{
			Name:      l.Name(),
			Aliases:   []string{"boogie"},
			Filenames: []string{"*.bpl"},
		},
		func() chroma.Rules {
			return chroma.Rules{
				"root": {},
			}
		},
	)
}

// Name returns the name of the lexer.
func (Boogie) Name() string {
	return heartbeat.LanguageBoogie.StringChroma()
}
