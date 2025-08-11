package conway

import (
	"fmt"
	"math"
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
		return Vector3{}, fmt.Errorf("insufficient vertices for normal calculation: %d", len(vertices))
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
	if length < 1e-12 {
		return Vector3{}, fmt.Errorf("degenerate face normal (length: %e)", length)
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

// allocateFaceSlice pre-allocates a face slice with known capacity
func allocateFaceSlice(capacity int) []*Face {
	return make([]*Face, 0, capacity)
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

// CalculateGeometryStats computes geometric statistics for a polyhedron
func (p *Polyhedron) CalculateGeometryStats() *GeometryStats {
	stats := &GeometryStats{}

	if len(p.Edges) == 0 || len(p.Faces) == 0 {
		return stats
	}

	// Edge statistics
	minEdgeLen := math.Inf(1)
	maxEdgeLen := 0.0
	totalEdgeLen := 0.0

	for _, edge := range p.Edges {
		length := edge.Length()
		if length < minEdgeLen {
			minEdgeLen = length
		}
		if length > maxEdgeLen {
			maxEdgeLen = length
		}
		totalEdgeLen += length
	}

	stats.MinEdgeLength = minEdgeLen
	stats.MaxEdgeLength = maxEdgeLen
	stats.AvgEdgeLength = totalEdgeLen / float64(len(p.Edges))

	// Face statistics
	minFaceArea := math.Inf(1)
	maxFaceArea := 0.0
	totalFaceArea := 0.0

	for _, face := range p.Faces {
		area := face.Area()
		if area < minFaceArea {
			minFaceArea = area
		}
		if area > maxFaceArea {
			maxFaceArea = area
		}
		totalFaceArea += area
	}

	stats.MinFaceArea = minFaceArea
	stats.MaxFaceArea = maxFaceArea
	stats.AvgFaceArea = totalFaceArea / float64(len(p.Faces))

	// Bounding box
	if len(p.Vertices) > 0 {
		first := true
		for _, vertex := range p.Vertices {
			pos := vertex.Position
			if first {
				stats.BoundingBox.Min = pos
				stats.BoundingBox.Max = pos
				first = false
			} else {
				if pos.X < stats.BoundingBox.Min.X {
					stats.BoundingBox.Min.X = pos.X
				}
				if pos.Y < stats.BoundingBox.Min.Y {
					stats.BoundingBox.Min.Y = pos.Y
				}
				if pos.Z < stats.BoundingBox.Min.Z {
					stats.BoundingBox.Min.Z = pos.Z
				}
				if pos.X > stats.BoundingBox.Max.X {
					stats.BoundingBox.Max.X = pos.X
				}
				if pos.Y > stats.BoundingBox.Max.Y {
					stats.BoundingBox.Max.Y = pos.Y
				}
				if pos.Z > stats.BoundingBox.Max.Z {
					stats.BoundingBox.Max.Z = pos.Z
				}
			}
		}
	}

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
