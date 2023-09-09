package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/dhowden/tag"
	"github.com/mayudev/dedupthingy/util"
)

type Result struct {
	res util.Metadata
	err error
}

var scannedFiles []string

func checkPaths(paths []string) error {
	for _, path := range paths {
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			return errors.New("invalid path: " + path)
		}
	}
	return nil
}

func walk(s string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if !d.IsDir() {
		scannedFiles = append(scannedFiles, s)
	}

	return nil
}

func scanFile(file string, ch chan<- Result) {

	f, err := os.Open(file)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\033[2K\rOn %s", file)
	m, err := tag.ReadFrom(f)

	if err == nil {
		ma := util.Metadata{
			Filename: file,
			Title:    m.Title(),
			Album:    m.Album(),
			Artist:   m.Artist(),
			Year:     m.Year(),
		}

		ch <- Result{err: nil, res: ma}

	} else {
		ch <- Result{err: errors.New("failed to read metadata"), res: util.Metadata{}}
	}

}

func findDuplicates(results []util.Metadata) {
	hashes := make(map[util.Metadata]string)
	comparator := util.NewComparator(MatchBy, CaseSensitive)

	f, err := os.Create("output.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	for _, v := range results {
		c := comparator.CreateComparator(v)

		if val, exists := hashes[c]; exists {
			fragment := fmt.Sprintf("%s\n%s\n", val, v.Filename)
			f.WriteString(fragment)
			fmt.Println("\033[1m\033[93m-> Duplicate found!\033[0m\n" + fragment)
		} else {
			hashes[c] = v.Filename
		}

	}

}

func runDeduplicate(paths []string) error {
	sem := make(chan struct{}, 200)

	err := checkPaths(paths)
	if err != nil {
		return err
	}

	for _, path := range paths {
		filepath.WalkDir(path, walk)
	}

	ch := make(chan Result)

	for _, v := range scannedFiles {
		go func(v string) {
			sem <- struct{}{}
			defer func() { <-sem }()
			scanFile(v, ch)
		}(v)
	}

	results := []util.Metadata{}

	for range scannedFiles {
		a := <-ch

		if a.err == nil {
			results = append(results, a.res)
		}
	}

	close(ch)

	fmt.Println()

	findDuplicates(results)

	return nil
}
