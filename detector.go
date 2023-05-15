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

type Options struct {
	IgnoreCase bool
	Noises     []rune
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

func New(opts Options) *Detector {
	d := new(Detector)

	trie := NewTrie()
	trie.IgnoreCase = opts.IgnoreCase
	if len(opts.Noises) > 0 {
		for _, r := range opts.Noises {
			trie.Noises[r] = struct{}{}
		}
	}

	d.trie = trie

	return d
}
