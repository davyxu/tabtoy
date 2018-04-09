package exprvm

import (
	"fmt"
	"strings"
)

type node struct {
	value interface{}
	prev  *node
}
type Stack struct {
	top    *node
	length int
}

func (self *Stack) String() string {

	var sb strings.Builder

	index := 0
	for i := self.top; i != nil; i = i.prev {
		sb.WriteString(fmt.Sprintf("[%d] %v\n", index, i.value))
		index++
	}

	return sb.String()
}

// Create a new stack
func NewStack() *Stack {
	return &Stack{nil, 0}
}

// Return the number of items in the stack
func (this *Stack) Len() int {
	return this.length
}

// View the top item on the stack
func (this *Stack) Peek() interface{} {
	if this.length == 0 {
		return nil
	}
	return this.top.value
}

// Pop the top item of the stack and return it
func (this *Stack) Pop() interface{} {
	if this.length == 0 {
		return nil
	}

	n := this.top
	this.top = n.prev
	this.length--
	return n.value
}

// Push a value onto the top of the stack
func (this *Stack) Push(value interface{}) {
	n := &node{value, this.top}
	this.top = n
	this.length++
}
