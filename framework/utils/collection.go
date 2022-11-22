package utils

type Collection struct {
	list interface{}
}

func NewStrCollection(list []string) *Collection {
	return &Collection{list: list}
}
