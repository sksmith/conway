package conway_test

import (
	"testing"

	"github.com/sksmith/conway/conway"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIntegrationBasicOperations tests basic operations work correctly together.
func TestIntegrationBasicOperations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		notation      string
		expectedValid bool
		expectedEuler int
		minVertices   int
		minEdges      int
		minFaces      int
	}{
		{"Tetrahedron", "T", true, 2, 4, 6, 4},
		{"Cube", "C", true, 2, 8, 12, 6},
		{"Octahedron", "O", true, 2, 6, 12, 8},
		{"Dodecahedron", "D", true, 2, 20, 30, 12},
		{"Icosahedron", "I", true, 2, 12, 30, 20},
		{"Dual Cube", "dC", true, 2, 6, 12, 8},
		{"Truncated Cube", "tC", true, 2, 14, 36, 14},
		{"Dual Icosahedron", "dI", true, 2, 20, 30, 12},
		{"Truncated Icosahedron", "tI", true, 2, 60, 90, 32},
		{"Ambo Cube", "aC", true, 2, 12, 24, 14},
		{"Kis Cube", "kC", true, 2, 14, 36, 24},
		{"Join Cube", "jC", true, 2, 12, 24, 14},
		{"Complex operation", "dtC", true, 2, 14, 36, 14},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			poly, err := conway.Parse(tt.notation)
			require.NoError(t, err, "Failed to parse notation: %s", tt.notation)
			require.NotNil(t, poly, "Polyhedron should not be nil")

			// Test basic validity.
			if tt.expectedValid {
				assert.True(t, poly.IsValid(), "Polyhedron should be valid: %s", poly.Stats())
			}

			// Test Euler characteristic.
			assert.Equal(t, tt.expectedEuler, poly.EulerCharacteristic(),
				"Euler characteristic mismatch for %s: %s", tt.notation, poly.Stats())

			// Test minimum counts (operations should not decrease complexity dramatically)
			assert.GreaterOrEqual(t, len(poly.Vertices), tt.minVertices,
				"Too few vertices for %s", tt.notation)
			assert.GreaterOrEqual(t, len(poly.Edges), tt.minEdges,
				"Too few edges for %s", tt.notation)
			assert.GreaterOrEqual(t, len(poly.Faces), tt.minFaces,
				"Too few faces for %s", tt.notation)
		})
	}
}

// TestIntegrationDualInvolution tests that dual is an involution (dd = identity).
func TestIntegrationDualInvolution(t *testing.T) {
	t.Parallel()

	seeds := []string{"T", "C", "O", "D", "I"}

	for _, seed := range seeds {
		t.Run("Dual_involution_"+seed, func(t *testing.T) {
			t.Parallel()

			original, err := conway.Parse(seed)
			require.NoError(t, err)

			// Apply dual twice.
			_, err = conway.Parse("d" + seed)
			require.NoError(t, err)

			dual2, err := conway.Parse("dd" + seed)
			require.NoError(t, err)

			// Should have same topology as original (vertices/faces swapped back)
			assert.Equal(t, len(original.Vertices), len(dual2.Vertices),
				"Dual involution failed for vertices in %s", seed)
			assert.Equal(t, len(original.Faces), len(dual2.Faces),
				"Dual involution failed for faces in %s", seed)
			assert.Equal(t, len(original.Edges), len(dual2.Edges),
				"Dual involution failed for edges in %s", seed)
		})
	}
}

