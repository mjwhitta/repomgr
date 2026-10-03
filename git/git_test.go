//nolint:godoclint // These are tests
package git_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mjwhitta/repomgr/git"
	"github.com/mjwhitta/repomgr/testhelper"
	assert "github.com/stretchr/testify/require"
)

var repo string = filepath.Join("testdata", "dummy")

func init() {
	// For coverage
	git.Debug = true
	git.SSHKey = "/tmp/fake.key"
}

func TestAdd(t *testing.T) {
	t.Run(
		"Fail",
		func(t *testing.T) {
			assert.NoError(t, testhelper.Setup(t, ".."))
			assert.Error(t, git.Add(repo, "test"))
		},
	)

	t.Run(
		"Success",
		func(t *testing.T) {
			assert.NoError(t, testhelper.Setup(t, ".."))
			assert.NoError(t, git.Add(repo, "dummyUntracked"))
		},
	)
}

func TestBranch(t *testing.T) {
	assert.NoError(t, testhelper.Setup(t, ".."))
	assert.Equal(t, "dummy", git.Branch(repo))
	assert.Equal(t, "main", git.Branch("")) // For coverage
}

func TestClean(t *testing.T) {
	t.Run(
		"FailNotRepo",
		func(t *testing.T) {
			assert.Error(t, git.Clean("/noexist"))
		},
	)

	t.Run(
		"Success",
		func(t *testing.T) {
			var changes map[string][]string
			var e error

			assert.NoError(t, testhelper.Setup(t, ".."))

			assert.NoError(t, git.Clean(repo))

			changes, e = git.Status(repo)
			assert.NoError(t, e)
			assert.Len(t, changes["added"], 1)
			assert.Len(t, changes["deleted"], 1)
			assert.Len(t, changes["modified"], 1)
			assert.Empty(t, changes["renamed"])
			assert.Empty(t, changes["untracked"])
		},
	)
}

func TestClone(t *testing.T) {
	t.Run(
		"FailDummyBranch",
		func(t *testing.T) {
			var e error

			assert.NoError(t, testhelper.Setup(t, ".."))

			e = git.Clone(
				"testdata",
				"github.com",
				"mjwhitta/repomgr",
				"dummy",
				"dummy",
			)
			assert.Error(t, e)
		},
	)

	t.Run(
		"FailDummyRepo",
		func(t *testing.T) {
			var e error

			assert.NoError(t, testhelper.Setup(t, ".."))

			e = git.Clone(
				"testdata",
				"github.com",
				"mjwhitta/dummy",
				"",
			)
			assert.Error(t, e)
		},
	)

	t.Run(
		"FailMainBranch",
		func(t *testing.T) {
			var e error

			assert.NoError(t, testhelper.Setup(t, ".."))

			e = git.Clone(
				"testdata",
				"github.com",
				"mjwhitta/repomgr",
				"dummy",
			)
			assert.Error(t, e)
		},
	)

	t.Run(
		"Success",
		func(t *testing.T) {
			var dir string = "testdata"
			var name string = "dummy"
			var path string = "mjwhitta/repomgr"
			var provider string = "github.com"

			assert.NoError(t, os.RemoveAll(repo))
			assert.NoError(t, git.Clone(dir, provider, path, name))

			// For coverage
			assert.NoError(t, git.Fetch(repo))
			assert.NoError(t, git.Pull(repo, ""))

			// For coverage, but only I can push to this repo, so I
			// can't assert.
			_ = git.Push(repo, false, "", "main")
		},
	)
}

func TestCommit(t *testing.T) {
	t.Run(
		"FailNoChanges",
		func(t *testing.T) {
			assert.NoError(t, testhelper.Setup(t, ".."))
			assert.NoError(t, git.Reset(repo)) // For coverage
			assert.NoError(t, git.ResetHard(repo, ""))
			assert.Error(t, git.Commit(repo, "test"))
		},
	)

	t.Run(
		"FailNoMessage",
		func(t *testing.T) {
			assert.NoError(t, testhelper.Setup(t, ".."))
			assert.Error(t, git.Commit(repo, ""))
		},
	)

	t.Run(
		"Success",
		func(t *testing.T) {
			var changes map[string][]string
			var e error

			assert.NoError(t, testhelper.Setup(t, ".."))

			assert.NoError(t, git.Commit(repo, "test"))

			changes, e = git.Status(repo)
			assert.NoError(t, e)
			assert.Empty(t, changes["added"])
			assert.Len(t, changes["deleted"], 1)
			assert.Len(t, changes["modified"], 1)
			assert.Empty(t, changes["renamed"])
			assert.Len(t, changes["untracked"], 1)
		},
	)
}

func TestExists(t *testing.T) {
	assert.NoError(t, testhelper.Setup(t, ".."))
	assert.True(t, git.Exists(repo))
}

func TestFetch(t *testing.T) {
	assert.NoError(t, testhelper.Setup(t, ".."))
	assert.Error(t, git.Fetch(repo))
}

