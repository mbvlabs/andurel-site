package services

import (
	"errors"
	"testing"

	"andurel-site/models"
)

func TestNextVersionPosition(t *testing.T) {
	if got := nextVersionPosition(nil); got != 0 {
		t.Fatalf("empty = %d, want 0", got)
	}

	got := nextVersionPosition([]models.DocumentationVersion{
		{Position: 0},
		{Position: 2},
		{Position: 1},
	})
	if got != 3 {
		t.Fatalf("next = %d, want 3", got)
	}
}

func TestApplyVersionOrder(t *testing.T) {
	existing := []models.DocumentationVersion{
		{ID: 10, Slug: "latest", Label: "latest", Position: 0},
		{ID: 11, Slug: "head", Label: "head", Position: 1},
		{ID: 12, Slug: "1.5.5", Label: "1.5.5", Position: 1},
	}

	updates, err := applyVersionOrder(existing, []int64{12, 10, 11})
	if err != nil {
		t.Fatalf("applyVersionOrder: %v", err)
	}
	if len(updates) != 3 {
		t.Fatalf("updates = %d, want 3", len(updates))
	}
	if updates[0].ID != 12 || updates[0].Position != 0 {
		t.Fatalf("first = %+v", updates[0])
	}
	if updates[1].ID != 10 || updates[1].Position != 1 {
		t.Fatalf("second = %+v", updates[1])
	}
	if updates[2].ID != 11 || updates[2].Position != 2 {
		t.Fatalf("third = %+v", updates[2])
	}

	if _, err := applyVersionOrder(existing, []int64{10, 11}); !errors.Is(err, ErrVersionOrderMismatch) {
		t.Fatalf("missing id error = %v", err)
	}
	if _, err := applyVersionOrder(existing, []int64{10, 11, 99}); !errors.Is(err, ErrVersionOrderMismatch) {
		t.Fatalf("unknown id error = %v", err)
	}
	if _, err := applyVersionOrder(existing, []int64{10, 10, 11}); !errors.Is(err, ErrVersionOrderMismatch) {
		t.Fatalf("duplicate id error = %v", err)
	}
}
