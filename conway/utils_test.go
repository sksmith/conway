package conway

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateGeometryStats(t *testing.T) {
	t.Run("EmptyPolyhedron", func(t *testing.T) {
		p := &Polyhedron{}
		stats := p.CalculateGeometryStats()

		assert.NotNil(t, stats)
		assert.Equal(t, 0.0, stats.MinEdgeLength)
		assert.Equal(t, 0.0, stats.MaxEdgeLength)
		assert.Equal(t, 0.0, stats.AvgEdgeLength)
		assert.Equal(t, 0.0, stats.MinFaceArea)
		assert.Equal(t, 0.0, stats.MaxFaceArea)
		assert.Equal(t, 0.0, stats.AvgFaceArea)
	})

	t.Run("NoEdges", func(t *testing.T) {
		p := &Polyhedron{
			Vertices: map[int]*Vertex{0: {ID: 0, Position: Vector3{0, 0, 0}}},
			Faces:    map[int]*Face{0: {}},
		}
		stats := p.CalculateGeometryStats()

		assert.NotNil(t, stats)
		assert.Equal(t, 0.0, stats.MinEdgeLength)
		assert.Equal(t, 0.0, stats.MaxEdgeLength)
		assert.Equal(t, 0.0, stats.AvgEdgeLength)
	})

	t.Run("NoFaces", func(t *testing.T) {
		v1 := &Vertex{ID: 0, Position: Vector3{0, 0, 0}}
		v2 := &Vertex{ID: 1, Position: Vector3{1, 0, 0}}
		e := &Edge{V1: v1, V2: v2}

		p := &Polyhedron{
			Vertices: map[int]*Vertex{0: v1, 1: v2},
			Edges:    map[int]*Edge{0: e},
		}
		stats := p.CalculateGeometryStats()

		// Function returns early if no faces, so all values are zero
		assert.NotNil(t, stats)
		assert.Equal(t, 0.0, stats.MinEdgeLength)
		assert.Equal(t, 0.0, stats.MaxEdgeLength)
		assert.Equal(t, 0.0, stats.AvgEdgeLength)
		assert.Equal(t, 0.0, stats.MinFaceArea)
		assert.Equal(t, 0.0, stats.MaxFaceArea)
		assert.Equal(t, 0.0, stats.AvgFaceArea)
	})

	t.Run("SingleVertex", func(t *testing.T) {
		v := &Vertex{ID: 0, Position: Vector3{1, 2, 3}}
		p := &Polyhedron{
			Vertices: map[int]*Vertex{0: v},
			Edges:    map[int]*Edge{}, // Empty but not nil
			Faces:    map[int]*Face{}, // Empty but not nil
		}
		stats := p.CalculateGeometryStats()

		// Function returns early if no edges or faces, so bounding box is not set
		assert.NotNil(t, stats)
		assert.Equal(t, Vector3{0, 0, 0}, stats.BoundingBox.Min)
		assert.Equal(t, Vector3{0, 0, 0}, stats.BoundingBox.Max)
	})

	t.Run("MultipleVertices", func(t *testing.T) {
		v1 := &Vertex{ID: 0, Position: Vector3{-1, -2, -3}}
		v2 := &Vertex{ID: 1, Position: Vector3{4, 5, 6}}
		v3 := &Vertex{ID: 2, Position: Vector3{0, 1, 2}}

		p := &Polyhedron{
			Vertices: map[int]*Vertex{0: v1, 1: v2, 2: v3},
			Edges:    map[int]*Edge{}, // Empty but not nil
			Faces:    map[int]*Face{}, // Empty but not nil
		}
		stats := p.CalculateGeometryStats()

		// Function returns early if no edges or faces, so bounding box is not set
		assert.NotNil(t, stats)
		assert.Equal(t, Vector3{0, 0, 0}, stats.BoundingBox.Min)
		assert.Equal(t, Vector3{0, 0, 0}, stats.BoundingBox.Max)
	})

	t.Run("ValidCube", func(t *testing.T) {
		// Create a simple cube for testing
		cube, err := Parse("C")
		assert.NoError(t, err)

		stats := cube.CalculateGeometryStats()
		assert.NotNil(t, stats)

		assert.Greater(t, stats.MinEdgeLength, 0.0)
		assert.Greater(t, stats.MaxEdgeLength, 0.0)
		assert.GreaterOrEqual(t, stats.MaxEdgeLength, stats.MinEdgeLength)
		assert.Greater(t, stats.AvgEdgeLength, 0.0)

		assert.Greater(t, stats.MinFaceArea, 0.0)
		assert.Greater(t, stats.MaxFaceArea, 0.0)
		assert.GreaterOrEqual(t, stats.MaxFaceArea, stats.MinFaceArea)
		assert.Greater(t, stats.AvgFaceArea, 0.0)

		// Cube should have uniform edge lengths and face areas
		assert.InDelta(t, stats.MinEdgeLength, stats.MaxEdgeLength, 1e-10)
		assert.InDelta(t, stats.MinFaceArea, stats.MaxFaceArea, 1e-10)
	})
}

