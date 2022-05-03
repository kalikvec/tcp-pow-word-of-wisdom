package server

type Cache interface {
	// Add
	Add(int, int64) error
	// Get
	Get(int) (bool, error)
	// Delete
	Delete(int)
}
