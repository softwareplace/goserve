package request

func (i *_impl[T]) Patch(config *Config) (*T, error) {
	return i.exec("PATCH", config)
}
