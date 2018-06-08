package placeholder

import (
	"strings"
)

type PlaceHolder struct {
	Replacer string
}

func NewPlaceHolder(replacer string) *PlaceHolder {
	return &PlaceHolder{
		Replacer: replacer,
	}
}

func (ph *PlaceHolder) Replace(oldNew ...string) string {
	r := strings.NewReplacer(oldNew...)
	return r.Replace(ph.Replacer)
}
