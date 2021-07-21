// +build !RELEASE

package perm

type FS int

const (
	PUBLIC       = 0777
	PUBLIC_VIEW  = 0777
	PUBLIC_WRITE = 0777
	PRIVATE      = 0777
)
