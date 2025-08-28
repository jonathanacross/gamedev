package tiled

import "fmt"

// --------------- Public interface -----------

type Property struct {
	Value interface{}
}

// A property set is just a map of key value pairs.
// The values are Typed, and must be one of bool, int, float64, string,
// according to the setup in Tiled.
type PropertySet map[string]Property

// Gets the property with the given name as a bool.
func (ps PropertySet) GetPropertyBool(name string) (bool, error) {
	return getProp[bool](ps, name)
}

// Gets the property with the given name as an int.
func (ps PropertySet) GetPropertyInt(name string) (int, error) {
	return getProp[int](ps, name)
}

// Gets the property with the given name as a float64.
func (ps PropertySet) GetPropertyFloat64(name string) (float64, error) {
	return getProp[float64](ps, name)
}

// Gets the property with the given name as a string.
func (ps PropertySet) GetPropertyString(name string) (string, error) {
	return getProp[string](ps, name)
}

// ----------  Tiled JSON structs --------------

type PropertiesJSON struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// -----------  Internal conversion functions -----------

// TODO: make these members private once package is refactored
func (p PropertiesJSON) IntValue() (int, bool) {
	return p.intValue()
}

func (p PropertiesJSON) BoolValue() (bool, bool) {
	return p.boolValue()
}

// IntValue attempts to convert the property's value to an int.
func (p PropertiesJSON) intValue() (int, bool) {
	if v, ok := p.Value.(float64); ok {
		return int(v), true
	}
	if val, ok := p.Value.(int); ok {
		return val, true
	}
	return 0, false
}

// BoolValue attempts to convert the property's value to a bool.
func (p PropertiesJSON) boolValue() (bool, bool) {
	if v, ok := p.Value.(bool); ok {
		return v, true
	}
	// Tiled can sometimes export booleans as 0 or 1.
	if v, ok := p.Value.(float64); ok {
		return v == 1, true
	}
	// Tiled can sometimes export booleans as 0 or 1.
	if v, ok := p.Value.(int); ok {
		return v == 1, true
	}
	return false, false
}

// float64Value attempts to convert the property's value to a float64
func (p PropertiesJSON) float64Value() (float64, bool) {
	if v, ok := p.Value.(float64); ok {
		return v, true
	}
	return 0.0, false
}

// stringValue attempts to convert the property's value to a string
func (p PropertiesJSON) stringValue() (string, bool) {
	if v, ok := p.Value.(string); ok {
		return v, true
	}
	return "", false
}

// getProp is a generic helper function that retrieves a property and asserts its type.
func getProp[T any](ps PropertySet, name string) (T, error) {
	var zero T

	prop, ok := ps[name]
	if !ok {
		return zero, fmt.Errorf("property '%s' not found", name)
	}

	value, ok := prop.Value.(T)
	if !ok {
		return zero, fmt.Errorf("property '%s' is not of type %T", name, zero)
	}

	return value, nil
}

func (ps PropertySet) addPropertyHelper(pj *PropertiesJSON, value interface{}, ok bool) error {
	if !ok {
		return fmt.Errorf("could not add property '%s'; not of type %s (actual type %T)", pj.Name, pj.Type, value)
	}
	ps[pj.Name] = Property{Value: value}
	return nil
}

func (ps PropertySet) addProperty(pj *PropertiesJSON) error {
	var value interface{}
	var ok bool

	switch pj.Type {
	case "bool":
		value, ok = pj.boolValue()
	case "int":
		value, ok = pj.intValue()
	case "float":
		value, ok = pj.float64Value()
	case "string":
		value, ok = pj.stringValue()
	default:
		return fmt.Errorf("unknown property type: %s for property %s", pj.Type, pj.Name)
	}

	if !ok {
		return fmt.Errorf("could not add property '%s'; not of type %s (actual type %T)", pj.Name, pj.Type, value)
	}
	ps[pj.Name] = Property{Value: value}

	return nil
}

// TODO: add a unit test
func GetProperties(pj []PropertiesJSON) (*PropertySet, error) {
	properties := PropertySet{}
	for _, p := range pj {
		err := properties.addProperty(&p)
		if err != nil {
			return nil, err
		}
	}
	return &properties, nil
}
