package api

type Book struct {
	AuthorName   string   `json:"author_name"`
	Title        string   `json:"title,omitempty"`
	PublishYear  int      `json:"first_publish_year"`
	Price        int      `json:"price"`
	EditionCount int      `json:"edition_count"`
	Language     []string `json:"language"`
	Contributor  string   `json:"contributor"`
}
