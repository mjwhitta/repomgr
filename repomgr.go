//nolint:wrapcheck // External package is in the same repo
package repomgr

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/mjwhitta/errors"
	"github.com/mjwhitta/pathname"
	"github.com/mjwhitta/repomgr/git"
)

// AddPush will add files, commit with the specified message, and
// push.
func AddPush(dir string, msg string, files ...string) error {
	if e := git.Add(dir, files...); e != nil {
		return e
	}

	if e := git.Commit(dir, msg); e != nil {
		return e
	}

	if e := git.Push(dir, false, ""); e != nil {
		return e
	}

	return nil
}

// CleanSlate will ensure the repo is back to a clean state (no
// changes, no additional files), optionally with all the newest
// changes.
func CleanSlate(
	dir string,
	fetch bool,
	remote string,
	branch ...string,
) error {
	if fetch {
		if e := git.Fetch(dir, remote); e != nil {
			return e
		}
	}

	if e := git.ResetHard(dir, remote, branch...); e != nil {
		return e
	}

	if e := git.Clean(dir); e != nil {
		return e
	}

	return nil
}

// DirName will return a directory name in the form of
// <owner>.<name>[.branch].
func DirName(r Repository) string {
	var branch string
	var name string
	var owner string

	_, owner, name, _ = git.RepoInfo(r.CloneURL())

	if branch = r.BranchName(); branch != "" {
		branch = strings.ReplaceAll(branch, "/", "__")
		branch = strings.ReplaceAll(branch, "\\", "__")

		name += "." + branch
	}

	return owner + "." + name
}

// Download will clone the repo into the provided directory. If the
// repo already exists, it will update it.
func Download(dir string, r Repository) error {
	var e error
	var path string
	var provider string

	//nolint:mnd // u=rwx,g=-,o=-
	if e = os.MkdirAll(dir, 0o700); e != nil {
		e = errors.Newf("directory %s cannot be created: %w", dir, e)
		return e
	}

	dir = filepath.Join(dir, DirName(r))
	provider, _, _, path = git.RepoInfo(r.CloneURL())

	if ok, _ := pathname.DoesExist(dir); !ok {
		if Logger != nil {
			_ = Logger.SubInfof("Cloning %s...", filepath.Base(dir))
		}

		e = git.Clone(
			filepath.Dir(dir),
			provider,
			path,
			filepath.Base(dir),
			r.BranchName(),
		)
	} else {
		if Logger != nil {
			_ = Logger.SubInfof("Updating %s...", filepath.Base(dir))
		}

		e = CleanSlate(dir, true, "", r.BranchName())
	}

	if e != nil {
		//nolint:wrapcheck // Submodule of this module
		return e
	}

	return nil
}

// DownloadAll will loop thru the provided repos and clone all of them
// into the provided directory.
func DownloadAll(dir string, repos []Repository) error {
	if len(repos) == 0 {
		return errors.New("no repos provided")
	}

	for _, repo := range repos {
		if e := Download(dir, repo); e != nil {
			return e
		}
	}

	return nil
}

// RollBack will reset the repo back to the initial commit and force
// push.
//
// WARNING: This will destroy data!
func RollBack(dir string, remote string, branch ...string) error {
	var e error

	if (len(branch) == 0) || (strings.TrimSpace(branch[0]) == "") {
		branch = []string{"HEAD"}
	}

	branch[0] = strings.TrimSpace(branch[0])
	remote = strings.TrimSpace(remote)

	if e = CleanSlate(dir, false, remote, branch[0]); e != nil {
		return e
	}

	for e == nil {
		e = git.ResetHard(dir, "", branch[0]+"~1")
	}

	if e = git.Push(dir, true, remote); e != nil {
		return e
	}

	return nil
}
