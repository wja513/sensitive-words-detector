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
		// basic utf-8
		{
			text:    "这篇文章真tmd傻X，脑残，tmd~",
			words:   []string{"tmd", "脑残", "傻X"},
			expects: true,
		},
		// basic utf-8
		{
			text:    "这篇文章真tmd傻X，脑残，tmd~",
			words:   []string{"五毛", "小粉红"},
			expects: false,
		},
	}
	for _, detectCase := range detectCases {
		assert.Equal(t, initTrie(detectCase.words).Contains(detectCase.text), detectCase.expects)
	}
}

func TestFindAll2(t *testing.T) {
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
			text:    "这篇文章真tmd傻X，脑残，tmd~",
			words:   []string{"tmd", "脑残", "傻X"},
			expects: []string{"tmd", "傻X", "脑残", "tmd"},
		},
	}
	for _, detectCase := range detectCases {
		assert.Equal(t, initTrie(detectCase.words).FindAll2(detectCase.text), detectCase.expects)
	}
}

func TestFindAll(t *testing.T) {
	var detectCases = []struct {
		text    string
		words   []string
		expects []string
	}{
		//// basic ascii
		//{
		//	text:    "his-his",
		//	words:   []string{"his"},
		//	expects: []string{"his"},
		//},
		// basic ascii
		{
			text: "傻X-傻X=傻X",
			//words: []string{"bb"},
			words:   []string{"傻X"},
			expects: []string{"傻X"},
		},
		//// basic ascii
		//{
		//	text:    "ahishers",
		//	words:   []string{"he", "she", "hers", "his"},
		//	expects: []string{"his", "she", "he", "hers"},
		//},
		//// basic utf-8
		//{
		//	text:    "这篇文章真tmd傻X，脑残，tmd~",
		//	words:   []string{"tmd", "脑残", "傻X"},
		//	expects: []string{"tmd", "傻X", "脑残"},
		//},
	}
	for _, detectCase := range detectCases {
		assert.Equal(t, initTrie(detectCase.words).FindAll(detectCase.text), detectCase.expects)
	}
}

func initTrie(words []string) (trie *Trie) {
	trie = New()
	for _, word := range words {
		trie.Insert(word)
	}

	return
}
