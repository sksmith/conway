// Package conway implements Conway polyhedron notation, a mathematical system for
// describing complex polyhedra through systematic transformations of simple seed shapes.
//
// Conway notation provides operations like dual (d), ambo (a), truncate (t), kis (k),
// and join (j) that can be applied to the five Platonic solids to generate a wide
// variety of interesting polyhedra. Operations can be chained together using a simple
// text notation.
//
// # Basic Usage
//
// The simplest way to use this library is through the Parse function:
//
//	// Create a truncated icosahedron (soccer ball)
//	soccerBall, err := conway.Parse("tI")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Printf("Soccer ball: %s\n", soccerBall.Stats())
//
// # Seed Shapes
//
// The library provides five Platonic solids as seed shapes:
//   - T: Tetrahedron (4 vertices, 6 edges, 4 faces)
//   - C: Cube (8 vertices, 12 edges, 6 faces)
//   - O: Octahedron (6 vertices, 12 edges, 8 faces)
//   - D: Dodecahedron (20 vertices, 30 edges, 12 faces)
//   - I: Icosahedron (12 vertices, 30 edges, 20 faces)
//
// # Operations
//
// Basic operations include:
//   - d: Dual - exchanges vertices and faces
//   - a: Ambo - rectification, creates new vertices at edge midpoints
//   - t: Truncate - vertex truncation, cuts off each vertex
//   - k: Kis - stellation, raises a pyramid on each face
//   - j: Join - dual of ambo
//
// Compound operations include:
//   - o: Ortho - double join (jj)
//   - e: Expand - double ambo (aa)
//   - g: Gyro - pentagonal rotation
//   - s: Snub - chiral snub operation
//
// # Advanced Usage
//
// For more control, operations can be applied manually:
//
//	cube := conway.Cube()
//	dual := conway.NewDual().Apply(cube)
//	truncated := conway.NewTruncate().Apply(dual)
//
// # Validation
//
// All generated polyhedra can be validated:
//
//	if err := polyhedron.ValidateComplete(); err != nil {
//		log.Printf("Invalid polyhedron: %v", err)
//	}
//
// The library ensures all operations preserve the topological validity
// of polyhedra, maintaining Euler's formula (V - E + F = 2) and other
// geometric invariants.
//
// # Thread Safety
//
// All operations are thread-safe and can be used concurrently.
// The Polyhedron type includes internal synchronization for safe
// concurrent access.
//
// # Performance
//
// The library is optimized for performance with large polyhedra:
//   - O(1) edge lookup using hash tables
//   - Lazy evaluation of computed properties
//   - Memory-efficient half-edge data structure
//   - Caching of expensive calculations
package conway
