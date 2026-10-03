package repomgr

// Repo is a struct containing required info to clone a git
// repository.
type Repo struct {
	Branch string
	URL    string
}

// Verify interface compliance at compile time.
var _ Repository = (*Repo)(nil)

// BranchName will return the branch name for the repo or an empty
// string, if there is no branch.
func (r *Repo) BranchName() string {
	if r == nil {
		return ""
	}

	return r.Branch
}

// CloneURL will return the URL used to clone the repo.
func (r *Repo) CloneURL() string {
	if r == nil {
		return ""
	}

	return r.URL
}
