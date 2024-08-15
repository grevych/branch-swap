package branchswapper

import (
	"bufio"
	"errors"
	"os"

	"github.com/urfave/cli/v2"
)

const Filename = ".branchswap"

type VCS interface {
	GetCurrentBranch() (string, error)
	CheckoutBranch(branch string) error
	GetLocalBranches() (map[string]struct{}, error)
}

type BranchSwapper struct {
	vcs      VCS
	filename string
	stack    []string
}

func NewBranchSwapper(vcs VCS) *BranchSwapper {
	return &BranchSwapper{
		vcs:      vcs,
		filename: Filename,
		stack:    make([]string, 0),
	}
}

// NewBranchSwapperWithFilename
func NewBranchSwapperWithFilename(vcs VCS, filename string) *BranchSwapper {
	s := NewBranchSwapper(vcs)
	s.filename = filename
	return s
}

func (s *BranchSwapper) Load() error {
	var file *os.File
	file, err := s.openFile()
	if err != nil {
		return err
	}
	defer file.Close()
	if err := s.loadBranches(file); err != nil {
		return err
	}

	return nil
}

func (s *BranchSwapper) Unload() error {
	var file *os.File
	file, err := s.openFile()
	if err != nil {
		return err
	}
	defer file.Close()
	if err := s.unloadBranches(file); err != nil {
		return err
	}

	return nil
}

func (s *BranchSwapper) GetStack() []string {
	stack := make([]string, len(s.stack))
	copy(stack, s.stack)
	return stack
}

func (s *BranchSwapper) Swap(branchName string) error {
	if err := s.stashCurrentBranch(); err != nil {
		return err
	}

	if branchName == "" {
		return nil
	}

	if err := s.moveToBranch(branchName); err != nil {
		return cli.Exit(err.Error(), 1)
	}
	return nil
}

func (s *BranchSwapper) SwapFromStack(branchIndex int) error {
	if err := s.stashCurrentBranch(); err != nil {
		return err
	}

	if err := s.moveToBranchWithIndex(branchIndex); err != nil {
		return cli.Exit(err.Error(), 1)
	}
	return nil
}

func (s *BranchSwapper) stashCurrentBranch() error {
	currentBranch, err := s.vcs.GetCurrentBranch()
	if err != nil {
		return err
	}
	for _, branch := range s.stack {
		if branch == currentBranch {
			return nil
		}
	}
	s.stack = append(s.stack, currentBranch)

	// write stack back to file
	if err := s.Unload(); err != nil {
		return err
	}
	return nil
}

func (s *BranchSwapper) moveToBranchWithIndex(index int) error {
	// handle error
	if index < 0 || index >= len(s.stack) {
		return errors.New("index out of range")
	}
	branch := s.stack[index]
	if err := s.vcs.CheckoutBranch(branch); err != nil {
		return err
	}
	s.stack = append(s.stack[:index], s.stack[index+1:]...)
	// write stack back to file
	if err := s.Unload(); err != nil {
		return err
	}
	return nil
}

func (s *BranchSwapper) moveToBranch(branch string) error {
	// handle error
	if err := s.vcs.CheckoutBranch(branch); err != nil {
		return err
	}

	for i, b := range s.stack {
		if b == branch {
			s.stack = append(s.stack[:i], s.stack[i+1:]...)
			// write stack back to file
			if err := s.Unload(); err != nil {
				return err
			}
			break
		}
	}

	return nil
}

// loadBranches loads branches from filename. If it
// doesn't exist, it creates it.
func (s *BranchSwapper) loadBranches(f *os.File) error {
	localBranches, err := s.vcs.GetLocalBranches()
	if err != nil {
		return err
	}
	// Continue with the rest of the code...
	branches := make(map[string]int, 0)
	scanner := bufio.NewScanner(f)
	for i := 0; scanner.Scan(); i++ {
		// Validate branch against git
		line := scanner.Text()
		if _, ok := localBranches[line]; !ok {
			// log invalid branch
			continue
		}
		// Process each line here
		if _, ok := branches[line]; !ok {
			branches[line] = i
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	s.stack = make([]string, len(branches))
	for branch, index := range branches {
		s.stack[index] = branch
	}

	// rewrite file after reading and validating
	s.unloadBranches(f)
	return nil
}

func (s *BranchSwapper) unloadBranches(f *os.File) error {
	f.Truncate(0)
	f.Seek(0, 0)
	for _, branch := range s.stack {
		if _, err := f.WriteString(branch + "\n"); err != nil {
			return err
		}
	}
	return nil
}

func (s *BranchSwapper) openFile() (*os.File, error) {
	file, err := os.OpenFile(s.filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	return file, nil
}
