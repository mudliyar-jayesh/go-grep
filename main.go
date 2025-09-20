package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type FileData struct {
	Path    string
	Content []byte
}

func findFiles(rootDir string, filePaths chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			filePaths <- path
		}
		return nil
	})
	if err != nil {
		log.Printf("Error walking directory: %v", err)
	}
	close(filePaths)
}

func readFile(filePaths <-chan string, fileContents chan<- FileData, wg *sync.WaitGroup) {
	defer wg.Done()
	for path := range filePaths {
		content, err := os.ReadFile(path)
		if err != nil {
			log.Printf("Warning: Could not read file %s : %v", path, err)
			continue
		}
		fileContents <- FileData{Path: path, Content: content}
	}
}

func indexContent(fileContents <-chan FileData, trie *Trie, wg *sync.WaitGroup) {
	defer wg.Done()
	for data := range fileContents {
		scanner := bufio.NewScanner(strings.NewReader(string(data.Content)))
		scanner.Split(bufio.ScanWords)

		for scanner.Scan() {
			word := strings.ToLower(scanner.Text())
			trie.Insert(word, data.Path)
		}
	}
}

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usuage: %s <diretory> <search-term>", os.Args[0])
	}

	rootDir := os.Args[1]
	searchTerm := strings.ToLower(os.Args[2])

	var wg sync.WaitGroup
	filePaths := make(chan string)
	fileContents := make(chan FileData)

	wg.Add(1)
	go findFiles(rootDir, filePaths, &wg)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go readFile(filePaths, fileContents, &wg)
	}

	var indexerWg sync.WaitGroup
	indexerWg.Add(1)
	trie := NewTrie()
	go indexContent(fileContents, trie, &indexerWg)

	wg.Wait()
	close(fileContents)
	indexerWg.Wait()

	log.Println("Index built successfully")

	results := trie.Search(searchTerm)
	if len(results) == 0 {
		fmt.Printf("\nNo results found for '%s'.\n", searchTerm)
		return
	}

	fmt.Printf("\nFound '%s' in the following files:\n", searchTerm)
	for _, path := range results {
		fmt.Printf("- %s\n", path)
	}
}
