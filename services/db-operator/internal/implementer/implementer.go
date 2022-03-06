package implementer

type Implementer struct {
	actions    []interface{}
	NextAction chan interface{}
}

func New() *Implementer {
	return &Implementer{}
}

func (i *Implementer) AddAction(action interface{}) {
	i.actions = append(i.actions, action)
}

func (i *Implementer) DoAction(action interface{}) {

}
