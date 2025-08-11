// Package conway implements Conway polyhedron notation, a mathematical system for
// describing complex polyhedra through systematic transformations of simple seed shapes.
//
// Conway notation provides operations like dual (d), ambo (a), truncate (t), kis (k),
// and join (j) that can be applied to the five Platonic solids to generate a wide
// variety of interesting polyhedra. Operations can be chained together using a simple
// text notation, for example "dtC" creates the dual of a truncated cube.
//
// The library uses a half-edge data structure for efficient polyhedron representation
// and includes comprehensive validation, caching of computed properties, and performance
// optimizations for large polyhedra.
//
// Example usage:
//
//	// Create a truncated icosahedron (soccer ball)
//	soccerBall, err := conway.Parse("tI")
//
//	// Create operations manually
//	cube := conway.Cube()
//	dual := conway.Dual(cube)
//
//	// Validate the result
//	if err := soccerBall.ValidateComplete(); err != nil {
//	    log.Fatal(err)
//	}
package conway

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
)

const (
	// halfScale is used for calculating midpoints
	halfScale = 0.5
)

// Vector3 represents a 3D vector with X, Y, and Z components.
// It provides basic vector operations including addition, subtraction,
// scaling, dot product, cross product, normalization, and distance calculation.
type Vector3 struct {
	X, Y, Z float64
}

// Add returns the vector sum of v and other.
func (v Vector3) Add(other Vector3) Vector3 {
	return Vector3{v.X + other.X, v.Y + other.Y, v.Z + other.Z}
}

// Sub returns the vector difference of v and other.
func (v Vector3) Sub(other Vector3) Vector3 {
	return Vector3{v.X - other.X, v.Y - other.Y, v.Z - other.Z}
}

// Scale returns the vector v scaled by scalar s.
func (v Vector3) Scale(s float64) Vector3 {
	return Vector3{v.X * s, v.Y * s, v.Z * s}
}

// Dot returns the dot product of v and other.
func (v Vector3) Dot(other Vector3) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

// Cross returns the cross product of v and other.
// The result is perpendicular to both input vectors.
func (v Vector3) Cross(other Vector3) Vector3 {
	return Vector3{
		X: v.Y*other.Z - v.Z*other.Y,
		Y: v.Z*other.X - v.X*other.Z,
		Z: v.X*other.Y - v.Y*other.X,
	}
}

// Length returns the Euclidean length (magnitude) of the vector.
func (v Vector3) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

// Normalize returns a unit vector in the same direction as v.
// If v has zero length, it returns v unchanged.
func (v Vector3) Normalize() Vector3 {
	l := v.Length()
	if l == 0 {
		return v
	}

	return v.Scale(1.0 / l)
}

// Distance returns the Euclidean distance between v and other.
func (v Vector3) Distance(other Vector3) float64 {
	return v.Sub(other).Length()
}

// Vertex represents a point in 3D space with connectivity information.
// Each vertex maintains references to all edges and faces that contain it,
// enabling efficient traversal of the polyhedron's topology.
type Vertex struct {
	ID       int           // Unique identifier within the polyhedron
	Position Vector3       // 3D coordinates of the vertex
	Edges    map[int]*Edge // All edges incident to this vertex
	Faces    map[int]*Face // All faces containing this vertex
}

// NewVertex creates a new vertex with the given ID and position.
// The edge and face maps are initialized as empty.
func NewVertex(id int, pos Vector3) *Vertex {
	return &Vertex{
		ID:       id,
		Position: pos,
		Edges:    make(map[int]*Edge),
		Faces:    make(map[int]*Face),
	}
}

// Degree returns the number of edges incident to this vertex.
// In a valid polyhedron, each vertex must have degree >= 3.
func (v *Vertex) Degree() int {
	return len(v.Edges)
}

