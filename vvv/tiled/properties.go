package tiled

import (
	"fmt"
	"strconv"
)

// Gets the property with the given name as a bool.
// TODO: simplify names to just GetBool, GetInt, etc.
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

// GetProperties converts a slice of Tiled properties into a PropertySet map.
func GetProperties(tiledProperties []tiledProperty) (PropertySet, error) {
	ps := make(PropertySet, len(tiledProperties))
	for _, p := range tiledProperties {
		value, err := parseValue(p)
		if err != nil {
			return nil, err
		}
		ps[p.Name] = Property{Value: value}
	}
	return ps, nil
}

// parseValue handles type assertion and conversion for a single tiledProperty.
func parseValue(p tiledProperty) (interface{}, error) {
	switch p.Type {
	case "int":
		return parseInt(p.Value)
	case "bool":
		return parseBool(p.Value)
	case "float":
		return parseFloat(p.Value)
	case "string":
		return parseString(p.Value)
	default:
		return p.Value, nil
	}
}

// parseInt handles the various types an integer property might be.
func parseInt(v interface{}) (int, error) {
	if i, ok := v.(int); ok {
		return i, nil
	}
	if f, ok := v.(float64); ok {
		return int(f), nil
	}
	if s, ok := v.(string); ok {
		parsed, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, err
		}
		return int(parsed), nil
	}
	return 0, fmt.Errorf("value is not a number")
}

// parseBool handles the various types a boolean property might be.
func parseBool(v interface{}) (bool, error) {
	if b, ok := v.(bool); ok {
		return b, nil
	}
	return false, fmt.Errorf("value is not a boolean")
}

// parseFloat handles the various types a float property might be.
func parseFloat(v interface{}) (float64, error) {
	if f, ok := v.(float64); ok {
		return f, nil
	}
	if i, ok := v.(int); ok {
		return float64(i), nil
	}
	if s, ok := v.(string); ok {
		parsed, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, err
		}
		return parsed, nil
	}
	return 0, fmt.Errorf("value is not a number")
}

// parseString handles the various types a string property might be.
func parseString(v interface{}) (string, error) {
	if s, ok := v.(string); ok {
		return s, nil
	}
	return "", fmt.Errorf("value is not a string")
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
