package document

type Document struct {
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Category    string `json:"category"`
	GroupSlug   string `json:"groupSlug,omitempty"`
	GroupTitle  string `json:"groupTitle,omitempty"`
	Order       int    `json:"order"`
	WordCount   int    `json:"wordCount"`
	ReadingTime int    `json:"readingTime"`
	Content     string `json:"content,omitempty"`
}
