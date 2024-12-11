package object

type Enviroment struct {
	store map[string]Object
}

func NewEnviroment() *Enviroment {
	return &Enviroment{
		store: make(map[string]Object),
	}
}

func (e *Enviroment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]

	return obj, ok
}

func (e *Enviroment) Set(name string, val Object) {
	e.store[name] = val
}
