package main

// Initializes and asserts system preparedness for calls to Facetime. Should be called before processing any user request.
func initSession() error {
	if processLock {
		return ErrBusy
	}
	processLock = true
	err := openFacetime()

	return err
}

// Dispose and closes lock for the active request. Should always be called after request processing is complete.
func closeSession() {
	processLock = false
}
