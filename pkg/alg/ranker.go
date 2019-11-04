package alg

import (
	"math"
	"sort"
)

// Data .
type Data interface {
	Score() float64
}

// Ranker .
type Ranker interface {
	Rank(score float64) int
}

// NewRanker .
func NewRanker(d []float64, delim int) Ranker {
	r := &StdRanker{
		indices: make(map[float64]*scoreAcc),
		// scores:  make(scores, len(d)),
		delim: delim,
		pow10: math.Pow10(delim),
	}

	for _, s := range d {
		s = r.truncFloat64(s)
		if _, ok := r.indices[s]; !ok {
			r.indices[s] = new(scoreAcc)
			r.indices[s].score = s
		}
		r.indices[s].cnt++
	}

	r.scores = make(scores, len(r.indices))
	idx := 0
	for _, acc := range r.indices {
		acc.index = idx
		r.scores[idx] = acc
		idx++
	}

	// log.Println(r.scores, r.indices)
	sort.Sort(r.scores)
	// log.Println(r.scores, r.indices)

	return r
}

type scoreAcc struct {
	score float64 // score
	cnt   int     // cnt
	index int     // index
}

type scores []*scoreAcc

func (s scores) Len() int           { return len(s) }
func (s scores) Less(i, j int) bool { return s[i].score > s[j].score } // 降序排序
func (s scores) Swap(i, j int) {
	s[i].index = j
	s[j].index = i
	s[i], s[j] = s[j], s[i]
}

// StdRanker .
type StdRanker struct {
	indices map[float64]*scoreAcc
	scores  scores
	delim   int
	pow10   float64
}

// Rank .
func (r *StdRanker) Rank(score float64) int {
	score = r.truncFloat64(score)
	v, ok := r.indices[score]
	if !ok {
		// true: no such score
		return -1
	}

	_ = v

	rank := v.index + 1
	lastIdx := v.index - 1
	if lastIdx < 0 {
		// true: v.index = 0
		// do nothing
	} else {
		rank += r.scores[lastIdx].cnt
	}

	return rank
}

func (r *StdRanker) truncFloat64(f float64) float64 {
	return math.Trunc((f+0.5/r.pow10)*r.pow10) / r.pow10
}
