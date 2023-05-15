package detector

import (
	"bufio"
	"io"
	"sync"
)

type Detector struct {
	mu   sync.RWMutex
	trie *Trie
}

func (d *Detector) Load(rd io.Reader) error {
	buf := bufio.NewScanner(bufio.NewReader(rd))
	for buf.Scan() {
		d.trie.Insert(buf.Text())
	}

	return nil
}

func (d *Detector) AddWord(word string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.trie.Insert(word)
}

func (d *Detector) Detect(text string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.trie.Check(text)
}

func (d *Detector) Search(text string) []string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.trie.Search(text)
}

func (d *Detector) Match(text string) []Result {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.trie.Match(text)
}

func (d *Detector) Filter(text string, replace ...string) string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.trie.Filter(text, replace...)
}

func New(opts ...func(detector *Detector)) *Detector {
	d := new(Detector)
	trie := NewTrie()
	d.trie = trie
	for _, opt := range opts {
		opt(d)
	}

	return d
}

func WithIgnoreCase(d *Detector) {
	d.trie.IgnoreCase = true
}

func WithNosies(noises string) func(*Detector) {
	return func(d *Detector) {
		if len(noises) > 0 {
			for _, r := range noises {
				d.trie.Noises[r] = struct{}{}
			}
		}
	}
}
