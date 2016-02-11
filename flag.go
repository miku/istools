package istools

import (
	"flag"
	"fmt"
	"strings"
)

// Tagged is just a pair of strings. A values and some associated tag.
type Tagged struct {
	Tag   string
	Value string
}

// TagSlice is a list of tagged options.
type TagSlice []Tagged

// String returns the list of tags.
func (v *TagSlice) String() string {
	return fmt.Sprintf("%v", *v)
}

// Set fills up the slice gradually.
func (v *TagSlice) Set(value string) error {
	f := taggedFlag{}
	if err := f.Set(value); err != nil {
		return err
	}
	*v = append(*v, f.Tagged)
	return nil
}

// taggedFlag satisfies the Value interface.
type taggedFlag struct {
	Tagged
}

// String returns the tag and value separated by colon.
func (v *taggedFlag) String() string {
	return fmt.Sprintf("%s:%s", v.Tag, v.Value)
}

// Set tries to parse the value.
func (v *taggedFlag) Set(s string) error {
	var parts = strings.Split(s, ":")
	if len(parts) != 2 {
		return fmt.Errorf("format must be [TAG]:[PATH]")
	}
	v.Tagged.Tag, v.Tagged.Value = parts[0], parts[1]
	return nil
}

// TaggedFlag can be used for parsing flags.
func TaggedFlag(name string, value Tagged, usage string) *Tagged {
	f := taggedFlag{value}
	flag.CommandLine.Var(&f, name, usage)
	return &f.Tagged
}
