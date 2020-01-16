package main

import "testing"

func TestInit(t *testing.T) {
	var want error
	if got := Init(); got != want {
		t.Errorf("Init() = %q, want %q", got, want)
	}
}

func TestGetInfo(t *testing.T) {
	var want error
	if _, got := GetInfo(); got != want {
		t.Errorf("GetInfo() = %q, want %q", got, want)
	}
}
