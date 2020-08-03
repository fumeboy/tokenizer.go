package tokenizer

type ierr interface {
	display() string
}

type find_sep_err struct {
	i int // 文本进行位置
	have_token *token// 已捕捉的最新一个token
	will_token *token // 将捕捉的token
	next_seps []*separator // 接下来应该遇到的分割符
}

func (f find_sep_err) display() string {
	panic("implement me")
}

type token_check_error struct {
	i int // 文本进行位置
	text string // 受检测文本
	token *token
}

func (this *token_check_error) display() string {
	return ""
}
