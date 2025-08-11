package main

import (
	"fmt"
	"math"

	"github.com/sksmith/conway/conway"
)

func main() {
	fmt.Println("Conway Polyhedron Notation Library - Advanced Examples")
	fmt.Println("======================================================")

	fmt.Println("\n1. Systematic exploration of operations on all seeds:")
	seeds := []string{"T", "C", "O", "D", "I"}
	operations := []string{"d", "a", "t", "k", "j"}

	for _, seed := range seeds {
		fmt.Printf("\n   Operations on %s:\n", seed)
		original := conway.GetSeed(seed)
		fmt.Printf("     %s: %s\n", seed, original.Stats())

		for _, op := range operations {
			notation := op + seed
			poly, err := conway.Parse(notation)
			if err != nil {
				fmt.Printf("     Error with %s: %v\n", notation, err)
				continue
			}
			fmt.Printf("     %s: %s\n", notation, poly.Stats())
		}
	}

	fmt.Println("\n2. Compound operations:")
	compounds := []struct {
		notation string
		name     string
	}{
		{"oC", "Ortho Cube (jjC)"},
		{"eT", "Expand Tetrahedron (aaT)"},
		{"gC", "Gyro Cube"},
		{"sO", "Snub Octahedron"},
	}

	for _, compound := range compounds {
		poly, err := conway.Parse(compound.notation)
		if err != nil {
			fmt.Printf("   Error with %s: %v\n", compound.name, err)
			continue
		}
		fmt.Printf("   %s: %s\n", compound.name, poly.Stats())
	}

	fmt.Println("\n3. Multi-step transformations:")
	transformations := []string{
		"tkdC",
		"dtakT",
		"jsakO",
		"egdtI",
	}

	for _, notation := range transformations {
		poly, err := conway.Parse(notation)
		if err != nil {
			fmt.Printf("   Error with %s: %v\n", notation, err)
			continue
		}
		fmt.Printf("   %s: %s (valid: %v)\n", notation, poly.Stats(), poly.IsValid())
	}

	fmt.Println("\n4. Geometric properties analysis:")
	analyzePolyhedron("Tetrahedron", conway.Tetrahedron())
	analyzePolyhedron("Truncated Cube", conway.MustParse("tC"))
	analyzePolyhedron("Kis Octahedron", conway.MustParse("kO"))

	fmt.Println("\n5. Operation validation and properties:")
	cube := conway.Cube()

	testDualInvolution(cube)
	testEulerCharacteristic(cube)

	fmt.Println("\n6. Parser capabilities:")
	parser := conway.NewParser()

	fmt.Println("   Available operations:")
	for symbol, name := range parser.GetAvailableOperations() {
		fmt.Printf("     %s: %s\n", symbol, name)
	}

	fmt.Println("\n   Available seeds:")
	for symbol, name := range parser.GetAvailableSeeds() {
		fmt.Printf("     %s: %s\n", symbol, name)
	}

	fmt.Println("\n7. Error handling:")
	invalidNotations := []string{"xT", "T x", "", "dX", "123"}

	for _, notation := range invalidNotations {
		_, err := conway.Parse(notation)
		if err != nil {
			fmt.Printf("   %s -> Error: %v\n", notation, err)
		}
	}
}

func analyzePolyhedron(name string, poly *conway.Polyhedron) {
	fmt.Printf("\n   Analysis of %s:\n", name)
	fmt.Printf("     %s\n", poly.Stats())
	fmt.Printf("     Valid: %v\n", poly.IsValid())

	centroid := poly.Centroid()
	fmt.Printf("     Centroid: (%.3f, %.3f, %.3f)\n", centroid.X, centroid.Y, centroid.Z)

	maxDist := 0.0
	minDist := math.Inf(1)
	for _, v := range poly.Vertices {
		dist := v.Position.Length()
		if dist > maxDist {
			maxDist = dist
		}
		if dist < minDist {
			minDist = dist
		}
	}
	fmt.Printf("     Vertex distances: min=%.3f, max=%.3f\n", minDist, maxDist)

	totalSurfaceArea := 0.0
	for _, f := range poly.Faces {
		totalSurfaceArea += f.Area()
	}
	fmt.Printf("     Total surface area: %.3f\n", totalSurfaceArea)

	vertexDegrees := make(map[int]int)
	for _, v := range poly.Vertices {
		degree := v.Degree()
		vertexDegrees[degree]++
	}

	fmt.Print("     Vertex degree distribution: ")
	for degree, count := range vertexDegrees {
		fmt.Printf("%d-valent:%d ", degree, count)
	}
	fmt.Println()

	faceSizes := make(map[int]int)
	for _, f := range poly.Faces {
		size := f.Degree()
		faceSizes[size]++
	}

	fmt.Print("     Face size distribution: ")
	for size, count := range faceSizes {
		fmt.Printf("%d-gon:%d ", size, count)
	}
	fmt.Println()
}

func testDualInvolution(poly *conway.Polyhedron) {
	fmt.Printf("\n   Testing dual involution on %s:\n", poly.Name)
	fmt.Printf("     Original: %s\n", poly.Stats())

	dual := conway.Dual(poly)
	fmt.Printf("     Dual: %s\n", dual.Stats())

	doubleDual := conway.Dual(dual)
	fmt.Printf("     Double dual: %s\n", doubleDual.Stats())

	originalValid := poly.IsValid()
	dualValid := dual.IsValid()
	doubleDualValid := doubleDual.IsValid()

	fmt.Printf("     Validity: orig=%v, dual=%v, double=%v\n",
		originalValid, dualValid, doubleDualValid)

	if len(poly.Vertices) == len(doubleDual.Vertices) &&
		len(poly.Faces) == len(doubleDual.Faces) {
		fmt.Println("     ✓ Dual involution property satisfied")
	} else {
		fmt.Println("     ✗ Dual involution property violated")
	}
}

func testEulerCharacteristic(poly *conway.Polyhedron) {
	fmt.Printf("\n   Testing Euler characteristic preservation:\n")

	operations := []struct {
		name string
		op   func(*conway.Polyhedron) *conway.Polyhedron
	}{
		{"original", func(p *conway.Polyhedron) *conway.Polyhedron { return p }},
		{"dual", conway.Dual},
		{"ambo", conway.Ambo},
		{"truncate", conway.Truncate},
		{"kis", conway.Kis},
	}

	for _, op := range operations {
		result := op.op(poly)
		chi := result.EulerCharacteristic()
		fmt.Printf("     %s: χ = %d %s\n", op.name, chi,
			map[bool]string{true: "✓", false: "✗"}[chi == 2])
	}
}
