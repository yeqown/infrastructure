package fs_test

import (
	"testing"

	"github.com/yeqown/infrastructure/pkg/fs"
)

func Test_ListFiles(t *testing.T) {
	fns := fs.ListFiles(".", fs.IgnoreDirFilter())
	t.Log(fns)
}
