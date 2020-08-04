package tokenizer

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

func (ctx *context) match(p *path, end []*separator) ierr {
	i, l := 0, len(p.nodes)
	var word itoken
	var next = func() []*separator{
		if i != l-1 {
			return p.nodes[i+1].l()
		}
		return end
	}
	for ; i < l; i++ {
		word = p.nodes[i]
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
				if err := ctx.match(rp.repeated, []*separator{rp.sep, rp.right}); err != nil {
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
				for m, n := 0, len(s); m < n; m++ { // #2
					spm[s[m]] = r.paths[j]
				}
				seps = append(seps, s...)
			}
			if r.optional { // #1
				if pp := next(); pp != nil {
					seps = append(seps, pp...)
				}
			}
			ret_sep := split(ctx, seps...)
			if ret_sep == nil {
				if !r.optional {
					return ctx.err(seps)
				}
			} else {
				ctx.borrow = ret_sep
				if spm[ret_sep] == nil { // #1

				}else{
					if err := ctx.match(spm[ret_sep], next()); err != nil { // #2
						return err
					}
				}
			}
		}
	}
	return nil
}

func split(ctx *context, sep ...*separator) *separator {
	if ctx.borrow != nil {
		b := ctx.borrow
		ctx.borrow = nil
		return b
	}
	text := ctx.text[ctx.i:]
	var i = 0
	var j = 0
	var l = len(text)
	var n = len(sep)

	var separators = []*separator{}
	var if_notblank *separator
	var if_blank *separator
	for m := 0; m < n; m++ {
		switch sep[m].text {
		case NotBlank:
			if_notblank = sep[m]
			break
		case " ":
			if_blank = sep[m]
			break
		default:
			separators = append(separators, sep[m])
		}
	}
	n = len(separators)

	var success = func() {
		ctx.i += i
		if ctx.reg_token != nil { //捕获
			ctx.add(text[:j])
		}
	}
	if ctx.reg_token == nil {
		if if_blank != nil && text[0] == ' ' {
			for ; i < l; i++ {
				if text[i] != ' ' {
					break
				}
			}
			success()
			return if_blank
		}
		for ; i < l; i++ {
			if text[i] != ' ' {
				break
			}
		}
		if i == l {
			return nil
		}
		j = i
		for m := 0; m < n; m++ {
			s := separators[m]
			if i+len(s.text) <= l && equal(text[i:len(s.text)+i],[]byte(s.text)) {
				i = i+len(s.text)
				success()
				return s
			}
		}
		if if_notblank != nil {
			success()
			return if_notblank
		}
	} else {
		for ; i < l; i++ {
			if if_blank != nil && text[i] == ' '{
				j = i
				for ; i < l; i++ {
					if text[i] != ' ' {
						break
					}
				}
				success()
				return if_blank
			}
			for m := 0; m < n; m++ {
				s := separators[m]
				if text[i] == s.text[0]{
					if i+len(s.text) < l && equal(text[i:len(s.text)+i],[]byte(s.text)) {
						j = i
						i = i+len(s.text)
						success()
						return s
					}
				}
			}
		}
	}
	return nil
}

func equal(a []byte, b []byte) bool {
	if len(a) != len(b){
		return false
	}
	for i,l := 0,len(a);i<l;i++{
		if a[i] != b[i]{
			return false
		}
	}
	return true
}