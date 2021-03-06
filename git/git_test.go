package git

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	TESTDIR     = "testDir"
	TESTFILE    = "testFile"
	TESTMESSAGE = "testMessage"
	TESTDATE    = "2005-04-07T22:13:13"
)

func cleanup(directory string, t *testing.T) {
	if err := os.RemoveAll(directory); err != nil {
		t.Errorf("Failed to remove the directory (%s) during cleanup", directory)
	}
}

func TestCreation(t *testing.T) {
	testDir := getTestDir()
	git := Repo{testDir}
	if git.Location != testDir {
		t.Errorf("git.Location == %s; wanted %s", git.Location, testDir)
	}
}

func TestInit(t *testing.T) {
	testDir := getTestDir()
	git := Repo{testDir}
	git.Init()
	if !exists(filepath.Join(testDir, ".git")) {
		t.Error("After git init, there should be a .git dir.")
	}
	cleanup(testDir, t)
}

func TestAdd(t *testing.T) {
	testDir := getTestDir()
	git := Repo{testDir}
	git.Init()
	git.Add(createTestFile())
	if !exists(filepath.Join(testDir, ".git", "index")) {
		t.Error("After the first git add, there should be an index file.")
	}
	cleanup(testDir, t)
}

func TestCommit(t *testing.T) {
	testDir := getTestDir()
	git := Repo{testDir}
	git.Init()
	testFile := createTestFile()
	git.Add(testFile)
	git.Commit(TESTMESSAGE, TESTDATE)
	log := filepath.Join(testDir, ".git", "logs", "refs", "heads", "master")

	if !containsMessage(log) {
		t.Error("After commiting the commit message should be in the logs.")
	}
	cleanup(testDir, t)
}

func containsMessage(logPath string) bool {
	logFile, _ := os.Open(logPath)
	defer logFile.Close()
	logScanner := bufio.NewScanner(logFile)
	for logScanner.Scan() {
		firstLine := logScanner.Text()
		if strings.Contains(firstLine, TESTMESSAGE) {
			return true
		}
		break
	}
	return false
}

func exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func createTestFile() string {
	testFile := filepath.Join(getTestDir(), TESTFILE)
	file, _ := os.Create(testFile)
	file.Close()
	return testFile
}

func getTestDir() string {
	currentDir, _ := os.Getwd()
	return filepath.Join(currentDir, TESTDIR)
}
