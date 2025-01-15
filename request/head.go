package request

func (i *_impl[T]) Head(config *Config) (*T, error) {
	return i.exec("HEAD", config)
}
