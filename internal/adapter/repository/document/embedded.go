package document

import (
	"context"
	"errors"
	"io/fs"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
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
			INSERT INTO documents (slug, title, category, group_slug, group_title, sort_order, word_count, reading_time, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
			ON CONFLICT (slug) DO UPDATE SET
				title = EXCLUDED.title,
				category = EXCLUDED.category,
				group_slug = EXCLUDED.group_slug,
				group_title = EXCLUDED.group_title,
				sort_order = EXCLUDED.sort_order,
				word_count = EXCLUDED.word_count,
				reading_time = EXCLUDED.reading_time,
				updated_at = NOW()
		`, document.Slug, document.Title, document.Category, document.GroupSlug, document.GroupTitle, document.Order, document.WordCount, document.ReadingTime); err != nil {
			return err
		}
	}
	return nil
}

func (r *EmbeddedRepository) listFromEmbedded() ([]domain.Document, error) {
	documents := make([]domain.Document, 0)
	groupTitles := make(map[string]string)

	err := fs.WalkDir(assets.Documents, "documents", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(d.Name()) != ".md" {
			return nil
		}

		relPath := strings.TrimPrefix(path, "documents/")
		relPath = strings.TrimPrefix(relPath, "/")

		raw, err := assets.Documents.ReadFile(path)
		if err != nil {
			return err
		}
		raw = stripBOM(raw)

		dir := filepath.Dir(relPath)
		var groupSlug, groupTitle string
		var order int

		if dir != "." && dir != "" {
			groupSlug = strings.ReplaceAll(dir, "\\", "/")
			if _, exists := groupTitles[groupSlug]; !exists {
				groupTitles[groupSlug] = groupTitleFor(groupSlug)
			}
			groupTitle = groupTitles[groupSlug]
			order = parseOrder(d.Name())
		}

		slug := strings.ReplaceAll(strings.TrimSuffix(relPath, ".md"), "/", "__")
		slug = strings.ReplaceAll(slug, "\\", "__")

		documents = append(documents, domain.Document{
			Slug:        slug,
			Title:       extractTitle(raw, d.Name()),
			Category:    categoryFor(relPath),
			GroupSlug:   groupSlug,
			GroupTitle:  groupTitle,
			Order:       order,
			WordCount:   wordCount(raw),
			ReadingTime: readingTime(raw),
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Slice(documents, func(i, j int) bool {
		if documents[i].GroupSlug != documents[j].GroupSlug {
			return documents[i].GroupSlug < documents[j].GroupSlug
		}
		if documents[i].Order != documents[j].Order {
			return documents[i].Order < documents[j].Order
		}
		return documents[i].Title < documents[j].Title
	})

	return documents, nil
}

func (r *EmbeddedRepository) listFromDatabase(ctx context.Context) ([]domain.Document, error) {
	rows, err := r.db.Query(ctx, `
		SELECT slug, title, category, COALESCE(group_slug, ''), COALESCE(group_title, ''), COALESCE(sort_order, 0), word_count, reading_time
		FROM documents
		ORDER BY group_slug ASC, sort_order ASC, title ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	documents, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Document, error) {
		var document domain.Document
		err := row.Scan(&document.Slug, &document.Title, &document.Category, &document.GroupSlug, &document.GroupTitle, &document.Order, &document.WordCount, &document.ReadingTime)
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
	relPath := strings.ReplaceAll(slug, "__", "/") + ".md"
	path := "documents/" + relPath

	raw, err := assets.Documents.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		path = "documents/" + slug + ".md"
		raw, err = assets.Documents.ReadFile(path)
	}
	if errors.Is(err, fs.ErrNotExist) {
		return domain.Document{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Document{}, err
	}
	raw = stripBOM(raw)

	dir := filepath.Dir(relPath)
	var groupSlug, groupTitle string
	var order int

	if dir != "." && dir != "" && !strings.HasPrefix(relPath, slug) {
		groupSlug = strings.ReplaceAll(dir, "\\", "/")
		groupTitle = groupTitleFor(groupSlug)
		order = parseOrder(filepath.Base(relPath))
	}

	return domain.Document{
		Slug:        slug,
		Title:       extractTitle(raw, filepath.Base(path)),
		Category:    categoryFor(relPath),
		GroupSlug:   groupSlug,
		GroupTitle:  groupTitle,
		Order:       order,
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
	case strings.Contains(name, "python"):
		return "backend"
	default:
		return "backend"
	}
}

var leadingDigitsRegex = regexp.MustCompile(`^(\d+)`)

func parseOrder(filename string) int {
	base := strings.ToLower(filename)
	if strings.HasPrefix(base, "readme") {
		return 0
	}
	matches := leadingDigitsRegex.FindStringSubmatch(base)
	if len(matches) > 1 {
		if val, err := strconv.Atoi(matches[1]); err == nil {
			return val
		}
	}
	return 99
}

func groupTitleFor(groupSlug string) string {
	readmePath := "documents/" + groupSlug + "/README.md"
	if raw, err := assets.Documents.ReadFile(readmePath); err == nil {
		title := extractTitle(stripBOM(raw), groupSlug)
		if title != "" && title != groupSlug {
			return title
		}
	}
	parts := strings.Split(groupSlug, "/")
	last := parts[len(parts)-1]
	last = strings.ReplaceAll(last, "_", " ")
	last = strings.ReplaceAll(last, "-", " ")
	return strings.Title(last)
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

func stripBOM(b []byte) []byte {
	if len(b) >= 3 && b[0] == 0xEF && b[1] == 0xBB && b[2] == 0xBF {
		return b[3:]
	}
	return b
}
