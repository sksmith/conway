package conway

type OrthoOp struct{}

func (o OrthoOp) Symbol() string {
	return "o"
}

func (o OrthoOp) Name() string {
	return "ortho"
}

func (o OrthoOp) Apply(p *Polyhedron) *Polyhedron {
	return Join(Join(p))
}

type ExpandOp struct{}

func (e ExpandOp) Symbol() string {
	return "e"
}

func (e ExpandOp) Name() string {
	return "expand"
}

func (e ExpandOp) Apply(p *Polyhedron) *Polyhedron {
	return Ambo(Ambo(p))
}

type GyroOp struct{}

func (g GyroOp) Symbol() string {
	return "g"
}

func (g GyroOp) Name() string {
	return "gyro"
}

func (g GyroOp) Apply(p *Polyhedron) *Polyhedron {
	return Dual(Ambo(p))
}

type SnubOp struct{}

func (s SnubOp) Symbol() string {
	return "s"
}

func (s SnubOp) Name() string {
	return "snub"
}

func (s SnubOp) Apply(p *Polyhedron) *Polyhedron {
	return Dual(Gyro(p))
}

func Ortho(p *Polyhedron) *Polyhedron {
	op := OrthoOp{}
	return op.Apply(p)
}

func Expand(p *Polyhedron) *Polyhedron {
	op := ExpandOp{}
	return op.Apply(p)
}

func Gyro(p *Polyhedron) *Polyhedron {
	op := GyroOp{}
	return op.Apply(p)
}

func Snub(p *Polyhedron) *Polyhedron {
	op := SnubOp{}
	return op.Apply(p)
}
