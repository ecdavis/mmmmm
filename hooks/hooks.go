package hooks

type Hook func(...interface{}) error

var hooks = make(map[string][]Hook)

func Add(name string, hook Hook) {
	_, ok := hooks[name]
	if !ok {
		hooks[name] = make([]Hook, 0)
	}
	hooks[name] = append(hooks[name], hook)
}

func Run(name string, args ...interface{}) error {
	for _, hook := range(hooks[name]) {
		err := hook(args...)
		if err != nil {
			return err
		}
	}
	return nil
}
