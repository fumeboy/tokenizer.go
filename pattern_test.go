package tokenizer

import (
	"fmt"
	"testing"
)

func TestT(t *testing.T) {
	var (
		function_name = Path(
			Sep(NotBlank),
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

		text = []byte("let apple(a,b,c,'abc','a b c')")
		ctx = NEWcontext(nil,text)
	)

	fmt.Println(sentence_.Match(ctx),ctx.matched)
}
