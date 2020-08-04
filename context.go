package tokenizer

import "fmt"

type Context interface {}

type context struct {
	text []byte
	i int
	borrow *separator

	have_token *token
	reg_token *token
	matched map[string][][]byte
	ctx Context



	error ierr
}

func NEWcontext(ctx Context, text []byte) *context {
	return &context{
		text:       text,
		i:          0,
		have_token: nil,
		reg_token:  nil,
		matched: map[string][][]byte{},
		ctx:        ctx,
	}
}


func (this *context) add(data []byte) ierr {
	if this.reg_token.check != nil && !this.reg_token.check(data) {
		return &token_check_error{
			i:     this.i,
			text:  string(data),
			token: this.reg_token,
		}
	}else{
		if this.matched[this.reg_token.identifier] == nil{
			this.matched[this.reg_token.identifier] = [][]byte{}
		}
		this.matched[this.reg_token.identifier] = append(this.matched[this.reg_token.identifier], data)
		this.have_token = this.reg_token
		this.reg_token = nil
		return nil
	}
}
func (this *context) get(id string)[][]byte {
	r := this.matched[id]
	delete(this.matched, id)
	return r
}

func (this *context) err(n []*separator) ierr{
	return &find_sep_err{
		i:          this.i,
		have_token: this.have_token,
		will_token: this.reg_token,
		next_seps:  n,
	}
}

func (this *context) RUN(path *path) ierr {
	if err := this.match(path,nil); err != nil {
		this.error = err
		return err
	}
	if this.reg_token != nil{
		this.add(this.text[this.i:])
	}
	return  nil
}

func (this *context) Display() {
	if this.error == nil {
		fmt.Println("text:",`"`+string(this.text)+`"`)
		for k,v := range this.matched{
			fmt.Println("  token <"+k+">","=>")
			for _,vv := range v{
				fmt.Println("  - ", string(vv))
			}
		}
		fmt.Println("====")
	}else{
		this.error.display()
		for k,v := range this.matched{
			fmt.Println("  token <"+k+">","=>")
			for _,vv := range v{
				fmt.Println("  - ", string(vv))
			}
		}
	}
}
