package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	flag "github.com/spf13/pflag"
)

var (
	version = "dev"
	commit  = "none"
)

func main() {

	if len(commit) > 7 {
		commit = commit[:7]
	}

	var help bool
	var inputDir string
	var outputDir string
	var caseSensitive bool

	flag.StringVarP(&inputDir, "input", "i", "", "Input directory. Specify the directory where the text files to be searched are located.")
	flag.StringVarP(&outputDir, "output", "o", "", "(optional) The directory to output text files matched by the search.")
	flag.BoolVarP(&caseSensitive, "cs", "c", false, "(optional) The search is case-sensitive. By default, it is not case-sensitive.")
	flag.BoolVarP(&help, "help", "h", false, "Help.")
	flag.Parse()
	flag.CommandLine.SortFlags = false
	flag.Usage = func() {
		fmt.Printf("txtsearch v%s (%s)\n", version, commit)
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] WORDS...\nflags\n", os.Args[0])
		flag.PrintDefaults()
	}

	if help {
		flag.Usage()
		os.Exit(0)
	}

	words := flag.Args()

	if len(words) == 0 || inputDir == "" {
		flag.Usage()
		os.Exit(1)
	}

	result, err := search(inputDir, words, caseSensitive)

	if err != nil {
		fmt.Println("\nError: ", err)
		os.Exit(1)
	}

	fmt.Printf("total: %d  match: %d\n", result.totalCount, len(result.matchFiles))

	if outputDir != "" {
		// 出力ディレクトリが指定されていた場合
		// -> 一致したファイルを出力
		_, err := os.Stat(outputDir)
		if os.IsNotExist(err) {
			os.Mkdir(outputDir, 0777)
		} else {
			fmt.Println("\nError: ", err)
			os.Exit(1)
		}

		err = copyFiles(result.matchFiles, outputDir)
		if err != nil {
			fmt.Println("\nError: ", err)
			os.Exit(1)
		}

		fmt.Printf("copy the matched files to %s\n", outputDir)
	} else {
		for _, matchFile := range result.matchFiles {
			fmt.Println(matchFile)
		}
	}

	if err != nil {
		fmt.Println("\nError: ", err)
		os.Exit(1)
	}
}

type SearchResult struct {
	totalCount int
	matchFiles []string
}

func search(inputDir string, words []string, caseSensitive bool) (SearchResult, error) {

	filePaths, err := getFiles(inputDir)
	if err != nil {
		return SearchResult{}, err
	}

	matchFiles := []string{}
	for _, filePath := range filePaths {

		match, err := match(filePath, words, caseSensitive)
		if err != nil {
			return SearchResult{}, err
		}
		if match {
			matchFiles = append(matchFiles, filePath)
		}
	}

	return SearchResult{
		totalCount: len(filePaths),
		matchFiles: matchFiles,
	}, nil
}

func match(path string, words []string, caseSensitive bool) (bool, error) {

	if !caseSensitive {
		for i, word := range words {
			words[i] = strings.ToLower(word)
		}
	}

	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	r := bufio.NewReader(f)

	eof := false
	for !eof {

		line, err := r.ReadString('n')
		if err == io.EOF {
			// EOFの場合でも末尾までの文字が返ってくるのでフラグ立ててそのまま継続
			eof = true
		} else if err != nil {
			return false, err
		}

		if !caseSensitive {
			line = strings.ToLower(line)
		}

		for _, word := range words {
			if strings.Contains(line, word) {
				return true, nil
			}
		}
	}

	return false, nil
}

func getFiles(baseDir string) ([]string, error) {

	fileInfos, err := ioutil.ReadDir(baseDir)
	if err != nil {
		return nil, err
	}

	var filePaths []string
	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			filePaths = append(filePaths, filepath.Join(baseDir, fileInfo.Name()))
		}
	}

	sort.Strings(filePaths)
	return filePaths, nil
}

func copyFiles(srcFiles []string, outputDir string) error {

	for _, srcFile := range srcFiles {
		err := copyFile(srcFile, filepath.Join(outputDir, filepath.Base(srcFile)))
		if err != nil {
			return err
		}
	}

	return nil
}

func copyFile(src string, dest string) error {

	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()

	df, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer df.Close()

	_, err = io.Copy(sf, df)
	return err
}
