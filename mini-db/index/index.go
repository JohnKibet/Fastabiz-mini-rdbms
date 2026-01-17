package index

type Index interface {
	Insert(key any, rowID int) error
	Get(key any) (int, bool)
	Delete(key any)
}

// Insert → update index
// Delete → remove from index
// Lookup → O(1)
