package lang_test

import (
	"testing"

	"github.com/yeqown/infrastructure/pkg/lang"
)

type sortSource struct {
	First  int
	Second int
}

func lessByFirst(ssa, ssb interface{}) bool {
	return ssa.(*sortSource).First < ssb.(*sortSource).First
}

func lessBySecond(ssa, ssb interface{}) bool {
	return ssa.(*sortSource).Second < ssb.(*sortSource).Second
}

func Test_MultiSorter(t *testing.T) {
	data := []interface{}{
		&sortSource{1, 2},
		&sortSource{0, 2},
		&sortSource{1, 0},
		&sortSource{2, 1},
		&sortSource{3, 2},
		&sortSource{4, 3},
		&sortSource{5, 4},
	}

	lang.OrderedBy(lessByFirst, lessBySecond).Sort(data)

	for _, v := range data {
		t.Log(v)
	}

	// output:
	// multi_sorter_test.go:36: &{0 2}
	// multi_sorter_test.go:36: &{1 0}
	// multi_sorter_test.go:36: &{1 2}
	// multi_sorter_test.go:36: &{2 1}
	// multi_sorter_test.go:36: &{3 2}
	// multi_sorter_test.go:36: &{4 3}
	// multi_sorter_test.go:36: &{5 4}

	lang.OrderedBy(lessBySecond, lessByFirst).Sort(data)
	for _, v := range data {
		t.Log(v)
	}
	// output:
	// multi_sorter_test.go:50: &{1 0}
	// multi_sorter_test.go:50: &{2 1}
	// multi_sorter_test.go:50: &{0 2}
	// multi_sorter_test.go:50: &{1 2}
	// multi_sorter_test.go:50: &{3 2}
	// multi_sorter_test.go:50: &{4 3}
	// multi_sorter_test.go:50: &{5 4}
}
