package internal

type Item struct {
	Name string
	Path string
}

type Model struct {
	Items    []Item
	Selected int
}

func NewModel() Model {
	return Model{
		Items:    []Item{},
		Selected: 0,
	}
}

