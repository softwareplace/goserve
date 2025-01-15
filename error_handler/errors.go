package error_handler

func Handler(try func(), catch func(err any)) {
	defer func() {
		if r := recover(); r != nil {
			catch(r)
		}
	}()
	try()
}