func TestEdgeLookup(t *testing.T) {
	t.Run("BasicOperations", func(t *testing.T) {
		el := NewEdgeLookup()

		v1 := &Vertex{ID: 1}
		v2 := &Vertex{ID: 2}
		edge := &Edge{V1: v1, V2: v2}

		// Test Add
		el.Add(edge)

		// Test Find
		found := el.Find(1, 2)
		assert.Equal(t, edge, found)

		// Test symmetric lookup
		found = el.Find(2, 1)
		assert.Equal(t, edge, found)

		// Test Remove
		el.Remove(edge)
		found = el.Find(1, 2)
		assert.Nil(t, found)
	})

	t.Run("NonexistentEdge", func(t *testing.T) {
		el := NewEdgeLookup()
		found := el.Find(99, 100)
		assert.Nil(t, found)
	})
}

func TestCalculateFaceNormal(t *testing.T) {
	t.Run("InsufficientVertices", func(t *testing.T) {
		vertices := []*Vertex{
			{Position: Vector3{0, 0, 0}},
			{Position: Vector3{1, 0, 0}},
		}

		_, err := calculateFaceNormal(vertices)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insufficient vertices")
	})

	t.Run("ValidTriangle", func(t *testing.T) {
		vertices := []*Vertex{
			{Position: Vector3{0, 0, 0}},
			{Position: Vector3{1, 0, 0}},
			{Position: Vector3{0, 1, 0}},
		}

		normal, err := calculateFaceNormal(vertices)
		assert.NoError(t, err)
		assert.InDelta(t, 0.0, normal.X, 1e-10)
		assert.InDelta(t, 0.0, normal.Y, 1e-10)
		assert.InDelta(t, 1.0, normal.Z, 1e-10)
	})

	t.Run("DegenerateFace", func(t *testing.T) {
		// All vertices in a line
		vertices := []*Vertex{
			{Position: Vector3{0, 0, 0}},
			{Position: Vector3{1, 0, 0}},
			{Position: Vector3{2, 0, 0}},
		}

		_, err := calculateFaceNormal(vertices)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "degenerate face normal")
	})
}

func TestEnsureCounterClockwise(t *testing.T) {
	t.Run("InsufficientVertices", func(t *testing.T) {
		vertices := []*Vertex{
			{Position: Vector3{0, 0, 0}},
			{Position: Vector3{1, 0, 0}},
		}
		center := Vector3{0, 0, 0}

		result := ensureCounterClockwise(vertices, center)
		assert.Equal(t, vertices, result)
	})

	t.Run("ValidTriangle", func(t *testing.T) {
		vertices := []*Vertex{
			{Position: Vector3{1, 0, 0}},
			{Position: Vector3{0, 1, 0}},
			{Position: Vector3{0, 0, 1}},
		}
		center := Vector3{0, 0, 0}

		result := ensureCounterClockwise(vertices, center)
		assert.NotNil(t, result)
		assert.Len(t, result, 3)
	})
}

func TestMakeEdgeKey(t *testing.T) {
	t.Run("ConsistentOrdering", func(t *testing.T) {
		key1 := makeEdgeKey(1, 2)
		key2 := makeEdgeKey(2, 1)
		assert.Equal(t, key1, key2)
	})

	t.Run("DifferentVertices", func(t *testing.T) {
		key1 := makeEdgeKey(1, 2)
		key2 := makeEdgeKey(1, 3)
		assert.NotEqual(t, key1, key2)
	})
}

func TestAllocateSliceFunctions(t *testing.T) {
	t.Run("AllocateVertexSlice", func(t *testing.T) {
		vertices := allocateVertexSlice(10)
		assert.NotNil(t, vertices)
		assert.Equal(t, 0, len(vertices))
		assert.Equal(t, 10, cap(vertices))
	})

	t.Run("AllocateEdgeSlice", func(t *testing.T) {
		edges := allocateEdgeSlice(20)
		assert.NotNil(t, edges)
		assert.Equal(t, 0, len(edges))
		assert.Equal(t, 20, cap(edges))
	})
}
