package conway

type Operation interface {
	Apply(p *Polyhedron) *Polyhedron
	Symbol() string
	Name() string
}
