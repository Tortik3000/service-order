package entity

type Category struct {
	ID        string
	Name      string
	SortOrder int32
}

type MenuItem struct {
	ID          string
	CategoryID  string
	Name        string
	Description string
	Price       float64
	Active      bool
}
