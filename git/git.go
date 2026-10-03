package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/mjwhitta/errors"
	"github.com/mjwhitta/log"
)

// Add will add provided files to git's index.
func Add(dir string, files ...string) error {
	for _, fn := range files {
		if _, e := Git(dir, "add", fn); e != nil {
			return errors.Newf("failed to add %s: %w", fn, e)
		}
	}

	return nil
}

// Branch will return the branch name for the current HEAD ref of the
// provided repo.
func Branch(dir string) (s string) {
	s, _ = Git(dir, "status", "-b", "--porcelain")
	s, _, _ = strings.Cut(s, "\n")
	s, _, _ = strings.Cut(s, "...")

	return strings.TrimSpace(strings.TrimPrefix(s, "##"))
}

// Clean will remove untracked files from the working tree.
func Clean(dir string) error {
	if _, e := Git(dir, "clean", "-d", "-f", "-x"); e != nil {
		return errors.Newf("failed to clean: %w", e)
	}

	return nil
}

// Clone will clone a git repo in the specified directory.
func Clone(dir, provider, path, name string, branch ...string) error {
	var args []string = []string{"clone"}
	var url string

	provider = strings.TrimSpace(provider)
	path = strings.TrimSpace(path)

	// Only support clone via SSH at this time
	url = "git@" + provider + ":" + path + ".git"

	if name = strings.TrimSpace(name); name == "" {
		name = filepath.Base(path)
		name = strings.TrimSuffix(name, ".git")
	}

	if isSet(branch...) {
		args = append(args, "-b", strings.TrimSpace(branch[0]))
	}

	args = append(args, url, name)

	if _, e := Git(dir, args...); e != nil {
		if isSet(branch...) {
			return errors.Newf(
				"failed to clone branch %s for %s: %w",
				branch[0],
				url,
				e,
			)
		}

		return errors.Newf("failed to clone %s: %w", url, e)
	}

	return nil
}

// Commit will record changes to the repository.
func Commit(dir string, msg string) error {
	if msg = strings.TrimSpace(msg); msg == "" {
		return errors.New("no commit message provided")
	}

	if _, e := Git(dir, "commit", "-m", msg); e != nil {
		return errors.Newf("failed to commit: %w", e)
	}

	return nil
}

// Exists will return whether or not a Git repo exists in the
// specified directory.
func Exists(dir string) bool {
	return Hash(dir) != ""
}

// Fetch will download objects and refs from another repository.
func Fetch(dir string, remote ...string) error {
	remote = []string{setOrDefault("origin", remote...)}

	if _, e := Git(dir, "fetch", remote[0]); e != nil {
		return e
	}

	return nil
}

// Git will called the git command line tool.
func Git(dir string, args ...string) (string, error) {
	var b []byte
	var cmd *exec.Cmd
	var e error
	var gitOpts []string
	var ssh string = sshCmd()
	var sshOpts []string = []string{ // No windows support
		// "-o",
		// "ControlMaster=auto",
		// "-o",
		// "ControlPath=\"TODO/%C\"",
		// "-o",
		// "ControlPersist=yes",
	}

	if SSHKey != "" {
		sshOpts = append([]string{"-i", SSHKey}, sshOpts...)
	}

	if dir = strings.TrimSpace(dir); dir != "" {
		gitOpts = []string{"-C", dir}
	}

	gitOpts = append(gitOpts, "-c", "core.filemode=false")

	gitOpts = append(
		gitOpts,
		"-c",
		"core.sshCommand="+ssh+" "+strings.Join(sshOpts, " "),
	)

	if Debug {
		if dir == "" {
			log.Debugf("git %s", strings.Join(args, " "))
		} else {
			log.Debugf("git -C %s %s", dir, strings.Join(args, " "))
		}
	}

	// Add all the crazy git options after we've logged the command
	args = append(gitOpts, args...)

	//nolint:gosec // G204 - by design
	cmd = exec.Command("git", args...)
	if b, e = cmd.CombinedOutput(); e != nil {
		//nolint:wrapcheck // I want to return the actual error
		return "", e
	}

	if Debug {
		log.Msg(strings.TrimSpace(string(b)))
	}

	return strings.TrimSpace(string(b)), nil
}

// Hash will return the hash for the current HEAD ref of the provided
// repo.
func Hash(dir string) (s string) {
	s, _ = Git(dir, "rev-parse", "HEAD")
	return s
}

// Init will initialize a new git repo with the provided name in the
// provided directory. If no name is provided, it is assumed that the
// provided directory is full path to the repo.
func Init(dir string, name ...string) error {
	var e error

	name = []string{setOrDefault(".", name...)}
	e = errors.Newf("invalid repo name %s", name[0])

	switch {
	case name[0] == "..":
		return e
	case strings.Contains(name[0], "/"):
		return e
	case strings.Contains(name[0], "\\"):
		return e
	}

	dir = filepath.Clean(filepath.Join(dir, name[0]))
	name[0] = filepath.Base(dir)

	//nolint:mnd // u=rwx,g=-,o=-
	if e := os.MkdirAll(dir, 0o700); e != nil {
		return errors.Newf("failed to create repo: %w", e)
	}

	if _, e := Git(dir, "init"); e != nil {
		return errors.Newf("failed to init %s: %w", name[0], e)
	}

	return nil
}

