package main

import (
	"os"
	"strconv"
	"testing"
)

// TestParseGitRepository ...
// doc: https://www.npmjs.com/package/github-username-regex
func TestParseGitRepository(t *testing.T) {

	devenv, _ := strconv.ParseBool(os.Getenv("DEVENV"))

	// Test1: expected result SUCCESS
	testData1 := map[int]string{1: "apache", 2: ""}
	_, err := ParseGitRepository(testData1[1], testData1[2], devenv)
	if err != nil {
		t.Error(testData1)
		t.Errorf("%+v: SHOULD BE VALID. Program not agree....\nerr:%+v", testData1, err)

	}

	// Test2: expected result FAIL
	testData2 := map[int]string{1: "apache", 2: ""}
	_, err = ParseGitRepository(testData2[1], testData2[2], devenv)
	if err != nil {
		t.Errorf("%+v: not valid argument\nerr:%+v", testData2, err)
	}

	// Test3: expected result FAIL
	testData3 := map[int]string{1: "apache", 2: "kafka/hello"}
	_, err = ParseGitRepository(testData3[1], testData3[2], devenv)
	if err != nil {
		t.Errorf("%+v: not valid argument\nerr:%+v", testData3, err)
	}
}
