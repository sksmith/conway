package conway_test

import (
	"testing"

	"github.com/sksmith/conway/conway"
	"github.com/stretchr/testify/assert"
)

func TestCalculateGeometryStats(t *testing.T) {
	t.Parallel()

	t.Run("EmptyPolyhedron", func(t *testing.T) {
		t.Parallel()

		p := &conway.Polyhedron{}
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
		t.Parallel()

		p := &conway.Polyhedron{
			Vertices: map[int]*conway.Vertex{0: {ID: 0, Position: conway.Vector3{0, 0, 0}}},
			Faces:    map[int]*conway.Face{0: {}},
		}
		stats := p.CalculateGeometryStats()

		assert.NotNil(t, stats)
		assert.Equal(t, 0.0, stats.MinEdgeLength)
		assert.Equal(t, 0.0, stats.MaxEdgeLength)
		assert.Equal(t, 0.0, stats.AvgEdgeLength)
	})

	t.Run("NoFaces", func(t *testing.T) {
		t.Parallel()

		v1 := &conway.Vertex{ID: 0, Position: conway.Vector3{0, 0, 0}}
		v2 := &conway.Vertex{ID: 1, Position: conway.Vector3{1, 0, 0}}
		e := &conway.Edge{V1: v1, V2: v2}

		p := &conway.Polyhedron{
			Vertices: map[int]*conway.Vertex{0: v1, 1: v2},
			Edges:    map[int]*conway.Edge{0: e},
		}
		stats := p.CalculateGeometryStats()

		// Function returns early if no faces, so all values are zero.
		assert.NotNil(t, stats)
		assert.Equal(t, 0.0, stats.MinEdgeLength)
		assert.Equal(t, 0.0, stats.MaxEdgeLength)
		assert.Equal(t, 0.0, stats.AvgEdgeLength)
		assert.Equal(t, 0.0, stats.MinFaceArea)
		assert.Equal(t, 0.0, stats.MaxFaceArea)
		assert.Equal(t, 0.0, stats.AvgFaceArea)
	})

	t.Run("SingleVertex", func(t *testing.T) {
		t.Parallel()

		v := &conway.Vertex{ID: 0, Position: conway.Vector3{1, 2, 3}}
		p := &conway.Polyhedron{
			Vertices: map[int]*conway.Vertex{0: v},
			Edges:    map[int]*conway.Edge{}, // Empty but not nil
			Faces:    map[int]*conway.Face{}, // Empty but not nil
		}
		stats := p.CalculateGeometryStats()

		// Function returns early if no edges or faces, so bounding box is not set.
		assert.NotNil(t, stats)
		assert.Equal(t, conway.Vector3{0, 0, 0}, stats.BoundingBox.Min)
		assert.Equal(t, conway.Vector3{0, 0, 0}, stats.BoundingBox.Max)
	})

	t.Run("MultipleVertices", func(t *testing.T) {
		t.Parallel()

		v1 := &conway.Vertex{ID: 0, Position: conway.Vector3{-1, -2, -3}}
		v2 := &conway.Vertex{ID: 1, Position: conway.Vector3{4, 5, 6}}
		v3 := &conway.Vertex{ID: 2, Position: conway.Vector3{0, 1, 2}}

		p := &conway.Polyhedron{
			Vertices: map[int]*conway.Vertex{0: v1, 1: v2, 2: v3},
			Edges:    map[int]*conway.Edge{}, // Empty but not nil
			Faces:    map[int]*conway.Face{}, // Empty but not nil
		}
		stats := p.CalculateGeometryStats()

		// Function returns early if no edges or faces, so bounding box is not set.
		assert.NotNil(t, stats)
		assert.Equal(t, conway.Vector3{0, 0, 0}, stats.BoundingBox.Min)
		assert.Equal(t, conway.Vector3{0, 0, 0}, stats.BoundingBox.Max)
	})

	t.Run("ValidCube", func(t *testing.T) {
		t.Parallel()

		// Create a simple cube for testing.
		cube, err := conway.Parse("C")
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

		// Cube should have uniform edge lengths and face areas.
		assert.InDelta(t, stats.MinEdgeLength, stats.MaxEdgeLength, 1e-10)
		assert.InDelta(t, stats.MinFaceArea, stats.MaxFaceArea, 1e-10)
	})
}

