package model

type ProductFilter struct {
	Name     string
	IsActive bool
}

type OrderFilter struct {
	OrderReference string
}
