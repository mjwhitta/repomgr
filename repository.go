package repomgr

// Repository is an interface ensuring all repomgr operations are data
// agnostic.
type Repository interface {
	BranchName() string
	CloneURL() string
}
