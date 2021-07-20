// +build RELEASE

package perm

type FS int

const (
	PUBLIC       = 0777
	PUBLIC_VIEW  = 0744
	PUBLIC_WRITE = 0766
	PRIVATE      = 0700
)
