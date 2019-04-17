// Package mirrors handles managing mirrors in the running application
package mirrors

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/vcs"
	"github.com/jeffzhangme/glide/msg"
	gpath "github.com/jeffzhangme/glide/path"
)

var mirrors map[string]*mirror

func init() {
	mirrors = make(map[string]*mirror)
}

type mirror struct {
	Repo, Vcs string
}

// Get retrieves information about an mirror. It returns.
// - bool if found
// - new repo location
// - vcs type
func Get(k string) (bool, string, string) {
	o, f := mirrors[k]
	if !f {
		o = new(mirror)
		if strings.HasPrefix(k, "https://golang.org/x") {
			msg.Debug("Setting %s mirror to %s", k, "https://github.com/golang")
			o.Repo = strings.Replace(k, "golang.org/x", "github.com/golang", 1)
			o.Vcs = string(vcs.Git)
		} else if strings.HasPrefix(k, "https://google.golang.org/grpc") {
			msg.Debug("Setting %s mirror to %s", k, "https://github.com/grpc/grpc-go")
			o.Repo = strings.Replace(k, "google.golang.org/grpc", "github.com/grpc/grpc-go", 1)
			o.Vcs = string(vcs.Git)
		} else if strings.HasPrefix(k, "https://google.golang.org/genproto") {
			msg.Debug("Setting %s mirror to %s", k, "https://github.com/google/go-genproto")
			o.Repo = strings.Replace(k, "google.golang.org/genproto", "github.com/google/go-genproto", 1)
			o.Vcs = string(vcs.Git)
		} else if strings.HasPrefix(k, "https://google.golang.org/api") {
			msg.Debug("Setting %s mirror to %s", k, "https://github.com/googleapis/google-api-go-client")
			o.Repo = strings.Replace(k, "google.golang.org/api", "github.com/googleapis/google-api-go-client", 1)
			o.Vcs = string(vcs.Git)
		} else {
			return false, "", ""
		}
	}

	return true, o.Repo, o.Vcs
}

// Load pulls the mirrors into memory
func Load() error {
	home := gpath.Home()

	op := filepath.Join(home, "mirrors.yaml")

	var ov *Mirrors
	if _, err := os.Stat(op); os.IsNotExist(err) {
		msg.Debug("No mirrors.yaml file exists")
		ov = &Mirrors{
			Repos: make(MirrorRepos, 0),
		}
		return nil
	} else if err != nil {
		ov = &Mirrors{
			Repos: make(MirrorRepos, 0),
		}
		return err
	}

	var err error
	ov, err = ReadMirrorsFile(op)
	if err != nil {
		return fmt.Errorf("Error reading existing mirrors.yaml file: %s", err)
	}

	msg.Info("Loading mirrors from mirrors.yaml file")
	for _, o := range ov.Repos {
		msg.Debug("Found mirror: %s to %s (%s)", o.Original, o.Repo, o.Vcs)
		no := &mirror{
			Repo: o.Repo,
			Vcs:  o.Vcs,
		}
		mirrors[o.Original] = no
	}

	return nil
}
