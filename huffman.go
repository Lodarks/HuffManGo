package main

import (
	"container/heap"
)

// Düğüm yapısı
type Node struct {
	char       byte
	freq       int
	left       *Node
	right      *Node
	isInternal bool
	order      int
}

// Öncelik kuyruğu için gerekli metotlar
type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	if pq[i].freq == pq[j].freq {
		return pq[i].order < pq[j].order
	}
	return pq[i].freq < pq[j].freq
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Node)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// DynamicHuffman yapısı
type DynamicHuffman struct {
	root       *Node
	nodeMap    map[byte]*Node
	pq         PriorityQueue
	orderCount int
}

// Yeni bir DynamicHuffman nesnesi oluştur
func NewDynamicHuffman() *DynamicHuffman {
	return &DynamicHuffman{
		nodeMap:    make(map[byte]*Node),
		pq:         make(PriorityQueue, 0),
		orderCount: 0,
	}
}

// Karakter güncelleme
func (dh *DynamicHuffman) UpdateTree(char byte) {
	if node, exists := dh.nodeMap[char]; exists {
		// Var olan karakteri güncelle
		node.freq++
		heap.Init(&dh.pq)
	} else {
		// Yeni karakter ekle
		newNode := &Node{
			char:       char,
			freq:       1,
			isInternal: false,
			order:      dh.orderCount,
		}
		dh.orderCount++
		dh.nodeMap[char] = newNode
		heap.Push(&dh.pq, newNode)
	}

	// Ağacı yeniden düzenle
	dh.rebuildTree()
}

// Ağacı yeniden oluştur
func (dh *DynamicHuffman) rebuildTree() {
	// Öncelik kuyruğunu kopyala
	tempPQ := make(PriorityQueue, len(dh.pq))
	copy(tempPQ, dh.pq)
	heap.Init(&tempPQ)

	// Yeni ağacı oluştur
	for tempPQ.Len() > 1 {
		left := heap.Pop(&tempPQ).(*Node)
		right := heap.Pop(&tempPQ).(*Node)

		internal := &Node{
			freq:       left.freq + right.freq,
			left:       left,
			right:      right,
			isInternal: true,
			order:      dh.orderCount,
		}
		dh.orderCount++

		heap.Push(&tempPQ, internal)
	}

	if tempPQ.Len() > 0 {
		dh.root = heap.Pop(&tempPQ).(*Node)
	}
}

// Karakterin kodunu bul
func (dh *DynamicHuffman) GetCode(char byte) string {
	if node, exists := dh.nodeMap[char]; exists {
		code := ""
		current := node
		path := make([]byte, 0)

		// Kökten yaprağa kadar yolu bul
		for current != dh.root {
			parent := dh.findParent(current)
			if parent == nil {
				break
			}
			if parent.left == current {
				path = append(path, '0')
			} else {
				path = append(path, '1')
			}
			current = parent
		}

		// Yolu tersine çevir
		for i := len(path) - 1; i >= 0; i-- {
			code += string(path[i])
		}
		return code
	}
	return ""
}

// Düğümün ebeveynini bul
func (dh *DynamicHuffman) findParent(node *Node) *Node {
	var search func(*Node) *Node
	search = func(current *Node) *Node {
		if current == nil {
			return nil
		}
		if current.left == node || current.right == node {
			return current
		}
		if left := search(current.left); left != nil {
			return left
		}
		return search(current.right)
	}
	return search(dh.root)
}
