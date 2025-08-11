package conway_test

import (
	"testing"

	"github.com/sksmith/conway/conway"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderFacesAroundVertex(t *testing.T) {
	t.Parallel()

	t.Run("EmptyFaces", func(t *testing.T) {
		t.Parallel()

		vertex := &conway.Vertex{ID: 0, Faces: map[int]*conway.Face{}}
		result := conway.OrderFacesAroundVertex(vertex)

		assert.Empty(t, result)
	})

	t.Run("SingleFace", func(t *testing.T) {
		t.Parallel()

		vertex := &conway.Vertex{ID: 0}
		face := &conway.Face{ID: 1}

		vertex.Faces = map[int]*conway.Face{1: face}

		result := conway.OrderFacesAroundVertex(vertex)

		assert.Len(t, result, 1)
		assert.Equal(t, face, result[0])
	})

	t.Run("TwoFaces", func(t *testing.T) {
		t.Parallel()

		vertex := &conway.Vertex{ID: 0}
		face1 := &conway.Face{ID: 1}
		face2 := &conway.Face{ID: 2}

		vertex.Faces = map[int]*conway.Face{1: face1, 2: face2}

		result := conway.OrderFacesAroundVertex(vertex)

		assert.Len(t, result, 2)
		// With 2 or fewer faces, should return them as-is.
		assert.Contains(t, result, face1)
		assert.Contains(t, result, face2)
	})

	t.Run("ThreeFacesWithSharedEdges", func(t *testing.T) {
		t.Parallel()

		// Create a vertex with 3 faces that share edges.
		vertex := &conway.Vertex{ID: 0}
		v1 := &conway.Vertex{ID: 1}
		v2 := &conway.Vertex{ID: 2}
		v3 := &conway.Vertex{ID: 3}

		// Create edges.
		edge1 := &conway.Edge{ID: 1, V1: vertex, V2: v1}
		edge2 := &conway.Edge{ID: 2, V1: vertex, V2: v2}
		edge3 := &conway.Edge{ID: 3, V1: vertex, V2: v3}
		sharedEdge1 := &conway.Edge{ID: 4, V1: v1, V2: v2}
		sharedEdge2 := &conway.Edge{ID: 5, V1: v2, V2: v3}

		// Create faces that share edges.
		face1 := &conway.Face{
			ID:       1,
			Vertices: []*conway.Vertex{vertex, v1, v2},
			Edges:    []*conway.Edge{edge1, sharedEdge1, edge2},
		}
		face2 := &conway.Face{
			ID:       2,
			Vertices: []*conway.Vertex{vertex, v2, v3},
			Edges:    []*conway.Edge{edge2, sharedEdge2, edge3},
		}
		face3 := &conway.Face{
			ID:       3,
			Vertices: []*conway.Vertex{vertex, v3, v1},
			Edges:    []*conway.Edge{edge3, {ID: 6}, edge1},
		}

		vertex.Faces = map[int]*conway.Face{1: face1, 2: face2, 3: face3}
		vertex.Edges = map[int]*conway.Edge{1: edge1, 2: edge2, 3: edge3}

		// Set up back-references for edge faces.
		edge1.Faces = map[int]*conway.Face{1: face1, 3: face3}
		edge2.Faces = map[int]*conway.Face{1: face1, 2: face2}
		edge3.Faces = map[int]*conway.Face{2: face2, 3: face3}

		result := conway.OrderFacesAroundVertex(vertex)

		assert.Len(t, result, 3)

		// All faces should be present.
		assert.Contains(t, result, face1)
		assert.Contains(t, result, face2)
		assert.Contains(t, result, face3)
	})

	t.Run("DisconnectedFaces", func(t *testing.T) {
		t.Parallel()

		// Create faces that don't share edges (fallback case)
		vertex := &conway.Vertex{ID: 0}
		face1 := &conway.Face{ID: 1}
		face2 := &conway.Face{ID: 2}
		face3 := &conway.Face{ID: 3}

		vertex.Faces = map[int]*conway.Face{1: face1, 2: face2, 3: face3}
		vertex.Edges = map[int]*conway.Edge{} // No edges - should trigger fallback

		result := conway.OrderFacesAroundVertex(vertex)
		assert.Len(t, result, 3)

		// All faces should be present (though order may be arbitrary)
		assert.Contains(t, result, face1)
		assert.Contains(t, result, face2)
		assert.Contains(t, result, face3)
	})

	t.Run("ValidCubeVertex", func(t *testing.T) {
		t.Parallel()

		// Test with a real cube vertex.
		cube := conway.Cube()

		require.NotNil(t, cube)
		require.True(t, len(cube.Vertices) > 0)

		// Get a vertex from the cube.
		var testVertex *conway.Vertex

		for _, v := range cube.Vertices {
			if len(v.Faces) >= 3 {
				testVertex = v
				break
			}
		}

		require.NotNil(t, testVertex)

		result := conway.OrderFacesAroundVertex(testVertex)

		// Should return all faces for this vertex.
		assert.Len(t, result, len(testVertex.Faces))

		// All returned faces should be in the original set.
		for _, face := range result {
			_, exists := testVertex.Faces[face.ID]

			assert.True(t, exists, "Face should be associated with the test vertex")
		}
	})

	t.Run("ValidTetrahedronVertex", func(t *testing.T) {
		t.Parallel()

		// Test with a tetrahedron vertex (simpler case)
		tetra := conway.Tetrahedron()

		require.NotNil(t, tetra)

		// Get a vertex from the tetrahedron.
		var testVertex *conway.Vertex

		for _, v := range tetra.Vertices {
			testVertex = v
			break
		}

		require.NotNil(t, testVertex)

		result := conway.OrderFacesAroundVertex(testVertex)

		// Should return all faces for this vertex.
		assert.Len(t, result, len(testVertex.Faces))

		// Verify each face is associated with our vertex.
		for _, face := range result {
			_, exists := testVertex.Faces[face.ID]

			assert.True(t, exists)
		}
	})

	t.Run("LargeNumberOfFaces", func(t *testing.T) {
		t.Parallel()
		// Test with a complex polyhedron like dodecahedron.
		dodeca := conway.Dodecahedron()
		require.NotNil(t, dodeca)

		// Get a vertex with multiple faces.
		var testVertex *conway.Vertex

		for _, v := range dodeca.Vertices {
			if len(v.Faces) > 2 {
				testVertex = v
				break
			}
		}

		require.NotNil(t, testVertex)

		result := conway.OrderFacesAroundVertex(testVertex)

		// Should return all faces and they should all be unique.
		assert.Len(t, result, len(testVertex.Faces))

		// Check for duplicates.
		seen := make(map[int]bool)
		for _, face := range result {
			assert.False(t, seen[face.ID], "Found duplicate face ID %d", face.ID)
			seen[face.ID] = true
		}
	})

	t.Run("ComplexPolyhedron", func(t *testing.T) {
		t.Parallel()

		// Test with icosahedron.
		icosa := conway.Icosahedron()

		require.NotNil(t, icosa)

		// Test a few vertices.
		count := 0

		for _, v := range icosa.Vertices {
			if count >= 3 { // Test first 3 vertices
				break
			}

			result := conway.OrderFacesAroundVertex(v)

			// Should return all faces for this vertex.
			assert.Len(t, result, len(v.Faces))

			// All faces should be associated with the vertex.
			for _, face := range result {
				_, exists := v.Faces[face.ID]

				assert.True(t, exists, "Face should be associated with vertex %d", v.ID)
			}

			count++
		}
	})
}

func TestFindEdgeIndex(t *testing.T) {
	t.Parallel()

	t.Run("EdgeExists", func(t *testing.T) {
		t.Parallel()

		edge1 := &conway.Edge{ID: 1}
		edge2 := &conway.Edge{ID: 2}
		edge3 := &conway.Edge{ID: 3}

		face := &conway.Face{
			Edges: []*conway.Edge{edge1, edge2, edge3},
		}

		assert.Equal(t, 0, conway.FindEdgeIndex(face, edge1))
		assert.Equal(t, 1, conway.FindEdgeIndex(face, edge2))
		assert.Equal(t, 2, conway.FindEdgeIndex(face, edge3))
	})

	t.Run("EdgeDoesNotExist", func(t *testing.T) {
		t.Parallel()

		edge1 := &conway.Edge{ID: 1}
		edge2 := &conway.Edge{ID: 2}
		nonExistentEdge := &conway.Edge{ID: 99}

		face := &conway.Face{
			Edges: []*conway.Edge{edge1, edge2},
		}

		assert.Equal(t, -1, conway.FindEdgeIndex(face, nonExistentEdge))
	})

	t.Run("EmptyFace", func(t *testing.T) {
		t.Parallel()

		edge := &conway.Edge{ID: 1}
		face := &conway.Face{Edges: []*conway.Edge{}}

		assert.Equal(t, -1, conway.FindEdgeIndex(face, edge))
	})
}

func TestDualOpApply(t *testing.T) {
	t.Parallel()

	t.Run("ValidTetrahedron", func(t *testing.T) {
		t.Parallel()

		tetra := conway.Tetrahedron()

		require.NotNil(t, tetra)

		dualOp := conway.DualOp{}
		result := dualOp.Apply(tetra)

		assert.NotNil(t, result)
		assert.True(t, result.IsValid())
		assert.Equal(t, 2, result.EulerCharacteristic())
		assert.Equal(t, "dTetrahedron", result.Name)

		// Dual of tetrahedron is also a tetrahedron.
		assert.Equal(t, len(tetra.Vertices), len(result.Faces))
		assert.Equal(t, len(tetra.Faces), len(result.Vertices))
	})

	t.Run("ValidCube", func(t *testing.T) {
		t.Parallel()

		cube := conway.Cube()

		require.NotNil(t, cube)

		dualOp := conway.DualOp{}
		result := dualOp.Apply(cube)

		assert.NotNil(t, result)
		assert.True(t, result.IsValid())
		assert.Equal(t, 2, result.EulerCharacteristic())
		assert.Equal(t, "dCube", result.Name)

		// Dual of cube is octahedron.
		assert.Equal(t, len(cube.Vertices), len(result.Faces))
		assert.Equal(t, len(cube.Faces), len(result.Vertices))
	})

	t.Run("ValidOctahedron", func(t *testing.T) {
		t.Parallel()

		octa := conway.Octahedron()

		require.NotNil(t, octa)

		dualOp := conway.DualOp{}
		result := dualOp.Apply(octa)

		assert.NotNil(t, result)
		assert.True(t, result.IsValid())
		assert.Equal(t, 2, result.EulerCharacteristic())
		assert.Equal(t, "dOctahedron", result.Name)

		// Dual of octahedron is cube.
		assert.Equal(t, len(octa.Vertices), len(result.Faces))
		assert.Equal(t, len(octa.Faces), len(result.Vertices))
	})
}

func TestDualFunction(t *testing.T) {
	t.Parallel()
	t.Run("ConvenienceFunction", func(t *testing.T) {
		t.Parallel()

		cube := conway.Cube()
		require.NotNil(t, cube)

		// Test the convenience function.
		result := conway.Dual(cube)

		assert.NotNil(t, result)
		assert.True(t, result.IsValid())
		assert.Equal(t, "dCube", result.Name)
	})
}

func TestDualOpMethods(t *testing.T) {
	t.Parallel()
	t.Run("Symbol", func(t *testing.T) {
		t.Parallel()

		op := conway.DualOp{}
		assert.Equal(t, "d", op.Symbol())
	})

	t.Run("Name", func(t *testing.T) {
		t.Parallel()

		op := conway.DualOp{}
		assert.Equal(t, "dual", op.Name())
	})
}
