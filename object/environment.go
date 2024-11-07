package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]Object),
	}
}

// NewEnclosedEnvironment extends an environment with its own things.
// If something is not found it will look in the extending environment.
// Basically like a Matrojska Doll of Environments until it just can't find
// it anywhere and errors out.
// Without this we would potentially be overwriting existing arguments and values.
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// TODO: doesn't this need locks?
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}

	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