func TestEdgeLookup(t *testing.T) {
	t.Parallel()

	t.Run("BasicOperations", func(t *testing.T) {
		t.Parallel()

		el := conway.NewEdgeLookup()

		v1 := &conway.Vertex{ID: 1}
		v2 := &conway.Vertex{ID: 2}
		edge := &conway.Edge{V1: v1, V2: v2}

		// Test Add.
		el.Add(edge)

		// Test Find.
		found := el.Find(1, 2)
		assert.Equal(t, edge, found)

		// Test symmetric lookup.
		found = el.Find(2, 1)
		assert.Equal(t, edge, found)

		// Test Remove.
		el.Remove(edge)
		found = el.Find(1, 2)
		assert.Nil(t, found)
	})

	t.Run("NonexistentEdge", func(t *testing.T) {
		t.Parallel()

		el := conway.NewEdgeLookup()
		found := el.Find(99, 100)
		assert.Nil(t, found)
	})
}

func TestCalculateFaceNormal(t *testing.T) {
	t.Parallel()

	t.Run("InsufficientVertices", func(t *testing.T) {
		t.Parallel()

		vertices := []*conway.Vertex{
			{Position: conway.Vector3{0, 0, 0}},
			{Position: conway.Vector3{1, 0, 0}},
		}

		_, err := conway.CalculateFaceNormal(vertices)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insufficient vertices")
	})

	t.Run("ValidTriangle", func(t *testing.T) {
		t.Parallel()

		vertices := []*conway.Vertex{
			{Position: conway.Vector3{0, 0, 0}},
			{Position: conway.Vector3{1, 0, 0}},
			{Position: conway.Vector3{0, 1, 0}},
		}

		normal, err := conway.CalculateFaceNormal(vertices)
		assert.NoError(t, err)
		assert.InDelta(t, 0.0, normal.X, 1e-10)
		assert.InDelta(t, 0.0, normal.Y, 1e-10)
		assert.InDelta(t, 1.0, normal.Z, 1e-10)
	})

	t.Run("DegenerateFace", func(t *testing.T) {
		t.Parallel()

		// All vertices in a line.
		vertices := []*conway.Vertex{
			{Position: conway.Vector3{0, 0, 0}},
			{Position: conway.Vector3{1, 0, 0}},
			{Position: conway.Vector3{2, 0, 0}},
		}

		_, err := conway.CalculateFaceNormal(vertices)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "degenerate face normal")
	})
}

func TestEnsureCounterClockwise(t *testing.T) {
	t.Parallel()

	t.Run("InsufficientVertices", func(t *testing.T) {
		t.Parallel()

		vertices := []*conway.Vertex{
			{Position: conway.Vector3{0, 0, 0}},
			{Position: conway.Vector3{1, 0, 0}},
		}
		center := conway.Vector3{0, 0, 0}

		result := conway.EnsureCounterClockwise(vertices, center)
		assert.Equal(t, vertices, result)
	})

	t.Run("ValidTriangle", func(t *testing.T) {
		t.Parallel()

		vertices := []*conway.Vertex{
			{Position: conway.Vector3{1, 0, 0}},
			{Position: conway.Vector3{0, 1, 0}},
			{Position: conway.Vector3{0, 0, 1}},
		}
		center := conway.Vector3{0, 0, 0}

		result := conway.EnsureCounterClockwise(vertices, center)
		assert.NotNil(t, result)
		assert.Len(t, result, 3)
	})
}

func TestMakeEdgeKey(t *testing.T) {
	t.Parallel()

	t.Run("ConsistentOrdering", func(t *testing.T) {
		t.Parallel()

		key1 := conway.MakeEdgeKey(1, 2)
		key2 := conway.MakeEdgeKey(2, 1)
		assert.Equal(t, key1, key2)
	})

	t.Run("DifferentVertices", func(t *testing.T) {
		t.Parallel()

		key1 := conway.MakeEdgeKey(1, 2)
		key2 := conway.MakeEdgeKey(1, 3)
		assert.NotEqual(t, key1, key2)
	})
}

// Allocation functions (allocateVertexSlice, allocateEdgeSlice) are now tested
// indirectly through polyhedron construction operations rather than directly.
