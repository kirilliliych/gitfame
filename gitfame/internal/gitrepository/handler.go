package gitrepository

import (
	"os/exec"
	"strings"
)

type Handler string

func (handler *Handler) GetFileNamesList(revision string) ([]string, error) {
	cmd := exec.Command("git", "ls-tree", "-r", "--name-only", revision)
	cmd.Dir = string(*handler)
	fileNames, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	if len(fileNames) == 0 {
		return []string{}, nil
	}
	return strings.Split(strings.TrimSpace(string(fileNames)), "\n"), nil
}

func (handler *Handler) GetBlameInfo(fileName string, revision string) (string, error) {
	cmd := exec.Command("git", "blame", fileName, "--porcelain", revision)
	cmd.Dir = string(*handler)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func (handler *Handler) GetLog(fileName string, revision string) (string, error) {
	cmd := exec.Command("git", "log", revision, "-1", "--pretty=format:%H %an", "--", fileName)
	cmd.Dir = string(*handler)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
