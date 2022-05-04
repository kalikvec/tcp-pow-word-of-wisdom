package server

type Cache interface {
	// Add - add value by key
	Add(int, int64) error
	// Get - check whether key exists
	Get(int) (bool, error)
	// Delete - delete value by key
	Delete(int)
}
