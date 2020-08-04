原创的第二种分词器

程序一样很简单

它的特点是 低脑力成本

----

./pattern_test.go

编写下面这段结构体，就可以匹配 `let apple(a,b,c,'abc','a b c')` 和 `let apple`

```go
// Path  串行匹配
// Sep separator 分割符
// Token  要捕获的文本
// Branch 并行匹配
// Option 可选并行匹配
// Repeat 重复匹配
var (
    function_name = Path(
        Sep(NotBlank),// NotBlank 就是正则表达式中使用的 \S
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


输入是 `let apple(a,b,c,'abc','a b c')` ，那么结果是：

`"function_name"` => `apple`

`"param"` => `a`, `b`, `c`, `abc`, `a b c`

解释：

`let apple(a,b,c,'abc','a b c')` 的结构是 `let <function_name> (<param>,<param>...)`

我们的目的是从中取出 `<function_name>` 和 `<param>` 对应的值

而我们需要做的就是使用 Path、Sep、Branch、Option、Repeat 这些工具来描述出它的句型
