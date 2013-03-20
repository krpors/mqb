package main

import (
	"testing"
)

var lines = []string{
	"# This is a comment and should generate an error",
	"ACARS_UPLINK_TEST/RefAddr/0/Content=7",
	"ACARS_UPLINK_TEST/RefAddr/0/Type=VER",
	"THIS_HAS_AN_INCORRECTAMOUNTOFSLASHES/0/Type=VER",
}

func TestDefinitions(t *testing.T) {
	d := NewDefinition()
	d.Name = "doesnt matter in this test"
	_, r, n, v, _ := ParseSingleLine(lines[1])
	d.UpdatePropertyMap(r, n, v)

	_, r, n, v, _ = ParseSingleLine(lines[2])
	d.UpdatePropertyMap(r, n, v)

	if d.get("VER") != "7" {
		t.Errorf("Expected value 7 for key \"VER\"")
	}
}

func TestParseSingleLine(t *testing.T) {
	// assert line 0 
	j, r, n, v, err := ParseSingleLine(lines[0])
	if err == nil {
		t.Errorf("Expected error")
	}

	// assert line 1
	j, r, n, v, err = ParseSingleLine(lines[1])
	if err != nil {
		t.Errorf("Unexpected error: '%s'", err)
	}
	if j != "ACARS_UPLINK_TEST" {
		t.Errorf("Unexpected: '%s'", j)
	}
	if r != "0" {
		t.Errorf("Unexpected: '%s'", r)
	}
	if n != "Content" {
		t.Errorf("Unexpected: '%s'", n)
	}
	if v != "7" {
		t.Errorf("Unexpected: '%s'", v)
	}

	// assert line 3
	j, r, n, v, err = ParseSingleLine(lines[3])
	if err == nil {
		t.Errorf("Expected error")
	}
}
