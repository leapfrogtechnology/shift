package utils

import "testing"

// TestIsAlphaNumeric : command : go test -v ./core/utils/string
func TestIsAlphaNumeric(t *testing.T) {

	emptyResult := IsAlphaNumeric("")

	if emptyResult != false {
		t.Errorf("TestIsAlphaNumeric(\"\") Failed Expected %v, got %v", false, emptyResult)
	} else {
		t.Logf("TestIsAlphaNumeric(\"\") Success Expected %v, got %v", true, emptyResult)
	}

	alphaNumbericVal := IsAlphaNumeric("anish123")

	if alphaNumbericVal != true {
		t.Errorf("TestIsAlphaNumeric(\"\") Failed Expected %v, got %v", false, alphaNumbericVal)
	} else {
		t.Logf("TestIsAlphaNumeric(\"\") Success Expected %v, got %v", true, alphaNumbericVal)
	}

	alphaVal := IsAlphaNumeric("anish")

	if alphaNumbericVal != true {
		t.Errorf("TestIsAlphaNumeric(\"\") Failed Expected %v, got %v", false, alphaVal)
	} else {
		t.Logf("TestIsAlphaNumeric(\"\") Success Expected %v, got %v", true, alphaVal)
	}
}
