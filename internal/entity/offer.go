package entity

type Offer struct {
	OfferId   int // todo точно ли нужно это в энтити?
	Name      string
	Price     float64
	Available bool
}
