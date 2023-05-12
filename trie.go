package detector

import (
	"strings"
	"unicode/utf8"
)

type Result struct {
	CharStart  int // 1-based, left closed right closed
	CharEnd    int
	ByteStart  int // zero-based, left closed right open
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
		if child, ok := parent.Children[r]; ok {
			if child.End {
				return true
			}
			parent = child
		} else {
			parent = trie.Root
		}
	}

	return false
}

func (trie *Trie) Filter(text string, replace ...string) string {
	results := trie.Match(text)
	if len(results) == 0 {
		return text
	}

	var (
		sb     = strings.Builder{}
		l      = len(text)
		merged = []Result{
			{
				CharStart: results[0].CharStart,
				CharEnd:   results[0].CharEnd,
				ByteStart: results[0].ByteStart,
				ByteEnd:   results[0].ByteEnd,
			},
		}
		rep = "*"
	)
	if len(replace) > 0 {
		rep = replace[0]
	}

	for cur, j := 0, 0; cur < len(results)-1; cur++ {
		next := cur + 1
		if results[cur].ByteEnd >= results[next].ByteStart && results[cur].ByteEnd <= results[next].ByteEnd {
			merged[j].CharEnd = results[next].CharEnd
			merged[j].ByteEnd = results[next].ByteEnd
		} else {
			merged = append(merged, Result{
				CharStart: results[next].CharStart,
				CharEnd:   results[next].CharEnd,
				ByteStart: results[next].ByteStart,
				ByteEnd:   results[next].ByteEnd,
			})
			j++
		}
	}

	//fmt.Println(merged)
	var pos int
	for _, res := range merged {
		sb.WriteString(text[pos:res.ByteStart])
		sb.WriteString(strings.Repeat(rep, res.CharEnd-res.CharStart+1))
		pos = res.ByteEnd
	}
	if pos < l {
		sb.WriteString(text[pos:])
	}

	return sb.String()
}

func (trie *Trie) FindAll(text string) []string {
	results := trie.Match(text)

	// unique
	m := make(map[string]struct{}, len(results))
	words := make([]string, 0, len(results))
	for _, res := range results {
		if _, ok := m[res.HitWord]; !ok {
			m[res.HitWord] = struct{}{}
			words = append(words, res.HitWord)
		}
	}

	return words
}

func (trie *Trie) Match(text string) (results []Result) {
	var (
		parent               = trie.Root
		sb                   = strings.Builder{}
		start                = 0
		pos                  = 0
		cstart               = 0 // utf-8 char start pos
		ncmatched            = 0 // matched char counter
		l                    = len(text)
		firstMatchedRuneSize int
	)

	for pos < l {
		r, size := utf8.DecodeRuneInString(text[pos:])

		if child, ok := parent.Children[r]; ok {
			sb.WriteString(text[pos : pos+size])
			ncmatched++
			if firstMatchedRuneSize == 0 {
				firstMatchedRuneSize = size
			}
			if child.End {
				result := Result{
					CharStart: cstart + 1,
					CharEnd:   cstart + ncmatched,
					ByteStart: start,
					ByteEnd:   pos + size,
					HitWord:   sb.String(),
				}
				result.MatchedStr = text[result.ByteStart:result.ByteEnd]
				results = append(results, result)
			}
			pos += size
			parent = child
		} else {
			parent = trie.Root
			if firstMatchedRuneSize > 0 {
				start += firstMatchedRuneSize
				firstMatchedRuneSize = 0
				ncmatched = 0
			} else {
				start += size
			}
			pos = start
			cstart++
			sb.Reset()
		}
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
