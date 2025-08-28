package porter

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/uwu-tools/magex/mgx"
	"github.com/uwu-tools/magex/pkg"
	"github.com/uwu-tools/magex/pkg/downloads"
	"github.com/uwu-tools/magex/shx"
	"github.com/uwu-tools/magex/xplat"
)

var (
	// DefaultPorterVersion is the default version of Porter that is installed when it's not present
	DefaultPorterVersion = "v1.3.0"
)

// Install the default version of porter, if porter isn't already installed
func EnsurePorter() {
	home := GetPorterHome()
	clientPath := filepath.Join(home, "porter"+xplat.FileExt())
	if ok, _ := pkg.IsCommandAvailable(clientPath, "--version", ""); !ok {
		EnsurePorterAt(DefaultPorterVersion)
	}
}

// Install the specified version of porter
func EnsurePorterAt(version string) {
	home := GetPorterHome()
	runtimesDir := filepath.Join(home, "runtimes")
	mgx.Must(os.MkdirAll(runtimesDir, 0770))

	var forceDownloadRuntime bool

	clientPath := filepath.Join(home, "porter"+xplat.FileExt())
	if clientFound, _ := pkg.IsCommandAvailable(clientPath, "--version", version); !clientFound {
		// When we download a new client, always download a new runtime
		forceDownloadRuntime = true

		log.Printf("Porter %s not found at %s\n", version, clientPath)
		log.Println("Installing porter into", home)
		opts := downloads.DownloadOptions{
			UrlTemplate: "https://cdn.porter.sh/{{.VERSION}}/porter-{{.GOOS}}-{{.GOARCH}}{{.EXT}}",
			Name:        "porter",
			Version:     version,
			Ext:         xplat.FileExt(),
		}
		if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
			// we don't yet publish arm64 for porter
			opts.UrlTemplate = "https://cdn.porter.sh/{{.VERSION}}/porter-darwin-amd64"
		}
		mgx.Must(downloads.Download(home, opts))
	}

	// Only check if the runtime file _exists_ but don't try to check the version
	// The binary isn't runnable in all cases because it will always be a linux binary, and the client could be windows or mac
	// So we can't really check the version here...
	runtimePath := filepath.Join(home, "runtimes", "porter-runtime")
	if _, err := os.Stat(runtimePath); forceDownloadRuntime || os.IsNotExist(err) {
		opts := downloads.DownloadOptions{
			UrlTemplate: "https://cdn.porter.sh/{{.VERSION}}/porter-linux-amd64",
			Name:        "porter-runtime",
			Version:     version,
		}
		mgx.Must(downloads.Download(runtimesDir, opts))
	}
}

type InstallMixinOptions struct {
	Name    string
	URL     string
	Feed    string
	Version string
}

// EnsureMixin installs the specified mixin.
func EnsureMixin(mixin InstallMixinOptions) error {
	home := GetPorterHome()
	mixinDir := filepath.Join(home, "mixins", mixin.Name)
	if _, err := os.Stat(mixinDir); err == nil {
		log.Println("Mixin already installed:", mixin.Name)
		return nil
	}

	log.Println("Installing mixin:", mixin.Name)
	if mixin.Version == "" {
		mixin.Version = DefaultMixinVersion
	}
	var source string
	if mixin.Feed != "" {
		source = "--feed-url=" + mixin.Feed
	} else {
		source = "--url=" + mixin.URL
	}

	porterPath := filepath.Join(home, "porter")
	return shx.Run(porterPath, "mixin", "install", mixin.Name, "--version", mixin.Version, source)
}