func TestInit(t *testing.T) {
	t.Run(
		"FailBadName",
		func(t *testing.T) {
			assert.NoError(t, testhelper.Setup(t, ".."))
			assert.Error(t, git.Init("testdata", "dummy/test"))
			assert.Error(t, git.Init("testdata", "dummy\\testtest"))
			assert.Error(t, git.Init("testdata", ".."))
		},
	)

	t.Run(
		"FailMkdir",
		func(t *testing.T) {
			assert.Error(t, git.Init("/noexist", "dummy"))
		},
	)

	t.Run(
		"FailNotRepo",
		func(t *testing.T) {
			assert.Error(t, git.Init("/"))
		},
	)

	t.Run(
		"Success",
		func(t *testing.T) {
			assert.NoError(t, testhelper.Setup(t, ".."))
			assert.NoError(t, git.Init("testdata", "dummy"))
			assert.NoError(t, git.Init("testdata/dummy"))
		},
	)
}

func TestNormalizeRepo(t *testing.T) {
	t.Run(
		"HTTPS",
		func(t *testing.T) {
			assert.Equal(
				t,
				"github.com/mjwhitta/repomgr",
				git.NormalizeRepo(
					"https://github.com/mjwhitta/repomgr",
				),
			)
		},
	)

	t.Run(
		"Invalid",
		func(t *testing.T) {
			assert.Equal(t, "mjwhitta", git.NormalizeRepo("mjwhitta"))
		},
	)

	t.Run(
		"Long",
		func(t *testing.T) {
			assert.Equal(
				t,
				"github.com/mjwhitta/path/to/repo",
				git.NormalizeRepo(
					"github.com/mjwhitta/path/to/repo",
				),
			)
		},
	)

	t.Run(
		"NoProvider",
		func(t *testing.T) {
			assert.Equal(
				t,
				"github.com/mjwhitta/repomgr",
				git.NormalizeRepo("mjwhitta/repomgr"),
			)
		},
	)

	t.Run(
		"Provider",
		func(t *testing.T) {
			assert.Equal(
				t,
				"gitlab.com/mjwhitta/repomgr",
				git.NormalizeRepo("gitlab.com/mjwhitta/repomgr"),
			)
		},
	)

	t.Run(
		"SSH1",
		func(t *testing.T) {
			assert.Equal(
				t,
				"github.com/mjwhitta/repomgr",
				git.NormalizeRepo("git@github.com:mjwhitta/repomgr"),
			)
		},
	)

	t.Run(
		"SSH2",
		func(t *testing.T) {
			assert.Equal(
				t,
				"github.com/mjwhitta/repomgr",
				git.NormalizeRepo(
					"ssh://git@github.com:22/mjwhitta/repomgr",
				),
			)
		},
	)

	t.Run(
		"TrailingSlashes",
		func(t *testing.T) {
			assert.Equal(
				t,
				"github.com/mjwhitta/repomgr",
				git.NormalizeRepo("mjwhitta/repomgr///"),
			)

			assert.Equal(
				t,
				"github.com/mjwhitta/repomgr",
				git.NormalizeRepo("mjwhitta/repomgr///.git///"),
			)
		},
	)
}

func TestPath(t *testing.T) {
	assert.NoError(t, testhelper.Setup(t, ".."))
	assert.Equal(t, "notmjwhitta/dummy", git.Path(repo))
}

func TestProvider(t *testing.T) {
	assert.NoError(t, testhelper.Setup(t, ".."))
	assert.Equal(t, "github.com", git.Provider(repo))
}

func TestPull(t *testing.T) {
	assert.NoError(t, testhelper.Setup(t, ".."))
	assert.Error(t, git.Pull(repo, "origin", "main"))
}

func TestPush(t *testing.T) {
	assert.NoError(t, testhelper.Setup(t, ".."))
	assert.Error(t, git.Push(repo, true, ""))
}

