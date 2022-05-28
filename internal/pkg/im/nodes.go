package im

import (
	"sync"
)

type Node struct {
	len   int         // 节点数量
	nodes []*sync.Map // 节点列表
}

func NewNode(len int) *Node {
	return &Node{len: len, nodes: maps(len)}
}

func (n *Node) index(cid int64) int {
	return int(cid) % n.len
}

// 添加客户端
func (n *Node) add(c *Client) {
	n.nodes[n.index(c.cid)].Store(c.cid, c)
}

// 删除客户端
func (n *Node) del(c *Client) {
	n.nodes[n.index(c.cid)].Delete(c.cid)
}

// 获取客户端
func (n *Node) get(cid int64) (*Client, bool) {
	value, ok := n.nodes[n.index(cid)].Load(cid)
	if !ok {
		return nil, false
	}

	if client, ok := value.(*Client); ok {
		return client, true
	}

	return nil, false
}

// 判断客户端是否存在
func (n *Node) exist(cid int64) bool {
	if _, ok := n.nodes[n.index(cid)].Load(cid); ok {
		return true
	}

	return false
}

// each 遍历所有客户端
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
