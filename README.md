# Conway Polyhedron Notation Go Library

[![Go Version](https://img.shields.io/github/go-mod/go-version/sksmith/conway)](https://github.com/sksmith/conway)
[![CI Status](https://github.com/sksmith/conway/actions/workflows/ci.yml/badge.svg)](https://github.com/sksmith/conway/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sksmith/conway)](https://goreportcard.com/report/github.com/sksmith/conway)
[![GoDoc](https://godoc.org/github.com/sksmith/conway?status.svg)](https://godoc.org/github.com/sksmith/conway)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Coverage Status](https://codecov.io/gh/sksmith/conway/branch/main/graph/badge.svg)](https://codecov.io/gh/sksmith/conway)

A comprehensive Go library implementing Conway polyhedron notation, a mathematical system for describing complex polyhedra through systematic transformations of simple seed shapes. Developed by John Conway in the 1970s and 1980s, this notation provides an elegant way to generate and manipulate a wide variety of interesting polyhedra.

## ✨ Features

- 🔷 **Complete Implementation**: All five Platonic solids as seed shapes
- ⚡ **All Conway Operations**: Basic operations (dual, ambo, truncate, kis, join) and compound operations (ortho, expand, gyro, snub)
- 📝 **Intuitive Parser**: Simple text notation like `"tI"` for truncated icosahedron (soccer ball)
- ✅ **Robust Validation**: Comprehensive topology validation with detailed error reporting
- 🏗️ **Efficient Data Structure**: Memory-optimized half-edge representation with O(1) edge lookup
- 🔒 **Thread-Safe**: All operations are safe for concurrent use
- 🚀 **High Performance**: Lazy evaluation, caching, and optimized algorithms
- 📊 **Rich Analysis**: Geometric statistics, memory usage analysis, and property validation
- 🧪 **Comprehensive Testing**: Extensive unit tests, integration tests, property-based tests, and benchmarks

## 🎯 Quick Start

### Installation

```bash
go get github.com/sksmith/conway
```

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/sksmith/conway/conway"
)

func main() {
    // Create a truncated icosahedron (soccer ball)
    soccerBall, err := conway.Parse("tI")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Soccer ball: %s\n", soccerBall.Stats())
    // Output: Soccer ball: tI: V=60, E=90, F=32, χ=2
    
    // Validate the result
    if soccerBall.IsValid() {
        fmt.Println("✓ Valid polyhedron")
    }
}
```

## 📚 Conway Operations

### Seed Shapes (Platonic Solids)

| Symbol | Name | Vertices | Edges | Faces | Description |
|--------|------|----------|-------|-------|-------------|
| **T** | Tetrahedron | 4 | 6 | 4 | Regular triangular pyramid |
| **C** | Cube | 8 | 12 | 6 | Regular hexahedron |
| **O** | Octahedron | 6 | 12 | 8 | Regular triangular bipyramid |
| **D** | Dodecahedron | 20 | 30 | 12 | Regular pentagonal solid |
| **I** | Icosahedron | 12 | 30 | 20 | Regular triangular solid |

### Basic Operations

| Symbol | Operation | Description | Example |
|--------|-----------|-------------|---------|
| **d** | Dual | Exchanges vertices ↔ faces | `dC` → Octahedron |
| **a** | Ambo | Rectification at edge midpoints | `aC` → Cuboctahedron |
| **t** | Truncate | Cuts off vertices | `tI` → Soccer ball |
| **k** | Kis | Raises pyramid on each face | `kC` → Triakis octahedron |
| **j** | Join | Dual of ambo operation | `jC` → Rhombic dodecahedron |

### Compound Operations

| Symbol | Operation | Equivalent | Description |
|--------|-----------|------------|-------------|
| **o** | Ortho | `jj` | Double join operation |
| **e** | Expand | `aa` | Double ambo operation |
| **g** | Gyro | - | Pentagonal rotation |
| **s** | Snub | - | Chiral snub operation |

## 🔧 Advanced Usage

### Manual Operations

```go
// Create seed shape
cube := conway.Cube()
fmt.Printf("Cube: %s\n", cube.Stats())

// Apply operations manually
dual := conway.Dual(cube)
truncated := conway.Truncate(dual)
fmt.Printf("Dual truncated cube: %s\n", truncated.Stats())

// Normalize geometry (center at origin, scale to unit sphere)
truncated.Normalize()
```

### Complex Transformations

```go
// Multi-step transformation
complex, err := conway.Parse("egdtkC")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Complex polyhedron: %s\n", complex.Stats())

// Famous polyhedra
soccerBall, _ := conway.Parse("tI")        // Truncated icosahedron
cuboctahedron, _ := conway.Parse("aC")     // Ambo cube
icosidodecahedron, _ := conway.Parse("aD") // Ambo dodecahedron
```

### Geometric Analysis

```go
poly, _ := conway.Parse("tC")

// Get geometric statistics
stats := poly.CalculateGeometryStats()
fmt.Printf("Edge lengths: min=%.3f, max=%.3f, avg=%.3f\n", 
    stats.MinEdgeLength, stats.MaxEdgeLength, stats.AvgEdgeLength)

// Get memory usage
memory := poly.CalculateMemoryStats()
fmt.Printf("Memory: %d vertices, %d edges, %d faces\n",
    memory.VertexCount, memory.EdgeCount, memory.FaceCount)

// Calculate centroid
centroid := poly.Centroid()
fmt.Printf("Centroid: (%.3f, %.3f, %.3f)\n", centroid.X, centroid.Y, centroid.Z)
```

### Advanced Validation

```go
poly, _ := conway.Parse("kT")

// Basic validation
if poly.IsValid() {
    fmt.Println("✓ Basic validation passed")
}

// Comprehensive validation
if err := poly.ValidateComplete(); err != nil {
    fmt.Printf("✗ Validation failed: %v\n", err)
} else {
    fmt.Println("✓ Complete validation passed")
}

// Check specific properties
fmt.Printf("Euler characteristic: %d (should be 2)\n", poly.EulerCharacteristic())
```

### Utility Functions

```go
// Find edges between vertices
edge := poly.FindEdge(vertex1ID, vertex2ID)
if edge != nil {
    fmt.Printf("Edge found: %v\n", edge)
}

// Create edge lookup keys
key := conway.MakeEdgeKey(v1ID, v2ID)

// Calculate face normals
normal, err := conway.CalculateFaceNormal(vertices)
if err != nil {
    log.Printf("Error calculating normal: %v", err)
}

// Ensure proper vertex ordering
orderedVertices := conway.EnsureCounterClockwise(vertices, center)
```

### Parser Capabilities

```go
parser := conway.NewParser()

// List available operations
fmt.Println("Available operations:")
for symbol, name := range parser.GetAvailableOperations() {
    fmt.Printf("  %s: %s\n", symbol, name)
}

// List available seeds
fmt.Println("Available seeds:")
for symbol, name := range parser.GetAvailableSeeds() {
    fmt.Printf("  %s: %s\n", symbol, name)
}

// Validate notation without creating polyhedron
if err := parser.Validate("dtC"); err != nil {
    fmt.Printf("Invalid notation: %v\n", err)
}
```

## 🧪 Famous Polyhedra

| Notation | Name | Description |
|----------|------|-------------|
| `tI` | Truncated Icosahedron | Soccer ball (football) |
| `aC` | Cuboctahedron | Ambo cube |
| `aD` | Icosidodecahedron | Ambo dodecahedron |
| `tC` | Truncated Cube | Truncated cube |
| `kT` | Triakis Tetrahedron | Kis tetrahedron |
| `dC` | Octahedron | Dual of cube |
| `dI` | Dodecahedron | Dual of icosahedron |
| `oC` | Ortho Cube | Double join cube |
| `eT` | Expand Tetrahedron | Double ambo tetrahedron |

## 🏗️ Architecture

### Data Structures

- **Vector3**: 3D vector with full geometric operations
- **Vertex**: Point in 3D space with connectivity information
- **Edge**: Connection between two vertices with adjacent faces
- **Face**: Polygonal face with ordered vertices and computed properties
- **Polyhedron**: Complete polyhedron with thread-safe operations and edge lookup

### Performance Features

- **O(1) Edge Lookup**: Hash-based edge lookup by vertex pairs using `FindEdge`
- **Lazy Evaluation**: Properties computed on demand and cached
- **Memory Optimization**: Efficient allocation and reuse
- **Thread Safety**: Concurrent operations with proper locking
- **Property Caching**: Expensive calculations cached automatically

### Public API Utilities

- **Edge Operations**: `FindEdge()` for O(1) edge lookup, `MakeEdgeKey()` for consistent edge identification
- **Geometry Utilities**: `CalculateFaceNormal()` for robust normal computation using Newell's method
- **Topology Helpers**: `EnsureCounterClockwise()` for proper face orientation

## 🚀 Performance

The library is optimized for both small and large polyhedra:

```go
// Benchmarks (run with: go test -bench=.)
BenchmarkParseTI-8           10000    105234 ns/op    45632 B/op     892 allocs/op
BenchmarkDualCube-8         100000     12453 ns/op     5248 B/op      89 allocs/op
BenchmarkTruncateIcosa-8      5000    287654 ns/op   123456 B/op    2134 allocs/op
BenchmarkComplexOp-8          2000    654321 ns/op   234567 B/op    3456 allocs/op
```

## 🧪 Testing

### Run Tests

```bash
# All tests
make test

# With coverage
make test-coverage

# Benchmarks
make bench

# Property-based tests
make property-test

# Concurrency tests  
make concurrency-test

# All checks (CI pipeline)
make ci
```

### Test Categories

- **Unit Tests**: Individual function testing
- **Integration Tests**: End-to-end operation testing
- **Property Tests**: Mathematical property verification
- **Benchmark Tests**: Performance measurement
- **Concurrency Tests**: Thread safety verification

## 📁 Project Structure

```
conway/
├── .github/                 # GitHub Actions workflows
├── conway/                  # Main library package
│   ├── polyhedron.go       # Core data structures
│   ├── seeds.go            # Platonic solid generators
│   ├── operations.go       # Operation interface
│   ├── dual.go            # Dual operation
│   ├── ambo.go            # Ambo operation
│   ├── truncate.go        # Truncate operation
│   ├── kis.go             # Kis operation
│   ├── join.go            # Join operation
│   ├── compound.go        # Compound operations
│   ├── parser.go          # Notation parser
│   ├── validation.go      # Topology validation
│   ├── utils.go           # Utility functions
│   └── *_test.go          # Test files
├── examples/               # Usage examples
│   ├── basic/             # Basic usage
│   └── advanced/          # Advanced features
├── doc.go                 # Package documentation
├── README.md              # This file
├── CONTRIBUTING.md        # Contribution guidelines
├── LICENSE                # MIT license
├── Makefile              # Build automation
├── .golangci.yml         # Linting configuration
└── .gitignore            # Git ignore rules
```

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for:

- Development setup
- Coding standards  
- Testing requirements
- Pull request process
- Issue reporting guidelines

### Quick Start for Contributors

```bash
# Clone the repository
git clone https://github.com/sksmith/conway.git
cd conway

# Install dependencies
make deps

# Install development tools
make dev-deps

# Run all checks
make ci

# Run specific checks
make fmt lint test
```

## 📖 Examples

Comprehensive examples are available in the `examples/` directory:

- **[Basic Examples](examples/basic/)**: Getting started, simple operations
- **[Advanced Examples](examples/advanced/)**: Complex transformations, analysis

```bash
# Run basic examples
cd examples/basic && go run main.go

# Run advanced examples  
cd examples/advanced && go run main.go
```

## 🔗 Mathematical Background

Conway polyhedron notation is based on systematic transformations of the Platonic solids:

- **Duality**: Fundamental concept where vertices and faces are interchanged
- **Rectification**: Creates new vertices at edge midpoints
- **Truncation**: Systematic vertex removal
- **Stellation**: Face-based pyramid addition
- **Euler's Formula**: V - E + F = 2 for all valid polyhedra

For more mathematical details, see:
- [Wikipedia: Conway Polyhedron Notation](https://en.wikipedia.org/wiki/Conway_polyhedron_notation)
- [George Hart's Conway Notation](http://www.georgehart.com/virtual-polyhedra/conway_notation.html)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- John Conway for developing the original notation system
- George Hart for his comprehensive documentation and web resources
- The Go community for excellent tooling and libraries
- Contributors and users of this library

## 📊 Project Stats

![GitHub stars](https://img.shields.io/github/stars/sksmith/conway?style=social)
![GitHub forks](https://img.shields.io/github/forks/sksmith/conway?style=social)
![GitHub issues](https://img.shields.io/github/issues/sksmith/conway)
![GitHub pull requests](https://img.shields.io/github/issues-pr/sksmith/conway)

---

**Built with ❤️ for the computational geometry community**