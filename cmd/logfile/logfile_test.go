package logfile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/zeedas/zeedas-cli/cmd/logfile"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadParams(t *testing.T) {
	tmpFile, err := os.CreateTemp(t.TempDir(), "")
	require.NoError(t, err)

	defer tmpFile.Close()

	dir, _ := filepath.Split(tmpFile.Name())

	logFile, err := os.Create(filepath.Join(dir, "zeedas.log"))
	require.NoError(t, err)

	defer logFile.Close()

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	tests := map[string]struct {
		ViperLogFile       string
		ViperLogFileConfig string
		ViperLogFileOld    string
		ViperToStdout      bool
		EnvVar             string
		ViperDebug         bool
		ViperDebugConfig   bool
		Expected           logfile.Params
	}{
		"log file and verbose set": {
			ViperDebug: true,
			Expected: logfile.Params{
				File:    filepath.Join(home, ".zeedas", "zeedas.log"),
				Verbose: true,
			},
		},
		"log file and verbose from config": {
			ViperDebugConfig: true,
			Expected: logfile.Params{
				File:    filepath.Join(home, ".zeedas", "zeedas.log"),
				Verbose: true,
			},
		},
		"log file flag takes precedence": {
			ViperLogFile:       tmpFile.Name(),
			ViperLogFileConfig: "otherfolder/zeedas.config.log",
			ViperLogFileOld:    "otherfolder/zeedas.old.log",
			Expected: logfile.Params{
				File: tmpFile.Name(),
			},
		},
		"log file deprecated flag takes precedence": {
			ViperLogFileConfig: "otherfolder/zeedas.config.log",
			ViperLogFileOld:    tmpFile.Name(),
			Expected: logfile.Params{
				File: tmpFile.Name(),
			},
		},
		"log file from config": {
			ViperLogFileConfig: tmpFile.Name(),
			Expected: logfile.Params{
				File: tmpFile.Name(),
			},
		},
		"log file from ZEEDAS_HOME": {
			EnvVar: dir,
			Expected: logfile.Params{
				File: filepath.Join(dir, "zeedas.log"),
			},
		},
		"log file from home dir": {
			Expected: logfile.Params{
				File: filepath.Join(home, ".zeedas", "zeedas.log"),
			},
		},
		"log to stdout": {
			ViperToStdout: true,
			Expected: logfile.Params{
				File:     filepath.Join(home, ".zeedas", "zeedas.log"),
				ToStdout: true,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			v := viper.New()
			v.Set("log-file", test.ViperLogFile)
			v.Set("logfile", test.ViperLogFileOld)
			v.Set("log-to-stdout", test.ViperToStdout)
			v.Set("settings.log_file", test.ViperLogFileConfig)
			v.Set("settings.debug", test.ViperDebug)
			v.Set("verbose", test.ViperDebugConfig)

			err := os.Setenv("ZEEDAS_HOME", test.EnvVar)
			require.NoError(t, err)

			defer os.Unsetenv("ZEEDAS_HOME")

			params, err := logfile.LoadParams(v)
			require.NoError(t, err)

			assert.Equal(t, test.Expected, params)
		})
	}
}
