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
	check := func(wantVersion string) {
		UseBinForPorterHome()
		defer os.RemoveAll("bin")

		EnsurePorter()
		require.FileExists(t, filepath.Join("bin", "porter"+xplat.FileExt()), "expected the porter client to be in bin")
		assert.FileExists(t, filepath.Join("bin", "runtimes", "porter-runtime"), "expected the porter runtime to be in bin")

		ok, err := pkg.IsCommandAvailable("porter", wantVersion, "--version")
		require.NoError(t, err)
		assert.True(t, ok, "could not resolve the desired porter version")
	}

	t.Run("EnsurePorter - DefaultVersion", func(t *testing.T) {
		check(DefaultPorterVersion)
	})

	t.Run("EnsurePorterAt - CustomVersion", func(t *testing.T) {
		check("v1.0.0-alpha.1")
	})
}
