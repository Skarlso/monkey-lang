package object

type Environment struct {
	store map[string]Object
}

func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]Object),
	}
}

// TODO: doesn't this need locks?
func (e *Environment) Get(name string) (Object, bool) {
	v, ok := e.store[name]

	return v, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
