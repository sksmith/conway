package conway_test

import (
	"math"
	"testing"

	"github.com/sksmith/conway/conway"
)

func TestVector3Operations(t *testing.T) {
	t.Parallel()

	v1 := conway.Vector3{1, 2, 3}
	v2 := conway.Vector3{4, 5, 6}

	add := v1.Add(v2)
	if add.X != 5 || add.Y != 7 || add.Z != 9 {
		t.Errorf("Add failed: got %v, expected {5, 7, 9}", add)
	}

	sub := v2.Sub(v1)
	if sub.X != 3 || sub.Y != 3 || sub.Z != 3 {
		t.Errorf("Sub failed: got %v, expected {3, 3, 3}", sub)
	}

	scale := v1.Scale(2)
	if scale.X != 2 || scale.Y != 4 || scale.Z != 6 {
		t.Errorf("Scale failed: got %v, expected {2, 4, 6}", scale)
	}

	dot := v1.Dot(v2)
	if dot != 32 {
		t.Errorf("Dot failed: got %f, expected 32", dot)
	}

	cross := conway.Vector3{1, 0, 0}.Cross(conway.Vector3{0, 1, 0})
	if cross.X != 0 || cross.Y != 0 || cross.Z != 1 {
		t.Errorf("Cross failed: got %v, expected {0, 0, 1}", cross)
	}
}

func TestPolyhedronBasics(t *testing.T) {
	t.Parallel()

	p := conway.NewPolyhedron("test")

	v1 := p.AddVertex(conway.Vector3{0, 0, 0})
	v2 := p.AddVertex(conway.Vector3{1, 0, 0})
	v3 := p.AddVertex(conway.Vector3{0, 1, 0})

	p.AddEdge(v1, v2)
	p.AddEdge(v2, v3)
	p.AddEdge(v3, v1)

	f := p.AddFace([]*conway.Vertex{v1, v2, v3})

	if len(p.Vertices) != 3 {
		t.Errorf("Expected 3 vertices, got %d", len(p.Vertices))
	}

	if len(p.Edges) != 3 {
		t.Errorf("Expected 3 edges, got %d", len(p.Edges))
	}

	if len(p.Faces) != 1 {
		t.Errorf("Expected 1 face, got %d", len(p.Faces))
	}

	if v1.Degree() != 2 {
		t.Errorf("Expected vertex degree 2, got %d", v1.Degree())
	}

	if f.Degree() != 3 {
		t.Errorf("Expected face degree 3, got %d", f.Degree())
	}
}

func TestPolyhedronClone(t *testing.T) {
	t.Parallel()

	original := conway.Tetrahedron()
	clone := original.Clone()

	if len(clone.Vertices) != len(original.Vertices) {
		t.Errorf("Clone vertex count mismatch: %d vs %d",
			len(clone.Vertices), len(original.Vertices))
	}

	if len(clone.Edges) != len(original.Edges) {
		t.Errorf("Clone edge count mismatch: %d vs %d",
			len(clone.Edges), len(original.Edges))
	}

	if len(clone.Faces) != len(original.Faces) {
		t.Errorf("Clone face count mismatch: %d vs %d",
			len(clone.Faces), len(original.Faces))
	}

	if !clone.IsValid() {
		t.Error("Cloned polyhedron is not valid")
	}
}

func TestEulerCharacteristic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		poly     func() *conway.Polyhedron
		expected int
	}{
		{"conway.Tetrahedron", conway.Tetrahedron, 2},
		{"conway.Cube", conway.Cube, 2},
		{"conway.Octahedron", conway.Octahedron, 2},
		{"conway.Dodecahedron", conway.Dodecahedron, 2},
		{"conway.Icosahedron", conway.Icosahedron, 2},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			p := test.poly()

			chi := p.EulerCharacteristic()
			if chi != test.expected {
				t.Errorf("%s Euler characteristic: got %d, expected %d",
					test.name, chi, test.expected)
			}
		})
	}
}

func TestVertexRemoval(t *testing.T) {
	t.Parallel()

	p := conway.Tetrahedron()
	originalVertexCount := len(p.Vertices)
	originalEdgeCount := len(p.Edges)
	originalFaceCount := len(p.Faces)

	for _, v := range p.Vertices {
		p.RemoveVertex(v)
		break
	}

	if len(p.Vertices) != originalVertexCount-1 {
		t.Errorf("Expected %d vertices after removal, got %d",
			originalVertexCount-1, len(p.Vertices))
	}

	if len(p.Edges) >= originalEdgeCount {
		t.Error("Expected fewer edges after vertex removal")
	}

	if len(p.Faces) >= originalFaceCount {
		t.Error("Expected fewer faces after vertex removal")
	}
}