func TestRepoInfo(t *testing.T) {
	t.Run(
		"HTTPS",
		func(t *testing.T) {
			var name string
			var owner string
			var path string
			var provider string

			provider, owner, name, path = git.RepoInfo(
				"https://github.com/mjwhitta/repomgr",
			)

			assert.Equal(t, "github.com", provider)
			assert.Equal(t, "mjwhitta", owner)
			assert.Equal(t, "repomgr", name)
			assert.Equal(t, "mjwhitta/repomgr", path)
		},
	)

	t.Run(
		"Invalid",
		func(t *testing.T) {
			var name string
			var owner string
			var path string
			var provider string

			provider, owner, name, path = git.RepoInfo("mjwhitta")

			assert.Empty(t, provider)
			assert.Empty(t, owner)
			assert.Empty(t, name)
			assert.Empty(t, path)
		},
	)

	t.Run(
		"Long",
		func(t *testing.T) {
			var name string
			var owner string
			var path string
			var provider string

			provider, owner, name, path = git.RepoInfo(
				"github.com/mjwhitta/path/to/repo",
			)

			assert.Equal(t, "github.com", provider)
			assert.Equal(t, "mjwhitta", owner)
			assert.Equal(t, "repo", name)
			assert.Equal(t, "mjwhitta/path/to/repo", path)
		},
	)

	t.Run(
		"NoProvider",
		func(t *testing.T) {
			var name string
			var owner string
			var path string
			var provider string

			provider, owner, name, path = git.RepoInfo(
				"mjwhitta/repomgr",
			)

			assert.Equal(t, "github.com", provider)
			assert.Equal(t, "mjwhitta", owner)
			assert.Equal(t, "repomgr", name)
			assert.Equal(t, "mjwhitta/repomgr", path)
		},
	)

	t.Run(
		"Provider",
		func(t *testing.T) {
			var name string
			var owner string
			var path string
			var provider string

			provider, owner, name, path = git.RepoInfo(
				"gitlab.com/mjwhitta/repomgr",
			)

			assert.Equal(t, "gitlab.com", provider)
			assert.Equal(t, "mjwhitta", owner)
			assert.Equal(t, "repomgr", name)
			assert.Equal(t, "mjwhitta/repomgr", path)
		},
	)

	t.Run(
		"SSH1",
		func(t *testing.T) {
			var name string
			var owner string
			var path string
			var provider string

			provider, owner, name, path = git.RepoInfo(
				"git@github.com:mjwhitta/repomgr",
			)

			assert.Equal(t, "github.com", provider)
			assert.Equal(t, "mjwhitta", owner)
			assert.Equal(t, "repomgr", name)
			assert.Equal(t, "mjwhitta/repomgr", path)
		},
	)

	t.Run(
		"SSH2",
		func(t *testing.T) {
			var name string
			var owner string
			var path string
			var provider string

			provider, owner, name, path = git.RepoInfo(
				"ssh://git@github.com:22/mjwhitta/repomgr",
			)

			assert.Equal(t, "github.com", provider)
			assert.Equal(t, "mjwhitta", owner)
			assert.Equal(t, "repomgr", name)
			assert.Equal(t, "mjwhitta/repomgr", path)
		},
	)

	t.Run(
		"TrailingSlashes",
		func(t *testing.T) {
			var name string
			var owner string
			var path string
			var provider string

			provider, owner, name, path = git.RepoInfo(
				"mjwhitta/repomgr///",
			)

			assert.Equal(t, "github.com", provider)
			assert.Equal(t, "mjwhitta", owner)
			assert.Equal(t, "repomgr", name)
			assert.Equal(t, "mjwhitta/repomgr", path)
		},
	)
}

func TestReset(t *testing.T) {
	assert.Error(t, git.Reset("/noexist"))
}

func TestResetHard(t *testing.T) {
	assert.NoError(t, testhelper.Setup(t, ".."))
	assert.Error(t, git.ResetHard(repo, "origin"))
}

func TestStatus(t *testing.T) {
	t.Run(
		"FailNotRepo",
		func(t *testing.T) {
			var e error

			_, e = git.Status("/noexist")
			assert.Error(t, e)
		},
	)

	t.Run(
		"Success",
		func(t *testing.T) {
			var changes map[string][]string
			var e error

			assert.NoError(t, testhelper.Setup(t, ".."))

			changes, e = git.Status(repo)
			assert.NoError(t, e)
			assert.Len(t, changes["added"], 1)
			assert.Len(t, changes["deleted"], 1)
			assert.Len(t, changes["modified"], 1)
			assert.Empty(t, changes["renamed"])
			assert.Len(t, changes["untracked"], 1)

			e = git.Add(repo, "dummyAdded", "dummyRemoved")
			assert.NoError(t, e)

			changes, e = git.Status(repo)
			assert.NoError(t, e)
			assert.Empty(t, changes["added"])
			assert.Empty(t, changes["deleted"])
			assert.Len(t, changes["modified"], 1)
			assert.Len(t, changes["renamed"], 1)
			assert.Len(t, changes["untracked"], 1)
		},
	)
}

func TestURL(t *testing.T) {
	t.Run(
		"FailNotRepo",
		func(t *testing.T) {
			assert.Empty(t, git.URL("/noexist"))
		},
	)

	t.Run(
		"NoRemote",
		func(t *testing.T) {
			assert.NoError(t, os.RemoveAll(repo))
			assert.NoError(t, git.Init(repo))
			assert.Empty(t, git.URL(repo))
		},
	)

	t.Run(
		"PushURL",
		func(t *testing.T) {
			var args []string
			var e error
			var url string = "git@github.com:notmjwhitta/dummy.git"

			assert.NoError(t, testhelper.Setup(t, ".."))

			args = []string{"remote", "set-url", "origin"}
			_, e = git.Git(repo, append(args, "origin")...)
			assert.NoError(t, e)

			args = []string{"remote", "set-url", "origin", "--push"}
			_, e = git.Git(repo, append(args, url)...)
			assert.NoError(t, e)

			assert.Equal(t, url, git.URL(repo))
		},
	)
}
