package conway_test

import (
	"testing"

	"github.com/sksmith/conway/conway"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTruncateOpApply(t *testing.T) {
	t.Parallel()

	t.Run("EmptyPolyhedron", func(t *testing.T) {
		t.Parallel()

		empty := &conway.Polyhedron{
			Name:     "empty",
			Vertices: map[int]*conway.Vertex{},
			Edges:    map[int]*conway.Edge{},
			Faces:    map[int]*conway.Face{},
		}

		truncateOp := conway.TruncateOp{}
		result := truncateOp.Apply(empty)

		assert.NotNil(t, result)
		assert.Equal(t, "tempty", result.Name)
		assert.Equal(t, 0, len(result.Vertices))
		assert.Equal(t, 0, len(result.Edges))
		assert.Equal(t, 0, len(result.Faces))
	})

	t.Run("SingleTriangle", func(t *testing.T) {
		t.Parallel()

		// Create a simple triangle polyhedron.
		triangle := conway.NewPolyhedron("triangle")

		// Add vertices.
		v1 := triangle.AddVertex(conway.Vector3{0, 0, 0})
		v2 := triangle.AddVertex(conway.Vector3{1, 0, 0})
		v3 := triangle.AddVertex(conway.Vector3{0.5, 1, 0})

		// Add face.
		triangle.AddFace([]*conway.Vertex{v1, v2, v3})

		truncateOp := conway.TruncateOp{}
		result := truncateOp.Apply(triangle)

		assert.NotNil(t, result)
		assert.Equal(t, "ttriangle", result.Name)
		// Truncating a triangle should create multiple vertices.
		assert.Greater(t, len(result.Vertices), 3)
	})

	t.Run("ValidTetrahedron", func(t *testing.T) {
		t.Parallel()

		tetra := conway.Tetrahedron()
		require.NotNil(t, tetra)

		truncateOp := conway.TruncateOp{}
		result := truncateOp.Apply(tetra)

		assert.NotNil(t, result)
		assert.True(t, result.IsValid())
		assert.Equal(t, 2, result.EulerCharacteristic())
		assert.Equal(t, "tTetrahedron", result.Name)

		// Truncated tetrahedron should have more vertices than original.
		assert.Greater(t, len(result.Vertices), len(tetra.Vertices))
		assert.Greater(t, len(result.Faces), len(tetra.Faces))
	})

	t.Run("ValidCube", func(t *testing.T) {
		t.Parallel()

		cube := conway.Cube()
		require.NotNil(t, cube)

		truncateOp := conway.TruncateOp{}
		result := truncateOp.Apply(cube)

		assert.NotNil(t, result)
		assert.True(t, result.IsValid())
		assert.Equal(t, 2, result.EulerCharacteristic())
		assert.Equal(t, "tCube", result.Name)

		// Truncated cube should have specific properties.
		assert.Greater(t, len(result.Vertices), len(cube.Vertices))
		assert.Greater(t, len(result.Faces), len(cube.Faces))
	})

	t.Run("ValidOctahedron", func(t *testing.T) {
		t.Parallel()

		octa := conway.Octahedron()
		require.NotNil(t, octa)

		truncateOp := conway.TruncateOp{}
		result := truncateOp.Apply(octa)

		assert.NotNil(t, result)
		assert.True(t, result.IsValid())
		assert.Equal(t, 2, result.EulerCharacteristic())
		assert.Equal(t, "tOctahedron", result.Name)
	})

	t.Run("ValidDodecahedron", func(t *testing.T) {
		t.Parallel()

		dodeca := conway.Dodecahedron()
		require.NotNil(t, dodeca)

		truncateOp := conway.TruncateOp{}
		result := truncateOp.Apply(dodeca)

		assert.NotNil(t, result)
		assert.True(t, result.IsValid())
		assert.Equal(t, 2, result.EulerCharacteristic())
		assert.Equal(t, "tDodecahedron", result.Name)

		// Truncated dodecahedron (soccer ball) is a well-known shape.
		assert.Greater(t, len(result.Vertices), 50)
		assert.Greater(t, len(result.Faces), 25)
	})

	t.Run("ValidIcosahedron", func(t *testing.T) {
		t.Parallel()

		icosa := conway.Icosahedron()
		require.NotNil(t, icosa)

		truncateOp := conway.TruncateOp{}
		result := truncateOp.Apply(icosa)

		assert.NotNil(t, result)
		assert.True(t, result.IsValid())
		assert.Equal(t, 2, result.EulerCharacteristic())
		assert.Equal(t, "tIcosahedron", result.Name)
	})

	t.Run("EdgeVertexKeyGeneration", func(t *testing.T) {
		t.Parallel()

		// Test the edge vertex key generation function.
		key1 := conway.EdgeVertexKey(1, 2)
		key2 := conway.EdgeVertexKey(2, 1)
		key3 := conway.EdgeVertexKey(1, 2)

		assert.NotEqual(t, key1, key2) // Different order should give different keys
		assert.Equal(t, key1, key3)    // Same parameters should give same key
		assert.Equal(t, "1_2", key1)
	})

	t.Run("TruncationFactorBehavior", func(t *testing.T) {
		t.Parallel()

		// Test that truncation creates new vertices at the expected positions.
		tetra := conway.Tetrahedron()
		require.NotNil(t, tetra)

		truncateOp := conway.TruncateOp{}
		result := truncateOp.Apply(tetra)

		// Check that all vertices are at reasonable positions (not at original positions)
		for _, vertex := range result.Vertices {
			pos := vertex.Position
			length := pos.Length()

			// Vertices should be at reasonable distances (not 0, not too far)
			assert.Greater(t, length, 0.001)
			assert.Less(t, length, 10.0)
		}
	})

	t.Run("FaceVertexValidation", func(t *testing.T) {
		t.Parallel()

		// Test that all generated faces have at least 3 vertices.
		cube := conway.Cube()
		require.NotNil(t, cube)

		truncateOp := conway.TruncateOp{}
		result := truncateOp.Apply(cube)

		for _, face := range result.Faces {
			assert.GreaterOrEqual(t, len(face.Vertices), 3,
				"Face should have at least 3 vertices")
			assert.Equal(t, len(face.Vertices), len(face.Edges),
				"Face should have equal vertices and edges")
		}
	})

	t.Run("VertexDegreeConsistency", func(t *testing.T) {
		t.Parallel()

		// Test vertex degree consistency after truncation.
		octa := conway.Octahedron()
		require.NotNil(t, octa)

		truncateOp := conway.TruncateOp{}
		result := truncateOp.Apply(octa)

		for _, vertex := range result.Vertices {
			assert.GreaterOrEqual(t, vertex.Degree(), 2,
				"Vertex should have degree >= 2")
			assert.Equal(t, len(vertex.Edges), len(vertex.Faces),
				"Vertex edges and faces should match for manifold")
		}
	})
}

func TestTruncateFunction(t *testing.T) {
	t.Parallel()

	t.Run("ConvenienceFunction", func(t *testing.T) {
		t.Parallel()

		cube := conway.Cube()
		require.NotNil(t, cube)

		// Test the convenience function.
		result := conway.Truncate(cube)

		assert.NotNil(t, result)
		assert.True(t, result.IsValid())
		assert.Equal(t, "tCube", result.Name)
	})
}

func TestTruncateOpMethods(t *testing.T) {
	t.Parallel()

	t.Run("Symbol", func(t *testing.T) {
		t.Parallel()

		op := conway.TruncateOp{}
		assert.Equal(t, "t", op.Symbol())
	})

	t.Run("Name", func(t *testing.T) {
		t.Parallel()

		op := conway.TruncateOp{}
		assert.Equal(t, "truncate", op.Name())
	})
}
