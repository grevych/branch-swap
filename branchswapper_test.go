// go:build brnchswppr_test
package branchswapper

import (
	"bufio"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockVCS struct {
	currentBranch string
	branches      map[string]struct{}
	err           error
}

func (m *mockVCS) GetLocalBranches() (map[string]struct{}, error) {
	return m.branches, m.err
}

func (m *mockVCS) CheckoutBranch(branch string) error {
	m.currentBranch = branch
	return m.err
}

func (m *mockVCS) GetCurrentBranch() (string, error) {
	return m.currentBranch, m.err
}

type mockVCSWithErrorOnGetCurrentBranch struct {
	*mockVCS
}

func (m *mockVCSWithErrorOnGetCurrentBranch) GetCurrentBranch() (string, error) {
	return "", fmt.Errorf("current branch error")
}

type mockVCSWithErrorOnCheckout struct {
	*mockVCS
}

func (m *mockVCSWithErrorOnCheckout) CheckoutBranch(branch string) error {
	return fmt.Errorf("checkout error")
}

func setupTestFile(t *testing.T) (string, func(t *testing.T)) {
	filename := fmt.Sprintf(".brnchswppr-%d.test", time.Now().UnixNano())
	_, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Failed to create file %s: %v", filename, err)
	}

	return filename, func(t *testing.T) {
		if err := os.Remove(filename); err != nil {
			t.Errorf("Failed to remove file %s: %v", filename, err)
		}
	}
}

func setupTestBranches(t *testing.T, filename string, branches []string) {
	file, err := os.OpenFile(filename, os.O_WRONLY, 0600)
	if err != nil {
		t.Fatalf("Failed to open file %s: %v", filename, err)
	}

	for _, branch := range branches {
		if _, err := file.WriteString(branch + "\n"); err != nil {
			t.Fatalf("Failed to write to file: %v", err)
		}
	}

	file.Close()
}

func TestLoad(t *testing.T) {
	filename, teardown := setupTestFile(t)
	defer teardown(t)

	setupTestBranches(t, filename, []string{"foo", "bar"})

	vcs := &mockVCS{branches: map[string]struct{}{"foo": {}, "bar": {}}}
	bs := NewBranchSwapperWithFilename(vcs, filename)

	err := bs.Load()

	assert.Nil(t, err)
	assert.Equal(t, 2, len(bs.stack))
	assert.Equal(t, "foo", bs.stack[0])
	assert.Equal(t, "bar", bs.stack[1])
}

func TestLoad_WithBranchesNotInVCS(t *testing.T) {
	filename, teardown := setupTestFile(t)
	defer teardown(t)

	setupTestBranches(t, filename, []string{"foo", "bar", "baz"})

	vcs := &mockVCS{branches: map[string]struct{}{"foo": {}, "bar": {}}}
	bs := NewBranchSwapperWithFilename(vcs, filename)

	err := bs.Load()

	assert.Nil(t, err)
	assert.Equal(t, 2, len(bs.stack))
	assert.Equal(t, "foo", bs.stack[0])
	assert.Equal(t, "bar", bs.stack[1])
}

func TestLoad_WithEmptyFile(t *testing.T) {
	filename, teardown := setupTestFile(t)
	defer teardown(t)

	vcs := &mockVCS{branches: map[string]struct{}{"foo": {}, "bar": {}}}
	bs := NewBranchSwapperWithFilename(vcs, filename)

	err := bs.Load()

	assert.Nil(t, err)
	assert.Equal(t, 0, len(bs.stack))
}

func TestUnload(t *testing.T) {
	filename, teardown := setupTestFile(t)
	defer teardown(t)

	bs := NewBranchSwapperWithFilename(nil, filename)
	bs.stack = []string{"foo", "bar"}

	err := bs.Unload()

	assert.Nil(t, err)

	file, err := os.OpenFile(filename, os.O_RDONLY, 0600)
	if err != nil {
		t.Fatalf("Failed to open file %s: %v", filename, err)
	}

	branches := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		branches = append(branches, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		t.Errorf("Failed to read file: %v", err)
	}

	assert.Equal(t, 2, len(branches))
	assert.Equal(t, "foo", branches[0])
	assert.Equal(t, "bar", branches[1])
}

func TestUnload_TruncatesFileBeforeWriting(t *testing.T) {
	filename, teardown := setupTestFile(t)
	defer teardown(t)

	setupTestBranches(t, filename, []string{"foo", "bar", "baz"})

	bs := NewBranchSwapperWithFilename(nil, filename)
	bs.stack = []string{"foo", "bar"}

	err := bs.Unload()

	assert.Nil(t, err)

	file, err := os.OpenFile(filename, os.O_RDONLY, 0600)
	if err != nil {
		t.Fatalf("Failed to open file %s: %v", filename, err)
	}

	branches := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		branches = append(branches, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		t.Errorf("Failed to read file: %v", err)
	}

	assert.Equal(t, 2, len(branches))
	assert.Equal(t, "foo", branches[0])
	assert.Equal(t, "bar", branches[1])
}

