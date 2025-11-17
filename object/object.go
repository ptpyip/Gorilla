package object

import "fmt"

type ObjectType string

const (
	NONE = "NONE"
	BOOL = "BOOL"
	INT  = "INT"
)

type Object interface {
	GetType() ObjectType
	Inspect() string
}

type None struct {
}

func (noneObj *None) GetType() ObjectType {
	return NONE
}

func (noneObj *None) Inspect() string {
	return "None"
}

type Bool struct {
	Value bool
}

func (boolObj *Bool) GetType() ObjectType {
	return BOOL
}

func (boolObj *Bool) Inspect() string {
	return fmt.Sprintf("%t", boolObj.Value)
}

type Int struct {
	Value int64
}

func (intObj *Int) GetType() ObjectType {
	return INT
}

func (intObj *Int) Inspect() string {
	return fmt.Sprintf("%d", intObj.Value)
}
