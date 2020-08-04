package tokenizer

import (
	"testing"
)

func Test1(t *testing.T) {
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

		text = []byte("let apple")
		ctx = NEWcontext(nil,text)

		text2 = []byte("let apple(a,b,c,'abc','a b c')")
		ctx2 = NEWcontext(nil,text2)
	)
	ctx.RUN(sentence_)
	ctx.Display()
	ctx2.RUN(sentence_)
	ctx2.Display()
}
