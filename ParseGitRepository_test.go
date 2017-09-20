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

	// Test4: expected result FAIL
	testData4 := map[int]string{1: "", 2: "kafka/hello"}
	_, err = ParseGitRepository(testData4[1], testData4[2], devenv)
	if err != nil {
		t.Errorf("%+v: not valid argument\nerr:%+v", testData4, err)
	}

	// Test1: expected result SUCCESS
	testData5 := map[int]string{1: "apache", 2: "kafka"}
	responseData, err1 := ParseGitRepository(testData5[1], testData5[2], true)
	if err1 != nil {
		t.Error(testData1)
		t.Errorf("%+v: Wierd stuff.\nerr:%+v", testData1, err)
	}

	if responseData.Repository == "" {
		t.Errorf("%+v: is nil", responseData.Repository)
	}

	if responseData.Owner == "" {
		t.Errorf("%+v: is nil", responseData.Owner)
	}

	if responseData.Committer == "" {
		t.Errorf("%+v: is nil", responseData.Committer)
	}

	// Test6: expected result SUCCESS
	testData5 = map[int]string{1: "apache", 2: "kafka"}
	responseData, err1 = ParseGitRepository(testData5[1], testData5[2], false)
	if err1 != nil {
		t.Error(testData1)
		t.Errorf("%+v: Wierd stuff.\nerr:%+v", testData1, err)
	}

	if responseData.Repository == "" {
		t.Errorf("%+v: is nil", responseData.Repository)
	}

	if responseData.Owner == "nil" {
		t.Errorf("%+v: is nil", responseData.Owner)
	}

	if responseData.Committer == "nil" {
		t.Errorf("%+v: is nil", responseData.Committer)
	}
}
