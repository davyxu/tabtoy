package exprvm

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

func Compile(src string) (*Chunk, error) {

	node, err := parser.ParseExpr(src)
	if err != nil {
		return nil, err
	}

	ast.Print(nil, node)

	ck := newChunk()

	if err := walkTree(ck, node); err != nil {
		return nil, err
	}

	return ck, nil
}

func walkTree(ck *Chunk, node ast.Node) (err error) {

	switch thisNode := node.(type) {
	case *ast.BinaryExpr:
		if err = walkTree(ck, thisNode.X); err != nil {
			return
		}

		if err = walkTree(ck, thisNode.Y); err != nil {
			return
		}

		switch thisNode.Op {
		case token.ADD:
			ck.AddCode(Opcode_Add)
		case token.SUB:
			ck.AddCode(Opcode_Sub)
		case token.MUL:
			ck.AddCode(Opcode_Mul)
		case token.QUO:
			ck.AddCode(Opcode_Div)
		}
	case *ast.BasicLit: // 字面量/常数

		switch thisNode.Kind {
		case token.INT:

			v, err := strconv.Atoi(thisNode.Value)
			if err != nil {
				return err
			}

			ck.AddCodeOperand(Opcode_Push, v)
		default:
			return ErrUnknownOperandType
		}
	case *ast.ParenExpr:
		return walkTree(ck, thisNode.X)
	case *ast.UnaryExpr:
		walkTree(ck, thisNode.X)
		ck.AddCode(Opcode_Minus)
	case *ast.Ident: // 变量/常量
	default:
		return ErrUnknownExpression
	}

	return nil
}
