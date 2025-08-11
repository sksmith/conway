package conway

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderFacesAroundVertex(t *testing.T) {
	t.Run("EmptyFaces", func(t *testing.T) {
		vertex := &Vertex{ID: 0, Faces: map[int]*Face{}}
		result := orderFacesAroundVertex(vertex)
		assert.Empty(t, result)
	})

	t.Run("SingleFace", func(t *testing.T) {
		vertex := &Vertex{ID: 0}
		face := &Face{ID: 1}
		vertex.Faces = map[int]*Face{1: face}

		result := orderFacesAroundVertex(vertex)
		assert.Len(t, result, 1)
		assert.Equal(t, face, result[0])
	})

	t.Run("TwoFaces", func(t *testing.T) {
		vertex := &Vertex{ID: 0}
		face1 := &Face{ID: 1}
		face2 := &Face{ID: 2}
		vertex.Faces = map[int]*Face{1: face1, 2: face2}

		result := orderFacesAroundVertex(vertex)
		assert.Len(t, result, 2)
		// With 2 or fewer faces, should return them as-is
		assert.Contains(t, result, face1)
		assert.Contains(t, result, face2)
	})

	t.Run("ThreeFacesWithSharedEdges", func(t *testing.T) {
		// Create a vertex with 3 faces that share edges
		vertex := &Vertex{ID: 0}
		v1 := &Vertex{ID: 1}
		v2 := &Vertex{ID: 2}
		v3 := &Vertex{ID: 3}

		// Create edges
		edge1 := &Edge{ID: 1, V1: vertex, V2: v1}
		edge2 := &Edge{ID: 2, V1: vertex, V2: v2}
		edge3 := &Edge{ID: 3, V1: vertex, V2: v3}
		sharedEdge1 := &Edge{ID: 4, V1: v1, V2: v2}
		sharedEdge2 := &Edge{ID: 5, V1: v2, V2: v3}

		// Create faces that share edges
		face1 := &Face{
			ID:       1,
			Vertices: []*Vertex{vertex, v1, v2},
			Edges:    []*Edge{edge1, sharedEdge1, edge2},
		}
		face2 := &Face{
			ID:       2,
			Vertices: []*Vertex{vertex, v2, v3},
			Edges:    []*Edge{edge2, sharedEdge2, edge3},
		}
		face3 := &Face{
			ID:       3,
			Vertices: []*Vertex{vertex, v3, v1},
			Edges:    []*Edge{edge3, {ID: 6}, edge1},
		}

		vertex.Faces = map[int]*Face{1: face1, 2: face2, 3: face3}
		vertex.Edges = map[int]*Edge{1: edge1, 2: edge2, 3: edge3}

		// Set up back-references for edge faces
		edge1.Faces = map[int]*Face{1: face1, 3: face3}
		edge2.Faces = map[int]*Face{1: face1, 2: face2}
		edge3.Faces = map[int]*Face{2: face2, 3: face3}

		result := orderFacesAroundVertex(vertex)
		assert.Len(t, result, 3)

		// All faces should be present
		assert.Contains(t, result, face1)
		assert.Contains(t, result, face2)
		assert.Contains(t, result, face3)
	})

	t.Run("DisconnectedFaces", func(t *testing.T) {
		// Create faces that don't share edges (fallback case)
		vertex := &Vertex{ID: 0}
		face1 := &Face{ID: 1}
		face2 := &Face{ID: 2}
		face3 := &Face{ID: 3}

		vertex.Faces = map[int]*Face{1: face1, 2: face2, 3: face3}
		vertex.Edges = map[int]*Edge{} // No edges - should trigger fallback

		result := orderFacesAroundVertex(vertex)
		assert.Len(t, result, 3)

		// All faces should be present (though order may be arbitrary)
		assert.Contains(t, result, face1)
		assert.Contains(t, result, face2)
		assert.Contains(t, result, face3)
	})

	t.Run("ValidCubeVertex", func(t *testing.T) {
		// Test with a real cube vertex
		cube := Cube()
		require.NotNil(t, cube)
		require.True(t, len(cube.Vertices) > 0)

		// Get a vertex from the cube
		var testVertex *Vertex
		for _, v := range cube.Vertices {
			if len(v.Faces) >= 3 {
				testVertex = v
				break
			}
		}
		require.NotNil(t, testVertex)

		result := orderFacesAroundVertex(testVertex)

		// Should return all faces for this vertex
		assert.Len(t, result, len(testVertex.Faces))

		// All returned faces should be in the original set
		for _, face := range result {
			_, exists := testVertex.Faces[face.ID]
			assert.True(t, exists, "Face should be associated with the test vertex")
		}
	})

	t.Run("ValidTetrahedronVertex", func(t *testing.T) {
		// Test with a tetrahedron vertex (simpler case)
		tetra := Tetrahedron()
		require.NotNil(t, tetra)

		// Get a vertex from the tetrahedron
		var testVertex *Vertex
		for _, v := range tetra.Vertices {
			testVertex = v
			break
		}
		require.NotNil(t, testVertex)

		result := orderFacesAroundVertex(testVertex)

		// Should return all faces for this vertex
		assert.Len(t, result, len(testVertex.Faces))

		// Verify each face is associated with our vertex
		for _, face := range result {
			_, exists := testVertex.Faces[face.ID]
			assert.True(t, exists)
		}
	})

	t.Run("LargeNumberOfFaces", func(t *testing.T) {
		// Test with a complex polyhedron like dodecahedron
		dodeca := Dodecahedron()
		require.NotNil(t, dodeca)

		// Get a vertex with multiple faces
		var testVertex *Vertex
		for _, v := range dodeca.Vertices {
			if len(v.Faces) > 2 {
				testVertex = v
				break
			}
		}
		require.NotNil(t, testVertex)

		result := orderFacesAroundVertex(testVertex)

		// Should return all faces and they should all be unique
		assert.Len(t, result, len(testVertex.Faces))

		// Check for duplicates
		seen := make(map[int]bool)
		for _, face := range result {
			assert.False(t, seen[face.ID], "Found duplicate face ID %d", face.ID)
			seen[face.ID] = true
		}
	})

	t.Run("ComplexPolyhedron", func(t *testing.T) {
		// Test with icosahedron
		icosa := Icosahedron()
		require.NotNil(t, icosa)

		// Test a few vertices
		count := 0
		for _, v := range icosa.Vertices {
			if count >= 3 { // Test first 3 vertices
				break
			}

			result := orderFacesAroundVertex(v)

			// Should return all faces for this vertex
			assert.Len(t, result, len(v.Faces))

			// All faces should be associated with the vertex
			for _, face := range result {
				_, exists := v.Faces[face.ID]
				assert.True(t, exists, "Face should be associated with vertex %d", v.ID)
			}

			count++
		}
	})
}

