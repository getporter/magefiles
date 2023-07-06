package porter

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/carolynvs/magex/pkg"
	"github.com/carolynvs/magex/xplat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnsurePorter(t *testing.T) {
	tmp, err := os.MkdirTemp("", "magefiles")
	require.NoError(t, err)
	defer os.RemoveAll(tmp)

	const wantVersion = "v1.0.0-beta.1"
	UsePorterHome(tmp)
	EnsurePorterAt(wantVersion)
	EnsurePorter() // This should not change the installed version in bin!

	require.FileExists(t, filepath.Join(tmp, "porter"+xplat.FileExt()), "expected the porter client to be in bin")
	assert.FileExists(t, filepath.Join(tmp, "runtimes", "porter-runtime"), "expected the porter runtime to be in bin")

	ok, err := pkg.IsCommandAvailable("porter", "--version", wantVersion)
	require.NoError(t, err)
	assert.True(t, ok, "could not resolve the desired porter version")
}

func TestEnsurePorterAt(t *testing.T) {
	testcases := []struct {
		name        string
		wantVersion string
	}{
		{name: "default version", wantVersion: DefaultPorterVersion},
		{name: "custom version", wantVersion: "v1.0.0-alpha.10"},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tmp := t.TempDir()

			UsePorterHome(tmp)
			EnsurePorterAt(tc.wantVersion)
			require.FileExists(t, filepath.Join(tmp, "porter"+xplat.FileExt()), "expected the porter client to be in bin")
			assert.FileExists(t, filepath.Join(tmp, "runtimes", "porter-runtime"), "expected the porter runtime to be in bin")

			ok, err := pkg.IsCommandAvailable("porter", "--version", tc.wantVersion)
			require.NoError(t, err)
			assert.True(t, ok, "could not resolve the desired porter version")
		})
	}

	// when the runtime binary already exists, leave it
	t.Run("runtime binary only downloaded when client is stale", func(t *testing.T) {
		tmp := t.TempDir()

		UsePorterHome(tmp)
		EnsurePorter()

		porterPath := filepath.Join(tmp, "porter"+xplat.FileExt())
		runtimePath := filepath.Join(tmp, "runtimes", "porter-runtime")
		origPorterStat, err := os.Stat(porterPath)
		require.NoError(t, err, "failed to stat the porter binary")
		origRuntimeStat, err := os.Stat(runtimePath)
		require.NoError(t, err, "failed to stat the porter-runtime binary")

		// Nothing should be downloaded
		EnsurePorterAt(DefaultPorterVersion)
		newPorterStat, err := os.Stat(porterPath)
		require.NoError(t, err, "failed to stat the porter binary")
		require.Equal(t, origPorterStat.ModTime(), newPorterStat.ModTime(), "expected the porter binary to not be re-downloaded")
		newRuntimeStat, err := os.Stat(runtimePath)
		require.NoError(t, err, "failed to stat the porter-runtime binary")
		require.Equal(t, origRuntimeStat.ModTime(), newRuntimeStat.ModTime(), "expected the porter-runtime binary to not be re-downloaded")

		// Both should be re-downloaded
		EnsurePorterAt("v1.0.0-rc.1")
		newPorterStat, err = os.Stat(porterPath)
		require.NoError(t, err, "failed to stat the porter binary")
		require.Less(t, origPorterStat.ModTime(), newPorterStat.ModTime(), "expected the porter binary to be re-downloaded")
		newRuntimeStat, err = os.Stat(runtimePath)
		require.NoError(t, err, "failed to stat the porter-runtime binary")
		require.Less(t, origRuntimeStat.ModTime(), newRuntimeStat.ModTime(), "expected the porter-runtime binary to be re-downloaded")
	})
}
