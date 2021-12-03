package im

import (
	"sync"
)

type Node struct {
	len   int
	nodes []*sync.Map
}

func NewNode(len int) *Node {
	return &Node{len: len, nodes: maps(len)}
}

func (n *Node) index(cid int64) int {
	return getMapIndex(cid, n.len)
}

func (n *Node) add(c *Client) {
	n.nodes[n.index(c.cid)].Store(c.cid, c)
}

func (n *Node) del(c *Client) {
	n.nodes[n.index(c.cid)].Delete(c.cid)
}

func (n *Node) get(cid int64) (*Client, bool) {
	value, ok := n.nodes[n.index(cid)].Load(cid)
	if !ok {
		return nil, false
	}

	if client, ok := value.(*Client); ok {
		return client, true
	} else {
		return nil, false
	}
}

func (n *Node) exist(cid int64) bool {
	if _, ok := n.nodes[n.index(cid)].Load(cid); ok {
		return true
	}

	return false
}

func (n *Node) each(fn func(c *Client)) {
	for _, node := range n.nodes {
		node.Range(func(key, value interface{}) bool {
			if client, ok := value.(*Client); ok {
				fn(client)
			}

			return true
		})
	}
}
