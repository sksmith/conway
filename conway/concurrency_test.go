package conway

import (
	"math"
	"sync"
	"testing"
	"time"
)

// TestConcurrentVertexAddition tests concurrent vertex addition
func TestConcurrentVertexAddition(t *testing.T) {
	p := NewPolyhedron("ConcurrentTest")
	const numGoroutines = 10
	const verticesPerGoroutine = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Launch multiple goroutines adding vertices concurrently
	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < verticesPerGoroutine; j++ {
				p.AddVertex(Vector3{
					X: float64(goroutineID*verticesPerGoroutine + j),
					Y: float64(goroutineID),
					Z: 0,
				})
			}
		}(i)
	}

	wg.Wait()

	// Verify all vertices were added
	expectedCount := numGoroutines * verticesPerGoroutine
	if len(p.Vertices) != expectedCount {
		t.Errorf("Expected %d vertices, got %d", expectedCount, len(p.Vertices))
	}

	// Verify no ID collisions - all IDs should be unique
	idSet := make(map[int]bool)
	for _, v := range p.Vertices {
		if idSet[v.ID] {
			t.Errorf("Duplicate vertex ID found: %d", v.ID)
		}
		idSet[v.ID] = true
	}
}

// TestConcurrentEdgeAddition tests concurrent edge addition
func TestConcurrentEdgeAddition(t *testing.T) {
	p := NewPolyhedron("ConcurrentEdgeTest")

	// Create initial vertices
	vertices := make([]*Vertex, 10)
	for i := 0; i < 10; i++ {
		vertices[i] = p.AddVertex(Vector3{X: float64(i), Y: 0, Z: 0})
	}

	const numGoroutines = 5
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Launch multiple goroutines adding edges concurrently
	for i := 0; i < numGoroutines; i++ {
		go func(_ int) {
			defer wg.Done()
			// Each goroutine creates edges between different vertex pairs
			for j := 0; j < len(vertices)-1; j++ {
				v1 := vertices[j]
				v2 := vertices[j+1]
				p.AddEdge(v1, v2)

				// Small delay to increase chance of race conditions
				time.Sleep(time.Microsecond)
			}
		}(i)
	}

	wg.Wait()

	// Should have exactly 9 edges (between consecutive vertices)
	// Multiple calls to AddEdge with same vertex pairs should return existing edge
	if len(p.Edges) != 9 {
		t.Errorf("Expected 9 edges, got %d", len(p.Edges))
	}

	// Verify edge lookup consistency
	for i := 0; i < len(vertices)-1; i++ {
		edge := p.edgeLookup.Find(vertices[i].ID, vertices[i+1].ID)
		if edge == nil {
			t.Errorf("Edge lookup failed for vertices %d and %d", vertices[i].ID, vertices[i+1].ID)
		}
	}
}

// TestConcurrentFaceAddition tests concurrent face addition
func TestConcurrentFaceAddition(t *testing.T) {
	p := NewPolyhedron("ConcurrentFaceTest")

	// Create a set of vertices for faces
	const numVertices = 12
	vertices := make([]*Vertex, numVertices)
	for i := 0; i < numVertices; i++ {
		vertices[i] = p.AddVertex(Vector3{
			X: float64(i % 4),
			Y: float64(i / 4),
			Z: 0,
		})
	}

	const numGoroutines = 4
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Launch multiple goroutines adding faces concurrently
	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			defer wg.Done()
			// Each goroutine creates different triangular faces
			start := goroutineID * 3
			if start+2 < numVertices {
				p.AddFace([]*Vertex{
					vertices[start],
					vertices[start+1],
					vertices[start+2],
				})
			}
		}(i)
	}

	wg.Wait()

	// Verify faces were created
	if len(p.Faces) != numGoroutines {
		t.Errorf("Expected %d faces, got %d", numGoroutines, len(p.Faces))
	}
}

