// Package chashing didn't think about migrating work
package chashing

import (
	"errors"
	"fmt"
	"hash"
	"hash/crc32"
	"sort"
	"sync"
)

type ConsistentHash interface {
	AddNode(n Node) error
	RemoveNode(n Node) error

	HashKey(key []byte) Node
}

// Node represents a node in the consistent hashing ring.
type Node interface {
	Identity() []byte
	Replicas() uint8
}

var _ ConsistentHash = (*ConsistentHashing)(nil)

const (
	defaultVirtualNodeReplicas = 100
)

type ConsistentHashing struct {
	lock sync.RWMutex

	hash32 hash.Hash32 // hash32 algorithm
	nodes  []Node      // all actual nodes
	// TODO: merge nodeIds and nodeReplicas
	nodeIds       map[string]uint32 // all node identities for duplicated check; nodeId:nodeIdx
	nodesReplicas map[uint32]uint8  // all actual nodes setting; nodeIdx: replicas
	ring          []uint32          // ring of virtual nodes;    virtualNodeHash
	ringMapping   map[uint32]uint32 // ring mapping to actual nodes; virtualNode: nodeIdx
}

// New constructs a ConsistentHashing instance to manage nodes and
// locate key to target node.
func New(h hash.Hash32, initialNodes ...Node) (*ConsistentHashing, error) {
	if len(initialNodes) == 0 {
		return nil, errors.New("zero initialNodes to init")
	}
	if h == nil {
		h = crc32.NewIEEE()
	}

	n := len(initialNodes)

	ch := ConsistentHashing{
		hash32:        h,
		nodes:         make([]Node, 0, n),
		nodeIds:       make(map[string]uint32),
		nodesReplicas: make(map[uint32]uint8, n),
		ring:          make([]uint32, 0, n*defaultVirtualNodeReplicas),
		ringMapping:   make(map[uint32]uint32, n*defaultVirtualNodeReplicas),
	}

	for _, node := range initialNodes {
		if err := ch.addNode(node, false); err != nil {
			fmt.Printf("WARN: addNode(%s) failed: %v", node.Identity(), err.Error())
			continue
		}
	}

	// sort ring for binary search
	sort.Slice(ch.ring, func(i, j int) bool {
		return ch.ring[i] < ch.ring[j]
	})

	return &ch, nil
}

func (ch *ConsistentHashing) AddNode(n Node) error {
	return ch.addNode(n, true)
}

func (ch *ConsistentHashing) addNode(node Node, doSort bool) error {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	nodeId := node.Identity()
	if _, ok := ch.nodeIds[string(nodeId)]; ok {
		fmt.Println("WARN: duplicated node identity, skip it: ", string(nodeId))
		return errors.New("duplicated node identity")
	}

	nodeIdx := len(ch.nodes)
	ch.nodes = append(ch.nodes, node)

	// record node identity
	ch.nodeIds[string(nodeId)] = uint32(nodeIdx)

	// decide replicas for each node by setting or default
	replicas := node.Replicas()
	if replicas == 0 {
		replicas = defaultVirtualNodeReplicas
	}
	ch.nodesReplicas[uint32(nodeIdx)] = replicas

	// generate virtual nodes and ring mapping
	for i := uint8(0); i < replicas; i++ {
		vNodeKey := append(nodeId, '@', '@', byte(i))
		vNodeHash := crc32.ChecksumIEEE(vNodeKey)
		ch.ring = append(ch.ring, vNodeHash)
		ch.ringMapping[vNodeHash] = uint32(nodeIdx)
	}

	if doSort {
		// sort ring for binary search
		sort.Slice(ch.ring, func(i, j int) bool {
			return ch.ring[i] < ch.ring[j]
		})
	}

	return nil
}

func (ch *ConsistentHashing) RemoveNode(n Node) error {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	nodeId := n.Identity()
	if _, ok := ch.nodeIds[string(nodeId)]; !ok {
		return errors.New("node not found")
	}

	nodeIdx := ch.nodeIds[string(nodeId)]

	// remove actual node from nodes
	ch.nodes = append(ch.nodes[:nodeIdx], ch.nodes[nodeIdx+1:]...)
	// remove node identity from nodeIds
	delete(ch.nodeIds, string(nodeId))
	// remove replicas setting for each node by setting or default
	replicas := ch.nodesReplicas[nodeIdx]
	delete(ch.nodesReplicas, nodeIdx)

	// remove virtual nodes and ring mapping
	newRing := make([]uint32, 0, len(ch.ring)-int(replicas))
	for i := 0; i < len(ch.ring); i++ {
		if ch.ringMapping[ch.ring[i]] != nodeIdx {
			// keep the virtual node if it's not belong to the removed node
			newRing = append(newRing, ch.ring[i])
			continue
		}

		// remove the virtual node from ringMapping
		delete(ch.ringMapping, ch.ring[i])
	}

	return nil
}

func (ch *ConsistentHashing) HashKey(key []byte) Node {
	h := crc32.ChecksumIEEE(key)
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	// binary search to locate the first node that >= h
	idx := sort.Search(len(ch.ring), func(i int) bool {
		return ch.ring[i] >= h
	})

	// if idx == len(ch.ring), then idx = 0
	if idx == len(ch.ring) {
		idx = 0
	}

	// get the actual node index
	nodeIdx := ch.ringMapping[ch.ring[idx]]
	return ch.nodes[nodeIdx]
}

type builtinNode struct {
	id       []byte
	replicas uint8
}

func (b builtinNode) Identity() []byte { return b.id }
func (b builtinNode) Replicas() uint8  { return b.replicas }

func NewNode(id []byte, replicas uint8) Node {
	// empty id is not allowed
	if len(id) == 0 || allZero(id) {
		panic("empty id is not allowed")
	}

	return builtinNode{
		id:       id,
		replicas: replicas,
	}
}

func allZero(b []byte) bool {
	for _, v := range b {
		if v != 0 {
			return false
		}
	}
	return true
}
