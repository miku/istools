package istools

import "github.com/miku/span/finc"

type MatchAll struct{}

func (f MatchAll) Apply(is finc.IntermediateSchema) bool {
	return true
}
