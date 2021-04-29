package style

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spoke-d/clui/ui/color"
)

type Style struct {
	color color.RGB
}

// ParseStyle allows you to pass a string and attempt parse it into a style.
// Example:
//     color=rgb(0,0,0);
func ParseStyle(src string) (*Style, error) {
	lex := NewLexer(src)
	parser := NewParser(lex)
	ast, err := parser.Run()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	box, err := run(NewBox(&Style{}), ast)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return box.Root(), nil
}

func run(box *Box, e Expression) (*Box, error) {
	// Useful for debugging.
	fmt.Printf("%[1]T %[1]v\n", e)

	switch node := e.(type) {
	case *QueryExpression:
		for _, exp := range node.Expressions {
			var err error
			box, err = run(box, exp)
			if err != nil {
				return nil, errors.WithStack(err)
			}
		}

	case *ExpressionStatement:
		return run(box, node.Expression)

	case *InfixExpression:
		left, err := run(box, node.Left)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		switch node.Token.Type {
		case ASSIGN:
		default:
			return nil, errors.Errorf("unknown operator %q", node.Token.Literal)
		}

		return run(left, node.Right)

	case *Identifier:
		return box.Push(node.Token.Literal)

	case *CallExpression:
		fmt.Println("HERE")

	case *Integer:
		return box.Value(node.Value)

	case *Float:
		return box.Value(node.Value)

	case *String:
		return box.Value(node.Token.Literal)

	case *Bool:
		return box.Value(node.Value)

	case *Empty:
		return box.Pop()
	}
	return nil, RuntimeErrorf("Syntax Error: Unexpected expression %T", e)
}

type Box struct {
	root      *Style
	leaf      string
	traversed bool
}

func NewBox(root *Style) *Box {
	return &Box{root: root}
}

func (b *Box) Root() *Style {
	return b.root
}

func (b *Box) Push(name string) (*Box, error) {
	switch name {
	case "color":
		b.leaf = name
		b.traversed = true
	default:
		return nil, errors.Errorf("unknown field name %q", name)
	}
	return b, nil
}

func (b *Box) Value(value interface{}) (*Box, error) {
	if !b.traversed {
		return nil, errors.Errorf("unable to set value %q", value)
	}

	switch b.leaf {
	case "color":
		fmt.Println("!!", value)
	}

	return b.Pop()
}

func (b *Box) Pop() (*Box, error) {
	b.leaf = ""
	b.traversed = false

	return b, nil
}

// RuntimeError creates an invalid error.
type RuntimeError struct {
	err error
}

func (e *RuntimeError) Error() string {
	return e.err.Error()
}

// RuntimeErrorf defines a sentinel error for invalid index.
func RuntimeErrorf(msg string, args ...interface{}) error {
	return &RuntimeError{
		err: errors.Errorf("Runtime Error: "+msg, args...),
	}
}

// IsRuntimeError returns if the error is an ErrInvalidIndex error
func IsRuntimeError(err error) bool {
	err = errors.Cause(err)
	_, ok := err.(*RuntimeError)
	return ok
}
