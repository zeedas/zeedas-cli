package language

import (
	"fmt"
	"io"
	"os"
	fp "path/filepath"
	"runtime/debug"
	"sort"
	"strings"

	"github.com/zeedas/zeedas-cli/pkg/heartbeat"
	"github.com/zeedas/zeedas-cli/pkg/log"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/lexers"
	_ "github.com/alecthomas/chroma/lexers/a"        // not used directly
	_ "github.com/alecthomas/chroma/lexers/b"        // not used directly
	_ "github.com/alecthomas/chroma/lexers/c"        // not used directly
	_ "github.com/alecthomas/chroma/lexers/circular" // not used directly
	_ "github.com/alecthomas/chroma/lexers/d"        // not used directly
	_ "github.com/alecthomas/chroma/lexers/e"        // not used directly
	_ "github.com/alecthomas/chroma/lexers/f"        // not used directly
	_ "github.com/alecthomas/chroma/lexers/g"        // not used directly
	_ "github.com/alecthomas/chroma/lexers/h"        // not used directly
	_ "github.com/alecthomas/chroma/lexers/i"        // not used directly
	_ "github.com/alecthomas/chroma/lexers/j"        // not used directly
	_ "github.com/alecthomas/chroma/lexers/k"        // not used directly
	_ "github.com/alecthomas/chroma/lexers/l"        // not used directly
	_ "github.com/alecthomas/chroma/lexers/m"        // not used directly
	_ "github.com/alecthomas/chroma/lexers/n"        // not used directly
	_ "github.com/alecthomas/chroma/lexers/o"        // not used directly
	_ "github.com/alecthomas/chroma/lexers/p"        // not used directly
	_ "github.com/alecthomas/chroma/lexers/q"        // not used directly
	_ "github.com/alecthomas/chroma/lexers/r"        // not used directly
	_ "github.com/alecthomas/chroma/lexers/s"        // not used directly
	_ "github.com/alecthomas/chroma/lexers/t"        // not used directly
	_ "github.com/alecthomas/chroma/lexers/v"        // not used directly
	_ "github.com/alecthomas/chroma/lexers/w"        // not used directly
	_ "github.com/alecthomas/chroma/lexers/x"        // not used directly
	_ "github.com/alecthomas/chroma/lexers/y"        // not used directly
	_ "github.com/alecthomas/chroma/lexers/z"        // not used directly
	"github.com/danwakefield/fnmatch"
)

// Max file size supporting reading from file. Default is 512Kb.
const maxFileSize = 512000

// detectChromaCustomized returns the best by filename matching lexer. Best lexer is determined
// by customized priority.
// This is a modified implementation of chroma.lexers.internal.api:Match().
func detectChromaCustomized(filepath string) (heartbeat.Language, float32, bool) {
	_, file := fp.Split(filepath)
	filename := fp.Base(file)
	matched := chroma.PrioritisedLexers{}

	// First, try primary filename matches.
	for _, lexer := range lexers.Registry.Lexers {
		config := lexer.Config()
		for _, glob := range config.Filenames {
			if fnmatch.Match(glob, filename, 0) || fnmatch.Match(glob, strings.ToLower(filename), 0) {
				matched = append(matched, lexer)
			}
		}
	}

	if len(matched) > 0 {
		bestLexer, weight := selectByCustomizedPriority(filepath, matched)

		language, ok := heartbeat.ParseLanguageFromChroma(bestLexer.Config().Name)
		if !ok {
			log.Warnf("failed to parse language from chroma lexer name %q", bestLexer.Config().Name)
			return heartbeat.LanguageUnknown, 0, false
		}

		return language, weight, true
	}

	// Next, try filename aliases.
	for _, lexer := range lexers.Registry.Lexers {
		config := lexer.Config()
		for _, glob := range config.AliasFilenames {
			if fnmatch.Match(glob, filename, 0) {
				matched = append(matched, lexer)
			}
		}
	}

	if len(matched) > 0 {
		bestLexer, weight := selectByCustomizedPriority(filepath, matched)

		language, ok := heartbeat.ParseLanguageFromChroma(bestLexer.Config().Name)
		if !ok {
			log.Warnf("failed to parse language from chroma lexer name %q", bestLexer.Config().Name)
			return heartbeat.LanguageUnknown, 0, false
		}

		return language, weight, true
	}

	// Finally, try matching by file content.
	head, err := fileHead(filepath)
	if err != nil {
		log.Warnf("failed to load head from file %q: %s", filepath, err)
		return heartbeat.LanguageUnknown, 0, false
	}

	if len(head) == 0 {
		return heartbeat.LanguageUnknown, 0, false
	}

	if lexer := analyse(string(head)); lexer != nil {
		language, ok := heartbeat.ParseLanguageFromChroma(lexer.Config().Name)
		if !ok {
			log.Warnf("failed to parse language from chroma lexer name %q", lexer.Config().Name)
			return heartbeat.LanguageUnknown, 0, false
		}

		return language, 0, true
	}

	return heartbeat.LanguageUnknown, 0, false
}

