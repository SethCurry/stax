package stax

import (
	"errors"
	"fmt"
	"strings"
)

type Type string
type Supertype string
type Subtype string

type ParsedTypes struct {
	Types      []Type
	Supertypes []Supertype
	Subtypes   []Subtype
}

func ParseTypeLine(typeLine string) (ParsedTypes, error) {
	ret := ParsedTypes{
		Types:      []Type{},
		Supertypes: []Supertype{},
		Subtypes:   []Subtype{},
	}

	splitTypes := strings.Split(typeLine, "—")
	if len(splitTypes) == 0 {
		return ret, errors.New("cannot parse an empty string")
	}
	if len(splitTypes) > 2 {
		return ret, fmt.Errorf("type line contains too many dashes: %s", typeLine)
	}

	mainTypes := strings.Split(splitTypes[0], " ")
	for _, mainType := range mainTypes {
		if mainType == "" {
			continue
		}

		if IsType(mainType) {
			ret.Types = append(ret.Types, Type(mainType))
		} else {
			ret.Supertypes = append(ret.Supertypes, Supertype(mainType))
		}
	}

	if len(splitTypes) == 2 {
		subtypes := strings.Split(splitTypes[1], " ")
		foundSubtypes := []string{}
		for _, subtype := range subtypes {
			if subtype != "" {
				foundSubtypes = append(foundSubtypes, subtype)
			}
		}

		ret.Subtypes = make([]Subtype, len(foundSubtypes))

		for k, v := range foundSubtypes {
			ret.Subtypes[k] = Subtype(v)
		}
	}

	return ret, nil
}
