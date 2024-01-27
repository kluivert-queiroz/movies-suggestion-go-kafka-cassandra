package entity

type Suggestion struct {
	UserId string            `cql:"user_id"`
	Title  string            `cql:"title"`
	Movies []MovieSuggestion `cql:"movies"`
}

type MovieSuggestion struct {
	ID    string  `cql:"id"`
	Title string  `cql:"title"`
	Score float32 `cql:"score"`
}
