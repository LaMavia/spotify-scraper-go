package main

// HandleError : Handles an `error != nil`
func HandleError(err error) bool {
	if err != nil {
		panic(err)
	}

	return true
}
