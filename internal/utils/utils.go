package utils

import "github.com/sahilm/fuzzy"

func FilterMinimumScoreMatches(matches []fuzzy.Match, minScore int64) []fuzzy.Match {
	for i, v := range matches {
		if v.Score < int(minScore) {
			return matches[:i]
		}
	}

	return matches
}
