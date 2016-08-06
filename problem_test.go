package main

import (
	"fmt"
	"testing"
)

func TestStartTranformation(t *testing.T) {

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

	startline := 3
	rotate := false
	flip := true

	// fmt.Printf("%s\n", prob.String())

	var dr drone
	if !dr.search(&prob, startline, rotate, flip) {
		t.Fail()
	}

	fmt.Println("finish:", dr.String())

}
