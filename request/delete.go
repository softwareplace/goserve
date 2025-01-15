package request

func (i *_impl[T]) Delete(config *Config) (*T, error) {
	return i.exec("DELETE", config)
}
