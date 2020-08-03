package tokenizer

const (
	LF       = "\n"
	NotBlank = `\S`
)

type itoken interface {
	l() []*separator
}

type path struct {
	nodes []itoken
}

func Path(nodes ...itoken) *path {
	if len(nodes) == 0{
		panic("path 的长度应大于 1")
	}
	if _, ok := nodes[0].(*token); ok{
		panic("path 的首个项不能是 token")
	}
	return &path{nodes: nodes}
}

type token struct {
	identifier string
	check func([]byte)bool
}
func Token(text string) *token {
	if len(text) == 0 {
		panic("非法的标识符")
	}
	return &token{identifier: text}
}
func (this *token) Check(e func([]byte)bool) *token{
	this.check = e
	return this
}

type separator struct {
	text       string
}

func Sep(text string) *separator {
	if len(text) == 0 {
		panic("非法的分割符")
	}
	return &separator{text: text}
}

type repeat struct {
	left     *separator
	repeated *path
	sep      *separator
	right    *separator
}

func Repeat(left string, repeated *path, sep string, right string) *repeat {
	r := &repeat{
		left:     nil,
		repeated: repeated,
		sep:      nil,
		right:    nil,
	}
	if left != "" {
		r.left = Sep(left)
	}
	if right != "" {
		r.right = Sep(right)
	}else{
		r.right = Sep(NotBlank)
	}
	if sep != "" {
		r.sep = Sep(sep)
	}else{
		r.sep = Sep(" ")
	}
	return r
}

type branch struct {
	paths []*path
	optional bool
}

func Branch(paths ...*path) *branch {
	// branch 的各个 path 的首项不能相同
	return &branch{paths: paths}
}
func Option(paths ...*path) *branch {
	// branch 的各个 path 的首项不能相同
	return &branch{paths: paths, optional: true}
}