// TestIntegrationParserEdgeCases tests parser with edge cases.
func TestIntegrationParserEdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		notation    string
		expectError bool
	}{
		{"Empty string", "", true},
		{"Whitespace only", "   ", true},
		{"Invalid seed", "X", true},
		{"No seed", "dt", true},
		{"Invalid operation", "xC", true},
		{"Valid single seed", "T", false},
		{"Valid single operation", "dT", false},
		{"Valid complex", "dtkaC", false},
		{"Whitespace handling", "  dT  ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			poly, err := conway.Parse(tt.notation)
			if tt.expectError {
				assert.Error(t, err, "Should have failed for: %s", tt.notation)
				assert.Nil(t, poly, "Polyhedron should be nil on error")
			} else {
				assert.NoError(t, err, "Should not have failed for: %s", tt.notation)
				assert.NotNil(t, poly, "Polyhedron should not be nil on success")
				assert.True(t, poly.IsValid(), "Result should be valid")
			}
		})
	}
}

// TestIntegrationConcurrentOperations tests thread safety.
func TestIntegrationConcurrentOperations(t *testing.T) {
	t.Parallel()

	const numGoroutines = 10

	const numOperations = 100

	// Channel to collect results.
	results := make(chan *conway.Polyhedron, numGoroutines*numOperations)
	errors := make(chan error, numGoroutines*numOperations)

	// Launch multiple goroutines performing operations.
	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			for j := 0; j < numOperations; j++ {
				// Vary the operation based on goroutine ID and iteration.
				operations := []string{"T", "dT", "tC", "aO", "kI", "dD"}
				notation := operations[(goroutineID+j)%len(operations)]

				poly, err := conway.Parse(notation)
				if err != nil {
					errors <- err
					continue
				}

				// Perform some operations on the polyhedron.
				_ = poly.Clone()
				_ = poly.IsValid()
				_ = poly.Stats()
				_ = poly.EulerCharacteristic()

				results <- poly
			}
		}(i)
	}

	// Collect results.
	successCount := 0
	errorCount := 0

	for i := 0; i < numGoroutines*numOperations; i++ {
		select {
		case poly := <-results:
			assert.NotNil(t, poly)
			assert.True(t, poly.IsValid())

			successCount++
		case err := <-errors:
			assert.NoError(t, err, "Unexpected error in concurrent operations")

			errorCount++
		}
	}

	assert.Equal(t, numGoroutines*numOperations, successCount,
		"All operations should succeed")
	assert.Equal(t, 0, errorCount, "No errors should occur")
}

// TestIntegrationTopologyPreservation tests that operations preserve manifold properties.
func TestIntegrationTopologyPreservation(t *testing.T) {
	t.Parallel()

	operations := []string{"d", "a", "t", "k", "j"}
	seeds := []string{"T", "C", "O"}

	for _, seed := range seeds {
		for _, op := range operations {
			t.Run(seed+"_"+op, func(t *testing.T) {
				t.Parallel()

				notation := op + seed
				poly, err := conway.Parse(notation)
				require.NoError(t, err, "Failed to parse: %s", notation)

				// Check manifold properties.
				for _, edge := range poly.Edges {
					assert.LessOrEqual(t, len(edge.Faces), 2,
						"Edge should have at most 2 adjacent faces in %s", notation)
					assert.Greater(t, len(edge.Faces), 0,
						"Edge should have at least 1 adjacent face in %s", notation)
				}

				for _, vertex := range poly.Vertices {
					assert.GreaterOrEqual(t, vertex.Degree(), 3,
						"Vertex should have degree >= 3 in %s", notation)
				}

				for _, face := range poly.Faces {
					assert.GreaterOrEqual(t, len(face.Vertices), 3,
						"Face should have at least 3 vertices in %s", notation)
					assert.Equal(t, len(face.Vertices), len(face.Edges),
						"Face should have equal number of vertices and edges in %s", notation)
				}
			})
		}
	}
}

