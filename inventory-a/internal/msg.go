package internal

type Msg string

const (
	LoadItems    Msg = "load_items"
	AddItem      Msg = "add_item"
	RemoveItem   Msg = "remove_item"
	ClearItems   Msg = "clear_items"
	SelectItem   Msg = "select_item"
)

