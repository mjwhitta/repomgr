//nolint:godoclint // These are tests
package repomgr_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mjwhitta/log"
	"github.com/mjwhitta/repomgr"
	"github.com/mjwhitta/repomgr/git"
	"github.com/mjwhitta/repomgr/testhelper"
	assert "github.com/stretchr/testify/require"
)

var repo string = filepath.Join("testdata", "dummy")

func TestAddPush(t *testing.T) {
	t.Run(
		"FailAdd",
		func(t *testing.T) {
			assert.NoError(t, testhelper.Setup(t, "."))
			assert.Error(t, repomgr.AddPush(repo, "test", "test"))
		},
	)

	t.Run(
		"FailCommit",
		func(t *testing.T) {
			assert.NoError(t, testhelper.Setup(t, "."))
			assert.Error(t, repomgr.AddPush(repo, ""))
		},
	)

	t.Run(
		"FailPush",
		func(t *testing.T) {
			var e error

			assert.NoError(t, testhelper.Setup(t, "."))

			e = repomgr.AddPush(
				repo,
				"test",
				"dummy",
				"dummyRemoved",
				"dummyUntracked",
			)
			assert.Error(t, e)
		},
	)
}

func TestCleanSlate(t *testing.T) {
	t.Run(
		"FailFetch",
		func(t *testing.T) {
			assert.NoError(t, testhelper.Setup(t, "."))
			assert.Error(t, repomgr.CleanSlate(repo, true, ""))
		},
	)

	t.Run(
		"FailReset",
		func(t *testing.T) {
			assert.NoError(t, testhelper.Setup(t, "."))
			assert.Error(t, repomgr.CleanSlate(repo, false, "origin"))
		},
	)

	t.Run(
		"Success",
		func(t *testing.T) {
			assert.NoError(t, testhelper.Setup(t, "."))
			assert.NoError(t, repomgr.CleanSlate(repo, false, ""))
		},
	)
}

func TestDirName(t *testing.T) {
	var r *repomgr.Repo

	assert.Equal(t, ".", repomgr.DirName(r))

	r = &repomgr.Repo{URL: "git@github.com:mjwhitta/repomgr.git"}
	assert.Equal(t, "mjwhitta.repomgr", repomgr.DirName(r))

	r = &repomgr.Repo{
		Branch: "main",
		URL:    "git@github.com:mjwhitta/repomgr.git",
	}
	assert.Equal(t, "mjwhitta.repomgr.main", repomgr.DirName(r))

	r = &repomgr.Repo{
		Branch: "mw/feature/name",
		URL:    "git@github.com:mjwhitta/repomgr.git",
	}
	assert.Equal(
		t,
		"mjwhitta.repomgr.mw__feature__name",
		repomgr.DirName(r),
	)
}

func TestDownloadAll(t *testing.T) {
	t.Run(
		"FailClone",
		func(t *testing.T) {
			var repos []repomgr.Repository = []repomgr.Repository{
				&repomgr.Repo{
					URL: "git@github.com:mjwhitta/repomgr.git",
				},
			}

			assert.Error(t, repomgr.DownloadAll("/", repos))
		},
	)

	t.Run(
		"FailMkdir",
		func(t *testing.T) {
			var repos []repomgr.Repository = []repomgr.Repository{
				&repomgr.Repo{
					URL: "git@github.com:mjwhitta/repomgr.git",
				},
			}

			assert.Error(t, repomgr.DownloadAll("/noexist", repos))
		},
	)

	t.Run(
		"FailNoRepos",
		func(t *testing.T) {
			assert.Error(t, repomgr.DownloadAll("testdata", nil))
		},
	)

	t.Run(
		"Success",
		func(t *testing.T) {
			var e error
			var r *repomgr.Repo = &repomgr.Repo{
				URL: "git@github.com:mjwhitta/repomgr.git",
			}
			var repos []repomgr.Repository = []repomgr.Repository{r}

			// For coverage
			repomgr.Logger = log.NewMessenger()
			repomgr.Logger.Stdout = false

			e = os.RemoveAll(
				filepath.Join("testdata", "mjwhitta.repomgr"),
			)
			assert.NoError(t, e)

			// Once to clone
			e = repomgr.DownloadAll("testdata", repos)
			assert.NoError(t, e)

			// Twice to update
			e = repomgr.DownloadAll("testdata", repos)
			assert.NoError(t, e)
		},
	)
}

func TestRollBack(t *testing.T) {
	t.Run(
		"FailCleanSlate",
		func(t *testing.T) {
			assert.NoError(t, testhelper.Setup(t, "."))
			assert.Error(t, repomgr.RollBack("/noexist", ""))
		},
	)

	t.Run(
		"FailPush",
		func(t *testing.T) {
			var changes map[string][]string
			var e error

			assert.NoError(t, testhelper.Setup(t, "."))
			assert.Error(t, repomgr.RollBack(repo, ""))

			changes, e = git.Status(repo)
			assert.NoError(t, e)
			assert.Empty(t, changes["added"])
			assert.Empty(t, changes["deleted"])
			assert.Empty(t, changes["modified"])
			assert.Empty(t, changes["renamed"])
			assert.Empty(t, changes["untracked"])
		},
	)
}
