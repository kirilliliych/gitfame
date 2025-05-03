package gitstatistics

import (
	"errors"
	"gitlab.com/slon/shad-go/gitfame/internal/gitrepository"
	"gitlab.com/slon/shad-go/gitfame/internal/language"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type Handler struct {
	Repository   gitrepository.Handler
	Revision     string
	UseCommitter bool
	Extensions   []string
	Languages    []string
	Exclude      []string
	RestrictTo   []string
}

type fileStatistics struct {
	commits map[string]bool
	lines   int
}

func (handler *Handler) matchesExtensions(fileName string) bool {
	if len(handler.Extensions) == 0 {
		return true
	}

	fileExtension := filepath.Ext(fileName)
	for _, extension := range handler.Extensions {
		if extension == fileExtension {
			return true
		}
	}

	return false
}

func (handler *Handler) matchesLanguages(fileName string) bool {
	if len(handler.Languages) == 0 {
		return true
	}

	fileExtension := filepath.Ext(fileName)
	for _, lang := range handler.Languages {
		languageExtensions := language.GetLanguageExtensions(lang)
		for _, extension := range languageExtensions {
			if extension == fileExtension {
				return true
			}
		}
	}

	return false
}

func (handler *Handler) notMatchesExclude(fileName string) bool {
	if len(handler.Exclude) == 0 {
		return true
	}

	for _, globPattern := range handler.Exclude {
		matches, err := path.Match(globPattern, fileName)
		if (err != nil) || matches {
			return false
		}
	}

	return true
}

func (handler *Handler) matchesRestrictTo(fileName string) bool {
	if len(handler.RestrictTo) == 0 {
		return true
	}

	for _, globPattern := range handler.RestrictTo {
		matches, err := path.Match(globPattern, fileName)
		if err != nil {
			return false
		}
		if matches {
			return true
		}
	}

	return false
}

func (handler *Handler) matchesRequest(fileName string) bool {
	return handler.matchesExtensions(fileName) && handler.matchesLanguages(fileName) && handler.notMatchesExclude(fileName) && handler.matchesRestrictTo(fileName)
}

func (handler *Handler) getLogInfo(fileName string) (map[string]fileStatistics, error) {
	logText, err := handler.Repository.GetLog(fileName, handler.Revision)
	if err != nil {
		return nil, err
	}
	logByLines := strings.Split(logText, "\n")
	result := make(map[string]fileStatistics)
	for _, line := range logByLines {
		fields := strings.Fields(line)
		hashPrefix := fields[0] + " "
		name := strings.TrimPrefix(line, hashPrefix)
		currentCommits := make(map[string]bool)
		currentCommits[fields[0]] = true
		result[name] = fileStatistics{commits: currentCommits}
	}

	return result, nil
}

func (handler *Handler) getFileStatistics(fileName string) (map[string]fileStatistics, error) {
	blameText, err := handler.Repository.GetBlameInfo(fileName, handler.Revision)
	if err != nil {
		return nil, err
	}
	if len(blameText) == 0 {
		return handler.getLogInfo(fileName)
	}

	commitLinesCountByHash := make(map[string]int)
	userByCommitHash := make(map[string]string)
	userPrivileges := "author"
	if handler.UseCommitter {
		userPrivileges = "committer"
	}
	var currentGroupHash string
	blameByLines := strings.Split(blameText, "\n")
	for _, line := range blameByLines {
		fields := strings.Fields(line)
		if len(fields) == 0 || fields[0] == "summary" {
			continue
		}
		const blameMainHeaderStringLength = 4
		if len(fields) == blameMainHeaderStringLength { // hash, line number in orig file, line number in final file, lines count in group
			currentGroupHash = fields[0]
			linesInGroupCount, _ := strconv.Atoi(fields[3])
			commitLinesCountByHash[currentGroupHash] += linesInGroupCount
			continue
		}
		if fields[0] == userPrivileges { // author ... / committer ...
			_, ok := userByCommitHash[currentGroupHash]
			if !ok {
				userByCommitHash[currentGroupHash] = strings.TrimPrefix(line, userPrivileges+" ")
			}
		}
	}

	result := make(map[string]fileStatistics)
	for hash, user := range userByCommitHash {
		curUserStatistics, ok := result[user]
		if !ok {
			curUserStatistics.commits = make(map[string]bool)
		}
		curUserStatistics.commits[hash] = true
		count, ok := commitLinesCountByHash[hash]
		if !ok {
			return nil, errors.New("no available lines count by commit hash")
		}
		curUserStatistics.lines += count
		result[user] = curUserStatistics
	}

	return result, nil
}

func (handler *Handler) GetStatistics() ([]UserStatistics, error) {
	fileNames, err := handler.Repository.GetFileNamesList(handler.Revision)
	if err != nil {
		return nil, err
	}

	userStatistics := make(map[string]userStatisticsCommitsSet)
	for _, fileName := range fileNames {
		if !handler.matchesRequest(fileName) {
			continue
		}
		nameToStatistics, err := handler.getFileStatistics(fileName)
		if err != nil {
			return nil, err
		}
		for name, statistics := range nameToStatistics {
			curUserStatistics, ok := userStatistics[name]
			if !ok {
				curUserStatistics.commits = make(map[string]bool)
			}
			curUserStatistics.lines += statistics.lines
			curUserStatistics.files++
			for commitHash := range statistics.commits {
				curUserStatistics.commits[commitHash] = true
			}
			userStatistics[name] = curUserStatistics
		}
	}

	var result []UserStatistics
	for name, statistics := range userStatistics {
		result = append(result, UserStatistics{Name: name, Lines: statistics.lines, Commits: len(statistics.commits), Files: statistics.files})
	}
	return result, nil
}
