package utils

import (
	"reflect"
	"testing"
)

func TestExcept_1(t *testing.T) {
	src := []string{"hello", "and", "earth", "other", "planets"}
	exp := []string{"other", "and", "planets"}
	want := []string{"hello", "earth"}
	got := Except(src, exp)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("\nwant: %v\n got: %v", want, got)
	}
}

func TestExcept_2(t *testing.T) {
	exp := []string{"hello", "and", "earth", "other", "planets"}
	src := []string{"other", "and", "planets"}
	want := []string{}
	got := Except(src, exp)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("\nwant: %v\n got: %v", want, got)
	}
}

func TestExcept_3(t *testing.T) {
	exp := []string{}
	src := []string{"hello", "world"}
	want := []string{"hello", "world"}
	got := Except(src, exp)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("\nwant: %v\n got: %v", want, got)
	}
}
