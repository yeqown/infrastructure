package fontutil_test

import (
	"testing"

	"github.com/yeqown/infrastructure/pkg/fontutil"
)

func Test_GetSysFontList(t *testing.T) {
	fonts := fontutil.GetSysFontList()
	t.Log(fonts)
}
