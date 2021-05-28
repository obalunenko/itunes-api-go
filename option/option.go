package option

import (
	"github.com/obalunenko/itunes-api-go/internal"
)

type LookupOption interface {
	Apply(params *internal.LookupParams)
}

type withCountry string

func (w withCountry) Apply(p *internal.LookupParams) {
	p.Country = string(w)
}

func WithCountry(country string) LookupOption {
	return withCountry(country)
}

type withMedia string

func (w withMedia) Apply(p *internal.LookupParams) {
	p.Media = string(w)
}

func WithMedia(media string) LookupOption {
	return withMedia(media)
}

type withID int

func (w withID) Apply(p *internal.LookupParams) {
	p.ID = int(w)
}

func WithID(id int) LookupOption {
	return withID(id)
}
