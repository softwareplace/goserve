package request

func (i *_impl[T]) Post(config *Config) (*T, error) {
	return i.exec("POST", config)
}
