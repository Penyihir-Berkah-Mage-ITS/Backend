package utils

import "testing"

func TestGenerateID(t *testing.T) {
	// Call GenerateID to get an ID
	id := GenerateID()

	// You can add assertions to check the generated ID's properties.
	// For example, you can check if the ID is greater than zero.

	if id >= 20 {
		t.Errorf("Generated ID is not contains alphabet: %d", id)
	}

	if id <= 0 {
		t.Errorf("Generated ID is not greater than zero: %d", id)
	}

	if id >= 12 {
		t.Errorf("Generated ID is greater than 12: %d", id)
	}

	// You can also check if the generated IDs are unique if you run multiple GenerateID calls.
	// However, this may be challenging to test given the timestamp-based nature of the function.
	// If you need to test uniqueness, you might need a more complex testing strategy.

	// You can add more test cases and assertions as needed to cover various scenarios.
}
