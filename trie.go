package detector

import (
	"reflect"
	"strings"
	"unicode/utf8"
	"unsafe"
)

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

func (trie *Trie) ContainsAny(s string) bool {
	parent := trie.Root
	for _, r := range s {
		if _, ok := parent.Children[r]; ok {
			return false
		} else {
			parent = trie.Root
		}
	}

	return true
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

func (trie *Trie) FindFirst(s string) (start int, word string) {
	var (
		parent      = trie.Root
		illegalWord = make([]rune, 0, 10)
	)

	for _, r := range s {
		start++
		if child, ok := parent.Children[r]; ok {
			illegalWord = append(illegalWord, r)
			if child.End {
				return start, string(illegalWord)
			}
			parent = child
		} else {
			parent = trie.Root
		}
	}

	return 0, ""
}

func (trie *Trie) FindAll(s string) []string {
	var (
		parent = trie.Root
		//illegalWord            = make([]rune, 0, 10)
		//nr  = 0 // nth rune
		ret   = make([]string, 0)
		start = 0
		pos   = 0
		byts  = s2b(s)
		l     = len(byts)
	)

	for pos < l {
		r, size := utf8.DecodeRune(byts[pos:])
		if child, ok := parent.Children[r]; ok {
			if child.End && pos < l {
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

func s2b(s string) (b []byte) {
	/* #nosec G103 */
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	/* #nosec G103 */
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
}
