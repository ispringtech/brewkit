package executor

type Args []string

func (args *Args) AddKV(key, value string) {
	*args = append(*args, key, value)
}

func (args *Args) AddArgs(arguments ...string) {
	*args = append(*args, arguments...)
}
