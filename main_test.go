package main

import "testing"

func TestGetRetailerPoints(t *testing.T) {

	got := getRetailerPoints("M&Ms")
	want := 3

	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
	got = getRetailerPoints(" $TARGET% ^ R")
	want = 7

	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
}

func TestGetRoundPoints(t *testing.T) {
	got := getRoundPoints("100")
	want := 50

	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
	got = getRoundPoints("60.00")
	want = 50

	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
	got = getRoundPoints("80.99")
	want = 0

	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
}
func TestGetMultiplePoints(t *testing.T) {
	got := getMultiplePoints("100")
	want := 25

	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
	got = getMultiplePoints("60.25")
	want = 25

	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
	got = getMultiplePoints("80.75")
	want = 25

	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
	got = getMultiplePoints("80.99")
	want = 0

	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
}
func TestGetItemPoints(t *testing.T) {
	got := getItemPoints(3)
	want := (3 / 2) * 5

	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
	got = getItemPoints(1)
	want = 0

	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
	got = getItemPoints(0)
	want = 0

	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
}
func TestEvauluateDescription(t *testing.T) {
	got := evaluateDescription("   Klarbrunn 12-PK 12 FL OZ", "12.00")
	want := 3

	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
	got = evaluateDescription("   Emils Cheese Pizza    ", "12.25")
	want = 3

	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
	got = evaluateDescription("   Klarbrunn 12-PK 12 FL OZ", "12.25")
	want = 3

	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
}
func TestGetDatePoints(t *testing.T) {
	got := getDatePoints("2022-01-01")
	want := 6

	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
	got = getDatePoints("2022-01-19")
	want = 6

	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
	got = getDatePoints("2022-01-10")
	want = 0

	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
}
func TestGetTimePoints(t *testing.T) {
	got := getTimePoints("13:01")
	want := 0

	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
	got = getTimePoints("14:00")
	want = 10

	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
	got = getTimePoints("14:33")
	want = 10

	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
	got = getTimePoints("16:00")
	want = 0

	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}
}
