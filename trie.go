package detector

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type Result struct {
	//UcharStart int // 1-based
	//UcharEnd   int
	ByteStart  int // zero-based
	ByteEnd    int
	HitWord    string
	MatchedStr string
}

type Node struct {
	End      bool
	Children map[rune]*Node
}

type Trie struct {
	Root *Node
}

func (trie *Trie) Insert(s string) {
	if s == "" {
		return
	}

	parent := trie.Root
	for _, r := range s {
		if child, ok := parent.Children[r]; ok {
			parent = child
		} else {
			child := newNode()
			parent.Children[r] = child
			parent = child
		}
	}

	parent.End = true
}

func (trie *Trie) Contains(s string) bool {
	parent := trie.Root
	for _, r := range s {
		if _, ok := parent.Children[r]; ok {
			return true
		}
		parent = trie.Root
	}

	return false
}

func (trie *Trie) Replace(s string, replace string) string {
	var (
		sb      = strings.Builder{}
		parent  = trie.Root
		nrune   int
		changed bool
	)
	for _, r := range s {
		if child, ok := parent.Children[r]; ok {
			nrune++
			if child.End {
				changed = true
				sb.WriteString(strings.Repeat(replace, nrune))
			}
			parent = child
		} else {
			sb.WriteRune(r)
			parent = trie.Root
			nrune = 0
		}
	}

	if changed {
		return sb.String()
	}

	return s
}

func (trie *Trie) Filter(s string) string {
	var (
		sb      = strings.Builder{}
		parent  = trie.Root
		changed bool
	)
	for _, r := range s {
		if child, ok := parent.Children[r]; ok {
			if child.End {
				changed = true
			}
			parent = child
		} else {
			sb.WriteRune(r)
			parent = trie.Root
		}
	}

	if changed {
		return sb.String()
	}

	return s
}

func (trie *Trie) FindFirst(s string) (nth int, word string) {
	var (
		parent      = trie.Root
		illegalWord = make([]rune, 0, 10)
	)

	for _, r := range s {
		nth++
		if child, ok := parent.Children[r]; ok {
			illegalWord = append(illegalWord, r)
			if child.End {
				return nth, string(illegalWord)
			}
			parent = child
		} else {
			parent = trie.Root
		}
	}

	return 0, ""
}

func (trie *Trie) FindAll2(text string) []string {
	var (
		parent = trie.Root
		ret    = make([]string, 0)
		start  = 0
		pos    = 0
		byts   = s2b(text)
		l      = len(byts)
	)

	for pos < l {
		r, size := utf8.DecodeRune(byts[pos:])
		if child, ok := parent.Children[r]; ok {
			if child.End {
				ret = append(ret, string(byts[start:pos+size]))
			}
			parent = child
			pos += size
		} else {
			parent = trie.Root
			start += size
			pos = start
		}
	}

	return ret
}

func (trie *Trie) FindAll(text string) []string {
	results := trie.Match(text)
	fmt.Println(results)
	m := make(map[string]struct{}, len(results))
	s2 := make([]string, 0, len(results))
	for _, res := range results {
		if _, ok := m[res.HitWord]; !ok {
			m[res.HitWord] = struct{}{}
			s2 = append(s2, res.HitWord)
		}
	}

	return s2
}

func (trie *Trie) Match(text string) (results []Result) {
	var (
		parent = trie.Root
		b      = strings.Builder{}
		start  = 0
		pos    = 0
		rstart = 0 // rune start pos
		rend   = 0 // rune end pos
		byts   = s2b(text)
		l      = len(byts)
		e      bool
	)

	for pos < l {
		r, size := utf8.DecodeRune(byts[pos:])
		if child, ok := parent.Children[r]; ok {
			b.WriteRune(r)
			rend++
			if child.End {
				result := Result{
					//UcharStart: rstart + 1,
					//UcharEnd:   rstart + rend,
					ByteStart: start,
					ByteEnd:   pos + size,
					HitWord:   b.String(),
				}
				result.MatchedStr = string(byts[result.ByteStart:result.ByteEnd])
				results = append(results, result)
				e = true
			}
			pos += size
			parent = child
		} else {
			parent = trie.Root
			if e {
				pos -= size
				start = pos
				e = false
			} else {
				start += size
				pos = start
				rstart++
			}
			b.Reset()
		}

		//fmt.Println(pos, rstart, string(r))
	}

	return results
}

func newNode() *Node {
	node := new(Node)
	node.Children = make(map[rune]*Node)
	return node
}

func New() *Trie {
	trie := new(Trie)
	trie.Root = newNode()
	return trie
}
