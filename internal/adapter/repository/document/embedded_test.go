package document

import (
	"context"
	"testing"
)

func TestEmbeddedRepositoryListAndFindSubdirectories(t *testing.T) {
	repo := NewEmbeddedRepository()
	docs, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(docs) == 0 {
		t.Fatalf("List() returned empty documents")
	}

	foundGroup := false
	for _, doc := range docs {
		if doc.GroupSlug == "onphongvanpython" {
			foundGroup = true
			if doc.GroupTitle == "" {
				t.Errorf("expected non-empty GroupTitle for slug %s", doc.Slug)
			}
		}
	}
	if !foundGroup {
		t.Errorf("expected to find documents with GroupSlug 'onphongvanpython'")
	}

	// Test FindBySlug for a nested document
	slug := "onphongvanpython__03_database_sql_postgresql_mysql_mongodb"
	doc, err := repo.FindBySlug(context.Background(), slug)
	if err != nil {
		t.Fatalf("FindBySlug(%q) error = %v", slug, err)
	}
	if doc.GroupSlug != "onphongvanpython" {
		t.Errorf("expected GroupSlug 'onphongvanpython', got %q", doc.GroupSlug)
	}
	if doc.Content == "" {
		t.Errorf("expected non-empty Content")
	}
}
