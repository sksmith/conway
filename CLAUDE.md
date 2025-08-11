# Conway Polyhedron Notation Go Library

## Overview
This Go library implements Conway polyhedron notation, a system for describing polyhedra through a series of operations applied to seed shapes.

## Architecture

### Core Data Structures

#### Vertex
- `ID`: Unique identifier
- `Position`: 3D coordinates (x, y, z)
- `Edges`: References to connected edges
- `Faces`: References to adjacent faces

#### Edge
- `ID`: Unique identifier
- `V1, V2`: References to endpoint vertices
- `Faces`: References to adjacent faces (max 2)

#### Face
- `ID`: Unique identifier
- `Vertices`: Ordered list of vertices forming the face
- `Edges`: List of edges bounding the face
- `Normal`: Face normal vector

#### Polyhedron
- `Vertices`: Map of all vertices
- `Edges`: Map of all edges
- `Faces`: Map of all faces
- `Name`: Descriptive name
- Methods for topology validation and manipulation

### Seed Polyhedra
The library provides the five Platonic solids as seed shapes:
- **T (Tetrahedron)**: 4 vertices, 6 edges, 4 faces
- **C (Cube)**: 8 vertices, 12 edges, 6 faces
- **O (Octahedron)**: 6 vertices, 12 edges, 8 faces
- **D (Dodecahedron)**: 20 vertices, 30 edges, 12 faces
- **I (Icosahedron)**: 12 vertices, 30 edges, 20 faces

### Conway Operations

#### Basic Operations
1. **Dual (d)**: Exchanges vertices and faces
   - Each vertex becomes a face
   - Each face becomes a vertex
   - Edges are perpendicular bisected

2. **Ambo (a)**: Rectification
   - Creates new vertices at edge midpoints
   - Original vertices and faces are removed
   - Creates degree-4 vertices

3. **Truncate (t)**: Vertex truncation
   - Cuts off each vertex
   - Creates small faces at vertex locations
   - Original faces become larger

4. **Kis (k)**: Stellation
   - Raises a pyramid on each face
   - Adds new vertex at face center
   - Splits each face into triangles

5. **Join (j)**: Join operation
   - Dual of ambo: j = da
   - Creates quadrilateral faces

#### Compound Operations
6. **Ortho (o)**: jj (double join)
7. **Expand (e)**: aa (double ambo)
8. **Gyro (g)**: Pentagonal rotation
9. **Snub (s)**: Chiral snub operation

### Parser
The parser interprets Conway notation strings:
- Reads operations from left to right
- Applies them in reverse order (right to left)
- Example: "dtC" = dual(truncate(Cube))

### Implementation Details

#### Topology Management
- Half-edge data structure for efficient traversal
- Winged-edge representation for edge operations
- Face ordering maintains consistency (CCW from outside)

#### Operation Implementation
Each operation is implemented as:
```go
type Operation interface {
    Apply(p *Polyhedron) *Polyhedron
    Symbol() string
    Name() string
}
```

#### Validation
- Euler characteristic check: V - E + F = 2
- Manifold validation (2 faces per edge max)
- Vertex degree validation (min 3)
- Face planarity check for non-triangular faces

### Testing Strategy

1. **Unit Tests**
   - Each operation tested individually
   - Verify topological properties
   - Check Euler formula preservation

2. **Integration Tests**
   - Compound operations
   - Parser validation
   - Known polyhedra verification

3. **Property-Based Tests**
   - Dual involution: dd(P) = P
   - Commutative properties where applicable
   - Topology preservation

### Performance Considerations
- Lazy evaluation for complex operations
- Caching of computed properties (normals, centroids)
- Efficient memory management for large polyhedra

## Usage Examples

```go
// Create a truncated icosahedron (soccer ball)
p := conway.Parse("tI")

// Create a snub cube
p := conway.Parse("sC")

// Complex example
p := conway.Parse("dgtkC")  // dual(gyro(truncate(kis(Cube))))
```

## File Structure
```
conway/
├── polyhedron.go      # Core data structures
├── seeds.go           # Platonic solid generators
├── operations.go      # Operation implementations
├── dual.go           # Dual operation
├── ambo.go           # Ambo operation
├── truncate.go       # Truncate operation
├── kis.go            # Kis operation
├── join.go           # Join operation
├── compound.go       # Compound operations (ortho, expand, etc.)
├── parser.go         # Notation parser
├── validation.go     # Topology validation
├── utils.go          # Helper functions
└── *_test.go         # Test files for each module
```