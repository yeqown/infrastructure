package chashing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsistentHashing(t *testing.T) {
	ch, err := New(nil)
	assert.NoError(t, err)

	// Test AddNode
	node1 := NewNode([]byte("node1"), 20)
	err = ch.AddNode(node1)
	assert.NoError(t, err)

	// Test HashKey
	node := ch.HashKey([]byte("key1"))
	assert.Equal(t, node1, node)

	// Test RemoveNode
	err = ch.RemoveNode(node1)
	assert.NoError(t, err)

	// Test node1 is removed
	node = ch.HashKey([]byte("key1"))
	assert.NotEqual(t, node1, node)
}