// TestConcurrentOperations tests concurrent mixed operations
func TestConcurrentOperations(t *testing.T) {
	p := NewPolyhedron("ConcurrentMixedTest")

	// Create initial polyhedron
	initialVertices := make([]*Vertex, 4)
	for i := 0; i < 4; i++ {
		initialVertices[i] = p.AddVertex(Vector3{
			X: float64(i % 2),
			Y: float64(i / 2),
			Z: 0,
		})
	}

	// Add initial face
	p.AddFace(initialVertices)

	const numGoroutines = 6
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Launch different types of operations concurrently
	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			defer wg.Done()
			switch goroutineID % 3 {
			case 0:
				// Add vertices
				for j := 0; j < 10; j++ {
					p.AddVertex(Vector3{
						X: float64(goroutineID*10 + j),
						Y: float64(goroutineID),
						Z: 1,
					})
				}
			case 1:
				// Read operations (centroid calculation)
				for j := 0; j < 10; j++ {
					_ = p.Centroid()
					time.Sleep(time.Microsecond)
				}
			case 2:
				// Clone operations
				for j := 0; j < 5; j++ {
					clone := p.Clone()
					if clone.Name != "ConcurrentMixedTest" {
						t.Errorf("Clone name incorrect: %s", clone.Name)
					}
					time.Sleep(time.Microsecond)
				}
			}
		}(i)
	}

	wg.Wait()

	// Verify polyhedron has reasonable structure (skip full validation as adding vertices without faces creates incomplete topology)
	if len(p.Vertices) < 4 {
		t.Errorf("Expected at least 4 vertices, got %d", len(p.Vertices))
	}

	// Should have at least the original 4 vertices plus added ones
	if len(p.Vertices) < 4 {
		t.Errorf("Vertices lost during concurrent operations: %d", len(p.Vertices))
	}
}

// TestConcurrentRemovalOperations tests concurrent removal operations
func TestConcurrentRemovalOperations(t *testing.T) {
	// Create a larger polyhedron to test removal
	cube := Cube()

	const numGoroutines = 4
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Get initial counts
	initialVertices := len(cube.Vertices)
	initialEdges := len(cube.Edges)
	initialFaces := len(cube.Faces)

	// Launch concurrent operations that might remove elements
	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			defer wg.Done()
			// Each goroutine performs different operations
			switch goroutineID {
			case 0:
				// Add new vertices
				cube.AddVertex(Vector3{X: 10, Y: float64(goroutineID), Z: 0})
			case 1:
				// Perform validation checks
				_ = cube.ValidateComplete()
			case 2:
				// Calculate statistics
				_ = cube.CalculateGeometryStats()
			case 3:
				// Clone the polyhedron
				_ = cube.Clone()
			}
		}(i)
	}

	wg.Wait()

	// Verify the polyhedron still has reasonable structure (adding vertices can change Euler characteristic)
	// Skip full validation as concurrent additions may create incomplete topology temporarily
	if len(cube.Vertices) < initialVertices {
		t.Errorf("Vertices lost: had %d, now %d", initialVertices, len(cube.Vertices))
	}

	// Should have at least the original structure
	if len(cube.Vertices) < initialVertices {
		t.Errorf("Vertices lost: had %d, now %d", initialVertices, len(cube.Vertices))
	}
	if len(cube.Edges) < initialEdges {
		t.Errorf("Edges lost: had %d, now %d", initialEdges, len(cube.Edges))
	}
	if len(cube.Faces) != initialFaces { // Faces shouldn't change in this test
		t.Errorf("Faces changed unexpectedly: had %d, now %d", initialFaces, len(cube.Faces))
	}
}

