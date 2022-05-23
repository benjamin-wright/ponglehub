package rules

type RuleFail struct {
	response string
	log      string
}

func (r *RuleFail) Response() string {
	return r.response
}

func (r *RuleFail) Log() string {
	return r.log
}
