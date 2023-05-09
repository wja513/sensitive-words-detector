package detector

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContains(t *testing.T) {
	var detectCases = []struct {
		text    string
		words   []string
		expects bool
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
		assert.Equal(t, initTrie(detectCase.words).Contains(detectCase.text), detectCase.expects)
	}
}

func TestFindAll(t *testing.T) {
	var detectCases = []struct {
		text    string
		words   []string
		expects []string
	}{
		// basic ascii
		{
			text:    "ahishers",
			words:   []string{"he", "she", "hers", "his"},
			expects: []string{"his", "she", "he", "hers"},
		},
		// basic utf-8
		{
			text:    "这篇文章真tmd傻X，脑残，tmd瞎逼带节奏~",
			words:   []string{"脑残", "tmd", "傻X"},
			expects: []string{"tmd", "傻X", "脑残"},
		},
	}
	for _, detectCase := range detectCases {
		assert.Equal(t, initTrie(detectCase.words).FindAll(detectCase.text), detectCase.expects)
	}
}

func TestMatch(t *testing.T) {
	var detectCases = []struct {
		text    string
		words   []string
		expects []Result
	}{
		// basic ascii
		{
			text:  "ahishers",
			words: []string{"he", "she", "hers", "his"},
			expects: []Result{
				{
					UcharStart: 2,
					UcharEnd:   4,
					ByteStart:  1,
					ByteEnd:    4,
					HitWord:    "his",
					MatchedStr: "his",
				},
				{
					UcharStart: 4,
					UcharEnd:   6,
					ByteStart:  3,
					ByteEnd:    6,
					HitWord:    "she",
					MatchedStr: "she",
				},
				{
					UcharStart: 5,
					UcharEnd:   6,
					ByteStart:  4,
					ByteEnd:    6,
					HitWord:    "he",
					MatchedStr: "he",
				},
				{
					UcharStart: 5,
					UcharEnd:   8,
					ByteStart:  4,
					ByteEnd:    8,
					HitWord:    "hers",
					MatchedStr: "hers",
				},
			},
		},
		// basic utf-8
		{
			text:  "这篇文章真tmd傻X脑残tmd瞎逼带节奏~",
			words: []string{"tmd", "脑残", "傻X"},
			expects: []Result{
				{
					UcharStart: 6,
					UcharEnd:   8,
					ByteStart:  15,
					ByteEnd:    18,
					HitWord:    "tmd",
					MatchedStr: "tmd",
				},
				{
					UcharStart: 9,
					UcharEnd:   10,
					ByteStart:  18,
					ByteEnd:    22,
					HitWord:    "傻X",
					MatchedStr: "傻X",
				},
				{
					UcharStart: 11,
					UcharEnd:   12,
					ByteStart:  22,
					ByteEnd:    28,
					HitWord:    "脑残",
					MatchedStr: "脑残",
				},
				{
					UcharStart: 13,
					UcharEnd:   15,
					ByteStart:  28,
					ByteEnd:    31,
					HitWord:    "tmd",
					MatchedStr: "tmd",
				},
			},
		},
	}
	for _, detectCase := range detectCases {
		assert.Equal(t, initTrie(detectCase.words).Match(detectCase.text), detectCase.expects)
	}
}

func initTrie(words []string) (trie *Trie) {
	trie = New()
	for _, word := range words {
		trie.Insert(word)
	}

	return
}
