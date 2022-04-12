package porter

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/carolynvs/magex/pkg"
	"github.com/carolynvs/magex/xplat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnsurePorter(t *testing.T) {
	testcases := []struct {
		name        string
		wantVersion string
	}{
		{name: "default version", wantVersion: DefaultPorterVersion},
		{name: "custom version", wantVersion: "v1.0.0-alpha.10"},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tmp, err := ioutil.TempDir("", "magefiles")
			require.NoError(t, err)
			defer os.RemoveAll(tmp)

			UsePorterHome(tmp)
			EnsurePorterAt(tc.wantVersion)
			require.FileExists(t, filepath.Join(tmp, "porter"+xplat.FileExt()), "expected the porter client to be in bin")
			assert.FileExists(t, filepath.Join(tmp, "runtimes", "porter-runtime"), "expected the porter runtime to be in bin")

			ok, err := pkg.IsCommandAvailable("porter", tc.wantVersion, "--version")
			require.NoError(t, err)
			assert.True(t, ok, "could not resolve the desired porter version")
		})
	}
}
