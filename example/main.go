package main

import (
	"fmt"
	swd "github.com/wja513/sensitive-words-detector"
	"os"
)

// main execute `go run example/main.go` to run the demo
func main() {
	d := swd.New(swd.Options{
		IgnoreCase: true,
		Noises:     []rune(" ~!@#$%^&*()_-+=?<>.—，。/\\|《》？;:：'‘；“¥·"),
	})

	pwd, _ := os.Getwd()
	f, err := os.Open(pwd + "/example/sensitive_words_dict.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	d.Load(f)

	text := "#@这$是#%一^&段包^&**含敏感词*#b和敏&*感#词A的文本@#"
	fmt.Println(d.Detect(text))
	fmt.Println(d.Search(text))
	fmt.Println(d.Match(text))
	fmt.Println(d.Filter(text))
}
