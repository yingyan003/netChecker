package sorter

type IResource interface {
	Sort()
	List()
}

type LessFunc func(d1, d2 IResource) bool

type Sorter struct {
	elems []IResource
	Lessf []LessFunc
}