// Edge represents a connection between two vertices.
// In a manifold polyhedron, each edge is shared by exactly two faces.
// The edge maintains references to both endpoints and adjacent faces.
type Edge struct {
	ID    int           // Unique identifier within the polyhedron
	V1    *Vertex       // First endpoint vertex
	V2    *Vertex       // Second endpoint vertex
	Faces map[int]*Face // Adjacent faces (typically 2 for manifold edges)
}

// NewEdge creates a new edge connecting vertices v1 and v2.
// The faces map is initialized as empty.
func NewEdge(id int, v1, v2 *Vertex) *Edge {
	return &Edge{
		ID:    id,
		V1:    v1,
		V2:    v2,
		Faces: make(map[int]*Face),
	}
}

// Midpoint returns the point halfway between the edge's endpoints.
func (e *Edge) Midpoint() Vector3 {
	return e.V1.Position.Add(e.V2.Position).Scale(halfScale)
}

// Length returns the Euclidean distance between the edge's endpoints.
func (e *Edge) Length() float64 {
	return e.V1.Position.Distance(e.V2.Position)
}

// OtherVertex returns the vertex at the other end of the edge from v.
// Returns nil if v is not an endpoint of this edge.
func (e *Edge) OtherVertex(v *Vertex) *Vertex {
	if e.V1.ID == v.ID {
		return e.V2
	}
	if e.V2.ID == v.ID {
		return e.V1
	}

	return nil
}

// Face represents a polygonal face of a polyhedron.
// Faces are defined by an ordered list of vertices (counter-clockwise when viewed
// from outside the polyhedron) and maintain references to their bounding edges.
// Computed properties like normal, centroid, and area are cached for performance.
type Face struct {
	ID       int       // Unique identifier within the polyhedron
	Vertices []*Vertex // Ordered vertices forming the face boundary (CCW from outside)
	Edges    []*Edge   // Edges bounding the face

	// Cached computed properties
	cachedNormal   *Vector3     // Cached face normal vector
	cachedCentroid *Vector3     // Cached face centroid
	cachedArea     *float64     // Cached face area
	mu             sync.RWMutex // Mutex for thread-safe access to cached properties
}

func NewFace(id int, vertices []*Vertex) *Face {
	return &Face{
		ID:       id,
		Vertices: vertices,
		Edges:    allocateEdgeSlice(len(vertices)), // Pre-allocate with expected capacity
	}
}

func (f *Face) Degree() int {
	return len(f.Vertices)
}

func (f *Face) Centroid() Vector3 {
	f.mu.RLock()
	if f.cachedCentroid != nil {
		defer f.mu.RUnlock()
		return *f.cachedCentroid
	}
	f.mu.RUnlock()

	// Need write lock to update cache
	f.mu.Lock()
	defer f.mu.Unlock()

	// Double-check pattern: check again under write lock
	if f.cachedCentroid != nil {
		return *f.cachedCentroid
	}

	if len(f.Vertices) == 0 {
		return Vector3{}
	}

	sum := Vector3{}
	for _, v := range f.Vertices {
		sum = sum.Add(v.Position)
	}
	centroid := sum.Scale(1.0 / float64(len(f.Vertices)))
	f.cachedCentroid = &centroid

	return centroid
}

func (f *Face) Normal() Vector3 {
	f.mu.RLock()
	if f.cachedNormal != nil {
		defer f.mu.RUnlock()
		return *f.cachedNormal
	}
	f.mu.RUnlock()

	// Need write lock to update cache
	f.mu.Lock()
	defer f.mu.Unlock()

	// Double-check pattern: check again under write lock
	if f.cachedNormal != nil {
		return *f.cachedNormal
	}

	if len(f.Vertices) < 3 {
		return Vector3{}
	}

	// Use robust Newell's method for normal calculation
	normal, err := calculateFaceNormal(f.Vertices)
	if err != nil {
		// Fallback to simple cross product for degenerate cases
		v1 := f.Vertices[1].Position.Sub(f.Vertices[0].Position)
		v2 := f.Vertices[2].Position.Sub(f.Vertices[0].Position)
		normal = v1.Cross(v2).Normalize()
	}

	f.cachedNormal = &normal

	return normal
}

