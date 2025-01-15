package request

func (i *_impl[T]) Put(config *Config) (*T, error) {
	return i.exec("PUT", config)
}
