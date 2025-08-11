package main

import (
	"fmt"
	"log"

	"github.com/sksmith/conway/conway"
)

func main() {
	fmt.Println("Conway Polyhedron Notation Library - Basic Examples")
	fmt.Println("==================================================")

	fmt.Println("\n1. Creating seed polyhedra:")
	seeds := []string{"T", "C", "O", "D", "I"}
	seedNames := []string{"Tetrahedron", "Cube", "Octahedron", "Dodecahedron", "Icosahedron"}

	for i, seed := range seeds {
		poly := conway.GetSeed(seed)
		fmt.Printf("   %s (%s): %s\n", seedNames[i], seed, poly.Stats())
	}

	fmt.Println("\n2. Basic operations on a cube:")
	cube := conway.Cube()
	fmt.Printf("   Original Cube: %s\n", cube.Stats())

	operations := []struct {
		name string
		op   func(*conway.Polyhedron) *conway.Polyhedron
	}{
		{"Dual", conway.Dual},
		{"Ambo", conway.Ambo},
		{"Truncate", conway.Truncate},
		{"Kis", conway.Kis},
		{"Join", conway.Join},
	}

	for _, op := range operations {
		result := op.op(cube)
		if !result.IsValid() {
			log.Printf("Warning: %s operation produced invalid polyhedron", op.name)
		}
		fmt.Printf("   %s of Cube: %s\n", op.name, result.Stats())
	}

	fmt.Println("\n3. Using Conway notation parser:")
	notations := []string{"tC", "dO", "kT", "aI", "jD"}

	for _, notation := range notations {
		poly, err := conway.Parse(notation)
		if err != nil {
			fmt.Printf("   Error parsing %s: %v\n", notation, err)
			continue
		}
		fmt.Printf("   %s: %s\n", notation, poly.Stats())
	}

	fmt.Println("\n4. Complex operations:")
	complexNotations := []string{"dtC", "akT", "jsO", "egI"}

	for _, notation := range complexNotations {
		poly, err := conway.Parse(notation)
		if err != nil {
			fmt.Printf("   Error parsing %s: %v\n", notation, err)
			continue
		}
		fmt.Printf("   %s: %s\n", notation, poly.Stats())
	}

	fmt.Println("\n5. Famous polyhedra:")
	famous := []struct {
		notation string
		name     string
	}{
		{"tI", "Truncated Icosahedron (Soccer Ball)"},
		{"aC", "Cuboctahedron"},
		{"aD", "Icosidodecahedron"},
		{"kT", "Triakis Tetrahedron"},
		{"dC", "Octahedron (dual of Cube)"},
	}

	for _, f := range famous {
		poly, err := conway.Parse(f.notation)
		if err != nil {
			fmt.Printf("   Error creating %s: %v\n", f.name, err)
			continue
		}
		fmt.Printf("   %s (%s): %s\n", f.name, f.notation, poly.Stats())
	}

	fmt.Println("\n6. Dual relationships:")
	cube = conway.Cube()
	octahedron := conway.Octahedron()

	dualCube := conway.Dual(cube)
	dualOctahedron := conway.Dual(octahedron)

	fmt.Printf("   Cube: %s\n", cube.Stats())
	fmt.Printf("   Dual of Cube: %s\n", dualCube.Stats())
	fmt.Printf("   Octahedron: %s\n", octahedron.Stats())
	fmt.Printf("   Dual of Octahedron: %s\n", dualOctahedron.Stats())

	doubleDualCube := conway.Dual(dualCube)
	fmt.Printf("   Double dual of Cube: %s\n", doubleDualCube.Stats())
}