func (f *Face) Area() float64 {
	f.mu.RLock()
	if f.cachedArea != nil {
		defer f.mu.RUnlock()
		return *f.cachedArea
	}
	f.mu.RUnlock()

	// Need write lock to update cache
	f.mu.Lock()
	defer f.mu.Unlock()

	// Double-check pattern: check again under write lock
	if f.cachedArea != nil {
		return *f.cachedArea
	}

	if len(f.Vertices) < 3 {
		return 0
	}

	area := 0.0
	for i := 1; i < len(f.Vertices)-1; i++ {
		v1 := f.Vertices[i].Position.Sub(f.Vertices[0].Position)
		v2 := f.Vertices[i+1].Position.Sub(f.Vertices[0].Position)
		area += v1.Cross(v2).Length() * halfScale
	}
	f.cachedArea = &area

	return area
}

// Polyhedron represents a 3D polyhedron using a half-edge data structure.
// It maintains maps of vertices, edges, and faces with their connectivity information.
// The structure includes optimizations like O(1) edge lookup and caching of computed
// properties for better performance with large polyhedra.
//
// All valid polyhedra satisfy Euler's formula: V - E + F = 2, where V is the number
// of vertices, E is the number of edges, and F is the number of faces.
//
// Thread-safe for concurrent operations.
type Polyhedron struct {
	Name       string          // Descriptive name (e.g., "Cube", "tI")
	Vertices   map[int]*Vertex // All vertices indexed by ID
	Edges      map[int]*Edge   // All edges indexed by ID
	Faces      map[int]*Face   // All faces indexed by ID
	nextID     int64           // Atomic counter for next available ID
	edgeLookup *EdgeLookup     // O(1) edge lookup by vertex pair
	mu         sync.RWMutex    // Read-write mutex for thread safety

	// Cached computed properties
	cachedCentroid *Vector3 // Cached polyhedron centroid
}

// NewPolyhedron creates a new empty polyhedron with the given name.
// All internal maps and the edge lookup structure are initialized.
func NewPolyhedron(name string) *Polyhedron {
	return &Polyhedron{
		Name:       name,
		Vertices:   make(map[int]*Vertex),
		Edges:      make(map[int]*Edge),
		Faces:      make(map[int]*Face),
		nextID:     0,
		edgeLookup: NewEdgeLookup(),
	}
}

func (p *Polyhedron) getNextID() int {
	return int(atomic.AddInt64(&p.nextID, 1))
}

// AddVertex creates a new vertex at the specified position and adds it to the polyhedron.
// Returns the created vertex. Invalidates cached properties.
// Thread-safe for concurrent access.
func (p *Polyhedron) AddVertex(pos Vector3) *Vertex {
	p.mu.Lock()
	defer p.mu.Unlock()

	v := NewVertex(p.getNextID(), pos)
	p.Vertices[v.ID] = v
	p.invalidateCache() // Invalidate cached centroid when vertices change

	return v
}

// AddEdge creates an edge between vertices v1 and v2, or returns the existing edge if it already exists.
// Uses O(1) lookup to check for existing edges. Updates vertex connectivity.
// Thread-safe for concurrent access.
func (p *Polyhedron) AddEdge(v1, v2 *Vertex) *Edge {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.addEdgeUnsafe(v1, v2)
}

// addEdgeUnsafe is the internal implementation of AddEdge without locking
func (p *Polyhedron) addEdgeUnsafe(v1, v2 *Vertex) *Edge {
	// Check if edge already exists using O(1) lookup
	if existing := p.edgeLookup.Find(v1.ID, v2.ID); existing != nil {
		return existing
	}

	e := NewEdge(p.getNextID(), v1, v2)
	p.Edges[e.ID] = e
	p.edgeLookup.Add(e)
	v1.Edges[e.ID] = e
	v2.Edges[e.ID] = e

	return e
}

