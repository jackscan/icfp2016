package main

import (
	"fmt"
	"testing"
)

func TestOffset(t *testing.T) {
	str := " " +
		"4\n" +
		"0,0\n" +
		"1,0\n" +
		"1/2,1/2\n" +
		"0,1/2\n"

	var p polygon
	p.parse(tokenize(str))

	var offset vec
	offset.x.SetString("3/2")
	offset.y.SetString("-1/2")

	p.translate(&offset)

	result := fmt.Sprintf("%s", p.String())

	expected := "" +
		"4:\n" +
		"	3/2, -1/2\n" +
		"	5/2, -1/2\n" +
		"	2, 0\n" +
		"	3/2, 0\n"

	if expected != result {
		t.Fail()
	}
}

func TestOverlap(t *testing.T) {
	str := " " +
		"4\n" +
		"0,0\n" +
		"1,0\n" +
		"1/2,1/2\n" +
		"0,1/2\n" +
		"4\n" +
		"0,0\n" +
		"0,-1\n" +
		"1,-1\n" +
		"1,0\n" +
		"4\n" +
		"1/4,0\n" +
		"1/2,-1/4\n" +
		"3/4,0\n" +
		"1/2,1/4\n"

	var tokens = tokenize(str)
	var p1, p2, p3 polygon
	pos := p1.parse(tokens)
	pos += p2.parse(tokens[pos:])
	p3.parse(tokens[pos:])

	if p1.overlaps(&p2) || p2.overlaps(&p1) {
		fmt.Println(p1.String(), p2.String())
		t.Fail()
	}

	if !p1.overlaps(&p1) || !p2.overlaps(&p2) {
		fmt.Println(p1.String(), p2.String())
		t.Fail()
	}

	if !p1.overlaps(&p3) || !p3.overlaps(&p1) {
		fmt.Println(p1.String(), p3.String())
		t.Fail()
	}

	if !p2.overlaps(&p3) || !p3.overlaps(&p2) {
		fmt.Println(p2.String(), p3.String())
		t.Fail()
	}
}

func TestUnitSquare(t *testing.T) {
	str := " " +
		"4\n" +
		"0,0\n" +
		"1,0\n" +
		"1/2,1/2\n" +
		"0,1/2\n" +
		"4\n" +
		"0,0\n" +
		"0,-1\n" +
		"1,-1\n" +
		"1,0\n" +
		"4\n" +
		"1/4,0\n" +
		"1/2,-1/4\n" +
		"3/4,0\n" +
		"1/2,1/4\n"

	var tokens = tokenize(str)
	var p1, p2, p3 polygon
	pos := p1.parse(tokens)
	pos += p2.parse(tokens[pos:])
	p3.parse(tokens[pos:])

	if !p1.inUnitSquare() {
		fmt.Println(p1.String())
		t.Fail()
	}

	if p2.inUnitSquare() {
		fmt.Println(p2.String())
		t.Fail()
	}

	if p3.inUnitSquare() {
		fmt.Println(p3.String())
		t.Fail()
	}

	var v1, v2 vec
	v1.x.SetString("1/2")
	v1.y.SetString("0")
	v2.x.SetString("3/2")
	v2.y.SetString("0")

	if !onUnitSquareEdge(&v1, &v2) {
		fmt.Println(v1.String(), v2.String())
		t.Fail()
	}

	v1.x.SetString("1/2")
	v1.y.SetString("1")
	v2.x.SetString("3/4")
	v2.y.SetString("1")

	if !onUnitSquareEdge(&v1, &v2) {
		fmt.Println(v1.String(), v2.String())
		t.Fail()
	}

	v1.x.SetString("0")
	v1.y.SetString("2/3")
	v2.x.SetString("0")
	v2.y.SetString("0")

	if !onUnitSquareEdge(&v1, &v2) {
		fmt.Println(v1.String(), v2.String())
		t.Fail()
	}

	v1.x.SetString("1")
	v1.y.SetString("2/3")
	v2.x.SetString("1")
	v2.y.SetString("0")

	if !onUnitSquareEdge(&v1, &v2) {
		fmt.Println(v1.String(), v2.String())
		t.Fail()
	}

	v1.x.SetString("0")
	v1.y.SetString("0")
	v2.x.SetString("1")
	v2.y.SetString("1/3")

	if onUnitSquareEdge(&v1, &v2) {
		fmt.Println(v1.String(), v2.String())
		t.Fail()
	}

	v1.x.SetString("1")
	v1.y.SetString("0")
	v2.x.SetString("1/2")
	v2.y.SetString("1/3")

	if onUnitSquareEdge(&v1, &v2) {
		fmt.Println(v1.String(), v2.String())
		t.Fail()
	}

	v1.x.SetString("0")
	v1.y.SetString("1")
	v2.x.SetString("1")
	v2.y.SetString("1/3")

	if onUnitSquareEdge(&v1, &v2) {
		fmt.Println(v1.String(), v2.String())
		t.Fail()
	}

	v1.x.SetString("1/2")
	v1.y.SetString("0")
	v2.x.SetString("1")
	v2.y.SetString("1")

	if onUnitSquareEdge(&v1, &v2) {
		fmt.Println(v1.String(), v2.String())
		t.Fail()
	}
}

func TestDot(t *testing.T) {
	var a, b vec
	a.x.SetString("3/2")
	a.y.SetString("-1/2")
	b.x.SetString("4")
	b.y.SetString("6")

	if fmt.Sprintf("%s", a.dot(&b).RatString()) != "3" {
		t.Fail()
	}
}

func TestVecMirror(t *testing.T) {
	var p, a, o vec
	o.x.SetString("3/2")
	o.y.SetString("-1/2")
	p.x.SetString("3")
	p.y.SetString("-1")
	a.x.SetString("3")
	a.y.SetString("3/2")

	p.mirror(&o, &a)

	if fmt.Sprintf("%s, %s", p.x.RatString(), p.y.RatString()) != "2, 1" {
		t.Fail()
	}

	a.set(1, 1)
	p.mirror(&o, &a)

	// fmt.Println("p: ", p.String())
	if fmt.Sprintf("%s, %s", p.x.RatString(), p.y.RatString()) != "3, 0" {
		t.Fail()
	}
}

func TestCounterclockwise(t *testing.T) {
	var a, b, c, d vec
	a.x.SetString("3/2")
	a.y.SetString("-1/2")
	b.x.SetString("5/2")
	b.y.SetString("-1")
	c.x.SetString("1/2")
	c.y.SetString("1")
	d.x.SetString("7/2")
	d.y.SetString("-3/2")

	if !(counterclockwise(&a, &b, &c) && counterclockwise(&b, &c, &a) && counterclockwise(&c, &a, &b)) {
		t.Fail()
	}

	if counterclockwise(&a, &c, &b) || counterclockwise(&c, &b, &a) || counterclockwise(&b, &a, &c) {
		t.Fail()
	}

	// test collinear
	if counterclockwise(&a, &b, &d) && counterclockwise(&b, &d, &a) && counterclockwise(&d, &a, &b) {
		t.Fail()
	}
}