func TestNormalization(t *testing.T) {
	t.Parallel()

	p := conway.Cube()

	for _, v := range p.Vertices {
		v.Position = v.Position.Add(conway.Vector3{10, 20, 30})
	}

	p.Normalize()

	newCentroid := p.Centroid()
	if newCentroid.Length() > 1e-10 {
		t.Errorf("Centroid not at origin after normalization: %v", newCentroid)
	}

	maxDist := 0.0

	for _, v := range p.Vertices {
		dist := v.Position.Length()
		if dist > maxDist {
			maxDist = dist
		}
	}

	if math.Abs(maxDist-1.0) > 1e-10 {
		t.Errorf("Max distance from origin not 1.0 after normalization: %f", maxDist)
	}
}

func TestRemoveEdge(t *testing.T) {
	t.Parallel()

	p := conway.NewPolyhedron("test")
	v1 := p.AddVertex(conway.Vector3{0, 0, 0})
	v2 := p.AddVertex(conway.Vector3{1, 0, 0})
	v3 := p.AddVertex(conway.Vector3{0, 1, 0})

	p.AddEdge(v1, v2)
	e2 := p.AddEdge(v2, v3)
	p.AddEdge(v3, v1)

	if len(p.Edges) != 3 {
		t.Errorf("Expected 3 edges, got %d", len(p.Edges))
	}

	p.RemoveEdge(e2)

	if len(p.Edges) != 2 {
		t.Errorf("Expected 2 edges after removal, got %d", len(p.Edges))
	}

	found := false

	for _, e := range p.Edges {
		if e == e2 {
			found = true
			break
		}
	}

	if found {
		t.Error("Removed edge still found in polyhedron")
	}
}

func TestRemoveFace(t *testing.T) {
	t.Parallel()

	p := conway.NewPolyhedron("test")
	v1 := p.AddVertex(conway.Vector3{0, 0, 0})
	v2 := p.AddVertex(conway.Vector3{1, 0, 0})
	v3 := p.AddVertex(conway.Vector3{0, 1, 0})
	v4 := p.AddVertex(conway.Vector3{0, 0, 1})

	f1 := p.AddFace([]*conway.Vertex{v1, v2, v3})
	p.AddFace([]*conway.Vertex{v1, v3, v4})

	if len(p.Faces) != 2 {
		t.Errorf("Expected 2 faces, got %d", len(p.Faces))
	}

	p.RemoveFace(f1)

	if len(p.Faces) != 1 {
		t.Errorf("Expected 1 face after removal, got %d", len(p.Faces))
	}

	found := false

	for _, f := range p.Faces {
		if f == f1 {
			found = true
			break
		}
	}

	if found {
		t.Error("Removed face still found in polyhedron")
	}
}

func TestNormalizeZeroLengthVector(t *testing.T) {
	t.Parallel()

	v := conway.Vector3{0, 0, 0}
	normalized := v.Normalize()

	if normalized.X != 0 || normalized.Y != 0 || normalized.Z != 0 {
		t.Errorf("Zero vector normalization should return zero vector, got %v", normalized)
	}
}

func TestEdgeOtherVertex(t *testing.T) {
	t.Parallel()

	p := conway.NewPolyhedron("test")
	v1 := p.AddVertex(conway.Vector3{0, 0, 0})
	v2 := p.AddVertex(conway.Vector3{1, 0, 0})
	v3 := p.AddVertex(conway.Vector3{0, 1, 0})

	edge := p.AddEdge(v1, v2)

	other := edge.OtherVertex(v1)
	if other != v2 {
		t.Error("OtherVertex should return v2 when given v1")
	}

	other = edge.OtherVertex(v2)
	if other != v1 {
		t.Error("OtherVertex should return v1 when given v2")
	}

	other = edge.OtherVertex(v3)
	if other != nil {
		t.Error("OtherVertex should return nil for vertex not in edge")
	}
}

func TestFaceNormalEdgeCases(t *testing.T) {
	t.Parallel()

	p := conway.NewPolyhedron("test")

	v1 := p.AddVertex(conway.Vector3{0, 0, 0})
	v2 := p.AddVertex(conway.Vector3{1, 0, 0})

	face := p.AddFace([]*conway.Vertex{v1, v2})
	normal := face.Normal()

	if normal.Length() != 0 {
		t.Error("Face with fewer than 3 vertices should have zero normal")
	}

	v3 := p.AddVertex(conway.Vector3{0, 0, 0})
	face2 := p.AddFace([]*conway.Vertex{v1, v2, v3})
	normal2 := face2.Normal()

	if normal2.Length() != 0 {
		t.Error("Degenerate face should have zero normal")
	}
}