// AddFace creates a new face from the given ordered vertices.
// Automatically creates edges between consecutive vertices and updates all connectivity.
// Vertices should be ordered counter-clockwise when viewed from outside the polyhedron.
// Thread-safe for concurrent access.
func (p *Polyhedron) AddFace(vertices []*Vertex) *Face {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Ensure proper winding order if we have a meaningful polyhedron center
	if len(p.Vertices) > 3 {
		center := p.calculateCentroidUnsafe()
		vertices = ensureCounterClockwise(vertices, center)
	}

	f := NewFace(p.getNextID(), vertices)
	p.Faces[f.ID] = f

	for i := 0; i < len(vertices); i++ {
		v1 := vertices[i]
		v2 := vertices[(i+1)%len(vertices)]

		e := p.addEdgeUnsafe(v1, v2)
		f.Edges = append(f.Edges, e)
		e.Faces[f.ID] = f

		v1.Faces[f.ID] = f
	}

	return f
}

// invalidateCache invalidates all cached properties
func (p *Polyhedron) invalidateCache() {
	p.cachedCentroid = nil
}

// invalidateFaceCache invalidates cached properties for a face
func (f *Face) invalidateFaceCache() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.cachedNormal = nil
	f.cachedCentroid = nil
	f.cachedArea = nil
}

// RemoveVertex removes a vertex from the polyhedron and all associated edges and faces.
// Thread-safe for concurrent access.
func (p *Polyhedron) RemoveVertex(v *Vertex) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.removeVertexUnsafe(v)
}

// removeVertexUnsafe is the internal implementation without locking
func (p *Polyhedron) removeVertexUnsafe(v *Vertex) {
	// Remove all associated edges (which will also clean up faces)
	for _, e := range v.Edges {
		p.removeEdgeUnsafe(e)
	}

	// Remove all associated faces
	for _, f := range v.Faces {
		p.removeFaceUnsafe(f)
	}

	delete(p.Vertices, v.ID)
	p.invalidateCache() // Invalidate cache when vertices are removed
}

// RemoveEdge removes an edge from the polyhedron and cleans up all references.
// Thread-safe for concurrent access.
func (p *Polyhedron) RemoveEdge(e *Edge) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.removeEdgeUnsafe(e)
}

// removeEdgeUnsafe is the internal implementation without locking
func (p *Polyhedron) removeEdgeUnsafe(e *Edge) {
	// Remove from vertices
	delete(e.V1.Edges, e.ID)
	delete(e.V2.Edges, e.ID)

	// Remove from faces
	for _, f := range e.Faces {
		for i, edge := range f.Edges {
			if edge.ID == e.ID {
				f.Edges = append(f.Edges[:i], f.Edges[i+1:]...)
				break
			}
		}
		// Also remove face reference from edge
		delete(e.Faces, f.ID)
	}

	// Remove from lookup table - critical for preventing memory leaks
	p.edgeLookup.Remove(e)
	delete(p.Edges, e.ID)
}

// RemoveFace removes a face from the polyhedron and cleans up all references.
// Thread-safe for concurrent access.
func (p *Polyhedron) RemoveFace(f *Face) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.removeFaceUnsafe(f)
}

// removeFaceUnsafe is the internal implementation without locking
func (p *Polyhedron) removeFaceUnsafe(f *Face) {
	// Remove face references from vertices
	for _, v := range f.Vertices {
		delete(v.Faces, f.ID)
	}

	// Remove face references from edges
	for _, e := range f.Edges {
		delete(e.Faces, f.ID)
	}

	delete(p.Faces, f.ID)
}

// EulerCharacteristic returns V - E + F for the polyhedron.
// For valid closed polyhedra, this should always equal 2.
func (p *Polyhedron) EulerCharacteristic() int {
	return len(p.Vertices) - len(p.Edges) + len(p.Faces)
}

