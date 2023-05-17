package detector

import (
	"strings"
	"unicode"
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

type TrieNode struct {
	End      bool
	Children map[rune]*TrieNode
}

type Trie struct {
	Root       *TrieNode
	IgnoreCase bool
	Noises     map[rune]struct{}
}

func (t *Trie) Insert(word string) {
	if word == "" {
		return
	}

	node := t.Root
	for _, r := range word {
		if t.IgnoreCase {
			r = unicode.ToLower(r)
		}
		if child, ok := node.Children[r]; ok {
			node = child
		} else {
			child := newNode()
			node.Children[r] = child
			node = child
		}
	}

	node.End = true
}

func (t *Trie) Check(text string) bool {
	node := t.Root
	for _, r := range text {
		if t.IgnoreCase {
			r = unicode.ToLower(r)
		}
		if child, ok := node.Children[r]; ok {
			if child.End {
				return true
			}
			node = child
		} else {
			node = t.Root
		}
	}

	return false
}

func (t *Trie) Filter(text string, replace ...string) string {
	results := t.Match(text)
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

func (t *Trie) Search(text string) []string {
	results := t.Match(text)

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

func (t *Trie) Match(text string) (results []Result) {
	var (
		node                 = t.Root
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
		if t.IgnoreCase {
			r = unicode.ToLower(r)
		}

		if child, ok := node.Children[r]; ok {
			sb.WriteRune(r)
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
			node = child
		} else {
			if firstMatchedRuneSize > 0 {
				if _, ok := t.Noises[r]; ok {
					pos += size
					ncmatched++
					continue
				}
			}

			node = t.Root
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

func (t *Trie) GetAllWords() (words []string) {
	stack := []*TrieNode{t.Root}
	prefixStack := []string{""}

	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		prefix := prefixStack[len(prefixStack)-1]
		prefixStack = prefixStack[:len(prefixStack)-1]

		if node.End {
			words = append(words, prefix)
		}

		for r, child := range node.Children {
			stack = append(stack, child)
			prefixStack = append(prefixStack, prefix+string(r))
		}
	}

	return
}

func (t *Trie) Delete(word string) {
	var (
		node  = t.Root
		stack []*TrieNode
	)

	for _, r := range word {
		if _, ok := node.Children[r]; !ok {
			return
		}
		stack = append(stack, node)
		node = node.Children[r]
	}
	node.End = false
	if len(node.Children) == 0 {
		for len(stack) > 0 {
			node = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			delete(node.Children, []rune(word)[len(stack)])
			if len(node.Children) > 0 || node.End {
				break
			}
		}
	}
}

func newNode() *TrieNode {
	node := new(TrieNode)
	node.Children = make(map[rune]*TrieNode)
	return node
}

func NewTrie() *Trie {
	trie := new(Trie)
	trie.Root = newNode()
	trie.Noises = make(map[rune]struct{})
	return trie
}
