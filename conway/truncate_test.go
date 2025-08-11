package conway

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTruncateOpApply(t *testing.T) {
	t.Run("EmptyPolyhedron", func(t *testing.T) {
		empty := &Polyhedron{
			Name:     "empty",
			Vertices: map[int]*Vertex{},
			Edges:    map[int]*Edge{},
			Faces:    map[int]*Face{},
		}

		truncateOp := TruncateOp{}
		result := truncateOp.Apply(empty)

		assert.NotNil(t, result)
		assert.Equal(t, "tempty", result.Name)
		assert.Equal(t, 0, len(result.Vertices))
		assert.Equal(t, 0, len(result.Edges))
		assert.Equal(t, 0, len(result.Faces))
	})

	t.Run("SingleTriangle", func(t *testing.T) {
		// Create a simple triangle polyhedron
		triangle := NewPolyhedron("triangle")

		// Add vertices
		v1 := triangle.AddVertex(Vector3{0, 0, 0})
		v2 := triangle.AddVertex(Vector3{1, 0, 0})
		v3 := triangle.AddVertex(Vector3{0.5, 1, 0})

		// Add face
		triangle.AddFace([]*Vertex{v1, v2, v3})

		truncateOp := TruncateOp{}
		result := truncateOp.Apply(triangle)

		assert.NotNil(t, result)
		assert.Equal(t, "ttriangle", result.Name)
		// Truncating a triangle should create multiple vertices
		assert.Greater(t, len(result.Vertices), 3)
	})

	t.Run("ValidTetrahedron", func(t *testing.T) {
		tetra := Tetrahedron()
		require.NotNil(t, tetra)

		truncateOp := TruncateOp{}
		result := truncateOp.Apply(tetra)

		assert.NotNil(t, result)
		assert.True(t, result.IsValid())
		assert.Equal(t, 2, result.EulerCharacteristic())
		assert.Equal(t, "tTetrahedron", result.Name)

		// Truncated tetrahedron should have more vertices than original
		assert.Greater(t, len(result.Vertices), len(tetra.Vertices))
		assert.Greater(t, len(result.Faces), len(tetra.Faces))
	})

	t.Run("ValidCube", func(t *testing.T) {
		cube := Cube()
		require.NotNil(t, cube)

		truncateOp := TruncateOp{}
		result := truncateOp.Apply(cube)

		assert.NotNil(t, result)
		assert.True(t, result.IsValid())
		assert.Equal(t, 2, result.EulerCharacteristic())
		assert.Equal(t, "tCube", result.Name)

		// Truncated cube should have specific properties
		assert.Greater(t, len(result.Vertices), len(cube.Vertices))
		assert.Greater(t, len(result.Faces), len(cube.Faces))
	})

	t.Run("ValidOctahedron", func(t *testing.T) {
		octa := Octahedron()
		require.NotNil(t, octa)

		truncateOp := TruncateOp{}
		result := truncateOp.Apply(octa)

		assert.NotNil(t, result)
		assert.True(t, result.IsValid())
		assert.Equal(t, 2, result.EulerCharacteristic())
		assert.Equal(t, "tOctahedron", result.Name)
	})

	t.Run("ValidDodecahedron", func(t *testing.T) {
		dodeca := Dodecahedron()
		require.NotNil(t, dodeca)

		truncateOp := TruncateOp{}
		result := truncateOp.Apply(dodeca)

		assert.NotNil(t, result)
		assert.True(t, result.IsValid())
		assert.Equal(t, 2, result.EulerCharacteristic())
		assert.Equal(t, "tDodecahedron", result.Name)

		// Truncated dodecahedron (soccer ball) is a well-known shape
		assert.Greater(t, len(result.Vertices), 50)
		assert.Greater(t, len(result.Faces), 25)
	})

	t.Run("ValidIcosahedron", func(t *testing.T) {
		icosa := Icosahedron()
		require.NotNil(t, icosa)

		truncateOp := TruncateOp{}
		result := truncateOp.Apply(icosa)

		assert.NotNil(t, result)
		assert.True(t, result.IsValid())
		assert.Equal(t, 2, result.EulerCharacteristic())
		assert.Equal(t, "tIcosahedron", result.Name)
	})

	t.Run("EdgeVertexKeyGeneration", func(t *testing.T) {
		// Test the edge vertex key generation function
		key1 := edgeVertexKey(1, 2)
		key2 := edgeVertexKey(2, 1)
		key3 := edgeVertexKey(1, 2)

		assert.NotEqual(t, key1, key2) // Different order should give different keys
		assert.Equal(t, key1, key3)    // Same parameters should give same key
		assert.Equal(t, "1_2", key1)
	})

	t.Run("TruncationFactorBehavior", func(t *testing.T) {
		// Test that truncation creates new vertices at the expected positions
		tetra := Tetrahedron()
		require.NotNil(t, tetra)

		truncateOp := TruncateOp{}
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
		// Test that all generated faces have at least 3 vertices
		cube := Cube()
		require.NotNil(t, cube)

		truncateOp := TruncateOp{}
		result := truncateOp.Apply(cube)

		for _, face := range result.Faces {
			assert.GreaterOrEqual(t, len(face.Vertices), 3,
				"Face should have at least 3 vertices")
			assert.Equal(t, len(face.Vertices), len(face.Edges),
				"Face should have equal vertices and edges")
		}
	})

	t.Run("VertexDegreeConsistency", func(t *testing.T) {
		// Test vertex degree consistency after truncation
		octa := Octahedron()
		require.NotNil(t, octa)

		truncateOp := TruncateOp{}
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
	t.Run("ConvenienceFunction", func(t *testing.T) {
		cube := Cube()
		require.NotNil(t, cube)

		// Test the convenience function
		result := Truncate(cube)

		assert.NotNil(t, result)
		assert.True(t, result.IsValid())
		assert.Equal(t, "tCube", result.Name)
	})
}

func TestTruncateOpMethods(t *testing.T) {
	t.Run("Symbol", func(t *testing.T) {
		op := TruncateOp{}
		assert.Equal(t, "t", op.Symbol())
	})

	t.Run("Name", func(t *testing.T) {
		op := TruncateOp{}
		assert.Equal(t, "truncate", op.Name())
	})
}