func TestGetStack(t *testing.T) {
	bs := NewBranchSwapper(nil)
	bs.stack = []string{"foo", "bar"}

	stack := bs.GetStack()

	assert.Equal(t, 2, len(stack))
	assert.Equal(t, "foo", stack[0])
	assert.Equal(t, "bar", stack[1])
}

func TestGetStack_ReturnsCopyOfStack(t *testing.T) {
	bs := NewBranchSwapper(nil)
	bs.stack = []string{"foo", "bar"}

	copy := bs.GetStack()

	bs.stack[0] += "baz"
	bs.stack[1] += "baz"

	assert.Equal(t, 2, len(copy))
	assert.Equal(t, "foo", copy[0])
	assert.Equal(t, "bar", copy[1])
}

func TestSwap(t *testing.T) {
	filename, teardown := setupTestFile(t)
	defer teardown(t)

	setupTestBranches(t, filename, []string{"bar"})

	vcs := &mockVCS{
		branches:      map[string]struct{}{"foo": {}, "bar": {}, "baz": {}},
		currentBranch: "foo",
	}
	bs := NewBranchSwapperWithFilename(vcs, filename)
	bs.stack = []string{"bar"}

	err := bs.Swap("baz")

	assert.Nil(t, err)
	assert.Equal(t, 2, len(bs.stack))
	assert.Equal(t, "bar", bs.stack[0])
	assert.Equal(t, "foo", bs.stack[1])
}

func TestSwap_UnloadsFile(t *testing.T) {
	filename, teardown := setupTestFile(t)
	defer teardown(t)

	setupTestBranches(t, filename, []string{"bar"})

	vcs := &mockVCS{
		branches:      map[string]struct{}{"foo": {}, "bar": {}, "baz": {}},
		currentBranch: "foo",
	}
	bs := NewBranchSwapperWithFilename(vcs, filename)
	bs.stack = []string{"bar"}

	err := bs.Swap("baz")

	assert.Nil(t, err)

	file, err := os.OpenFile(filename, os.O_RDONLY, 0600)
	if err != nil {
		t.Fatalf("Failed to open file %s: %v", filename, err)
	}

	branches := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		branches = append(branches, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		t.Errorf("Failed to read file: %v", err)
	}

	assert.Equal(t, 2, len(branches))
	assert.Equal(t, "bar", branches[0])
	assert.Equal(t, "foo", branches[1])
}

func TestSwap_WithCurrentBranchAlreadyInFile(t *testing.T) {
	filename, teardown := setupTestFile(t)
	defer teardown(t)

	setupTestBranches(t, filename, []string{"foo", "bar"})

	vcs := &mockVCS{
		branches:      map[string]struct{}{"foo": {}, "bar": {}, "baz": {}},
		currentBranch: "bar",
	}
	bs := NewBranchSwapperWithFilename(vcs, filename)
	bs.stack = []string{"foo", "bar"}

	err := bs.Swap("baz")

	assert.Nil(t, err)
	assert.Equal(t, 2, len(bs.stack))
	assert.Equal(t, "foo", bs.stack[0])
	assert.Equal(t, "bar", bs.stack[1])
}

// Current branch is not in VCS
func TestSwap_WithErrorOnGetCurrentBranch(t *testing.T) {
	filename, teardown := setupTestFile(t)
	defer teardown(t)

	setupTestBranches(t, filename, []string{"foo"})

	vcs := &mockVCSWithErrorOnCheckout{
		mockVCS: &mockVCS{
			branches:      map[string]struct{}{"foo": {}, "bar": {}},
			currentBranch: "bar",
		},
	}
	bs := NewBranchSwapperWithFilename(vcs, filename)
	bs.stack = []string{"foo"}

	err := bs.Swap("baz")

	assert.ErrorContains(t, err, "checkout error")
}

// Given branch not in VCS
func TestSwap_WithErrorOnCheckout(t *testing.T) {
	filename, teardown := setupTestFile(t)
	defer teardown(t)

	setupTestBranches(t, filename, []string{"foo"})

	vcs := &mockVCSWithErrorOnGetCurrentBranch{
		mockVCS: &mockVCS{
			branches:      map[string]struct{}{"foo": {}, "bar": {}},
			currentBranch: "bar",
		},
	}
	bs := NewBranchSwapperWithFilename(vcs, filename)
	bs.stack = []string{"foo"}

	err := bs.Swap("baz")

	assert.ErrorContains(t, err, "current branch error")
}

func TestSwap_WithGivenBranchNotInFile(t *testing.T) {
	filename, teardown := setupTestFile(t)
	defer teardown(t)

	setupTestBranches(t, filename, []string{"foo"})

	vcs := &mockVCS{
		branches:      map[string]struct{}{"foo": {}, "bar": {}},
		currentBranch: "foo",
	}
	bs := NewBranchSwapperWithFilename(vcs, filename)
	bs.stack = []string{"foo"}

	err := bs.Swap("bar")

	assert.Nil(t, err)
	assert.Equal(t, 1, len(bs.stack))
	assert.Equal(t, "foo", bs.stack[0])
}
