package project

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/wakatime/wakatime-cli/pkg/log"
)

// Subversion contains svn data.
type Subversion struct {
	// Filepath contains the entity path.
	Filepath string
}

// Detect gets information about the svn project for a given file.
func (s Subversion) Detect() (Result, bool, error) {
	binary, ok := findSvnBinary()
	if !ok {
		log.Debugln("svn binary not found")
		return Result{}, false, nil
	}

	var fp string

	// Take only the directory
	if fileExists(s.Filepath) {
		fp = filepath.Dir(s.Filepath)
	}

	// Find for .svn/wc.db file
	svnConfigFile, ok := FindFileOrDirectory(fp, filepath.Join(".svn", "wc.db"))
	if !ok {
		return Result{}, false, nil
	}

	info, ok, err := svnInfo(filepath.Join(svnConfigFile, "..", ".."), binary)
	if err != nil {
		return Result{}, false, fmt.Errorf("failed to get svn info: %s", err)
	}

	if !ok {
		return Result{}, false, nil
	}

	return Result{
		Project: resolveSvnInfo(info, "Repository Root"),
		Branch:  resolveSvnInfo(info, "URL"),
		Folder:  strings.ReplaceAll(info["Repository Root"], "\r", ""),
	}, true, nil
}

func svnInfo(fp string, binary string) (map[string]string, bool, error) {
	if runtime.GOOS == "darwin" && !hasXcodeTools() {
		return nil, false, nil
	}

	cmd := exec.Command(binary, "info", fp)
	out, err := cmd.Output()

	if err != nil {
		return nil, false, fmt.Errorf("error getting svn info: %s", err)
	}

	result := map[string]string{}

	for _, line := range strings.Split(string(out), "\n") {
		item := strings.Split(line, ": ")
		if len(item) == 2 {
			result[item[0]] = item[1]
		}
	}

	return result, true, nil
}

func findSvnBinary() (string, bool) {
	locations := []string{
		"svn",
		"/usr/bin/svn",
		"/usr/local/bin/svn",
	}

	for _, loc := range locations {
		cmd := exec.Command(loc, "--version") // nolint:gosec

		err := cmd.Run()
		if err != nil {
			log.Debugf("failed while calling %s --version: %s", loc, err)
			continue
		}

		return loc, true
	}

	return "", false
}

func hasXcodeTools() bool {
	cmd := exec.Command("/usr/bin/xcode-select", "-p")

	return cmd.Run() == nil
}

func resolveSvnInfo(info map[string]string, key string) string {
	if val, ok := info[key]; ok {
		parts := strings.Split(val, "/")
		last := parts[len(parts)-1]
		parts2 := strings.Split(last, "\\")
		last2 := parts2[len(parts2)-1]

		return strings.ReplaceAll(last2, "\r", "")
	}

	return ""
}

// ID returns its id.
func (Subversion) ID() DetectorID {
	return SubversionDetector
}
