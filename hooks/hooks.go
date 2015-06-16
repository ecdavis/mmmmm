package hooks

type Hook func(...interface{})

var hooks = make(map[string][]Hook)

func Add(name string, hook Hook) {
	list, ok := hooks[name]
	if !ok {
		list = make([]Hook, 0)
	}
	list = append(list, hook)
	hooks[name] = list
}

func Run(name string, args ...interface{}) {
	for _, hook := range(hooks[name]) {
		hook(args)
	}
}