func TestFindEdgeIndex(t *testing.T) {
	t.Run("EdgeExists", func(t *testing.T) {
		edge1 := &Edge{ID: 1}
		edge2 := &Edge{ID: 2}
		edge3 := &Edge{ID: 3}

		face := &Face{
			Edges: []*Edge{edge1, edge2, edge3},
		}

		assert.Equal(t, 0, findEdgeIndex(face, edge1))
		assert.Equal(t, 1, findEdgeIndex(face, edge2))
		assert.Equal(t, 2, findEdgeIndex(face, edge3))
	})

	t.Run("EdgeDoesNotExist", func(t *testing.T) {
		edge1 := &Edge{ID: 1}
		edge2 := &Edge{ID: 2}
		nonExistentEdge := &Edge{ID: 99}

		face := &Face{
			Edges: []*Edge{edge1, edge2},
		}

		assert.Equal(t, -1, findEdgeIndex(face, nonExistentEdge))
	})

	t.Run("EmptyFace", func(t *testing.T) {
		edge := &Edge{ID: 1}
		face := &Face{Edges: []*Edge{}}

		assert.Equal(t, -1, findEdgeIndex(face, edge))
	})
}

func TestDualOpApply(t *testing.T) {
	t.Run("ValidTetrahedron", func(t *testing.T) {
		tetra := Tetrahedron()
		require.NotNil(t, tetra)

		dualOp := DualOp{}
		result := dualOp.Apply(tetra)

		assert.NotNil(t, result)
		assert.True(t, result.IsValid())
		assert.Equal(t, 2, result.EulerCharacteristic())
		assert.Equal(t, "dTetrahedron", result.Name)

		// Dual of tetrahedron is also a tetrahedron
		assert.Equal(t, len(tetra.Vertices), len(result.Faces))
		assert.Equal(t, len(tetra.Faces), len(result.Vertices))
	})

	t.Run("ValidCube", func(t *testing.T) {
		cube := Cube()
		require.NotNil(t, cube)

		dualOp := DualOp{}
		result := dualOp.Apply(cube)

		assert.NotNil(t, result)
		assert.True(t, result.IsValid())
		assert.Equal(t, 2, result.EulerCharacteristic())
		assert.Equal(t, "dCube", result.Name)

		// Dual of cube is octahedron
		assert.Equal(t, len(cube.Vertices), len(result.Faces))
		assert.Equal(t, len(cube.Faces), len(result.Vertices))
	})

	t.Run("ValidOctahedron", func(t *testing.T) {
		octa := Octahedron()
		require.NotNil(t, octa)

		dualOp := DualOp{}
		result := dualOp.Apply(octa)

		assert.NotNil(t, result)
		assert.True(t, result.IsValid())
		assert.Equal(t, 2, result.EulerCharacteristic())
		assert.Equal(t, "dOctahedron", result.Name)

		// Dual of octahedron is cube
		assert.Equal(t, len(octa.Vertices), len(result.Faces))
		assert.Equal(t, len(octa.Faces), len(result.Vertices))
	})
}

func TestDualFunction(t *testing.T) {
	t.Run("ConvenienceFunction", func(t *testing.T) {
		cube := Cube()
		require.NotNil(t, cube)

		// Test the convenience function
		result := Dual(cube)

		assert.NotNil(t, result)
		assert.True(t, result.IsValid())
		assert.Equal(t, "dCube", result.Name)
	})
}

func TestDualOpMethods(t *testing.T) {
	t.Run("Symbol", func(t *testing.T) {
		op := DualOp{}
		assert.Equal(t, "d", op.Symbol())
	})

	t.Run("Name", func(t *testing.T) {
		op := DualOp{}
		assert.Equal(t, "dual", op.Name())
	})
}
