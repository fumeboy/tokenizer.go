package tokenizer

import "fmt"

func (this *path) l() []*separator {
	for i, l := 0, len(this.nodes); i < l; i++ {
		if _, ok := this.nodes[i].(*token); !ok {
			return this.nodes[i].l()
		}
	}
	return nil
}
func (this *token) l() []*separator { return nil }
func (this *branch) l() []*separator {
	var seps = []*separator{}
	for i, l := 0, len(this.paths); i < l; i++ {
		seps = append(seps, this.paths[i].l()...)
	}
	return seps
}
func (this *repeat) l() []*separator {
	if this.left == nil {
		return this.repeated.l()
	}
	return []*separator{this.left}
}
func (this *separator) l() []*separator {
	return []*separator{this}
}

func (this *path) Match(ctx *context) ierr {
	for i, l := 0, len(this.nodes); i < l; i++ {
		word := this.nodes[i]
		switch word.(type) {
		case *separator:
			s := split(ctx, word.(*separator))
			if s == nil {
				return ctx.err([]*separator{s})
			}
			break
		case *token:
			ctx.reg_token = word.(*token)
			break
		case *path:
			if err := word.(*path).Match(ctx); err != nil {
				return err
			}
			break
		case *repeat:
			rp := word.(*repeat)
			if rp.left != nil {
				s := split(ctx, rp.left)
				if s == nil {
					return ctx.err([]*separator{rp.left})
				}
			}
			var s *separator
			for {
				if err := rp.repeated.Match(ctx); err != nil {
					return err
				}
				s = split(ctx, rp.sep, rp.right)
				if s == nil {
					return ctx.err([]*separator{rp.sep, rp.right})
				} else if s == rp.sep {
					continue
				} else {
					break
				}
			}
			break
		case *branch:
			r := word.(*branch)
			var seps []*separator
			var spm = map[*separator]*path{}

			for j, k := 0, len(r.paths); j < k; j++ {
				s := r.paths[j].l()
				for m, n := 0, len(s); m < n; m++ {
					spm[s[m]] = r.paths[j]
				}
				seps = append(seps, s...)
			}
			ret_sep := split(ctx, seps...)
			if ret_sep == nil {
				if !r.optional{
					return ctx.err(seps)
				}
			}else{
				ctx.borrow = ret_sep
				if err := spm[ret_sep].Match(ctx); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func split(ctx *context, sep ...*separator) *separator {
	if ctx.borrow != nil{
		b := ctx.borrow
		ctx.borrow = nil
		return b
	}
	text := ctx.text[ctx.i:]
	var i = 0
	var j = 0
	var l = len(text)

	var t = 0
	var lt = len(sep)
	var this = sep[t]

	var failed = false
	var success = func() {
		if ctx.reg_token != nil { //捕获
			ctx.add(text[:j])
		}
	}
	var fail = func() {
		failed = true
	}
lalala:
	switch this.text {
	case " ": //blank, 匹配一或多个连续的空格
		for ; i < l; i++ {
			if text[i] == ' ' {
				j = i
				break
			}
		}
		for ; i < l; i++ {
			if text[i] != ' ' {
				break
			}
		}
		if i > 0 {
			//匹配到了空格
			success()
		} else {
			//没匹配到空格，但实际上这种情况很少吧
			fail()
		}
		break
	case NotBlank:
		for ; i < l; i++ {
			if text[i] != ' ' {
				j = i
				break
			}
		}
		if j == i {
			success()
		} else {
			fail()
		}
		break
	default:
		k, lk := 0, len(this.text)
		for ; i < l; i++ {
			if text[i] == this.text[0]{
				j = i
				break
			}
		}
		for ; i < l && k < lk; {
			if text[i] != this.text[k] {
				break
			}
			i++
			k++
		}
		if k == lk && k != 0 {
			success()
		} else {
			fail()
		}
	}
	if failed {
		t++
		if t < lt {
			i = 0
			j = 0
			this = sep[t]
			goto lalala
		} else {
			return nil
		}
	}
	fmt.Println(i, string(this.text))
	ctx.i += i
	return this
}
