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
		assert.Equal(t, detectCase.expects, initTrie(detectCase.words).Contains(detectCase.text))
	}
}

func TestFindAll(t *testing.T) {
	var detectCases = []struct {
		text    string
		words   []string
		expects []string
	}{
		// ascii
		{
			text:    "ahishers",
			words:   []string{"he", "she", "hers", "his"},
			expects: []string{"his", "she", "he", "hers"},
		},
		// utf-8
		{
			text:    "这篇文章真tmd傻X，脑残，tmd瞎逼带节奏~",
			words:   []string{"脑残", "tmd", "傻X"},
			expects: []string{"tmd", "傻X", "脑残"},
		},
	}
	for _, detectCase := range detectCases {
		assert.Equal(t, detectCase.expects, initTrie(detectCase.words).FindAll(detectCase.text))
	}
}

func TestFilter(t *testing.T) {
	var detectCases = []struct {
		text    string
		words   []string
		replace string
		expects string
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
			//replace: "",
			//expects: "这篇文章真，，瞎逼带节奏~",
		},
	}
	for _, detectCase := range detectCases {
		assert.Equal(t, detectCase.expects, initTrie(detectCase.words).Filter(detectCase.text, detectCase.replace))
	}
}

func TestMatch(t *testing.T) {
	var detectCases = []struct {
		text    string
		words   []string
		expects []Result
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
	}
	for _, detectCase := range detectCases {
		assert.Equal(t, detectCase.expects, initTrie(detectCase.words).Match(detectCase.text))
	}
}

func initTrie(words []string) (trie *Trie) {
	trie = New()
	for _, word := range words {
		trie.Insert(word)
	}

	return
}
