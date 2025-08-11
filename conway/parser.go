package conway

import (
	"errors"
	"fmt"
	"strings"
)

// Static errors for err113 compliance.
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
	parser := &Parser{
		operations: make(map[string]Operation),
	}

	parser.operations["d"] = DualOp{}
	parser.operations["a"] = AmboOp{}
	parser.operations["t"] = TruncateOp{}
	parser.operations["k"] = KisOp{}
	parser.operations["j"] = JoinOp{}
	parser.operations["o"] = OrthoOp{}
	parser.operations["e"] = ExpandOp{}
	parser.operations["g"] = GyroOp{}
	parser.operations["s"] = SnubOp{}

	return parser
}

func (p *Parser) Parse(notation string) (*Polyhedron, error) {
	notation = strings.TrimSpace(notation)
	if notation == "" {
		return nil, ErrEmptyNotation
	}

	seed, operations, err := p.parseNotation(notation)
	if err != nil {
		return nil, err
	}

	if seed == nil {
		return nil, ErrNoSeedPolyhedron
	}

	return p.applyOperations(seed, operations), nil
}

// parseNotation extracts seed and operations from notation string.
func (p *Parser) parseNotation(notation string) (*Polyhedron, []Operation, error) {
	var seed *Polyhedron

	var operations []Operation

	for i, char := range notation {
		symbol := string(char)

		if seed == nil {
			if parsedSeed := GetSeed(symbol); parsedSeed != nil {
				seed = parsedSeed
				continue
			}
		}

		if op, exists := p.operations[symbol]; exists {
			operations = append(operations, op)
			continue
		}

		if seed == nil && i == len(notation)-1 {
			if lastSeed := GetSeed(symbol); lastSeed != nil {
				seed = lastSeed
				continue
			}

			return nil, nil, fmt.Errorf("%w: %s", ErrUnknownSeedPolyhedron, symbol)
		}

		return nil, nil, fmt.Errorf("%w: %s at position %d", ErrUnknownOperation, symbol, i)
	}

	return seed, operations, nil
}

// applyOperations applies the operations to the seed polyhedron.
func (p *Parser) applyOperations(seed *Polyhedron, operations []Operation) *Polyhedron {
	result := seed.Clone()

	for i := len(operations) - 1; i >= 0; i-- {
		result = operations[i].Apply(result)
	}

	return result
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
