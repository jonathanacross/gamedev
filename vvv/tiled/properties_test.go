package tiled

import (
	"testing"
)

func TestGetPropertyBool(t *testing.T) {
	ps := make(PropertySet)
	ps["isAlive"] = Property{Value: true}
	ps["score"] = Property{Value: 123}

	tests := []struct {
		name        string
		propName    string
		expectedVal bool
		expectErr   bool
	}{
		{"success", "isAlive", true, false},
		{"not found", "isDead", false, true},
		{"type mismatch", "score", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := ps.GetPropertyBool(tt.propName)
			if (err != nil) != tt.expectErr {
				t.Fatalf("expected error: %v, got: %v", tt.expectErr, err)
			}
			if err == nil && val != tt.expectedVal {
				t.Errorf("expected value: %v, got: %v", tt.expectedVal, val)
			}
		})
	}
}

func TestGetPropertyInt(t *testing.T) {
	ps := make(PropertySet)
	ps["score"] = Property{Value: 123}
	ps["isAlive"] = Property{Value: true}   // Mismatch type for testing
	ps["floatVal"] = Property{Value: 45.67} // Another mismatch type

	tests := []struct {
		name        string
		propName    string
		expectedVal int
		expectErr   bool
	}{
		{"success - found int", "score", 123, false},
		{"error - not found", "level", 0, true},
		{"error - type mismatch (bool)", "isAlive", 0, true},
		{"error - type mismatch (float64)", "floatVal", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := ps.GetPropertyInt(tt.propName)

			if (err != nil) != tt.expectErr {
				t.Fatalf("expected error: %v, got: %v", tt.expectErr, err)
			}
			if err == nil && val != tt.expectedVal {
				t.Errorf("expected value: %v, got: %v", tt.expectedVal, val)
			}
		})
	}
}

func TestGetPropertyFloat64(t *testing.T) {
	ps := make(PropertySet)
	ps["health"] = Property{Value: 99.5}
	ps["name"] = Property{Value: "Player"} // Mismatch type

	tests := []struct {
		name        string
		propName    string
		expectedVal float64
		expectErr   bool
	}{
		{"success - found float64", "health", 99.5, false},
		{"error - not found", "mana", 0.0, true},
		{"error - type mismatch (string)", "name", 0.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := ps.GetPropertyFloat64(tt.propName)

			if (err != nil) != tt.expectErr {
				t.Fatalf("expected error: %v, got: %v", tt.expectErr, err)
			}
			if err == nil && val != tt.expectedVal {
				t.Errorf("expected value: %v, got: %v", tt.expectedVal, val)
			}
		})
	}
}

func TestGetPropertyString(t *testing.T) {
	ps := make(PropertySet)
	ps["name"] = Property{Value: "PlayerOne"}
	ps["isAlive"] = Property{Value: true} // Mismatch type

	tests := []struct {
		name        string
		propName    string
		expectedVal string
		expectErr   bool
	}{
		{"success - found string", "name", "PlayerOne", false},
		{"error - not found", "guild", "", true},
		{"error - type mismatch (bool)", "isAlive", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := ps.GetPropertyString(tt.propName)

			if (err != nil) != tt.expectErr {
				t.Fatalf("expected error: %v, got: %v", tt.expectErr, err)
			}
			if err == nil && val != tt.expectedVal {
				t.Errorf("expected value: %v, got: %v", tt.expectedVal, val)
			}
		})
	}
}

func TestGetProperties(t *testing.T) {
	tests := []struct {
		name              string
		inputPJ           []PropertiesJSON
		expectErr         bool
		expectedPropCount int
		expectedValues    map[string]interface{}
	}{
		{
			name: "Successfully get a mixed set of properties",
			inputPJ: []PropertiesJSON{
				{Name: "isVisible", Type: "bool", Value: true},
				{Name: "playerScore", Type: "int", Value: 100},
				{Name: "speed", Type: "float", Value: 1.23},
				{Name: "playerName", Type: "string", Value: "Alice"},
			},
			expectErr:         false,
			expectedPropCount: 4,
			expectedValues: map[string]interface{}{
				"isVisible":   true,
				"playerScore": 100,
				"speed":       1.23,
				"playerName":  "Alice",
			},
		},
		{
			name: "Fail on unknown property type",
			inputPJ: []PropertiesJSON{
				{Name: "propertyA", Type: "string", Value: "valid"},
				{Name: "propertyB", Type: "unknown", Value: 123}, // This one should cause the failure
				{Name: "propertyC", Type: "bool", Value: true},
			},
			expectErr:         true,
			expectedPropCount: 0,
			expectedValues:    nil,
		},
		{
			name: "Fail on type mismatch",
			inputPJ: []PropertiesJSON{
				{Name: "propertyA", Type: "int", Value: 10},
				{Name: "propertyB", Type: "bool", Value: "yes"}, // This one should cause the failure
				{Name: "propertyC", Type: "float", Value: 3.14},
			},
			expectErr:         true,
			expectedPropCount: 0,
			expectedValues:    nil,
		},
		{
			name:              "Empty slice",
			inputPJ:           []PropertiesJSON{},
			expectErr:         false,
			expectedPropCount: 0,
			expectedValues:    map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			propSet, err := GetProperties(tt.inputPJ)

			if tt.expectErr {
				if err == nil {
					t.Fatal("expected an error, but got nil")
				}
				if propSet != nil {
					t.Fatalf("expected nil PropertySet on error, but got %v", propSet)
				}
				return // Test is complete for this case
			}

			if err != nil {
				t.Fatalf("did not expect an error, but got: %v", err)
			}
			if propSet == nil {
				t.Fatal("expected a PropertySet, but got nil")
			}

			if len(*propSet) != tt.expectedPropCount {
				t.Errorf("expected %d properties, but got %d", tt.expectedPropCount, len(*propSet))
			}

			// Check that all expected properties and values exist
			for key, expectedVal := range tt.expectedValues {
				prop, ok := (*propSet)[key]
				if !ok {
					t.Errorf("expected property '%s' to exist, but it didn't", key)
					continue
				}

				if prop.Value != expectedVal {
					t.Errorf("for property '%s', expected value %v, but got %v", key, expectedVal, prop.Value)
				}
			}
		})
	}
}
