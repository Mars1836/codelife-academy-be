package progress

import (
	"context"
	"errors"
	"testing"

	domain "codelife-study-be/internal/domain/progress"
)

type fakeRepository struct {
	items map[string]domain.LearningProgress
}

func (f *fakeRepository) ListByUser(context.Context, string) ([]domain.LearningProgress, error) {
	result := make([]domain.LearningProgress, 0, len(f.items))
	for _, item := range f.items {
		result = append(result, item)
	}
	return result, nil
}

func (f *fakeRepository) FindByUserAndDocument(_ context.Context, _, slug string) (domain.LearningProgress, error) {
	item, ok := f.items[slug]
	if !ok {
		return domain.LearningProgress{}, domain.ErrNotFound
	}
	return item, nil
}

func (f *fakeRepository) Upsert(_ context.Context, _ string, item domain.LearningProgress) (domain.LearningProgress, error) {
	f.items[item.DocumentSlug] = item
	return item, nil
}

func TestUpdateCreatesAndMergesProgress(t *testing.T) {
	repository := &fakeRepository{items: map[string]domain.LearningProgress{}}
	service := New(repository)
	status := "studying"
	scroll := 320

	item, err := service.Update(context.Background(), "user-1", "go-basics", UpdateInput{
		Status:         &status,
		ScrollPosition: &scroll,
	})
	if err != nil {
		t.Fatal(err)
	}
	if item.Status != status || item.ScrollPosition != scroll || item.DocumentSlug != "go-basics" {
		t.Fatalf("unexpected progress: %#v", item)
	}

	note := "important"
	item, err = service.Update(context.Background(), "user-1", "go-basics", UpdateInput{Note: &note})
	if err != nil {
		t.Fatal(err)
	}
	if item.Status != status || item.ScrollPosition != scroll || item.Note != note {
		t.Fatalf("partial update lost existing data: %#v", item)
	}
}

func TestUpdateRejectsInvalidInput(t *testing.T) {
	service := New(&fakeRepository{items: map[string]domain.LearningProgress{}})
	invalidStatus := "deleted"
	_, err := service.Update(context.Background(), "user-1", "go-basics", UpdateInput{Status: &invalidStatus})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected invalid input, got %v", err)
	}
}
