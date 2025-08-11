package conway_test

import (
	"testing"

	"github.com/sksmith/conway/conway"
)

// BenchmarkPolyhedronCreation benchmarks the creation of seed polyhedra.
func BenchmarkPolyhedronCreation(b *testing.B) {
	benchmarks := []struct {
		name string
		fn   func() *conway.Polyhedron
	}{
		{"Tetrahedron", conway.Tetrahedron},
		{"Cube", conway.Cube},
		{"Octahedron", conway.Octahedron},
		{"Dodecahedron", conway.Dodecahedron},
		{"Icosahedron", conway.Icosahedron},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = bm.fn()
			}
		})
	}
}

// BenchmarkDualOperation benchmarks the dual operation on various polyhedra.
func BenchmarkDualOperation(b *testing.B) {
	benchmarks := []struct {
		name string
		poly *conway.Polyhedron
	}{
		{"Tetrahedron", conway.Tetrahedron()},
		{"Cube", conway.Cube()},
		{"Octahedron", conway.Octahedron()},
		{"Dodecahedron", conway.Dodecahedron()},
		{"Icosahedron", conway.Icosahedron()},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = conway.Dual(bm.poly)
			}
		})
	}
}

// BenchmarkAmboOperation benchmarks the ambo operation.
func BenchmarkAmboOperation(b *testing.B) {
	benchmarks := []struct {
		name string
		poly *conway.Polyhedron
	}{
		{"Tetrahedron", conway.Tetrahedron()},
		{"Cube", conway.Cube()},
		{"Octahedron", conway.Octahedron()},
		{"Dodecahedron", conway.Dodecahedron()},
		{"Icosahedron", conway.Icosahedron()},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = conway.Ambo(bm.poly)
			}
		})
	}
}

// BenchmarkTruncateOperation benchmarks the truncate operation.
func BenchmarkTruncateOperation(b *testing.B) {
	benchmarks := []struct {
		name string
		poly *conway.Polyhedron
	}{
		{"Tetrahedron", conway.Tetrahedron()},
		{"Cube", conway.Cube()},
		{"Octahedron", conway.Octahedron()},
		{"Dodecahedron", conway.Dodecahedron()},
		{"Icosahedron", conway.Icosahedron()},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = conway.Truncate(bm.poly)
			}
		})
	}
}

// BenchmarkKisOperation benchmarks the kis operation.
func BenchmarkKisOperation(b *testing.B) {
	benchmarks := []struct {
		name string
		poly *conway.Polyhedron
	}{
		{"Tetrahedron", conway.Tetrahedron()},
		{"Cube", conway.Cube()},
		{"Octahedron", conway.Octahedron()},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = conway.Kis(bm.poly)
			}
		})
	}
}

// BenchmarkJoinOperation benchmarks the join operation.
func BenchmarkJoinOperation(b *testing.B) {
	benchmarks := []struct {
		name string
		poly *conway.Polyhedron
	}{
		{"Tetrahedron", conway.Tetrahedron()},
		{"Cube", conway.Cube()},
		{"Octahedron", conway.Octahedron()},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = conway.Join(bm.poly)
			}
		})
	}
}

// BenchmarkParser benchmarks the notation parser.
func BenchmarkParser(b *testing.B) {
	testStrings := []string{
		"T", "C", "O", "D", "I",
		"dT", "aC", "tO", "kC", "jT",
		"dtC", "akO", "taC",
		"dtakC",
	}

	for _, str := range testStrings {
		b.Run(str, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = conway.Parse(str)
			}
		})
	}
}

// BenchmarkComplexOperations benchmarks complex operation chains.
func BenchmarkComplexOperations(b *testing.B) {
	cube := conway.Cube()

	b.Run("dtakC", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = conway.Parse("dtakC")
		}
	})

	b.Run("Manual_dtakC", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			k := conway.Kis(cube)
			a := conway.Ambo(k)
			t := conway.Truncate(a)
			_ = conway.Dual(t)
		}
	})
}

