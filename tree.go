package main

import (
	"fmt"
	// "math/rand"
)

// A Tree is a binary tree with integer values.
type Tree struct {
	Left  *Tree
	Value int
	Right *Tree
}

// New returns a new, random binary tree holding the values k, 2k, ..., 10k.
func New(k int) *Tree {
	var t *Tree
	for _, v := range []int{6, 4, 5, 2, 9, 8, 7, 3, 1} {
		t = insert(t, v)
	}
	return t
}

func insert(t *Tree, v int) *Tree {
	if t == nil {
		return &Tree{nil, v, nil}
	}
	if v < t.Value {
		t.Left = insert(t.Left, v)
	} else {
		t.Right = insert(t.Right, v)
	}
	return t
}

func (t *Tree) String() string {
	if t == nil {
		return "()"
	}
	s := ""
	if t.Left != nil {
		s += t.Left.String() + " "
	}
	s += fmt.Sprint(t.Value)
	if t.Right != nil {
		s += " " + t.Right.String()
	}
	return "(" + s + ")"
}

type Queue struct {
	nodes []*Tree
	head  int
	tail  int
	count int
}

// Push adds a node to the queue.
func (q *Queue) Push(n *Tree) {
	if q.head == q.tail && q.count > 0 {
		nodes := make([]*Tree, len(q.nodes)*2)
		copy(nodes, q.nodes[q.head:])
		copy(nodes[len(q.nodes)-q.head:], q.nodes[:q.head])
		q.head = 0
		q.tail = len(q.nodes)
		q.nodes = nodes
	}
	q.nodes[q.tail] = n
	q.tail = (q.tail + 1) % len(q.nodes)
	q.count++
}

// Pop removes and returns a node from the queue in first to last order.
func (q *Queue) Pop() *Tree {
	if q.count == 0 {
		return nil
	}
	node := q.nodes[q.head]
	q.head = (q.head + 1) % len(q.nodes)
	q.count--
	return node
}

// Stack is a basic LIFO stack that resizes as needed.
type Stack struct {
	nodes []*Tree
	count int
}

// Push adds a node to the stack.
func (s *Stack) Push(n *Tree) {
	if s.count >= len(s.nodes) {
		nodes := make([]*Tree, len(s.nodes)*2)
		copy(nodes, s.nodes)
		s.nodes = nodes
	}
	s.nodes[s.count] = n
	s.count++
}

// Pop removes and returns a node from the stack in last to first order.
func (s *Stack) Pop() *Tree {
	if s.count == 0 {
		return nil
	}
	node := s.nodes[s.count-1]
	s.count--
	return node
}

func (s *Stack) Top() *Tree {
	if s.count == 0 {
		return nil
	}
	return s.nodes[s.count-1]
}

// head -> left -> right
func (t *Tree) preOrder() {
	s := &Stack{nodes: make([]*Tree, 3)}
	for t != nil || s.count > 0 {
		if t != nil {
			fmt.Println(t.Value)
			s.Push(t)
			t = t.Left
		} else {
			t = s.Pop().Right
		}
	}
}

func (t *Tree) inOrder() {
	s := &Stack{nodes: make([]*Tree, 3)}

	for t != nil || s.count > 0 {
		for t != nil {
			s.Push(t)
			t = t.Left
		}
		if s.count > 0 {
			tmp := s.Pop()
			fmt.Println(tmp.Value)
			t = tmp.Right			
		}
		// fmt.Println(s.Pop)
	}
}

func (t *Tree) postOrder() {
	if t.Left != nil {
		t.Left.postOrder()
	}
	if t.Right != nil {
		t.Right.postOrder()
	}
	fmt.Println(t.Value)
}

func (t *Tree) postOrder2() {

	s := &Stack{nodes: make([]*Tree, 3)}
	var pre *Tree = nil
	s.Push(t)

	for s.count > 0 {
		cur := s.Top()

		if (cur.Left == nil && cur.Right == nil)||
			(pre != nil && (cur.Left == pre || cur.Right == pre)) {
			fmt.Println(cur.Value)
			s.Pop()
			pre = cur
		} else {
			if(cur.Right != nil) {
				s.Push(cur.Right)
			}
			if(cur.Left != nil) {
				s.Push(cur.Left)
			}
		}
	}

}

func main() {
	t := New(1)
	fmt.Println(t)
	t.postOrder2()
	// t.inOrder()
	// t.preOrder()
}
