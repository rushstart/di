package di

import (
	"github.com/rushstart/tid"
	"reflect"
)

func NewInput(c *Container, id tid.ID) (Input, error) {
	typ := id.Type()
	if IsInStruct(typ) {
		definitions := make([]Definition, 0, typ.NumField())
		for i := range typ.NumField() {
			field := typ.Field(i)
			fieldID := tid.FromType(field.Type, field.Tag.Get("tag"))
			definition, err := c.Definition(fieldID)
			if err != nil {
				return Input{}, err
			}
			definitions = append(definitions, definition)
		}
		return Input{id: id, definitions: definitions, isInStruct: true}, nil
	}

	definition, err := c.Definition(id)
	if err != nil {
		return Input{}, err
	}
	return Input{id: id, definitions: []Definition{definition}}, nil
}

type Input struct {
	id          tid.ID
	definitions []Definition
	isInStruct  bool
}

func (i Input) IsInStruct() bool {
	return i.isInStruct
}

func (i Input) InIDs() (ids []tid.ID) {
	for _, def := range i.definitions {
		ids = append(ids, def.id)
	}
	return ids
}

func (i Input) Resolve(sc Scope) (reflect.Value, error) {
	if i.IsInStruct() {
		st := reflect.New(i.id.Type()).Elem()
		for idx, def := range i.definitions {
			v, err := sc.ResolveDefinition(def)
			if err != nil {
				return reflect.Value{}, err
			}
			st.Field(idx).Set(v)
		}
		return st, nil
	}

	return sc.ResolveDefinition(i.definitions[0])
}
