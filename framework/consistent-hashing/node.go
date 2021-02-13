package chashing

import (
	"errors"
	"strconv"
)

var (
	// ErrEmptyNodeIdent .
	ErrEmptyNodeIdent = errors.New("empty node ident")
	// ErrDuplicatedNodeIdent .
	ErrDuplicatedNodeIdent = errors.New("duplicated node ident")
	// ErrNodeNotFound .
	ErrNodeNotFound = errors.New("node not found by ident")
)

// Node .
type Node interface {
	Ident() string
}

func validNode(n Node) error {
	if n.Ident() == "" {
		return ErrEmptyNodeIdent
	}
	return nil
}

type virtualNode struct {
	nodeIdent string
	idx       int
}

func newVirtualNode(nodeIdent string, idx int) virtualNode {
	return virtualNode{
		nodeIdent: nodeIdent,
		idx:       idx,
	}
}

func (vn virtualNode) Ident() string {
	return vn.nodeIdent + "@@@" + strconv.Itoa(vn.idx)
}
