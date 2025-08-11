package conway_test

import (
	"testing"

	"github.com/sksmith/conway/conway"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderEdgesAroundVertex(t *testing.T) {
	t.Parallel()
	t.Run("EmptyEdges", func(t *testing.T) {
		t.Parallel()

		vertex := &conway.Vertex{ID: 0, Edges: map[int]*conway.Edge{}}
		result := conway.OrderEdgesAroundVertex(vertex)
		assert.Empty(t, result)
	})

	t.Run("SingleEdge", func(t *testing.T) {
		t.Parallel()

		vertex := &conway.Vertex{ID: 0}
		edge := &conway.Edge{ID: 1, V1: vertex, V2: &conway.Vertex{ID: 1}}
		vertex.Edges = map[int]*conway.Edge{1: edge}

		result := conway.OrderEdgesAroundVertex(vertex)
		assert.Len(t, result, 1)
		assert.Equal(t, edge, result[0])
	})

	t.Run("TwoEdges", func(t *testing.T) {
		t.Parallel()

		vertex := &conway.Vertex{ID: 0}
		edge1 := &conway.Edge{ID: 1, V1: vertex, V2: &conway.Vertex{ID: 1}}
		edge2 := &conway.Edge{ID: 2, V1: vertex, V2: &conway.Vertex{ID: 2}}
		vertex.Edges = map[int]*conway.Edge{1: edge1, 2: edge2}

		result := conway.OrderEdgesAroundVertex(vertex)
		assert.Len(t, result, 2)
		// With 2 or fewer edges, should return them as-is.
		assert.Contains(t, result, edge1)
		assert.Contains(t, result, edge2)
	})

	t.Run("DisconnectedEdges", func(t *testing.T) {
		t.Parallel()
		// Create edges that don't share faces (fallback case)
		vertex := &conway.Vertex{ID: 0}
		edge1 := &conway.Edge{ID: 1, V1: vertex, V2: &conway.Vertex{ID: 1}}
		edge2 := &conway.Edge{ID: 2, V1: vertex, V2: &conway.Vertex{ID: 2}}
		edge3 := &conway.Edge{ID: 3, V1: vertex, V2: &conway.Vertex{ID: 3}}

		vertex.Edges = map[int]*conway.Edge{1: edge1, 2: edge2, 3: edge3}
		vertex.Faces = map[int]*conway.Face{} // No faces - should trigger fallback

		result := conway.OrderEdgesAroundVertex(vertex)
		assert.Len(t, result, 3)

		// All edges should be present (though order may be arbitrary)
		assert.Contains(t, result, edge1)
		assert.Contains(t, result, edge2)
		assert.Contains(t, result, edge3)
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
			if len(v.Edges) >= 3 {
				testVertex = v
				break
			}
		}

		require.NotNil(t, testVertex)

		result := conway.OrderEdgesAroundVertex(testVertex)

		// Should return all edges for this vertex.
		assert.Len(t, result, len(testVertex.Edges))

		// All edges should connect to the test vertex.
		for _, edge := range result {
			assert.True(t, edge.V1.ID == testVertex.ID || edge.V2.ID == testVertex.ID,
				"Edge should connect to the test vertex")
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

		result := conway.OrderEdgesAroundVertex(testVertex)

		// Should return all edges for this vertex.
		assert.Len(t, result, len(testVertex.Edges))

		// Verify each edge connects to our vertex.
		for _, edge := range result {
			assert.True(t, edge.V1.ID == testVertex.ID || edge.V2.ID == testVertex.ID)
		}
	})

	// Test edge cases and boundary conditions.
	t.Run("LargeNumberOfEdges", func(t *testing.T) {
		t.Parallel()
		// Test the function with many edges (stress test)
		octahedron := conway.Octahedron()
		require.NotNil(t, octahedron)

		// Get a vertex with multiple edges.
		var testVertex *conway.Vertex

		for _, v := range octahedron.Vertices {
			if len(v.Edges) > 3 {
				testVertex = v
				break
			}
		}

		require.NotNil(t, testVertex)

		result := conway.OrderEdgesAroundVertex(testVertex)

		// Should return all edges and they should all be unique.
		assert.Len(t, result, len(testVertex.Edges))

		// Check for duplicates.
		seen := make(map[int]bool)
		for _, edge := range result {
			assert.False(t, seen[edge.ID], "Found duplicate edge ID %d", edge.ID)
			seen[edge.ID] = true
		}
	})

	t.Run("ComplexPolyhedron", func(t *testing.T) {
		t.Parallel()
		// Test with a more complex polyhedron like dodecahedron.
		dodeca := conway.Dodecahedron()
		require.NotNil(t, dodeca)

		// Test a few vertices.
		count := 0
		for _, v := range dodeca.Vertices {
			if count >= 3 { // Test first 3 vertices
				break
			}

			result := conway.OrderEdgesAroundVertex(v)

			// Should return all edges for this vertex.
			assert.Len(t, result, len(v.Edges))

			// All edges should connect to the test vertex.
			for _, edge := range result {
				assert.True(t, edge.V1.ID == v.ID || edge.V2.ID == v.ID,
					"Edge should connect to vertex %d", v.ID)
			}

			count++
		}
	})
}
