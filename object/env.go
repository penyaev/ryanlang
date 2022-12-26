package object

type Environment struct {
	storage map[string]Object
	parent  *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		storage: map[string]Object{},
	}
}
func (e *Environment) Get(key string) (Object, bool) {
	environ := e
	for environ != nil {
		value, ok := environ.storage[key]
		if ok {
			return value, true
		}

		environ = environ.parent
	}

	return nil, false
}
func (e *Environment) GetFromCurrent(key string) (Object, bool) {
	value, ok := e.storage[key]
	return value, ok

}
func (e *Environment) Set(key string, value Object) {
	e.storage[key] = value
}
func (e *Environment) Replace(key string, value Object) bool {
	environ := e
	for environ != nil {
		_, ok := environ.storage[key]
		if ok {
			environ.storage[key] = value
			return true
		}

		environ = environ.parent
	}
	return false
}
func (e *Environment) Derive() *Environment {
	newEnv := NewEnvironment()
	newEnv.parent = e
	return newEnv
}
func (e *Environment) DeriveWith(key string, value Object) *Environment {
	newEnv := e.Derive()
	newEnv.Set(key, value)
	return newEnv
}
