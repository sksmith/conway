package conway

import (
	"fmt"
	"math"
)

// ValidationError represents an error in polyhedron validation
type ValidationError struct {
	Type    string
	Message string
}

func (ve ValidationError) Error() string {
	return fmt.Sprintf("%s validation error: %s", ve.Type, ve.Message)
}

// ValidateManifold checks if the polyhedron is a valid 2-manifold
// A valid 2-manifold requires:
// - Each edge connects exactly 2 faces (except boundary edges which have 1)
// - Each vertex's incident faces form a connected cycle
// Thread-safe for concurrent access.
func (p *Polyhedron) ValidateManifold() error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Check edge manifold property
	for _, edge := range p.Edges {
		faceCount := len(edge.Faces)
		if faceCount != 2 {
			if faceCount == 1 {
				// Boundary edge - this might be valid for open meshes
				// For closed polyhedra, all edges should have exactly 2 faces
				continue
			}
			return ValidationError{
				Type:    "Manifold",
				Message: fmt.Sprintf("Edge %d has %d faces (expected 2)", edge.ID, faceCount),
			}
		}
	}

	// Check vertex manifold property
	for _, vertex := range p.Vertices {
		if err := p.validateVertexManifold(vertex); err != nil {
			return err
		}
	}

	return nil
}

// validateVertexManifold checks if faces around a vertex form a proper manifold
func (p *Polyhedron) validateVertexManifold(vertex *Vertex) error {
	if len(vertex.Faces) < 3 {
		return ValidationError{
			Type:    "Manifold",
			Message: fmt.Sprintf("Vertex %d has only %d faces (minimum 3)", vertex.ID, len(vertex.Faces)),
		}
	}

	// Check that faces around vertex form a connected cycle
	// This is a complex check that requires face ordering
	orderedFaces := orderFacesAroundVertex(vertex)
	if len(orderedFaces) != len(vertex.Faces) {
		return ValidationError{
			Type:    "Manifold",
			Message: fmt.Sprintf("Vertex %d faces don't form a connected cycle", vertex.ID),
		}
	}

	return nil
}

// ValidatePlanarity checks if non-triangular faces are planar
// Thread-safe for concurrent access.
func (p *Polyhedron) ValidatePlanarity() error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	const tolerance = 1e-10

	for _, face := range p.Faces {
		if len(face.Vertices) <= 3 {
			// Triangular faces are always planar
			continue
		}

		if err := p.validateFacePlanarity(face, tolerance); err != nil {
			return err
		}
	}

	return nil
}

// validateFacePlanarity checks if a specific face is planar within tolerance
func (p *Polyhedron) validateFacePlanarity(face *Face, tolerance float64) error {
	if len(face.Vertices) < 4 {
		return nil // Triangular faces are always planar
	}

	// Calculate the plane from the first three vertices
	v0 := face.Vertices[0].Position
	v1 := face.Vertices[1].Position
	v2 := face.Vertices[2].Position

	// Calculate normal vector
	edge1 := v1.Sub(v0)
	edge2 := v2.Sub(v0)
	normal := edge1.Cross(edge2).Normalize()

	// Check if all other vertices lie in the same plane
	for i := 3; i < len(face.Vertices); i++ {
		vi := face.Vertices[i].Position
		// Calculate distance from point to plane
		dist := math.Abs(normal.Dot(vi.Sub(v0)))
		if dist > tolerance {
			return ValidationError{
				Type:    "Planarity",
				Message: fmt.Sprintf("Face %d vertex %d is %.2e units from face plane (tolerance: %.2e)", face.ID, i, dist, tolerance),
			}
		}
	}

	return nil
}

// ValidateWinding checks if all faces have consistent winding order (CCW from outside)
// Thread-safe for concurrent access.
func (p *Polyhedron) ValidateWinding() error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	centroid := p.calculateCentroidUnsafe()

	for _, face := range p.Faces {
		if err := p.validateFaceWinding(face, centroid); err != nil {
			return err
		}
	}

	return nil
}

// validateFaceWinding checks if a face has correct winding order
func (p *Polyhedron) validateFaceWinding(face *Face, polyhedronCentroid Vector3) error {
	if len(face.Vertices) < 3 {
		return ValidationError{
			Type:    "Winding",
			Message: fmt.Sprintf("Face %d has insufficient vertices for winding check", face.ID),
		}
	}

	faceNormal := face.Normal()
	faceCentroid := face.Centroid()

	// Vector from polyhedron center to face center
	outwardVector := faceCentroid.Sub(polyhedronCentroid).Normalize()

	// If face normal points outward, winding should be counter-clockwise when viewed from outside
	dotProduct := faceNormal.Dot(outwardVector)

	if dotProduct < -0.1 { // Allow some tolerance
		return ValidationError{
			Type:    "Winding",
			Message: fmt.Sprintf("Face %d has incorrect winding order (normal points inward)", face.ID),
		}
	}

	return nil
}

// ValidateTopology performs comprehensive topology validation
// Thread-safe for concurrent access.
func (p *Polyhedron) ValidateTopology() error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Check Euler characteristic
	euler := len(p.Vertices) - len(p.Edges) + len(p.Faces) // Calculate inline to avoid deadlock
	if euler != 2 {
		return ValidationError{
			Type:    "Topology",
			Message: fmt.Sprintf("Invalid Euler characteristic: %d (expected 2)", euler),
		}
	}

	// Check minimum vertex degree
	for _, vertex := range p.Vertices {
		if vertex.Degree() < 3 {
			return ValidationError{
				Type:    "Topology",
				Message: fmt.Sprintf("Vertex %d has degree %d (minimum 3)", vertex.ID, vertex.Degree()),
			}
		}
	}

	// Check minimum face degree
	for _, face := range p.Faces {
		if face.Degree() < 3 {
			return ValidationError{
				Type:    "Topology",
				Message: fmt.Sprintf("Face %d has degree %d (minimum 3)", face.ID, face.Degree()),
			}
		}
	}

	// Check edge-face connectivity
	for _, edge := range p.Edges {
		faceCount := len(edge.Faces)
		if faceCount == 0 || faceCount > 2 {
			return ValidationError{
				Type:    "Topology",
				Message: fmt.Sprintf("Edge %d has %d faces (expected 1 or 2)", edge.ID, faceCount),
			}
		}
	}

	return nil
}

// ValidateGeometry performs geometric validation checks
// Thread-safe for concurrent access.
func (p *Polyhedron) ValidateGeometry() error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Check for degenerate edges (zero length)
	const minEdgeLength = 1e-12
	for _, edge := range p.Edges {
		length := edge.Length()
		if length < minEdgeLength {
			return ValidationError{
				Type:    "Geometry",
				Message: fmt.Sprintf("Edge %d has degenerate length: %e", edge.ID, length),
			}
		}
	}

	// Check for degenerate faces (zero area)
	const minFaceArea = 1e-12
	for _, face := range p.Faces {
		area := face.Area()
		if area < minFaceArea {
			return ValidationError{
				Type:    "Geometry",
				Message: fmt.Sprintf("Face %d has degenerate area: %e", face.ID, area),
			}
		}
	}

	return nil
}

// ValidateComplete performs all validation checks
func (p *Polyhedron) ValidateComplete() error {
	if err := p.ValidateTopology(); err != nil {
		return err
	}

	if err := p.ValidateManifold(); err != nil {
		return err
	}

	if err := p.ValidatePlanarity(); err != nil {
		return err
	}

	if err := p.ValidateWinding(); err != nil {
		return err
	}

	if err := p.ValidateGeometry(); err != nil {
		return err
	}

	return nil
}
