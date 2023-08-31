package project_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/zeedas/zeedas-cli/pkg/project"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMercurial_Detect(t *testing.T) {
	fp := setupTestMercurial(t)

	m := project.Mercurial{
		Filepath: filepath.Join(fp, "zeedas-cli/src/pkg/file.go"),
	}

	result, detected, err := m.Detect()
	require.NoError(t, err)

	assert.True(t, detected)
	assert.Contains(t, result.Folder, fp)
	assert.Equal(t, project.Result{
		Project: "zeedas-cli",
		Branch:  "billing",
		Folder:  result.Folder,
	}, result)
}

func TestMercurial_Detect_BranchWithSlash(t *testing.T) {
	fp := setupTestMercurialBranchWithSlash(t)

	m := project.Mercurial{
		Filepath: filepath.Join(fp, "zeedas-cli/src/pkg/file.go"),
	}

	result, detected, err := m.Detect()
	require.NoError(t, err)

	assert.True(t, detected)
	assert.Contains(t, result.Folder, fp)
	assert.Equal(t, project.Result{
		Project: "zeedas-cli",
		Branch:  "feature/billing",
		Folder:  result.Folder,
	}, result)
}

func TestMercurial_Detect_NoBranch(t *testing.T) {
	fp := setupTestMercurialNoBranch(t)

	m := project.Mercurial{
		Filepath: filepath.Join(fp, "zeedas-cli/src/pkg/file.go"),
	}

	result, detected, err := m.Detect()
	require.NoError(t, err)

	assert.True(t, detected)
	assert.Contains(t, result.Folder, fp)
	assert.Equal(t, project.Result{
		Project: "zeedas-cli",
		Branch:  "default",
		Folder:  result.Folder,
	}, result)
}

func TestMercurial_ID(t *testing.T) {
	m := project.Mercurial{}

	assert.Equal(t, project.MercurialDetector, m.ID())
}

func setupTestMercurial(t *testing.T) (fp string) {
	tmpDir := t.TempDir()

	err := os.MkdirAll(filepath.Join(tmpDir, "zeedas-cli/src/pkg"), os.FileMode(int(0700)))
	require.NoError(t, err)

	tmpFile, err := os.Create(filepath.Join(tmpDir, "zeedas-cli/src/pkg/file.go"))
	require.NoError(t, err)

	defer tmpFile.Close()

	err = os.Mkdir(filepath.Join(tmpDir, "zeedas-cli/.hg"), os.FileMode(int(0700)))
	require.NoError(t, err)

	copyFile(t, "testdata/hg/branch", filepath.Join(tmpDir, "zeedas-cli/.hg/branch"))

	return tmpDir
}

func setupTestMercurialBranchWithSlash(t *testing.T) (fp string) {
	tmpDir := t.TempDir()

	err := os.MkdirAll(filepath.Join(tmpDir, "zeedas-cli/src/pkg"), os.FileMode(int(0700)))
	require.NoError(t, err)

	tmpFile, err := os.Create(filepath.Join(tmpDir, "zeedas-cli/src/pkg/file.go"))
	require.NoError(t, err)

	defer tmpFile.Close()

	err = os.Mkdir(filepath.Join(tmpDir, "zeedas-cli/.hg"), os.FileMode(int(0700)))
	require.NoError(t, err)

	copyFile(t, "testdata/hg/branch_with_slash", filepath.Join(tmpDir, "zeedas-cli/.hg/branch"))

	return tmpDir
}

func setupTestMercurialNoBranch(t *testing.T) (fp string) {
	tmpDir := t.TempDir()

	err := os.MkdirAll(filepath.Join(tmpDir, "zeedas-cli/src/pkg"), os.FileMode(int(0700)))
	require.NoError(t, err)

	tmpFile, err := os.Create(filepath.Join(tmpDir, "zeedas-cli/src/pkg/file.go"))
	require.NoError(t, err)

	defer tmpFile.Close()

	err = os.Mkdir(filepath.Join(tmpDir, "zeedas-cli/.hg"), os.FileMode(int(0700)))
	require.NoError(t, err)

	return tmpDir
}
