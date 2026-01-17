package index

import "errors"

type PKIndex struct {
	data map[any]int // Maps primary key value to row ID
}

func NewPKIndex() *PKIndex {
	return &PKIndex{
		data: make(map[any]int),
	}
}

func (i *PKIndex) Insert(key any, rowID int) error {
	if _, exists := i.data[key]; exists {
		return errors.New("duplicate primary key")
	}
	i.data[key] = rowID
	return nil
}

func (i *PKIndex) Get(key any) (int, bool) {
	rowID, ok := i.data[key]
	return rowID, ok
}

func (i *PKIndex) Delete(key any) {
	delete(i.data, key)
}