package interpol

import (
	"fmt"
	"testing"
)

// test data
var data1 = []string{"first", "second", "third"}      // youtube comment section
var data2 = []string{"Alphonse", "Gabriel", "Capone"} // good old Al

// file test helper function
func setupFile(t *testing.T, cmd string) (*Interpol, *InterpolatedString) {
	ip := NewInterpol()
	c, err := ip.Add(cmd)
	if err != nil {
		t.Fatalf("failed to add file: %v", err)
	}
	return ip, c
}

// testing linear
func TestFileLinearAll(t *testing.T) {
	ip, c := setupFile(t, "{{file, filename=data/file1 mode=linear count=0}}")

	testCompareAll(t, "file all-linear test",
		data1,
		testGetAll(ip, c))
}

func TestFileLinearSome(t *testing.T) {
	ip, c := setupFile(t, "{{file, filename=data/file1 mode=linear count=2}}")

	testCompareAll(t, "file some-linear test",
		data1[:2],
		testGetAll(ip, c))
}

func TestFileLinearMulti(t *testing.T) {
	ip, c := setupFile(t, "{{file, filename=data/file1 mode=linear count=5}}")

	testCompareAll(t, "file some.linear test",
		append(append(data1, data1[0]), data1[1]),
		testGetAll(ip, c))
}

// testing permutations
func TestFilePermAll(t *testing.T) {
	ip, c := setupFile(t, "{{file, filename=data/file1 mode=perm count=0}}")

	testOccurrence(t, "file perm-all", data1, 1, 1, 3, testGetAll(ip, c))
}

func TestFilePermSome(t *testing.T) {
	ip, c := setupFile(t, "{{file, filename=data/file1 mode=perm count=2}}")

	testOccurrence(t, "file perm-some", data1, 1, 1, 2, testGetAll(ip, c))
}

func TestFilePermMulti(t *testing.T) {
	ip, c := setupFile(t, "{{file, filename=data/file1 mode=perm count=5}}")

	testOccurrence(t, "file perm-multi", data1, 1, 2, 5, testGetAll(ip, c))
}

// testing random

func TestFileRandAll(t *testing.T) {
	ip, c := setupFile(t, "{{file, filename=data/file2 mode=rand count=0}}")
	testOccurrence(t, "file rand-all", data2, 0, 3, 3, testGetAll(ip, c))
}

func TestFileRandSome(t *testing.T) {
	ip, c := setupFile(t, "{{file, filename=data/file2 mode=rand count=2}}")
	testOccurrence(t, "file rand-some", data2, 0, 2, 2, testGetAll(ip, c))
}

func TestFileRandMulti(t *testing.T) {
	ip, c := setupFile(t, "{{file, filename=data/file2 mode=rand count=6}}")
	testOccurrence(t, "file rand-multi", data2, 0, 6, 6, testGetAll(ip, c))
}
