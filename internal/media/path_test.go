package media

import "testing"

func check(t *testing.T, got, want interface{}) {
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func shouldPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() { _ = recover() }()
	f()
	t.Errorf("should have panicked")
}

func TestPathArgs(t *testing.T) {
	p1 := NewPath("/")
	p2 := NewPath("/picture.jpg")
	p3 := NewPath("/users")
	p4 := NewPath("/users/picture.jpg")

	shouldPanic(t, func() { NewPath("") })
	check(t, p1.ToString(), "/")
	check(t, p2.ToString(), "/picture.jpg")
	check(t, p3.ToString(), "/users")
	check(t, p4.ToString(), "/users/picture.jpg")
}

func TestPathExtensions(t *testing.T) {
	p1 := NewPath("/")
	p2 := NewPath("/picture.jpg")
	p3 := NewPath("/users")
	p4 := NewPath("/users/picture.jpg")

	check(t, p1.Extension(), "")
	check(t, p2.Extension(), ".jpg")
	check(t, p3.Extension(), "")
	check(t, p4.Extension(), ".jpg")
}

func TestPathDirs(t *testing.T) {
	p1 := NewPath("/")
	p2 := NewPath("/picture.jpg")
	p3 := NewPath("/users")
	p4 := NewPath("/users/picture.jpg")

	check(t, p1.Dir(), "/")
	check(t, p2.Dir(), "/")
	check(t, p3.Dir(), "/users")
	check(t, p4.Dir(), "/users")
}

func TestPathBaseBasenames(t *testing.T) {
	p1 := NewPath("/")
	p2 := NewPath("/picture.jpg")
	p3 := NewPath("/users")
	p4 := NewPath("/users/picture.jpg")

	shouldPanic(t, func() { p1.Basename() })
	check(t, p2.Basename(), "picture")
	shouldPanic(t, func() { p3.Basename() })
	check(t, p4.Basename(), "picture")
}
