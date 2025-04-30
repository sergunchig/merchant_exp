package entity

type Offer struct {
	OfferId   int // todo точно ли нужно это в энтити? да, это импортируемуе значение
	Name      string
	Price     float64
	Available bool
}
