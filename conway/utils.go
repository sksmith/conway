package conway

import (
	"errors"
	"fmt"
	"math"
)

const (
	// lengthTolerance is the minimum length threshold for vector normalization
	lengthTolerance = 1e-12
)

// Static errors for err113 compliance
var (
	ErrInsufficientVertices = errors.New("insufficient vertices for normal calculation")
	ErrDegenerateFaceNormal = errors.New("degenerate face normal")
)

// EdgeLookup provides O(1) edge lookup by vertex pair
type EdgeLookup struct {
	edgeMap map[string]*Edge
}

// NewEdgeLookup creates a new edge lookup structure
func NewEdgeLookup() *EdgeLookup {
	return &EdgeLookup{
		edgeMap: make(map[string]*Edge),
	}
}

// makeEdgeKey creates a consistent key for vertex pairs
func makeEdgeKey(v1ID, v2ID int) string {
	if v1ID > v2ID {
		v1ID, v2ID = v2ID, v1ID // Ensure consistent ordering
	}

	return fmt.Sprintf("%d-%d", v1ID, v2ID)
}

// Add adds an edge to the lookup
func (el *EdgeLookup) Add(edge *Edge) {
	key := makeEdgeKey(edge.V1.ID, edge.V2.ID)
	el.edgeMap[key] = edge
}

// Find finds an edge between two vertices
func (el *EdgeLookup) Find(v1ID, v2ID int) *Edge {
	key := makeEdgeKey(v1ID, v2ID)
	return el.edgeMap[key]
}

// Remove removes an edge from the lookup
func (el *EdgeLookup) Remove(edge *Edge) {
	key := makeEdgeKey(edge.V1.ID, edge.V2.ID)
	delete(el.edgeMap, key)
}

// calculateFaceNormal computes face normal with proper error handling
func calculateFaceNormal(vertices []*Vertex) (Vector3, error) {
	if len(vertices) < 3 {
		return Vector3{}, fmt.Errorf("%w: %d", ErrInsufficientVertices, len(vertices))
	}

	// Use Newell's method for robust normal calculation
	normal := Vector3{}
	n := len(vertices)

	for i := 0; i < n; i++ {
		v1 := vertices[i].Position
		v2 := vertices[(i+1)%n].Position

		normal.X += (v1.Y - v2.Y) * (v1.Z + v2.Z)
		normal.Y += (v1.Z - v2.Z) * (v1.X + v2.X)
		normal.Z += (v1.X - v2.X) * (v1.Y + v2.Y)
	}

	length := normal.Length()
	if length < lengthTolerance {
		return Vector3{}, fmt.Errorf("%w (length: %e)", ErrDegenerateFaceNormal, length)
	}

	return normal.Scale(1.0 / length), nil
}

// ensureCounterClockwise ensures face vertices are in counter-clockwise order
// when viewed from outside the polyhedron
func ensureCounterClockwise(vertices []*Vertex, polyhedronCenter Vector3) []*Vertex {
	if len(vertices) < 3 {
		return vertices
	}

	normal, err := calculateFaceNormal(vertices)
	if err != nil {
		return vertices // Return original if we can't calculate normal
	}

	// Calculate face centroid
	centroid := Vector3{}
	for _, v := range vertices {
		centroid = centroid.Add(v.Position)
	}
	centroid = centroid.Scale(1.0 / float64(len(vertices)))

	// Vector from polyhedron center to face center
	outward := centroid.Sub(polyhedronCenter).Normalize()

	// If normal points inward, reverse vertex order
	if normal.Dot(outward) < 0 {
		reversed := make([]*Vertex, len(vertices))
		for i, v := range vertices {
			reversed[len(vertices)-1-i] = v
		}

		return reversed
	}

	return vertices
}

// allocateSliceWithCapacity pre-allocates a slice with known capacity
// This helps reduce memory reallocations during polyhedron construction
func allocateVertexSlice(capacity int) []*Vertex {
	return make([]*Vertex, 0, capacity)
}

// allocateEdgeSlice pre-allocates an edge slice with known capacity
func allocateEdgeSlice(capacity int) []*Edge {
	return make([]*Edge, 0, capacity)
}

