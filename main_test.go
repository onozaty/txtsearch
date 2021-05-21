package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestConvert_match_match(t *testing.T) {

	s := "word1"

	f, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(f.Name())

	result, err := match(f.Name(), []string{"word1"}, false)
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if !result {
		t.Fatal("failed test\n", result)
	}
}

func TestConvert_match_unmatch(t *testing.T) {

	s := "word1"

	f, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(f.Name())

	result, err := match(f.Name(), []string{"word2"}, false)
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if result {
		t.Fatal("failed test\n", result)
	}
}

func TestConvert_match_match_multiline(t *testing.T) {

	s := `word1
word2 word3
word4
`

	f, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(f.Name())

	result, err := match(f.Name(), []string{"word4"}, false)
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if !result {
		t.Fatal("failed test\n", result)
	}
}

func TestConvert_match_match_multiword(t *testing.T) {

	s := `word1
word2 word3
word4
`

	f, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(f.Name())

	result, err := match(f.Name(), []string{"wordx", "word2"}, false)
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if !result {
		t.Fatal("failed test\n", result)
	}
}

func TestConvert_match_match_casenotsensitive(t *testing.T) {

	s := `word1
Word2 word3
word4
`

	f, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(f.Name())

	result, err := match(f.Name(), []string{"worD2"}, false)
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if !result {
		t.Fatal("failed test\n", result)
	}
}

func TestConvert_match_match_casesensitive(t *testing.T) {

	s := `word1
Word2 word3
word4
`

	f, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(f.Name())

	result, err := match(f.Name(), []string{"worD2"}, true)
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if result {
		t.Fatal("failed test\n", result)
	}
}

func TestConvert_match_looongtext(t *testing.T) {

	// bufio.Scanner だと token too long が出るくらいの長さで
	s := strings.Repeat("x", 100000) + "word1"

	f, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(f.Name())

	result, err := match(f.Name(), []string{"word1"}, false)
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if !result {
		t.Fatal("failed test\n", result)
	}
}

func TestConvert_search(t *testing.T) {

	dir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.RemoveAll(dir)

	file1 := filepath.Join(dir, "1")
	err = os.WriteFile(file1, []byte("word1"), 0666)
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	file2 := filepath.Join(dir, "2")
	err = os.WriteFile(file2, []byte("word2"), 0666)
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	file3 := filepath.Join(dir, "3")
	err = os.WriteFile(file3, []byte("word3"), 0666)
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result, err := search(dir, []string{"word1", "WORD3"}, false)
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	expect := SearchResult{
		totalCount: 3,
		matchFiles: []string{
			filepath.Join(dir, "1"),
			filepath.Join(dir, "3"),
		},
	}

	if !reflect.DeepEqual(result, expect) {
		t.Fatal("failed test\n", result)
	}
}

func TestConvert_copyFiles(t *testing.T) {

	inDir, err := os.MkdirTemp("", "input")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.RemoveAll(inDir)

	file1 := filepath.Join(inDir, "1")
	err = os.WriteFile(file1, []byte("word1"), 0666)
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	file2 := filepath.Join(inDir, "2")
	err = os.WriteFile(file2, []byte("word2"), 0666)
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	file3 := filepath.Join(inDir, "3")
	err = os.WriteFile(file3, []byte("word3"), 0666)
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	outDir, err := os.MkdirTemp("", "ouput")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.RemoveAll(outDir)

	err = copyFiles([]string{file2, file3}, outDir)
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	// コピー先のフォルダ確認
	fileInfos, err := ioutil.ReadDir(outDir)
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if len(fileInfos) != 2 {
		t.Fatal("failed test\n", fileInfos)
	}
	for _, fileInfo := range fileInfos {
		if fileInfo.Size() == 0 {
			// 空じゃないことだけチェック
			t.Fatal("failed test\n", fileInfo)
		}
	}
}

func createTempFile(content string) (*os.File, error) {

	tempFile, err := os.CreateTemp("", "txt")
	if err != nil {
		return nil, err
	}

	_, err = tempFile.Write([]byte(content))
	if err != nil {
		return nil, err
	}

	return tempFile, nil
}
