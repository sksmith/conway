package conway

import (
	"errors"
	"fmt"
	"strings"
)

// Static errors for err113 compliance
var (
	ErrEmptyNotation         = errors.New("empty notation string")
	ErrNoSeedPolyhedron      = errors.New("no seed polyhedron found in notation")
	ErrUnknownSeedPolyhedron = errors.New("unknown seed polyhedron")
	ErrUnknownOperation      = errors.New("unknown operation")
)

type Parser struct {
	operations map[string]Operation
}

func NewParser() *Parser {
	p := &Parser{
		operations: make(map[string]Operation),
	}

	p.operations["d"] = DualOp{}
	p.operations["a"] = AmboOp{}
	p.operations["t"] = TruncateOp{}
	p.operations["k"] = KisOp{}
	p.operations["j"] = JoinOp{}
	p.operations["o"] = OrthoOp{}
	p.operations["e"] = ExpandOp{}
	p.operations["g"] = GyroOp{}
	p.operations["s"] = SnubOp{}

	return p
}

func (p *Parser) Parse(notation string) (*Polyhedron, error) {
	notation = strings.TrimSpace(notation)
	if notation == "" {
		return nil, ErrEmptyNotation
	}

	var seed *Polyhedron
	var operations []Operation

	for i, char := range notation {
		symbol := string(char)

		if seed == nil {
			seed = GetSeed(symbol)
			if seed != nil {
				continue
			}
		}

		op, exists := p.operations[symbol]
		if exists {
			operations = append(operations, op)
			continue
		}

		if seed == nil && i == len(notation)-1 {
			seed = GetSeed(symbol)
			if seed == nil {
				return nil, fmt.Errorf("%w: %s", ErrUnknownSeedPolyhedron, symbol)
			}

			continue
		}

		return nil, fmt.Errorf("%w: %s at position %d", ErrUnknownOperation, symbol, i)
	}

	if seed == nil {
		return nil, ErrNoSeedPolyhedron
	}

	result := seed.Clone()

	for i := len(operations) - 1; i >= 0; i-- {
		result = operations[i].Apply(result)
	}

	return result, nil
}

func (p *Parser) Validate(notation string) error {
	_, err := p.Parse(notation)
	return err
}

func (p *Parser) GetAvailableOperations() map[string]string {
	ops := make(map[string]string)
	for symbol, op := range p.operations {
		ops[symbol] = op.Name()
	}

	return ops
}

func (p *Parser) GetAvailableSeeds() map[string]string {
	return map[string]string{
		"T": "Tetrahedron",
		"C": "Cube",
		"O": "Octahedron",
		"D": "Dodecahedron",
		"I": "Icosahedron",
	}
}

func Parse(notation string) (*Polyhedron, error) {
	parser := NewParser()

	return parser.Parse(notation)
}

func MustParse(notation string) *Polyhedron {
	result, err := Parse(notation)
	if err != nil {
		panic(err)
	}

	return result
}
