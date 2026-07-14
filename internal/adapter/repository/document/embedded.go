package document

import (
	"context"
	"errors"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	domain "codelife-study-be/internal/domain/document"
	assets "codelife-study-be/src"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EmbeddedRepository struct {
	db *pgxpool.Pool
}

func NewEmbeddedRepository(db ...*pgxpool.Pool) *EmbeddedRepository {
	repository := &EmbeddedRepository{}
	if len(db) > 0 {
		repository.db = db[0]
	}
	return repository
}

func (r *EmbeddedRepository) List(ctx context.Context) ([]domain.Document, error) {
	if r.db != nil {
		return r.listFromDatabase(ctx)
	}
	return r.listFromEmbedded()
}

func (r *EmbeddedRepository) SyncMetadata(ctx context.Context) error {
	if r.db == nil {
		return nil
	}
	documents, err := r.listFromEmbedded()
	if err != nil {
		return err
	}
	for _, document := range documents {
		if _, err := r.db.Exec(ctx, `
			INSERT INTO documents (slug, title, category, word_count, reading_time, updated_at)
			VALUES ($1, $2, $3, $4, $5, NOW())
			ON CONFLICT (slug) DO UPDATE SET
				title = EXCLUDED.title,
				category = EXCLUDED.category,
				word_count = EXCLUDED.word_count,
				reading_time = EXCLUDED.reading_time,
				updated_at = NOW()
		`, document.Slug, document.Title, document.Category, document.WordCount, document.ReadingTime); err != nil {
			return err
		}
	}
	return nil
}

func (r *EmbeddedRepository) listFromEmbedded() ([]domain.Document, error) {
	entries, err := fs.ReadDir(assets.Documents, "documents")
	if err != nil {
		return nil, err
	}
	documents := make([]domain.Document, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}
		raw, err := assets.Documents.ReadFile("documents/" + entry.Name())
		if err != nil {
			return nil, err
		}
		documents = append(documents, domain.Document{
			Slug:        strings.TrimSuffix(entry.Name(), ".md"),
			Title:       extractTitle(raw, entry.Name()),
			Category:    categoryFor(entry.Name()),
			WordCount:   wordCount(raw),
			ReadingTime: readingTime(raw),
		})
	}
	sort.Slice(documents, func(i, j int) bool { return documents[i].Title < documents[j].Title })
	return documents, nil
}

func (r *EmbeddedRepository) listFromDatabase(ctx context.Context) ([]domain.Document, error) {
	rows, err := r.db.Query(ctx, `
		SELECT slug, title, category, word_count, reading_time
		FROM documents
		ORDER BY title ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	documents, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Document, error) {
		var document domain.Document
		err := row.Scan(&document.Slug, &document.Title, &document.Category, &document.WordCount, &document.ReadingTime)
		return document, err
	})
	if err != nil {
		return nil, err
	}
	if len(documents) == 0 {
		return r.listFromEmbedded()
	}
	return documents, nil
}

func (r *EmbeddedRepository) FindBySlug(_ context.Context, slug string) (domain.Document, error) {
	raw, err := assets.Documents.ReadFile("documents/" + slug + ".md")
	if errors.Is(err, fs.ErrNotExist) {
		return domain.Document{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Document{}, err
	}
	return domain.Document{
		Slug:        slug,
		Title:       extractTitle(raw, slug),
		Category:    categoryFor(slug),
		WordCount:   wordCount(raw),
		ReadingTime: readingTime(raw),
		Content:     string(raw),
	}, nil
}

func extractTitle(raw []byte, fallback string) string {
	for _, line := range strings.Split(string(raw), "\n") {
		line = strings.Trim(strings.TrimSpace(strings.TrimLeft(line, "#")), "*_ ")
		if line != "" {
			return line
		}
	}
	return strings.TrimSuffix(fallback, ".md")
}

func categoryFor(name string) string {
	name = strings.ToLower(name)
	switch {
	case strings.Contains(name, "acid"), strings.Contains(name, "postgresql"), strings.Contains(name, "redis"), strings.Contains(name, "locks"), strings.Contains(name, "storage"):
		return "database"
	case strings.Contains(name, "kafka"), strings.Contains(name, "log"), strings.Contains(name, "giao_tiep"), strings.Contains(name, "1tr_users"):
		return "architecture"
	case strings.Contains(name, "dsa"), strings.Contains(name, "solid"):
		return "algorithm-design"
	case strings.Contains(name, "owasp"):
		return "security"
	case strings.Contains(name, "docker"), strings.Contains(name, "kubernetes"), strings.Contains(name, "devops"):
		return "devops"
	default:
		return "backend"
	}
}

func wordCount(raw []byte) int {
	return len(strings.Fields(string(raw)))
}

func readingTime(raw []byte) int {
	count := wordCount(raw)
	minutes := (count + 199) / 200
	if minutes < 1 {
		return 1
	}
	return minutes
}
