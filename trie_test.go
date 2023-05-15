package detector

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheck(t *testing.T) {
	var detectCases = []struct {
		text       string
		words      []string
		ignoreCase bool
		nosies     []rune
		expects    bool
	}{
		{
			text:    "看，他真像个傻X^_^",
			words:   []string{"傻B", "傻X"},
			expects: true,
		},
		{
			text:    "看，他真像个傻X^_^",
			words:   []string{"傻B", "SB"},
			expects: false,
		},
	}
	for _, detectCase := range detectCases {
		assert.Equal(t, detectCase.expects, initTrie(detectCase.words, detectCase.ignoreCase, detectCase.nosies).Check(detectCase.text))
	}
}

func TestSearch(t *testing.T) {
	var detectCases = []struct {
		text       string
		words      []string
		ignoreCase bool
		nosies     []rune
		expects    []string
	}{
		// ascii
		{
			text:    "ahishers",
			words:   []string{"he", "she", "hers", "his"},
			expects: []string{"his", "she", "he", "hers"},
		},
		// utf-8
		{
			text:       "这篇文章真tmd傻X，脑残，tmd瞎逼带节奏~",
			words:      []string{"脑残", "tmd", "傻x"},
			ignoreCase: true,
			expects:    []string{"tmd", "傻x", "脑残"},
		},
		{
			text:    "#@这$是#%一^&段包^&**含敏感词*#3和敏&*感#词1的文本@#",
			words:   []string{"敏感词1", "敏感词2", "敏感词3"},
			nosies:  []rune("#@$%^*&"),
			expects: []string{"敏感词3", "敏感词1"},
		},
	}
	for _, detectCase := range detectCases {
		assert.Equal(t, detectCase.expects, initTrie(detectCase.words, detectCase.ignoreCase, detectCase.nosies).Search(detectCase.text))
	}
}

func TestFilter(t *testing.T) {
	var detectCases = []struct {
		text       string
		words      []string
		replace    string
		ignoreCase bool
		nosies     []rune
		expects    string
	}{
		{
			text:    "ahishers",
			words:   []string{"he", "she", "hers", "his"},
			replace: "*",
			expects: "a*******",
		},
		{
			text:    "这篇文章真tmd傻X，脑残，tmd瞎逼带节奏~",
			words:   []string{"脑残", "tmd", "傻X"},
			replace: "*",
			expects: "这篇文章真*****，**，***瞎逼带节奏~",
		},
		{
			text:    "#@这$是#%一^&段包^&**含敏感词*#3和敏&*感#词1的文本@#",
			words:   []string{"敏感词1", "敏感词2", "敏感词3"},
			nosies:  []rune("#@$%^*&"),
			replace: "*",
			expects: "#@这$是#%一^&段包^&**含******和*******的文本@#",
		},
	}
	for _, detectCase := range detectCases {
		assert.Equal(t, detectCase.expects, initTrie(detectCase.words, detectCase.ignoreCase, detectCase.nosies).Filter(detectCase.text, detectCase.replace))
	}
}

func TestMatch(t *testing.T) {
	var detectCases = []struct {
		text       string
		words      []string
		ignoreCase bool
		nosies     []rune
		expects    []Result
	}{
		// ascii
		{
			text:  "ahishers",
			words: []string{"he", "she", "hers", "his"},
			expects: []Result{
				{
					CharStart:  2,
					CharEnd:    4,
					ByteStart:  1,
					ByteEnd:    4,
					HitWord:    "his",
					MatchedStr: "his",
				},
				{
					CharStart:  4,
					CharEnd:    6,
					ByteStart:  3,
					ByteEnd:    6,
					HitWord:    "she",
					MatchedStr: "she",
				},
				{
					CharStart:  5,
					CharEnd:    6,
					ByteStart:  4,
					ByteEnd:    6,
					HitWord:    "he",
					MatchedStr: "he",
				},
				{
					CharStart:  5,
					CharEnd:    8,
					ByteStart:  4,
					ByteEnd:    8,
					HitWord:    "hers",
					MatchedStr: "hers",
				},
			},
		},
		// utf-8
		{
			text:  "这篇文章真tmd傻X脑残tmd瞎逼带节奏~",
			words: []string{"tmd", "脑残", "傻X"},
			expects: []Result{
				{
					CharStart:  6,
					CharEnd:    8,
					ByteStart:  15,
					ByteEnd:    18,
					HitWord:    "tmd",
					MatchedStr: "tmd",
				},
				{
					CharStart:  9,
					CharEnd:    10,
					ByteStart:  18,
					ByteEnd:    22,
					HitWord:    "傻X",
					MatchedStr: "傻X",
				},
				{
					CharStart:  11,
					CharEnd:    12,
					ByteStart:  22,
					ByteEnd:    28,
					HitWord:    "脑残",
					MatchedStr: "脑残",
				},
				{
					CharStart:  13,
					CharEnd:    15,
					ByteStart:  28,
					ByteEnd:    31,
					HitWord:    "tmd",
					MatchedStr: "tmd",
				},
			},
		},
		{
			text:       "#@这$是#%一^&段包^&**含敏感词*#b和敏&*感#词A的文本@#",
			words:      []string{"敏感词a", "敏感词2", "敏感词B"},
			ignoreCase: true,
			nosies:     []rune("#@$%^*&"),
			expects: []Result{
				{
					CharStart:  18,
					CharEnd:    23,
					ByteStart:  29,
					ByteEnd:    41,
					HitWord:    "敏感词b",
					MatchedStr: "敏感词*#b",
				},
				{
					CharStart:  25,
					CharEnd:    31,
					ByteStart:  44,
					ByteEnd:    57,
					HitWord:    "敏感词a",
					MatchedStr: "敏&*感#词A",
				},
			},
		},
	}
	for _, detectCase := range detectCases {
		assert.Equal(t, detectCase.expects, initTrie(detectCase.words, detectCase.ignoreCase, detectCase.nosies).Match(detectCase.text))
	}
}

func initTrie(words []string, ignoreCase bool, noises []rune) (trie *Trie) {
	trie = New()
	trie.IgnoreCase = ignoreCase
	trie.Noises = make(map[rune]struct{}, 0)
	for _, r := range noises {
		trie.Noises[r] = struct{}{}
	}

	for _, word := range words {
		trie.Insert(word)
	}

	return
}
