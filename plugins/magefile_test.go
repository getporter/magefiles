package plugins

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstallPlugin(t *testing.T) {
	magefile := NewMagefile("github.com/myplugin/test-plugin", "testplugin", "testdata/bin")

	// Change the porter home to a safe place for the test to write to
	require.NoError(t, os.MkdirAll("testdata/porter_home", 0770))
	os.Setenv("PORTER_HOME", "testdata/porter_home")
	defer os.Unsetenv("PORTER_HOME")

	magefile.Install()

	assert.DirExists(t, "testdata/porter_home/plugins/testplugin", "The plugin directory doesn't exist")
	assert.FileExists(t, "testdata/porter_home/plugins/testplugin/testplugin", "The plugin wasn't installed")
}
