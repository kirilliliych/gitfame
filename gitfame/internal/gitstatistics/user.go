package gitstatistics

type UserStatistics struct {
	Name    string `json:"name"`
	Lines   int    `json:"lines"`
	Commits int    `json:"commits"`
	Files   int    `json:"files"`
}

type userStatisticsCommitsSet struct {
	lines   int
	commits map[string]bool
	files   int
}
