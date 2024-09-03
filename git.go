package branchswapper

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type GitExecutor struct {
}

func NewGitExecutor() *GitExecutor {
	return &GitExecutor{}
}

func (g *GitExecutor) GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		if err2 := (exec.ExitError{}); errors.Is(err, &err2) {
			return "", fmt.Errorf("%s", err2.Stderr)

		}
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func (g *GitExecutor) CheckoutBranch(branch string) error {
	cmd := exec.Command("git", "checkout", branch)
	err := cmd.Run()
	if err != nil {
		if err2 := (exec.ExitError{}); errors.Is(err, &err2) {
			return fmt.Errorf("%s", err2.Stderr)

		}
		return err
	}
	return nil
}

func (g *GitExecutor) GetLocalBranches() (map[string]struct{}, error) {
	cmd := exec.Command("git", "branch")
	output, err := cmd.Output()
	if err != nil {
		if err2 := (exec.ExitError{}); errors.Is(err, &err2) {
			return nil, fmt.Errorf("%s", err2.Stderr)

		}
		return nil, err
	}
	branches := strings.Split(string(output), "\n")
	result := map[string]struct{}{}
	for _, branch := range branches {
		branch = strings.TrimSpace(branch)
		if branch == "" {
			continue
		}
		if branch[0] == '*' {
			result[branch[2:]] = struct{}{}
		} else {
			result[branch] = struct{}{}
		}
	}
	return result, nil
}
