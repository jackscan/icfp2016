package main

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func donTestPrintTimes(t *testing.T) {
	t0, err := time.Parse("Jan 2, 2006 03:04 (MST)", "Aug 6, 2016 07:00 (UTC)")
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 46; i++ {
		fmt.Println(t0.Unix())
		t0 = t0.Add(time.Hour)
	}
}

func TestBuildSolution1(t *testing.T) {

	// str := "" +
	// 	"6\n" +
	// 	"0,0\n" +
	// 	"1,0\n" +
	// 	"1,3/4\n" +
	// 	"0,1/4\n" +
	// 	"1,1\n" +
	// 	"0,1\n"

	str := `18
        0,1 1/3,1 2/3,1 13/16,1 1,1
        88/100,91/100
        0,3/4 1,3/4
        0,1/2 1,1/2
        0,1/4 1,1/4
        12/100,9/100
        0,0 3/16,0 1/3,0 2/3,0 1,0
        11
        3 0 6 1
        4 6 8 2 1
        5 8 10 5 3 2
        4 10 12 7 5
        5 12 14 15 9 7
        4 15 16 11 9
        3 16 17 11
        3 3 5 4
        3 5 7 4
        3 10 13 12
        3 13 14 12
    `

	tokens := tokenize(str)

	var s solution

	s.parseIncomplete(tokens)

	var v1, v2 vec

	v1.x.SetString("1")
	v1.y.SetString("3/4")
	v2.x.SetString("5/4")
	v2.y.SetString("0")

	v2.add(&v1)
	v1.set(0, 0)

	// fmt.Println(v1.String(), v2.String())
	s.mirror(&v1, &v2)

	v1.x.SetString("-3/20")
	v1.y.SetString("0")
	s.translate(&v1)
	// fmt.Println(s.String())

	s.mirrorAt(10, 14, []int{13})
	s.mirrorAt(3, 7, []int{4})
	s.mirrorAt(12, 7, []int{0, 1, 2, 3, 4, 5, 6, 8, 10})
	s.mirrorAt(10, 5, []int{0, 1, 2, 3, 6, 8})
	s.mirrorAt(8, 2, []int{0, 1, 6})
	s.mirrorAt(1, 6, []int{0})
	s.mirrorAt(15, 9, []int{11, 16, 17})
	s.mirrorAt(16, 11, []int{17})

	file, err := os.Create("problem1.txt")
	if err != nil {
		panic(err)
	}
	file.WriteString(s.StdString())
	file.Close()

	file, err = os.Create("problem2.txt")
	if err != nil {
		panic(err)
	}

	v1.x.SetString("0")
	v1.y.SetString("0")
	v2.x.SetString("1")
	v2.y.SetString("1")
	s.mirror(&v1, &v2)

	file.WriteString(s.StdString())
	file.Close()

}

func TestBuildSolution2(t *testing.T) {

	// str := "" +
	// 	"6\n" +
	// 	"0,0\n" +
	// 	"1,0\n" +
	// 	"1,3/4\n" +
	// 	"0,1/4\n" +
	// 	"1,1\n" +
	// 	"0,1\n"

	str := `26
        0,1 1/3,1 0,3/4
        2/3,1 0,1/2
        1,1 0,1/4
        1,3/4 0,0
        1,1/2 1/3,0
        1,1/4 2/3,0 1,0
        13/50,189/200
        13/50,139/200
        1/2,5/8
        1/2,3/8
        37/50,61/200
        37/50,11/200
        13/16,1
        88/100,91/100
        12/100,9/100
        3/16,0
        1/14,1
        13/14,0
        18
        3 24 14 1
        4 0 2 14 24
        4 2 4 15 14
        4 1 14 15 3
        4 4 6 16 15
        5 3 15 16 21 20
        4 6 22 17 16
        4 16 17 7 21
        5 17 22 23 10 18
        4 7 17 18 9
        4 10 12 19 18
        4 9 18 19 11
        3 12 25 19
        4 11 19 25 13
        3 5 20 21
        3 5 21 7
        3 6 8 22
        3 8 23 22
    `

	tokens := tokenize(str)

	var s solution

	s.parseIncomplete(tokens)

	// var v16, v17, v3, v4, v5, v6 vec
	//
	// v3.x.SetString("2/3")
	// v3.y.SetString("1")
	// v4.x.SetString("0")
	// v4.y.SetString("1/2")
	//
	// v3.sub(&v4)
	//
	// v16.x.SetString("1/2")
	// v16.y.SetString("5/8")
	//
	// fmt.Println(v4.String(), v3.String())
	// v16.mirror(&v4, &v3)
	// fmt.Println("v14: ", v16.String())
	//
	// v5.x.SetString("1")
	// v5.y.SetString("1")
	// v6.x.SetString("0")
	// v6.y.SetString("1/4")
	//
	// v5.sub(&v6)
	//
	// v17.x.SetString("1/2")
	// v17.y.SetString("3/8")
	//
	// v17.mirror(&v6, &v5)
	// fmt.Println("v15: ", v17.String())

	s.mirrorAt(1, 2, []int{0, 24})
	s.mirrorAt(3, 4, []int{0, 1, 2, 14, 24})
	s.mirrorAt(5, 6, []int{0, 1, 2, 3, 4, 14, 15, 20, 24})
	s.mirrorAt(7, 8, []int{9, 10, 11, 12, 13, 18, 19, 23, 25})
	s.mirrorAt(9, 10, []int{11, 12, 13, 19, 25})
	s.mirrorAt(11, 12, []int{13, 25})

	s.mirrorAt(16, 17, []int{1, 3, 5, 7, 9, 11, 13, 18, 19, 20, 21, 24})

	s.mirrorAt(7, 20, []int{5})
	s.mirrorAt(6, 23, []int{8})

	var v1, v2 vec

	v1.x.SetString("1")
	v1.y.SetString("3/4")
	v2.x.SetString("5/4")
	v2.y.SetString("0")

	v2.add(&v1)
	v1.set(0, 0)

	s.mirror(&v1, &v2)

	v1.x.SetString("-3/20")
	v1.y.SetString("0")
	s.translate(&v1)
	// fmt.Println(s.String())

	file, err := os.Create("problem3.txt")
	if err != nil {
		panic(err)
	}
	file.WriteString(s.StdString())
	file.Close()
	//
	// file, err = os.Create("problem2.txt")
	// if err != nil {
	// 	panic(err)
	// }
	//
	// v1.x.SetString("0")
	// v1.y.SetString("0")
	// v2.x.SetString("1")
	// v2.y.SetString("1")
	// s.mirror(&v1, &v2)
	//
	// file.WriteString(s.StdString())
	// file.Close()

}
