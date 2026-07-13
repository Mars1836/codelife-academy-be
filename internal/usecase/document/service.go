package document

import (
	"context"
	"strings"

	domain "codelife-study-be/internal/domain/document"
)

type Service struct {
	repository domain.Repository
	cache      domain.Cache
}

func New(repository domain.Repository, cache domain.Cache) *Service {
	return &Service{repository: repository, cache: cache}
}

func (s *Service) List(ctx context.Context) ([]domain.Document, error) {
	const key = "documents:list:v1"
	var documents []domain.Document
	if s.cache != nil && s.cache.Get(ctx, key, &documents) {
		return documents, nil
	}
	documents, err := s.repository.List(ctx)
	if err == nil && s.cache != nil {
		_ = s.cache.Set(ctx, key, documents)
	}
	return documents, err
}

func (s *Service) Get(ctx context.Context, slug string) (domain.Document, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" || strings.ContainsAny(slug, `/\\`) {
		return domain.Document{}, domain.ErrNotFound
	}
	key := "documents:item:v1:" + slug
	var doc domain.Document
	if s.cache != nil && s.cache.Get(ctx, key, &doc) {
		return doc, nil
	}
	doc, err := s.repository.FindBySlug(ctx, slug)
	if err == nil && s.cache != nil {
		_ = s.cache.Set(ctx, key, doc)
	}
	return doc, err
}
