package storage

// GenericStorage describes interface every other storage should embed
// and conform to as it contains essential things like context handling.
type GenericStorage interface {
	// Get returns data from storage by key.
	Get(string) string
	// GetDoneChan returns a channel which should be used to block execution
	// until storage's routines are completed.
	GetDoneChan() chan struct{}
	// Put puts passed data into storage.
	Put(map[string]string)
	// Start starts asynchronous things if needed.
	Start()
}