// BenchmarkValidation benchmarks validation operations.
func BenchmarkValidation(b *testing.B) {
	polyhedra := map[string]*conway.Polyhedron{
		"Tetrahedron":  conway.Tetrahedron(),
		"Cube":         conway.Cube(),
		"Octahedron":   conway.Octahedron(),
		"Dodecahedron": conway.Dodecahedron(),
		"Icosahedron":  conway.Icosahedron(),
	}

	for name, poly := range polyhedra {
		b.Run(name+"_IsValid", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = poly.IsValid()
			}
		})

		b.Run(name+"_ValidateComplete", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = poly.ValidateComplete()
			}
		})

		b.Run(name+"_ValidateTopology", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = poly.ValidateTopology()
			}
		})

		b.Run(name+"_ValidateManifold", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = poly.ValidateManifold()
			}
		})
	}
}

// BenchmarkEdgeLookup benchmarks the edge lookup optimization.
func BenchmarkEdgeLookup(b *testing.B) {
	// Create a polyhedron with many vertices for meaningful edge lookup benchmark.
	dodeca := conway.Dodecahedron()

	// Get some vertices for lookup testing.
	var v1, v2 *conway.Vertex
	for _, v := range dodeca.Vertices {
		if v1 == nil {
			v1 = v
		} else if v2 == nil {
			v2 = v
			break
		}
	}

	b.Run("EdgeLookup_Find", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = dodeca.FindEdge(v1.ID, v2.ID)
		}
	})

	b.Run("AddEdge_WithLookup", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			poly := conway.NewPolyhedron("test")
			vertex1 := poly.AddVertex(conway.Vector3{0, 0, 0})
			vertex2 := poly.AddVertex(conway.Vector3{1, 0, 0})
			_ = poly.AddEdge(vertex1, vertex2)
		}
	})
}

// BenchmarkGeometryCalculations benchmarks geometric computations.
func BenchmarkGeometryCalculations(b *testing.B) {
	polyhedra := map[string]*conway.Polyhedron{
		"Tetrahedron":  conway.Tetrahedron(),
		"Cube":         conway.Cube(),
		"Octahedron":   conway.Octahedron(),
		"Dodecahedron": conway.Dodecahedron(),
		"Icosahedron":  conway.Icosahedron(),
	}

	for name, poly := range polyhedra {
		b.Run(name+"_Centroid", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = poly.Centroid()
			}
		})

		b.Run(name+"_EulerCharacteristic", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = poly.EulerCharacteristic()
			}
		})

		b.Run(name+"_Normalize", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				testPoly := poly.Clone()
				testPoly.Normalize()
			}
		})

		b.Run(name+"_GeometryStats", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = poly.CalculateGeometryStats()
			}
		})
	}
}

// BenchmarkFaceOperations benchmarks face-related operations.
func BenchmarkFaceOperations(b *testing.B) {
	cube := conway.Cube()

	// Get a face for testing.
	var testFace *conway.Face
	for _, face := range cube.Faces {
		testFace = face
		break
	}

	b.Run("Face_Normal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = testFace.Normal()
		}
	})

	b.Run("Face_Area", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = testFace.Area()
		}
	})

	b.Run("Face_Centroid", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = testFace.Centroid()
		}
	})
}

// BenchmarkMemoryUsage benchmarks memory allocation patterns.
func BenchmarkMemoryUsage(b *testing.B) {
	b.Run("NewPolyhedron", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = conway.NewPolyhedron("test")
		}
	})

	b.Run("Clone", func(b *testing.B) {
		cube := conway.Cube()

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_ = cube.Clone()
		}
	})
}

// BenchmarkScalability benchmarks operations on increasingly complex polyhedra.
func BenchmarkScalability(b *testing.B) {
	// Create increasingly complex polyhedra through operations.
	base := conway.Cube()
	truncated := conway.Truncate(base)     // More complex
	compound := conway.Truncate(truncated) // Even more complex

	polyhedra := map[string]*conway.Polyhedron{
		"Simple":  base,
		"Medium":  truncated,
		"Complex": compound,
	}

	for name, poly := range polyhedra {
		b.Run(name+"_Dual", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = conway.Dual(poly)
			}
		})

		b.Run(name+"_Validation", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = poly.ValidateComplete()
			}
		})
	}
}
