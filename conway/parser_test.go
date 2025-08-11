package conway_test

import (
	"strings"
	"testing"

	"github.com/sksmith/conway/conway"
)

func TestParseValidNotations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		notation string
		valid    bool
	}{
		{"T", true},
		{"C", true},
		{"O", true},
		{"D", true},
		{"I", true},
		{"dT", true},
		{"aT", true},
		{"tC", true},
		{"kO", true},
		{"jD", true},
		{"oI", true},
		{"eC", true},
		{"gT", true},
		{"sO", true},
		{"dtC", true},
		{"akT", true},
		{"", false},
		{"X", false},
		{"dX", false},
		{"xT", false},
	}

	parser := conway.NewParser()

	for _, test := range tests {
		t.Run(test.notation, func(t *testing.T) {
			t.Parallel()

			result, err := parser.Parse(test.notation)

			if test.valid {
				if err != nil {
					t.Errorf("Expected valid notation %s, got error: %v", test.notation, err)
					return
				}

				if result == nil {
					t.Errorf("Expected result for valid notation %s", test.notation)
					return
				}

				if !result.IsValid() {
					t.Errorf("Result of %s is not a valid polyhedron", test.notation)
				}

				return
			}

			if err == nil {
				t.Errorf("Expected error for invalid notation %s", test.notation)
			}
		})
	}
}

func TestParseOperationOrder(t *testing.T) {
	t.Parallel()

	parser := conway.NewParser()

	result, err := parser.Parse("dtC")
	if err != nil {
		t.Fatalf("Failed to parse dtC: %v", err)
	}

	if !result.IsValid() {
		t.Error("Result of dtC is not valid")
	}

	if !strings.Contains(result.Name, "d") || !strings.Contains(result.Name, "t") {
		t.Errorf("Result name should contain operation symbols, got: %s", result.Name)
	}
}

func TestParserHelperMethods(t *testing.T) {
	t.Parallel()

	parser := conway.NewParser()

	t.Run("GetAvailableOperations", func(t *testing.T) {
		t.Parallel()

		ops := parser.GetAvailableOperations()
		expectedOps := []string{"d", "a", "t", "k", "j", "o", "e", "g", "s"}

		for _, op := range expectedOps {
			if _, exists := ops[op]; !exists {
				t.Errorf("Expected operation %s not found in available operations", op)
			}
		}
	})

	t.Run("GetAvailableSeeds", func(t *testing.T) {
		t.Parallel()

		seeds := parser.GetAvailableSeeds()
		expectedSeeds := []string{"T", "C", "O", "D", "I"}

		for _, seed := range expectedSeeds {
			if _, exists := seeds[seed]; !exists {
				t.Errorf("Expected seed %s not found in available seeds", seed)
			}
		}
	})

	t.Run("Validate", func(t *testing.T) {
		t.Parallel()

		if err := parser.Validate("dtC"); err != nil {
			t.Errorf("Validation failed for valid notation: %v", err)
		}

		if err := parser.Validate("xT"); err == nil {
			t.Error("Validation should have failed for invalid notation")
		}
	})
}

func TestGlobalParseFunction(t *testing.T) {
	t.Parallel()

	result, err := conway.Parse("tI")
	if err != nil {
		t.Fatalf("Global Parse failed: %v", err)
	}

	if !result.IsValid() {
		t.Error("Global Parse result is not valid")
	}
}

func TestMustParse(t *testing.T) {
	t.Parallel()

	t.Run("ValidNotation", func(t *testing.T) {
		t.Parallel()

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("MustParse panicked on valid notation: %v", r)
			}
		}()

		result := conway.MustParse("aC")
		if !result.IsValid() {
			t.Error("MustParse result is not valid")
		}
	})

	t.Run("InvalidNotation", func(t *testing.T) {
		t.Parallel()

		defer func() {
			if r := recover(); r == nil {
				t.Error("MustParse should have panicked on invalid notation")
			}
		}()

		conway.MustParse("xT")
	})
}

func TestComplexNotations(t *testing.T) {
	t.Parallel()

	complexNotations := []string{
		"dtC",
		"akT",
		"jsO",
		"egI",
		"tkdC",
		"aeT",
	}

	parser := conway.NewParser()

	for _, notation := range complexNotations {
		t.Run(notation, func(t *testing.T) {
			t.Parallel()

			result, err := parser.Parse(notation)
			if err != nil {
				t.Errorf("Failed to parse %s: %v", notation, err)
				return
			}

			if !result.IsValid() {
				t.Errorf("Result of %s is not valid: %s", notation, result.Stats())
			}

			if result.EulerCharacteristic() != 2 {
				t.Errorf("Result of %s has wrong Euler characteristic: %d",
					notation, result.EulerCharacteristic())
			}
		})
	}
}
