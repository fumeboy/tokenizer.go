原创的第二种分词器

程序一样很简单

它的特点是 低脑力成本

----

./pattern_test.go

编写下面这段结构体，就可以匹配 `let apple(a,b,c,'abc','a b c')`

```go
// Path  串行匹配
// Sep separator 分割符
// Token  要捕获的文本
// Branch 并行匹配
// Option 可选并行匹配
var (
    function_name = Path(
        Sep(NotBlank),
        Token("function_name"),
        Option(Path(Sep(" "))),
    )
    param = Branch(
        Path(
            Sep("'"),
            Token("param"),
            Sep("'"),
        ),
        Path(
            Sep(NotBlank),
            Token("param"),
            Option(Path(Sep(" "))),
        ),
    )

    sentence_ = Path(
        Sep("let"),
        function_name,
        Option(
            Path(Repeat("(", Path(param), ",", ")")),
        ),
    )
)
```

结果：


`"function_name"` => `apple`

`"param"` => `a`, `b`, `c`, `abc`, `a b c`