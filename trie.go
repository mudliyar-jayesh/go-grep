package main

type TrieNode struct {
	children    map[rune]*TrieNode
	isEndOfWord bool
	files       map[string]struct{}
}

func NewTN() *TrieNode {
	return &TrieNode{
		children: make(map[rune]*TrieNode),
		files:    make(map[string]struct{}),
	}
}

type Trie struct {
	root *TrieNode
}

func NewTrie() *Trie {
	return &Trie{
		root: NewTN(),
	}
}

func (t *Trie) Insert(word, filePath string) {
	current := t.root
	for _, char := range word {
		if _, ok := current.children[char]; !ok {
			current.children[char] = NewTN()
		}
		current = current.children[char]
	}

	current.isEndOfWord = true
	current.files[filePath] = struct{}{}
}

func (t *Trie) Search(prefix string) []string {
	current := t.root
	for _, char := range prefix {
		if node, ok := current.children[char]; ok {
			current = node
			continue
		}
		return nil
	}
	filesFound := make(map[string]struct{})
	t.collectFiles(current, filesFound)

	var results []string
	for file := range filesFound {
		results = append(results, file)
	}
	return results
}

func (t *Trie) collectFiles(node *TrieNode, files map[string]struct{}) {
	if node.isEndOfWord {
		for file := range node.files {
			files[file] = struct{}{}
		}
	}

	for _, childNode := range node.children {
		t.collectFiles(childNode, files)
	}
}
