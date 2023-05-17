# sensitive-words-detector

## Introduce
A sensitive words tool。 一个敏感词工具，看 https://github.com/importcjj/sensitive 年久失修，就重造了一个轮子。

## Usage
see example/main.go
```go
func main() {
    d := swd.NewWithOptions(&swd.Options{
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
	fmt.Println(d.Detect(text)) // false
	fmt.Println(d.Search(text)) // [敏感词b 敏感词a]
	fmt.Println(d.Match(text)) // [{18 23 29 41 敏感词b 敏感词*#b} {25 31 44 57 敏感词a 敏&*感#词A}]
    fmt.Println(d.Filter(text), "*") // #@这$是#%一^&段包^&**含******和*******的文本@#
}
```
## Benchmark
TODO

## Reference
- https://github.com/importcjj/sensitive
- https://github.com/toolgood/ToolGood.Words