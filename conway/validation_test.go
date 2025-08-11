package conway_test

import (
	"strings"
	"testing"

	"github.com/sksmith/conway/conway"
)

func TestValidationError(t *testing.T) {
	t.Parallel()

	err := conway.ValidationError{
		Type:    "test_error",
		Message: "This is a test error message",
	}

	errorStr := err.Error()
	if !strings.Contains(errorStr, "test_error") {
		t.Errorf("Error string should contain error type, got: %s", errorStr)
	}

	if !strings.Contains(errorStr, "This is a test error message") {
		t.Errorf("Error string should contain error message, got: %s", errorStr)
	}
}

func TestValidationErrorCases(t *testing.T) {
	t.Parallel()

	p := conway.NewPolyhedron("invalid")

	v1 := p.AddVertex(conway.Vector3{0, 0, 0})
	v2 := p.AddVertex(conway.Vector3{1, 0, 0})

	p.AddFace([]*conway.Vertex{v1, v2})

	err := p.ValidateComplete()
	if err == nil {
		t.Error("Expected validation error for invalid polyhedron")
		return
	}

	if err.Error() == "" {
		t.Error("Validation error should have non-empty error message")
	}
}

func TestValidateManifoldErrors(t *testing.T) {
	t.Parallel()

	p := conway.NewPolyhedron("test")

	v1 := p.AddVertex(conway.Vector3{0, 0, 0})
	v2 := p.AddVertex(conway.Vector3{1, 0, 0})
	v3 := p.AddVertex(conway.Vector3{0, 1, 0})
	v4 := p.AddVertex(conway.Vector3{0, 0, 1})

	p.AddFace([]*conway.Vertex{v1, v2, v3})
	p.AddFace([]*conway.Vertex{v1, v2, v4})
	p.AddFace([]*conway.Vertex{v2, v3, v4})
	p.AddFace([]*conway.Vertex{v1, v3, v4})
	p.AddFace([]*conway.Vertex{v1, v2, v3})

	err := p.ValidateManifold()
	if err == nil {
		t.Error("Expected manifold validation error for non-manifold polyhedron")
	}
}
