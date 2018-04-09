package exprvm

import (
	"fmt"
	"strings"
)

type Command struct {
	Type Opcode

	Operand []interface{}
}

func (self *Command) String() string {

	var sb strings.Builder

	sb.WriteString(self.Type.String())

	if len(self.Operand) > 0 {
		sb.WriteString(" ")
		for index, operand := range self.Operand {
			if index > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("%v", operand))
		}
	}

	return sb.String()
}

type Chunk struct {
	Commands []Command
}

func (self *Chunk) String() string {

	var sb strings.Builder

	for _, cmd := range self.Commands {

		sb.WriteString(fmt.Sprintf("%s\n", cmd.String()))
	}

	return sb.String()
}

func (self *Chunk) AddCode(t Opcode) {

	self.Commands = append(self.Commands, Command{Type: t})
}

func (self *Chunk) AddCodeOperand(t Opcode, operand ...interface{}) {

	self.Commands = append(self.Commands, Command{
		Type:    t,
		Operand: operand,
	})
}

func newChunk() *Chunk {
	return &Chunk{}
}
