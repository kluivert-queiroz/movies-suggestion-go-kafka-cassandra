package entity

type Movie struct {
	ID            string   `cql:"id"`
	Title         string   `cql:"title"`
	OriginalTitle string   `cql:"original_title"`
	TitleType     string   `cql:"title_type"`
	StartYear     int      `cql:"start_year"`
	Genres        []string `cql:"genres"`
	Score         float32  `cql:"score"`
}
