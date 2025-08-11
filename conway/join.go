package conway

type JoinOp struct{}

func (j JoinOp) Symbol() string {
	return "j"
}

func (j JoinOp) Name() string {
	return "join"
}

func (j JoinOp) Apply(p *Polyhedron) *Polyhedron {
	dual := Dual(p)

	ambo := Ambo(dual)

	return ambo
}

func Join(p *Polyhedron) *Polyhedron {
	op := JoinOp{}
	return op.Apply(p)
}
