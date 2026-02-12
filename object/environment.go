package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

// create new Environment

func NewEnvironment() *Environment { // only uppercase can be exported in Go !!!!
	return &Environment{store: make(map[string]Object), outer: nil}
}

// get from hashmap

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}
	return obj, ok
}

// set in hashmap

func (e *Environment) Set(name string, obj Object) Object {
	e.store[name] = obj // i mean value object
	return obj
}

// create new enclosed environment...

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}