// IsValid performs basic validation checks on the polyhedron.
// Returns true if the polyhedron satisfies basic topological requirements:
// - Euler characteristic equals 2
// - All vertices have degree >= 3
// - All edges have 1 or 2 adjacent faces
// - All faces have >= 3 vertices
// For more comprehensive validation, use ValidateComplete().
func (p *Polyhedron) IsValid() bool {
	if p.EulerCharacteristic() != 2 {
		return false
	}

	for _, v := range p.Vertices {
		if v.Degree() < 3 {
			return false
		}
	}

	for _, e := range p.Edges {
		if len(e.Faces) > 2 || len(e.Faces) == 0 {
			return false
		}
	}

	for _, f := range p.Faces {
		if f.Degree() < 3 {
			return false
		}
	}

	return true
}

// Clone creates a deep copy of the polyhedron.
// All vertices, edges, and faces are recreated with new IDs,
// but the geometric and topological structure is preserved.
// Thread-safe for concurrent access.
func (p *Polyhedron) Clone() *Polyhedron {
	p.mu.RLock()
	defer p.mu.RUnlock()

	newP := NewPolyhedron(p.Name)

	// Pre-allocate vertex map with known size
	vertexMap := make(map[int]*Vertex, len(p.Vertices))
	for _, v := range p.Vertices {
		newV := newP.AddVertex(v.Position)
		vertexMap[v.ID] = newV
	}

	for _, f := range p.Faces {
		// Pre-allocate slice with exact size needed
		newVertices := make([]*Vertex, len(f.Vertices))
		for i, v := range f.Vertices {
			newVertices[i] = vertexMap[v.ID]
		}
		newP.AddFace(newVertices)
	}

	return newP
}

func (p *Polyhedron) Centroid() Vector3 {
	p.mu.RLock()
	if p.cachedCentroid != nil {
		defer p.mu.RUnlock()
		return *p.cachedCentroid
	}
	p.mu.RUnlock()

	// Need write lock to update cache
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.calculateCentroidUnsafe()
}

// calculateCentroidUnsafe is the internal implementation without locking
func (p *Polyhedron) calculateCentroidUnsafe() Vector3 {
	// Double-check pattern: check again under write lock
	if p.cachedCentroid != nil {
		return *p.cachedCentroid
	}

	if len(p.Vertices) == 0 {
		return Vector3{}
	}

	sum := Vector3{}
	for _, v := range p.Vertices {
		sum = sum.Add(v.Position)
	}
	centroid := sum.Scale(1.0 / float64(len(p.Vertices)))
	p.cachedCentroid = &centroid

	return centroid
}

// Normalize centers the polyhedron at the origin and scales it so that
// the furthest vertex is at unit distance from the center.
// This is useful for consistent visualization and comparison of polyhedra.
// All cached properties are invalidated after normalization.
func (p *Polyhedron) Normalize() {
	centroid := p.Centroid()

	for _, v := range p.Vertices {
		v.Position = v.Position.Sub(centroid)
	}

	maxDist := 0.0
	for _, v := range p.Vertices {
		dist := v.Position.Length()
		if dist > maxDist {
			maxDist = dist
		}
	}

	if maxDist > 0 {
		scale := 1.0 / maxDist
		for _, v := range p.Vertices {
			v.Position = v.Position.Scale(scale)
		}
	}

	// Invalidate all cached properties after normalization
	p.invalidateCache()
	for _, f := range p.Faces {
		f.invalidateFaceCache()
	}
}

// Stats returns a string with basic polyhedron statistics including
// name, vertex count (V), edge count (E), face count (F), and Euler characteristic (χ).
func (p *Polyhedron) Stats() string {
	return fmt.Sprintf("%s: V=%d, E=%d, F=%d, χ=%d",
		p.Name, len(p.Vertices), len(p.Edges), len(p.Faces), p.EulerCharacteristic())
}