// analyse text content and return the "best" lexer.
// This is a copy of chroma.lexers.internal.api:Analyse().
func analyse(text string) chroma.Lexer {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("panicked: %v. Stack: %s", err, string(debug.Stack()))
		}
	}()

	var picked chroma.Lexer

	highest := float32(0.0)

	for _, lexer := range lexers.Registry.Lexers {
		if analyser, ok := lexer.(chroma.Analyser); ok {
			weight := analyser.AnalyseText(text)
			if weight > highest {
				picked = lexer
				highest = weight
			}
		}
	}

	return picked
}

// weightedLexer is a lexer with priority and weight.
type weightedLexer struct {
	chroma.Lexer
	Weight   float32
	Priority float32
}

// selectByCustomizedPriority selects the best matching lexer by customized priority evaluation.
func selectByCustomizedPriority(filepath string, lexers chroma.PrioritisedLexers) (chroma.Lexer, float32) {
	sort.Slice(lexers, func(i, j int) bool {
		icfg, jcfg := lexers[i].Config(), lexers[j].Config()

		// 1. by priority
		if icfg.Priority != jcfg.Priority {
			return icfg.Priority > jcfg.Priority
		}

		// 2. by name
		return strings.ToLower(icfg.Name) > strings.ToLower(jcfg.Name)
	})

	dir, _ := fp.Split(filepath)

	extensions, err := loadFolderExtensions(dir)
	if err != nil {
		log.Warnf("failed to load folder files extensions: %s", err)
		return lexers[0], 0
	}

	head, err := fileHead(filepath)
	if err != nil {
		log.Warnf("failed to load head from file %q: %s", filepath, err)
		return lexers[0], 0
	}

	var weighted []weightedLexer

	for _, lexer := range lexers {
		var weight float32

		if analyser, ok := lexer.(chroma.Analyser); ok {
			weight = analyser.AnalyseText(string(head))
		}

		cfg := lexer.Config()

		if p, ok := priority(cfg.Name); ok {
			weighted = append(weighted, weightedLexer{
				Lexer:    lexer,
				Priority: p,
				Weight:   weight,
			})

			continue
		}

		if cfg.Name == "Matlab" {
			weighted = append(weighted, weightedLexer{
				Lexer:    lexer,
				Priority: cfg.Priority,
				Weight:   matlabWeight(weight, extensions),
			})

			continue
		}

		if cfg.Name == "Objective-C" {
			weighted = append(weighted, weightedLexer{
				Lexer:    lexer,
				Priority: cfg.Priority,
				Weight:   objectiveCWeight(weight, extensions),
			})

			continue
		}

		weighted = append(weighted, weightedLexer{
			Lexer:    lexer,
			Priority: cfg.Priority,
			Weight:   weight,
		})
	}

	sort.Slice(weighted, func(i, j int) bool {
		// 1. by weight
		if weighted[i].Weight != weighted[j].Weight {
			return weighted[i].Weight > weighted[j].Weight
		}

		// 2. by priority
		if weighted[i].Priority != weighted[j].Priority {
			return weighted[i].Priority > weighted[j].Priority
		}

		// 3. name
		return weighted[i].Lexer.Config().Name > weighted[j].Lexer.Config().Name
	})

	return weighted[0].Lexer, weighted[0].Weight
}

// fileHead returns the first `maxFileSize` bytes of the file's content.
func fileHead(filepath string) ([]byte, error) {
	f, err := os.Open(filepath) // nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %s", err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Debugf("failed to close file '%s': %s", filepath, err)
		}
	}()

	data, err := io.ReadAll(io.LimitReader(f, maxFileSize))
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to read bytes from file: %s", err)
	}

	return data, nil
}

// objectiveCWeight determines the weight of objective-c by the provided same folder file extensions.
func objectiveCWeight(weight float32, extensions []string) float32 {
	var matFileExists bool

	for _, e := range extensions {
		if e == ".mat" {
			matFileExists = true
			break
		}
	}

	if matFileExists {
		weight -= 0.01
	} else {
		weight += 0.01
	}

	for _, e := range extensions {
		if e == ".h" {
			weight += 0.01
			break
		}
	}

	return weight
}

// matlabWeight determines the weight of matlab by the provided same folder file extensions.
func matlabWeight(weight float32, extensions []string) float32 {
	for _, e := range extensions {
		if e == ".mat" {
			weight += 0.01
			break
		}
	}

	var headerFileExists bool

	for _, e := range extensions {
		if e == ".h" {
			headerFileExists = true
			break
		}
	}

	if !headerFileExists {
		weight += 0.01
	}

	return weight
}
