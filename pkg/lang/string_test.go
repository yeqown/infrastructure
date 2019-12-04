package lang

import (
	"bytes"
	"strings"
	"testing"
)

func Test_String2Byte(t *testing.T) {
	var (
		s    = "teststring"
		tar1 = []byte(s)
	)

	tar2 := StringToBytes(s)
	if !bytes.Equal(tar1, tar2) {
		t.Errorf("want %v, got %v", tar1, tar2)
		t.FailNow()
	}
}

func Test_Byte2String(t *testing.T) {
	var (
		s1  = "teststring"
		tar = []byte(s1)
	)

	s2 := BytesToString(tar)
	if r := strings.Compare(s1, s2); r != 0 {
		t.Errorf("want %v, got %v", s1, s2)
		t.FailNow()
	}
}

func Benchmark_builtin_Convert(b *testing.B) {
	/*
		goos: darwin
		goarch: amd64
		pkg: github.com/yeqown/infrastructure/pkg/lang
		Benchmark_builtin_Convert-4   	39206592	        30.6 ns/op	      16 B/op	       1 allocs/op
		PASS
		ok  	github.com/yeqown/infrastructure/pkg/lang	3.715s
		Success: Benchmarks passed.
	*/
	var (
		s   = "teststring"
		tar = []byte{}
	)

	for i := 0; i < b.N; i++ {
		tar = []byte(s)
		_ = tar
	}
}

func Benchmark_unsafe_Convert(b *testing.B) {
	/*
		goos: darwin
		goarch: amd64
		pkg: github.com/yeqown/infrastructure/pkg/lang
		Benchmark_unsafe_Convert-4   	1000000000	         0.521 ns/op	       0 B/op	       0 allocs/op
		PASS
		ok  	github.com/yeqown/infrastructure/pkg/lang	0.896s
		Success: Benchmarks passed.
	*/
	var (
		s   = "teststring"
		tar = []byte{}
	)

	for i := 0; i < b.N; i++ {
		tar = StringToBytes(s)
		_ = tar
	}
}

func Benchmark_builtin_Convert_reverse(b *testing.B) {
	/*
		goos: darwin
		goarch: amd64
		pkg: github.com/yeqown/infrastructure/pkg/lang
		Benchmark_builtin_Convert_reverse-4   	43372372	        25.1 ns/op	      16 B/op	       1 allocs/op
		PASS
		ok  	github.com/yeqown/infrastructure/pkg/lang	2.719s
		Success: Benchmarks passed.
	*/
	var (
		s1  = "teststring"
		tar = []byte(s1)
		s2  string
	)

	for i := 0; i < b.N; i++ {
		s2 = string(tar)
		_ = s2
	}
}

func Benchmark_unsafe_Convert_reverse(b *testing.B) {
	/*
				goos: darwin
		goarch: amd64
		pkg: github.com/yeqown/infrastructure/pkg/lang
		Benchmark_unsafe_Convert_reverse-4   	1000000000	         0.349 ns/op	       0 B/op	       0 allocs/op
		PASS
		ok  	github.com/yeqown/infrastructure/pkg/lang	0.515s
		Success: Benchmarks passed.
	*/
	var (
		s1  = "teststring"
		tar = []byte(s1)
		s2  string
	)

	for i := 0; i < b.N; i++ {
		s2 = BytesToString(tar)
		_ = s2
	}
}