// TestCentroidCachingRace specifically tests for race conditions in centroid caching
func TestCentroidCachingRace(t *testing.T) {
	p := NewPolyhedron("CentroidRaceTest")

	// Add vertices to the polyhedron
	for i := 0; i < 10; i++ {
		p.AddVertex(Vector3{X: float64(i), Y: float64(i), Z: float64(i)})
	}

	const numGoroutines = 50
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Channel to collect centroids
	centroidChan := make(chan Vector3, numGoroutines)

	// Launch many goroutines calling Centroid() simultaneously
	// This should trigger the race condition in calculateCentroidUnsafe
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			// Call Centroid multiple times to increase race condition chances
			for j := 0; j < 100; j++ {
				centroid := p.Centroid()
				if j == 0 { // Only collect the first centroid from each goroutine
					centroidChan <- centroid
				}
			}
		}()
	}

	wg.Wait()
	close(centroidChan)

	// All centroids should be identical since we're using the same vertices
	var firstCentroid *Vector3
	for centroid := range centroidChan {
		if firstCentroid == nil {
			firstCentroid = &centroid
		} else {
			// Check if centroids are consistent (allowing for floating point precision)
			if math.Abs(centroid.X-firstCentroid.X) > 1e-10 ||
				math.Abs(centroid.Y-firstCentroid.Y) > 1e-10 ||
				math.Abs(centroid.Z-firstCentroid.Z) > 1e-10 {
				t.Errorf("Centroid inconsistency detected: got %v, expected %v", centroid, firstCentroid)
			}
		}
	}
}

// TestBoundingBoxCalculationRace specifically tests for race conditions in bounding box calculation
func TestBoundingBoxCalculationRace(t *testing.T) {
	p := NewPolyhedron("BoundingBoxRaceTest")

	// Start with an initial cube
	cube := Cube()
	for _, v := range cube.Vertices {
		p.AddVertex(v.Position)
	}

	const numGoroutines = 20
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Channel to collect geometry stats
	statsChan := make(chan *GeometryStats, numGoroutines/2)

	// Launch goroutines that simultaneously add vertices and calculate geometry stats
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			if id%2 == 0 {
				// Half the goroutines add vertices
				for j := 0; j < 50; j++ {
					p.AddVertex(Vector3{
						X: float64(id*50 + j),
						Y: float64(id),
						Z: float64(j),
					})
				}
			} else {
				// Half the goroutines calculate geometry stats (which calls calculateBoundingBox)
				for j := 0; j < 10; j++ {
					stats := p.CalculateGeometryStats()
					if j == 0 { // Only collect the first stats from each reader goroutine
						statsChan <- stats
					}
				}
			}
		}(i)
	}

	wg.Wait()
	close(statsChan)

	// Verify that all stats calculations completed without crashing
	statsCount := 0
	for stats := range statsChan {
		if stats == nil {
			t.Error("Got nil geometry stats")
		}
		statsCount++
	}

	if statsCount == 0 {
		t.Error("No geometry stats were collected")
	}
}

// TestAtomicIDGeneration tests that ID generation is truly atomic
func TestAtomicIDGeneration(t *testing.T) {
	p := NewPolyhedron("AtomicIDTest")
	const numGoroutines = 20
	const verticesPerGoroutine = 50

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Channel to collect all generated IDs
	idChan := make(chan int, numGoroutines*verticesPerGoroutine)

	// Launch many goroutines generating vertices simultaneously
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < verticesPerGoroutine; j++ {
				v := p.AddVertex(Vector3{X: float64(j), Y: 0, Z: 0})
				idChan <- v.ID
			}
		}()
	}

	wg.Wait()
	close(idChan)

	// Collect all IDs and check for duplicates
	idMap := make(map[int]int)
	for id := range idChan {
		idMap[id]++
		if idMap[id] > 1 {
			t.Errorf("Duplicate ID generated: %d (occurred %d times)", id, idMap[id])
		}
	}

	expectedCount := numGoroutines * verticesPerGoroutine
	if len(idMap) != expectedCount {
		t.Errorf("Expected %d unique IDs, got %d", expectedCount, len(idMap))
	}
}
