package document

type Document struct {
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Category    string `json:"category"`
	WordCount   int    `json:"wordCount"`
	ReadingTime int    `json:"readingTime"`
	Content     string `json:"content,omitempty"`
}
