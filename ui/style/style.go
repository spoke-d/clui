package style

import (
	"github.com/pkg/errors"
	"github.com/spoke-d/clui/ui/color"
)

type Style struct {
	color     color.RGB
	underline bool
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
	//fmt.Printf("%[1]T %[1]v\n", e)

	switch node := e.(type) {
	case *QueryExpression:
		for _, exp := range node.Expressions {
			var err error
			box, err = run(box, exp)
			if err != nil {
				return nil, errors.WithStack(err)
			}
		}
		return box, nil

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
		name, ok := node.Name.(*Identifier)
		if !ok {
			return nil, errors.Errorf("expected a identifier name %T", node.Name)
		}

		switch name.Token.Literal {
		case "rgb":
			color, err := parseRGBColor(node.Arguments)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			return box.Value(color)
		case "hsv":
			color, err := parseHSVColor(node.Arguments)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			return box.Value(color)
		case "hsl":
			color, err := parseHSLColor(node.Arguments)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			return box.Value(color)
		}
		return nil, errors.Errorf("unknown constructor %q", name)

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
	case "color", "underline":
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
		switch v := value.(type) {
		case color.RGB:
			b.root.color = v
		default:
			return nil, errors.Errorf("unable to set the color %q", value)
		}
	case "underline":
		switch v := value.(type) {
		case bool:
			b.root.underline = v
		default:
			return nil, errors.Errorf("unable to set the underline %q", value)
		}
	default:
		return nil, errors.Errorf("unknown property %q", b.leaf)
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

func parseRGBColor(args []Expression) (color.RGB, error) {
	if len(args) != 3 {
		return color.RGB{}, errors.Errorf("expected 3 arguments to construct a rgb color")
	}

	var rgb [3]float64
	for i, v := range args {
		switch arg := v.(type) {
		case *Integer:
			rgb[i] = float64(arg.Value)
		case *Float:
			rgb[i] = arg.Value
		}
	}

	return color.RGB{R: rgb[0], G: rgb[1], B: rgb[2]}, nil
}

func parseHSVColor(args []Expression) (color.RGB, error) {
	if len(args) != 3 {
		return color.RGB{}, errors.Errorf("expected 3 arguments to construct a hsv color")
	}

	var hsv [3]float64
	for i, v := range args {
		switch arg := v.(type) {
		case *Integer:
			hsv[i] = float64(arg.Value)
		case *Float:
			hsv[i] = arg.Value
		}
	}

	return color.HSV{H: hsv[0], S: hsv[1], V: hsv[2]}.RGB(), nil
}

func parseHSLColor(args []Expression) (color.RGB, error) {
	if len(args) != 3 {
		return color.RGB{}, errors.Errorf("expected 3 arguments to construct a hsl color")
	}

	var hsl [3]float64
	for i, v := range args {
		switch arg := v.(type) {
		case *Integer:
			hsl[i] = float64(arg.Value)
		case *Float:
			hsl[i] = arg.Value
		}
	}

	return color.HSL{H: hsl[0], S: hsl[1], L: hsl[2]}.RGB(), nil
}
