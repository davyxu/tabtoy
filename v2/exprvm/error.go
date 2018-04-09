package exprvm

import "github.com/pkg/errors"

var (
	ErrUnknownExpression  = errors.New("Unknown Expression")
	ErrUnknownOperandType = errors.New("Unknown Operand Type")
)
