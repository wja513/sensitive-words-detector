package detector

import (
	"strings"
	"unicode/utf8"
)

type Result struct {
	UcharStart int // 1-based, left closed right closed
	UcharEnd   int
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
		byts   = s2b(text)
		l      = len(text)
		next   = 1
		cur    = 0
		i      = 0
		merged = []Result{
			{
				UcharStart: results[cur].UcharStart,
				UcharEnd:   results[cur].UcharEnd,
				ByteStart:  results[cur].ByteStart,
				ByteEnd:    results[cur].ByteEnd,
			},
		}
		rep = "*"
	)
	if len(replace) > 0 {
		rep = replace[0]
	}

	for cur < len(results)-1 {
		if results[cur].ByteEnd >= results[next].ByteStart && results[cur].ByteEnd <= results[next].ByteEnd {
			merged[i].UcharEnd = results[next].UcharEnd
			merged[i].ByteEnd = results[next].ByteEnd

		} else {
			merged = append(merged, Result{
				UcharStart: results[next].UcharStart,
				UcharEnd:   results[next].UcharEnd,
				ByteStart:  results[next].ByteStart,
				ByteEnd:    results[next].ByteEnd,
			})
			i++
		}
		cur++
		next++
	}

	//fmt.Println(merged)
	var pos int
	for _, res := range merged {
		sb.Write(byts[pos:res.ByteStart])
		sb.WriteString(strings.Repeat(rep, res.UcharEnd-res.UcharStart+1))
		pos = res.ByteEnd
	}
	if pos < l {
		sb.Write(byts[pos:])
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
		byts                 = s2b(text)
		l                    = len(byts)
		firstMatchedRuneSize int
	)

	for pos < l {
		r, size := utf8.DecodeRune(byts[pos:])

		if child, ok := parent.Children[r]; ok {
			sb.Write(byts[pos : pos+size])
			ncmatched++
			if firstMatchedRuneSize == 0 {
				firstMatchedRuneSize = size
			}
			if child.End {
				result := Result{
					UcharStart: cstart + 1,
					UcharEnd:   cstart + ncmatched,
					ByteStart:  start,
					ByteEnd:    pos + size,
					HitWord:    sb.String(),
				}
				result.MatchedStr = string(byts[result.ByteStart:result.ByteEnd])
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
