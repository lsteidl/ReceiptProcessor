package main

import "testing"

func TestGetRetailerPoints(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"M&Ms", 3},
		{" $TARGET% ^ R", 7},
		{" % hello *&$ ()", 5},
	}
	for _, c := range cases {
		got := getRetailerPoints(c.in)
		if got != c.want {
			t.Errorf("getRetailerPoints(%s) == %d, want %d", c.in, got, c.want)
		}
	}
}

func TestGetRoundPoints(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"100", 50},
		{"60.00", 50},
		{"80.99", 0},
	}
	for _, c := range cases {
		got := getRoundPoints(c.in)
		if got != c.want {
			t.Errorf("getRoundPoints(%s) == %d, want %d", c.in, got, c.want)
		}
	}
}
func TestGetMultiplePoints(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"100", 25},
		{"60.25", 25},
		{"80.75", 25},
		{"80.99", 0},
	}
	for _, c := range cases {
		got := getMultiplePoints(c.in)
		if got != c.want {
			t.Errorf("getMultiplePoints(%s) == %d, want %d", c.in, got, c.want)
		}
	}

}
func TestGetItemPoints(t *testing.T) {
	cases := []struct {
		in   int
		want int
	}{
		{3, 5},
		{1, 0},
		{0, 0},
	}
	for _, c := range cases {
		got := getItemPoints(c.in)
		if got != c.want {
			t.Errorf("getItemPoints(%d) == %d, want %d", c.in, got, c.want)
		}
	}
}
func TestEvauluateDescription(t *testing.T) {
	cases := []struct {
		inDesc  string
		inPrice string
		want    int
	}{
		{" Klarbrunn 12-PK 12 FL OZ", "12.00", 3},
		{"  Emils Cheese Pizza    ", "12.25", 3},
		{" Klarbrunn 12-PK 12 FL OZ", "10.10", 3},
		{"ABC", "99", 20},
	}
	for _, c := range cases {
		got := evaluateDescription(c.inDesc, c.inPrice)
		if got != c.want {
			t.Errorf("evaluateDescription(%s, %s) == %d, want %d", c.inDesc, c.inPrice, got, c.want)
		}
	}

}
func TestGetDatePoints(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"2022-01-01", 6},
		{"1990-01-19", 6},
		{"2021-01-10", 0},
		{"80.99", 0},
	}
	for _, c := range cases {
		got := getDatePoints(c.in)
		if got != c.want {
			t.Errorf("getDatePoints(%s) == %d, want %d", c.in, got, c.want)
		}
	}

}
func TestGetTimePoints(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"13:01", 0},
		{"14:00", 10},
		{"14:33", 10},
		{"16:00", 0},
	}
	for _, c := range cases {
		got := getTimePoints(c.in)
		if got != c.want {
			t.Errorf("getTimePoints(%s) == %d, want %d", c.in, got, c.want)
		}
	}

}
