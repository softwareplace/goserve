package request

func (i *_impl[T]) Get(config *Config) (*T, error) {
	return i.exec("GET", config)
}
