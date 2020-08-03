package tokenizer

import (
	"fmt"
)

type p struct {
	a string
	b string
}
var pp p
func main()  {
	pp = p{
		a: "123",
		b: pp.a,
	}
	fmt.Println(pp)
}