// NormalizeRepo will normalize a repo's URL for convenient parsing.
// The output is of the form <provider>/<owner>/<repo>.
func NormalizeRepo(repo string) string {
	var m []string

	// Remove whitespace and trailing slashes
	repo = strings.TrimSpace(repo)
	for strings.HasSuffix(repo, "/") {
		repo = strings.TrimSuffix(repo, "/")
	}

	// Remove .git and trailing slashes
	repo = strings.TrimSuffix(repo, ".git")
	for strings.HasSuffix(repo, "/") {
		repo = strings.TrimSuffix(repo, "/")
	}

	// If no provider, assume GitHub
	if strings.Count(repo, "/") == 1 {
		if strings.Count(repo, ":") == 0 {
			repo = "github.com/" + repo
		}
	}

	if m = reRepo.FindStringSubmatch(repo); len(m) > 0 {
		return m[3] + "/" + m[5]
	}

	// Invalid repo, just return the input, I guess
	return repo
}

// Path will return the path (owner/repo) for the repo in the
// specified directory.
func Path(dir string) (s string) {
	_, s, _ = strings.Cut(NormalizeRepo(URL(dir)), "/")
	return s
}

// Provider will return the git provider for the repo in the specified
// directory.
func Provider(dir string) (s string) {
	s, _, _ = strings.Cut(NormalizeRepo(URL(dir)), "/")
	return s
}

// Pull will fetch from and integrate with another repository or a
// local branch.
func Pull(dir string, remote string, branch ...string) error {
	var args []string = []string{"pull"}

	if remote = strings.TrimSpace(remote); remote != "" {
		args = append(args, remote)
	}

	if isSet(branch...) {
		args = append(args, branch[0])
	}

	if _, e := Git(dir, args...); e != nil {
		return errors.Newf("failed to pull: %w", e)
	}

	return nil
}

// Push will update remote refs along with associated objects.
func Push(
	dir string,
	force bool,
	remote string,
	branch ...string,
) error {
	var args []string = []string{"push"}

	if force {
		args = append(args, "-f")
	}

	if remote = strings.TrimSpace(remote); remote == "" {
		remote = "origin"
	}

	args = append(args, remote)

	if isSet(branch...) {
		args = append(args, branch[0])
	}

	if _, e := Git(dir, args...); e != nil {
		return errors.Newf("failed to push: %w", e)
	}

	return nil
}

// RepoInfo will return the Git provider, the Git repo's owner, the
// Git repo's name, and the path.
func RepoInfo(repo string) (string, string, string, string) {
	var name string
	var owner string
	var path string
	var provider string

	//nolint:mnd // 2 is the min number of slashes for a valid repo
	if repo = NormalizeRepo(repo); strings.Count(repo, "/") < 2 {
		return "", "", "", ""
	}

	provider, path, _ = strings.Cut(repo, "/")

	name = filepath.Base(path)
	owner, _, _ = strings.Cut(path, "/")

	return provider, owner, name, path
}

// Reset will set HEAD or the index to a known state.
func Reset(dir string) error {
	if _, e := Git(dir, "reset"); e != nil {
		return errors.Newf("failed to reset: %w", e)
	}

	return nil
}

// ResetHard will set HEAD or the index to a known state.
func ResetHard(dir string, remote string, branch ...string) error {
	var e error

	branch = []string{setOrDefault("HEAD", branch...)}

	if remote = strings.TrimSpace(remote); remote != "" {
		remote += "/"
	}

	_, e = Git(dir, "reset", "--hard", remote+branch[0])
	if e != nil {
		return errors.Newf("failed to reset: %w", e)
	}

	return nil
}

// Status will return all files with modifications waiting to be
// pushed.
func Status(dir string) (map[string][]string, error) {
	var changes map[string][]string
	var e error
	var fn string
	var state string
	var status string

	if status, e = Git(dir, "status", "--porcelain"); e != nil {
		return nil, e
	}

	changes = map[string][]string{
		"added":     {},
		"deleted":   {},
		"modified":  {},
		"renamed":   {},
		"untracked": {},
	}

	for _, s := range strings.Split(status, "\n") {
		s = strings.TrimSpace(s)
		state, fn, _ = strings.Cut(s, " ")

		fn = strings.TrimSpace(fn)
		state = strings.TrimSpace(state)

		switch state {
		case "??":
			changes["untracked"] = append(changes["untracked"], fn)
		case "A":
			changes["added"] = append(changes["added"], fn)
		case "D":
			changes["deleted"] = append(changes["deleted"], fn)
		case "M":
			changes["modified"] = append(changes["modified"], fn)
		case "R":
			changes["renamed"] = append(changes["renamed"], fn)
		}
	}

	slices.Sort(changes["added"])
	slices.Sort(changes["deleted"])
	slices.Sort(changes["modified"])
	slices.Sort(changes["renamed"])
	slices.Sort(changes["untracked"])

	return changes, nil
}

// URL will return the URL for the repo.
func URL(dir string, remote ...string) string {
	var e error
	var remotes string
	var ss []string
	var url string

	remote = []string{setOrDefault("origin", remote...)}

	remotes, e = Git(dir, "remote", "show", "-n", remote[0])
	if e != nil {
		return ""
	}

	for _, s := range strings.Split(remotes, "\n") {
		s = strings.TrimSpace(s)

		if strings.Contains(s, " URL:") {
			ss = strings.Fields(s)

			if url = ss[len(ss)-1]; url == remote[0] {
				url = ""
			}

			if url != "" {
				break
			}
		}
	}

	return url
}