// GeometryStats provides statistical information about polyhedron geometry
type GeometryStats struct {
	MinEdgeLength float64
	MaxEdgeLength float64
	AvgEdgeLength float64
	MinFaceArea   float64
	MaxFaceArea   float64
	AvgFaceArea   float64
	BoundingBox   struct {
		Min, Max Vector3
	}
}

// calculateEdgeStats computes edge statistics for the polyhedron
func calculateEdgeStats(edges map[int]*Edge) (float64, float64, float64) {
	if len(edges) == 0 {
		return 0, 0, 0
	}

	minLength := math.Inf(1)
	maxLength := 0.0
	totalLength := 0.0

	for _, edge := range edges {
		length := edge.Length()
		if length < minLength {
			minLength = length
		}
		if length > maxLength {
			maxLength = length
		}
		totalLength += length
	}

	return minLength, maxLength, totalLength / float64(len(edges))
}

// calculateFaceStats computes face area statistics for the polyhedron
func calculateFaceStats(faces map[int]*Face) (float64, float64, float64) {
	if len(faces) == 0 {
		return 0, 0, 0
	}

	minFaceArea := math.Inf(1)
	maxFaceArea := 0.0
	totalFaceArea := 0.0

	for _, face := range faces {
		area := face.Area()
		if area < minFaceArea {
			minFaceArea = area
		}
		if area > maxFaceArea {
			maxFaceArea = area
		}
		totalFaceArea += area
	}

	return minFaceArea, maxFaceArea, totalFaceArea / float64(len(faces))
}

// updateBoundingBox updates min/max bounds with a new position
func updateBoundingBox(minBound, maxBound, pos *Vector3) {
	if pos.X < minBound.X {
		minBound.X = pos.X
	}
	if pos.Y < minBound.Y {
		minBound.Y = pos.Y
	}
	if pos.Z < minBound.Z {
		minBound.Z = pos.Z
	}
	if pos.X > maxBound.X {
		maxBound.X = pos.X
	}
	if pos.Y > maxBound.Y {
		maxBound.Y = pos.Y
	}
	if pos.Z > maxBound.Z {
		maxBound.Z = pos.Z
	}
}

// calculateBoundingBox computes the bounding box for all vertices
func calculateBoundingBox(vertices map[int]*Vertex) (Vector3, Vector3) {
	if len(vertices) == 0 {
		return Vector3{}, Vector3{}
	}

	var minBound, maxBound Vector3
	first := true
	for _, vertex := range vertices {
		pos := vertex.Position
		if first {
			minBound = pos
			maxBound = pos
			first = false
		} else {
			updateBoundingBox(&minBound, &maxBound, &pos)
		}
	}

	return minBound, maxBound
}

// CalculateGeometryStats computes geometric statistics for a polyhedron
func (p *Polyhedron) CalculateGeometryStats() *GeometryStats {
	p.mu.RLock()
	defer p.mu.RUnlock()

	stats := &GeometryStats{}

	if len(p.Edges) == 0 || len(p.Faces) == 0 {
		return stats
	}

	stats.MinEdgeLength, stats.MaxEdgeLength, stats.AvgEdgeLength = calculateEdgeStats(p.Edges)
	stats.MinFaceArea, stats.MaxFaceArea, stats.AvgFaceArea = calculateFaceStats(p.Faces)
	stats.BoundingBox.Min, stats.BoundingBox.Max = calculateBoundingBox(p.Vertices)

	return stats
}

// MemoryStats provides information about memory usage
type MemoryStats struct {
	VertexCount   int
	EdgeCount     int
	FaceCount     int
	TotalVertices int
	TotalEdges    int
	TotalFaces    int
}

// CalculateMemoryStats computes memory usage statistics
func (p *Polyhedron) CalculateMemoryStats() *MemoryStats {
	p.mu.RLock()
	defer p.mu.RUnlock()

	stats := &MemoryStats{
		VertexCount: len(p.Vertices),
		EdgeCount:   len(p.Edges),
		FaceCount:   len(p.Faces),
	}

	// Count total references
	for _, face := range p.Faces {
		stats.TotalVertices += len(face.Vertices)
		stats.TotalEdges += len(face.Edges)
	}

	for _, vertex := range p.Vertices {
		stats.TotalFaces += len(vertex.Faces)
	}

	return stats
}
