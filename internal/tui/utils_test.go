package tui

import (
  "testing"
)

func TestParseTime(t *testing.T) {
  ms := 210000
  want := "3:30"
  out := parseTime(ms)

  if out != want {
    t.Fatalf("Expected %v. Got %v", want, out)
  }
}
