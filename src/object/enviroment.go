package object

type Enviroment struct {
	store map[string]Object
	outer *Enviroment
}

func NewEnviroment() *Enviroment {
	return &Enviroment{
		store: make(map[string]Object),
		outer: nil,
	}
}

func NewOuterEnviroment(outer *Enviroment) *Enviroment {
	env := NewEnviroment()
	env.outer = outer

	return env
}

func (e *Enviroment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]

	if !ok && e.outer != nil {
		obj, ok := e.outer.Get(name)

		return obj, ok
	}

	return obj, ok
}

func (e *Enviroment) Set(name string, val Object) {
	e.store[name] = val
}