// TestIntegrationGeometryStats tests geometry statistics calculation.
func TestIntegrationGeometryStats(t *testing.T) {
	t.Parallel()

	poly, err := conway.Parse("tC")
	require.NoError(t, err)

	stats := poly.CalculateGeometryStats()
	require.NotNil(t, stats)

	// Basic sanity checks.
	assert.Greater(t, stats.MinEdgeLength, 0.0, "Min edge length should be positive")
	assert.Greater(t, stats.MaxEdgeLength, 0.0, "Max edge length should be positive")
	assert.GreaterOrEqual(t, stats.MaxEdgeLength, stats.MinEdgeLength,
		"Max edge length should be >= min edge length")
	assert.Greater(t, stats.AvgEdgeLength, 0.0, "Average edge length should be positive")

	assert.Greater(t, stats.MinFaceArea, 0.0, "Min face area should be positive")
	assert.Greater(t, stats.MaxFaceArea, 0.0, "Max face area should be positive")
	assert.GreaterOrEqual(t, stats.MaxFaceArea, stats.MinFaceArea,
		"Max face area should be >= min face area")
	assert.Greater(t, stats.AvgFaceArea, 0.0, "Average face area should be positive")

	// Bounding box should be reasonable.
	assert.Less(t, stats.BoundingBox.Min.X, stats.BoundingBox.Max.X,
		"Bounding box should have positive volume")
	assert.Less(t, stats.BoundingBox.Min.Y, stats.BoundingBox.Max.Y,
		"Bounding box should have positive volume")
	assert.Less(t, stats.BoundingBox.Min.Z, stats.BoundingBox.Max.Z,
		"Bounding box should have positive volume")
}

// TestIntegrationMemoryStats tests memory usage statistics.
func TestIntegrationMemoryStats(t *testing.T) {
	t.Parallel()

	poly, err := conway.Parse("kD")
	require.NoError(t, err)

	stats := poly.CalculateMemoryStats()
	require.NotNil(t, stats)

	// Basic sanity checks.
	assert.Equal(t, len(poly.Vertices), stats.VertexCount)
	assert.Equal(t, len(poly.Edges), stats.EdgeCount)
	assert.Equal(t, len(poly.Faces), stats.FaceCount)

	// Reference counts should be reasonable.
	assert.Greater(t, stats.TotalVertices, stats.VertexCount,
		"Total vertex references should exceed unique vertices")
	assert.Greater(t, stats.TotalEdges, stats.EdgeCount,
		"Total edge references should exceed unique edges")
	assert.Greater(t, stats.TotalFaces, stats.FaceCount,
		"Total face references should exceed unique faces")
}

// TestIntegrationNormalization tests polyhedron normalization.
func TestIntegrationNormalization(t *testing.T) {
	t.Parallel()

	poly, err := conway.Parse("C")
	require.NoError(t, err)

	// Store original stats.
	_ = poly.Centroid()
	_ = poly.CalculateGeometryStats()
	originalVertexCount := len(poly.Vertices)
	originalEdgeCount := len(poly.Edges)
	originalFaceCount := len(poly.Faces)

	// Normalize the polyhedron.
	poly.Normalize()

	// Check that it's centered at origin.
	newCentroid := poly.Centroid()
	assert.InDelta(t, 0.0, newCentroid.X, 1e-10, "Should be centered at origin")
	assert.InDelta(t, 0.0, newCentroid.Y, 1e-10, "Should be centered at origin")
	assert.InDelta(t, 0.0, newCentroid.Z, 1e-10, "Should be centered at origin")

	// Check that max distance is 1.
	maxDist := 0.0

	for _, v := range poly.Vertices {
		dist := v.Position.Length()
		if dist > maxDist {
			maxDist = dist
		}
	}

	assert.InDelta(t, 1.0, maxDist, 1e-10, "Max distance should be 1")

	// Topology should be preserved.
	assert.True(t, poly.IsValid(), "Normalization should preserve validity")
	assert.Equal(t, originalVertexCount, len(poly.Vertices),
		"Vertex count should be preserved")
	assert.Equal(t, originalEdgeCount, len(poly.Edges),
		"Edge count should be preserved")
	assert.Equal(t, originalFaceCount, len(poly.Faces),
		"Face count should be preserved")
}
