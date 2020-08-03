原创的第二种分词器

程序一样很简单

它的特点是 低脑力成本

----

./pattern_test.go

编写下面这段结构体，就可以匹配 `let apple(a,b,c,'abc','a b c')`

```go
var (
    function_name = Path(
        Sep(NotBlank), // NotBlank 即 \S, 是一种来自正则表达式规则的借用
        Token("function_name"),
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