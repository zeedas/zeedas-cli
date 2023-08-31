package heartbeat_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/zeedas/zeedas-cli/pkg/heartbeat"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithEntityModifier_XCodePlayground(t *testing.T) {
	tmpDir := t.TempDir()

	err := os.Mkdir(filepath.Join(tmpDir, "zeedas.playground"), os.FileMode(int(0700)))
	require.NoError(t, err)

	opt := heartbeat.WithEntityModifer()

	handle := opt(func(hh []heartbeat.Heartbeat) ([]heartbeat.Result, error) {
		assert.Equal(t, []heartbeat.Heartbeat{
			{
				Entity:     filepath.Join(tmpDir, "zeedas.playground", "Contents.swift"),
				EntityType: heartbeat.FileType,
			},
		}, hh)

		return []heartbeat.Result{
			{
				Status: 201,
			},
		}, nil
	})

	result, err := handle([]heartbeat.Heartbeat{
		{
			Entity:     filepath.Join(tmpDir, "zeedas.playground"),
			EntityType: heartbeat.FileType,
		},
	})
	require.NoError(t, err)

	assert.Equal(t, []heartbeat.Result{
		{
			Status: 201,
		},
	}, result)
}
