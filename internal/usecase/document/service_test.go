package document

import (
	"context"
	"testing"

	domain "codelife-study-be/internal/domain/document"
)

type fakeRepository struct{ documents []domain.Document }

func (f fakeRepository) List(context.Context) ([]domain.Document, error) { return f.documents, nil }
func (f fakeRepository) FindBySlug(_ context.Context, slug string) (domain.Document, error) {
	for _, document := range f.documents {
		if document.Slug == slug {
			return document, nil
		}
	}
	return domain.Document{}, domain.ErrNotFound
}

func TestGetRejectsPathTraversal(t *testing.T) {
	service := New(fakeRepository{}, nil)
	if _, err := service.Get(context.Background(), "../secret"); err != domain.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestListReturnsRepositoryDocuments(t *testing.T) {
	want := []domain.Document{{Slug: "redis", Title: "Redis", Category: "database", WordCount: 100, ReadingTime: 1}}
	service := New(fakeRepository{documents: want}, nil)
	got, err := service.List(context.Background())
	if err != nil || len(got) != 1 || got[0].Slug != want[0].Slug || got[0].Category != want[0].Category {
		t.Fatalf("unexpected result: %#v, %v", got, err)
	}
}
