package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupDCO(t *testing.T) {
	t.Run("git exists", func(t *testing.T) {
		tmp := t.TempDir()
		testDir := filepath.Join(tmp, "a/b")
		require.NoError(t, os.MkdirAll(testDir, 0755))
		gitDir := filepath.Join(tmp, ".git")
		require.NoError(t, os.Mkdir(gitDir, 0755))

		require.NoError(t, os.Chdir(testDir))
		require.NoError(t, SetupDCO())

		// test that the hook was created
		hookPath := filepath.Join(gitDir, "hooks/prepare-commit-msg")
		require.FileExists(t, hookPath)
		hookContents, err := os.ReadFile(hookPath)
		require.NoErrorf(t, err, "error reading %s", hookPath)
		assert.Equal(t, string(hookContents), prepareCommitMsg, "unexpected hook file contents found")
	})

	t.Run("git exists", func(t *testing.T) {
		tmp := t.TempDir()
		testPath := filepath.Join(tmp, "a/b")
		require.NoError(t, os.MkdirAll(testPath, 0755))
		// there is not .git directory

		require.NoError(t, os.Chdir(testPath))
		err := SetupDCO()
		require.ErrorContains(t, err, "could not find the repository root")
	})
}
