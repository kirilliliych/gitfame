package utility

import (
	"os"
	"sort"

	"gitlab.com/slon/shad-go/gitfame/internal/gitstatistics"
)

func sortLinesPrimarily(statistics *[]gitstatistics.UserStatistics) {
	sort.SliceStable(*statistics, func(i, j int) bool {
		if (*statistics)[i].Lines == (*statistics)[j].Lines {
			if (*statistics)[i].Commits == (*statistics)[j].Commits {
				if (*statistics)[i].Files == (*statistics)[j].Files {
					return (*statistics)[i].Name < (*statistics)[j].Name
				} else {
					return (*statistics)[i].Files > (*statistics)[j].Files
				}
			} else {
				return (*statistics)[i].Commits > (*statistics)[j].Commits
			}
		} else {
			return (*statistics)[i].Lines > (*statistics)[j].Lines
		}
	})
}

func sortCommitsPrimarily(statistics *[]gitstatistics.UserStatistics) {
	sort.SliceStable(*statistics, func(i, j int) bool {
		if (*statistics)[i].Commits == (*statistics)[j].Commits {
			if (*statistics)[i].Lines == (*statistics)[j].Lines {
				if (*statistics)[i].Files == (*statistics)[j].Files {
					return (*statistics)[i].Name < (*statistics)[j].Name
				} else {
					return (*statistics)[i].Files > (*statistics)[j].Files
				}
			} else {
				return (*statistics)[i].Lines > (*statistics)[j].Lines
			}
		} else {
			return (*statistics)[i].Commits > (*statistics)[j].Commits
		}
	})
}

func sortFilesPrimarily(statistics *[]gitstatistics.UserStatistics) {
	sort.SliceStable(*statistics, func(i, j int) bool {
		if (*statistics)[i].Files == (*statistics)[j].Files {
			if (*statistics)[i].Lines == (*statistics)[j].Lines {
				if (*statistics)[i].Commits == (*statistics)[j].Commits {
					return (*statistics)[i].Name < (*statistics)[j].Name
				} else {
					return (*statistics)[i].Commits > (*statistics)[j].Commits
				}
			} else {
				return (*statistics)[i].Lines > (*statistics)[j].Lines
			}
		} else {
			return (*statistics)[i].Files > (*statistics)[j].Files
		}
	})
}

func SortUserInfos(statistics []gitstatistics.UserStatistics, orderBy string) []gitstatistics.UserStatistics {
	result := statistics

	switch orderBy {
	case "lines":
		sortLinesPrimarily(&result)
	case "commits":
		sortCommitsPrimarily(&result)
	case "files":
		sortFilesPrimarily(&result)
	default:
		os.Exit(1)
	}

	return result
}
