package plugins

import (
	"fmt"
	"os"
	"path/filepath"

	"get.porter.sh/magefiles/ci"
	"get.porter.sh/magefiles/porter"
	"get.porter.sh/magefiles/releases"
	"github.com/carolynvs/magex/mgx"
	"github.com/carolynvs/magex/shx"
	"github.com/carolynvs/magex/xplat"
	"github.com/magefile/mage/mg"
)

// Magefile provides implementations for required mage targets needed by a Porter plugin.
type Magefile struct {
	// Pkg is the plugin's go module package name
	// For example, github.com/yourname/yourplugin
	Pkg string

	// PluginName is the name of the plugin binary
	PluginName string

	// OutDir is the path to the directory the plugin binary should be output.
	OutDir string
}

// NewMagefile creates a Magefile helper for a plugin.
func NewMagefile(pkg, pluginName, binDir string) Magefile {
	return Magefile{
		Pkg:        pkg,
		PluginName: pluginName,
		OutDir:     filepath.Join(binDir, "plugins", pluginName)}
}

var must = shx.CommandBuilder{StopOnError: true}

// ConfigureAgent sets up a CI worker to use Go and mage.
func (m Magefile) ConfigureAgent() {
	mgx.Must(ci.ConfigureAgent())
}

// Build the plugin
func (m Magefile) Build() {
	must.RunV("go", "mod", "tidy")
	releases.BuildPlugin(m.PluginName)
}

// XBuildAll cross-compiles the plugin before a release
func (m Magefile) XBuildAll() {
	releases.XBuildAll(m.Pkg, m.PluginName, m.OutDir)
	releases.PrepareMixinForPublish(m.PluginName)
}

// TestUnit runs unit tests
func (m Magefile) TestUnit() {
	v := ""
	if mg.Verbose() {
		v = "-v"
	}
	must.Command("go", "test", v, "./...").CollapseArgs().RunV()
}

// Test runs a full suite of tests
func (m Magefile) Test() {
	m.TestUnit()

	// Check that we can call `plugin version`
	m.Build()
	must.RunV(filepath.Join(m.OutDir, m.PluginName+xplat.FileExt()), "version")
}

// Publish the plugin and its plugin feed
func (m Magefile) Publish() {
	mg.SerialDeps(m.PublishBinaries, m.PublishPluginFeed)
}

// PublishBinaries uploads cross-compiled binaries to a GitHub release
// Requires PORTER_RELEASE_REPOSITORY to be set to github.com/USERNAME/REPO
func (m Magefile) PublishBinaries() {
	mg.SerialDeps(porter.UseBinForPorterHome, porter.EnsurePorter)
	releases.PreparePluginForPublish(m.PluginName)
	releases.PublishPlugin(m.PluginName)
}

// Publish a plugin feed
// Requires PORTER_PACKAGES_REMOTE to be set to git@github.com:USERNAME/REPO.git
func (m Magefile) PublishPluginFeed() {
	mg.SerialDeps(porter.UseBinForPorterHome, porter.EnsurePorter)
	releases.PublishPluginFeed(m.PluginName)
}

// TestPublish uploads release artifacts to a fork of the pluin's repo.
// If your plugin is officially hosted in a repository under your username, you will need to manually
// override PORTER_RELEASE_REPOSITORY and PORTER_PACKAGES_REMOTE to test out publishing safely.
func (m Magefile) TestPublish(username string) {
	pluginRepo := fmt.Sprintf("github.com/%s/%s-plugins", username, m.PluginName)
	pkgRepo := fmt.Sprintf("https://github.com/%s/packages.git", username)
	fmt.Printf("Publishing a release to %s and committing a plugin feed to %s\n", pluginRepo, pkgRepo)
	fmt.Printf("If you use different repository names, set %s and %s then call mage Publish instead.\n", releases.ReleaseRepository, releases.PackagesRemote)
	os.Setenv(releases.ReleaseRepository, pluginRepo)
	os.Setenv(releases.PackagesRemote, pkgRepo)

	m.Publish()
}

// Install the plugin
func (m Magefile) Install() {
	porterHome := porter.GetPorterHome()
	fmt.Printf("Installing the %s plugin into %s\n", m.PluginName, porterHome)

	os.MkdirAll(filepath.Join(porterHome, "plugins", m.PluginName), 0770)
	mgx.Must(shx.Copy(filepath.Join(m.OutDir, m.PluginName+xplat.FileExt()), filepath.Join(porterHome, "plugins", m.PluginName)))
}

// Clean removes generated build files
func (m Magefile) Clean() {
	os.RemoveAll("bin")
}
