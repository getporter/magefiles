package releases

import (
	"crypto/rand"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"get.porter.sh/magefiles/porter"
	"github.com/carolynvs/magex/mgx"
	"github.com/carolynvs/magex/shx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetReleaseAssets(t *testing.T) {
	tmp, err := ioutil.TempDir("", "magefiles")
	require.NoError(t, err)
	defer os.RemoveAll(tmp)

	mgx.Must(shx.Copy("testdata/mixins/v1.2.3/*", tmp, shx.CopyRecursive))

	gotFiles, err := getReleaseAssets(tmp)
	require.NoError(t, err)

	wantFiles := []string{
		filepath.Join(tmp, "mymixin-darwin-amd64"),
		filepath.Join(tmp, "mymixin-darwin-amd64.sha256sum"),
		filepath.Join(tmp, "mymixin-linux-amd64"),
		filepath.Join(tmp, "mymixin-linux-amd64.sha256sum"),
		filepath.Join(tmp, "mymixin-windows-amd64.exe"),
		filepath.Join(tmp, "mymixin-windows-amd64.exe.sha256sum"),
	}
	assert.Equal(t, wantFiles, gotFiles)

	// Read the existing checksum file with stale contents, and ensure it was updated
	gotChecksum, err := ioutil.ReadFile(filepath.Join(tmp, "mymixin-darwin-amd64.sha256sum"))
	require.NoError(t, err)
	wantCheckSum := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855  mymixin-darwin-amd64"
	assert.Equal(t, wantCheckSum, string(gotChecksum))
}

func TestAddChecksumExt(t *testing.T) {
	tests := []struct {
		input         string
		expectedAdded bool
		expected      string
	}{
		{
			input:         "porter.sh",
			expectedAdded: true,
			expected:      "porter.sh.sha256sum",
		},
		{
			input:         "porter.sh.sha256sum",
			expectedAdded: false,
			expected:      "porter.sh.sha256sum",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run("", func(t *testing.T) {
			output, added := AddChecksumExt(tt.input)
			assert.Equal(t, tt.expected, output)
			assert.Equal(t, tt.expectedAdded, added)
		})
	}

}

func TestAppendDataPath(t *testing.T) {
	data := make([]byte, 10)
	_, err := rand.Read(data)
	require.NoError(t, err)
	dataPath := "test/random"
	expected := hex.EncodeToString(data) + "  random"

	output := AppendDataPath(data, dataPath)
	require.Equal(t, expected, output)
}

func TestGenerateMixinFeed(t *testing.T) {
	tmp, err := ioutil.TempDir("", "magefiles")
	require.NoError(t, err)
	defer os.RemoveAll(tmp)

	// Install porter in our test bin
	tmpBin := filepath.Join(tmp, "bin")
	require.NoError(t, shx.Copy("../bin", tmpBin, shx.CopyRecursive), "failed to copy the porter bin into the test directory")
	porter.UsePorterHome(tmpBin)

	// Copy our atom feed template
	buildDir := filepath.Join(tmp, "build")
	require.NoError(t, os.Mkdir(buildDir, 0770))
	require.NoError(t, shx.Copy("testdata/atom-template.xml", buildDir))

	// Make a fake mixin release
	require.NoError(t, shx.Copy("testdata/mixins", tmpBin, shx.CopyRecursive))

	// Change into the tmp directory since the publish logic uses relative file paths
	origDir, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(tmp))
	defer os.Chdir(origDir)

	err = GenerateMixinFeed()
	require.NoError(t, err)

	assert.FileExists(t, filepath.Join(tmpBin, "mixins/.packages/mixins/atom.xml"), "expected a mixin feed")
}

func TestGeneratePluginFeed_PorterNotInstalled(t *testing.T) {
	tmp, err := ioutil.TempDir("", "magefiles")
	require.NoError(t, err)
	defer os.RemoveAll(tmp)

	// DO NOT INSTALL PORTER INTO THE BIN
	tmpBin := filepath.Join(tmp, "bin")

	// Copy our atom feed template
	buildDir := filepath.Join(tmp, "build")
	require.NoError(t, os.Mkdir(buildDir, 0770))
	require.NoError(t, shx.Copy("testdata/atom-template.xml", buildDir))

	// Make a fake mixin release
	require.NoError(t, shx.Copy("testdata/mixins", tmpBin, shx.CopyRecursive))

	// Change into the tmp directory since the publish logic uses relative file paths
	origDir, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(tmp))
	defer os.Chdir(origDir)

	err = GeneratePluginFeed()
	require.Errorf(t, err, "farts", "GeneratePluginFeed should fail when porter is not in the bin")
}
