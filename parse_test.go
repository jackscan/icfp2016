package main

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	str := " " +
		"1\n" +
		"4\n" +
		"0,0\n" +
		"1,0\n" +
		"1/2,1/2\n" +
		"0,1/2\n" +
		"5\n" +
		"0,0 1,0\n" +
		"1,0 1/2,1/2\n" +
		"1/2,1/2 0,1/2\n" +
		"0,1/2 0,0\n" +
		"0,0 1/2,1/2\n"

	var prob problem

	prob.parse(str)

	result := ""

	for _, p := range prob.polygons {
		result += fmt.Sprintf("%s", p.String())
	}

	result += fmt.Sprintf("%s", prob.skeleton.String())

	expected := "" +
		"4:\n" +
		"	0, 0\n" +
		"	1, 0\n" +
		"	1/2, 1/2\n" +
		"	0, 1/2\n" +
		"verts 4:\n" +
		"	0, 0\n" +
		"	1, 0\n" +
		"	1/2, 1/2\n" +
		"	0, 1/2\n" +
		"lines 5:\n" +
		"	0 - 1\n" +
		"	1 - 2\n" +
		"	2 - 3\n" +
		"	3 - 0\n" +
		"	0 - 2\n"

	if expected != result {
		fmt.Println(result)
		t.Fail()
	}
}
