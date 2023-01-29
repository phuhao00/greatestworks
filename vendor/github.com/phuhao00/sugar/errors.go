package sugar

//Try calls the function and return false in case of error.
func Try(callback func() error) (ok bool) {
	ok = true

	defer func() {
		if err := recover(); err != nil {
			ok = false
			return
		}
	}()

	err := callback()
	if err != nil {
		ok = false
	}

	return
}
