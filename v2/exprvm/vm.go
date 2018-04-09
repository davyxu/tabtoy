package exprvm

type Machine struct {
	DataStack *Stack
}

func (self *Machine) Run(chunk *Chunk) {

	for pc := 0; pc < len(chunk.Commands); {

		cmd := chunk.Commands[pc]

		if cmd.Type == Opcode_Exit {
			break
		}

		self.execute(cmd)

		pc++
	}

}

func (self *Machine) execute(cmd Command) (err error) {

	switch cmd.Type {
	case Opcode_Push:
		self.DataStack.Push(cmd.Operand[0])
	case Opcode_Add, Opcode_Sub, Opcode_Div, Opcode_Mul:
		v1 := self.DataStack.Pop()
		v2 := self.DataStack.Pop()

		result := arithOp(cmd.Type, v2, v1)

		self.DataStack.Push(result)
	case Opcode_Minus:
		v := self.DataStack.Pop()

		self.DataStack.Push(-v.(int))

	default:
		panic("Unknown opcode")
	}

	return nil
}

func arithOp(op Opcode, a, b interface{}) interface{} {

	an := a.(int)
	bn := b.(int)

	switch op {
	case Opcode_Add:
		return an + bn
	case Opcode_Sub:
		return an - bn
	case Opcode_Mul:
		return an * bn
	case Opcode_Div:
		return an / bn
	}

	return nil
}

func NewMachine() *Machine {
	return &Machine{
		DataStack: NewStack(),
	}
}
