// Package chashing didn't think about migrating work
package chashing

import (
	"errors"
	"hash"
	"strings"

	"github.com/yeqown/infrastructure/pkg/lang"
)

const (
	defaultMappingVirtualNodes = 100
)

// ConsistentHashing .
type ConsistentHashing struct {
	nodesNum  int               // nodes count
	vnodesNum int               // virtual nodes num
	h         hash.Hash32       // hash32 algthrim
	nodes     map[string]Node   // all actual nodes
	vnodes    map[string]Node   // all virtual nodes
	mapping   map[string]string // virtual to actual
	// TODO:
}

// NewConsistentHashing .
func NewConsistentHashing(h hash.Hash32, nodes ...Node) (*ConsistentHashing, error) {
	if len(nodes) == 0 {
		return nil, errors.New("zero nodes to init")
	}
	if h == nil {
		// true: set default hash.Hash
		// TODO:
		// h =
	}
	ch := ConsistentHashing{
		nodesNum: len(nodes),
		h:        h,
		nodes:    make(map[string]Node, len(nodes)),
		vnodes:   make(map[string]Node, len(nodes)*defaultMappingVirtualNodes),
		mapping:  make(map[string]string, len(nodes)*defaultMappingVirtualNodes),
	}

	for idx := range nodes {
		if err := ch.AddNode(nodes[idx]); err != nil {
			return nil, err
		}
	}

	return &ch, nil
}

// Locate .
func (ch *ConsistentHashing) Locate(key string) Node {
	pos := ch.pos(key)
	vnode := ch.locate(pos)
	return ch.mappingTo(vnode)
}

// calc hash position in 2^32
func (ch *ConsistentHashing) pos(key string) uint32 {
	ch.h.Reset()
	ch.h.Write(lang.StringToBytes(key))
	return ch.h.Sum32()
}

// locate
// TODO: pined same key into same node
// FIXME: if part of nodes is Migrating, should stop locate to target nodes
func (ch *ConsistentHashing) locate(pos uint32) Node {
	return nil
}

// AddNode .
// valid node; valid node ident duplicated;
// FIXME: actual you should migrate data from last node of new node
func (ch *ConsistentHashing) AddNode(n Node) error {
	if err := validNode(n); err != nil {
		return err
	}

	if _, ok := ch.nodes[n.Ident()]; ok {
		return ErrDuplicatedNodeIdent
	}

	return ch.addNode(n)
}

// add virtual nodes; add mapping
func (ch *ConsistentHashing) addNode(n Node) error {
	var (
		vn    virtualNode
		ident = n.Ident()
	)

	ch.nodes[ident] = n

	for i := 1; i <= defaultMappingVirtualNodes; i++ {
		vn = newVirtualNode(ident, i)
		ch.vnodes[vn.Ident()] = vn
		ch.mapping[vn.Ident()] = ident
	}
	ch.vnodesNum += defaultMappingVirtualNodes

	return nil
}

// RemoveNode .
// FIXME: actual you should migrate data from removed node
func (ch *ConsistentHashing) RemoveNode(n Node) error {
	var ident = n.Ident()
	if _, ok := ch.nodes[ident]; !ok {
		return ErrNodeNotFound
	}

	delete(ch.nodes, ident)

	for vIdent, aIdent := range ch.mapping {
		if strings.Compare(aIdent, ident) == 0 {
			// trur: equal
			delete(ch.vnodes, vIdent)
			delete(ch.mapping, vIdent)
		}
	}

	return nil
}

// mapping from virtual node to actual node
func (ch *ConsistentHashing) mappingTo(vNode Node) Node {
	id := ch.mapping[vNode.Ident()]
	return ch.nodes[id]
}
