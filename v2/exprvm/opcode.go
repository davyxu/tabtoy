package exprvm

type Opcode int

const (
	Opcode_Nop Opcode = iota
	Opcode_Push
	Opcode_Add
	Opcode_Sub
	Opcode_Mul
	Opcode_Div
	Opcode_Minus
	Opcode_Exit
)

func (self Opcode) String() string {
	switch self {
	case Opcode_Nop:
		return "Nop"
	case Opcode_Push:
		return "Push"
	case Opcode_Add:
		return "Add"
	case Opcode_Sub:
		return "Sub"
	case Opcode_Mul:
		return "Mul"
	case Opcode_Div:
		return "Div"
	case Opcode_Minus:
		return "Minus"
	case Opcode_Exit:
		return "Exit"
	}

	return "Unknown"
}
