package main

import (
	"fmt"
	"math/big"
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
	}
}

func TestOverlap2(t *testing.T) {
	str := "4  1/2,0  1/2,1/2  0,1  0,0\n" +
		"4  1,1  0,1  0,1/2  1/2,1/2"

	var tokens = tokenize(str)
	var p1, p2 polygon
	pos := p1.parse(tokens)
	p2.parse(tokens[pos:])

	if !p1.overlaps(&p2) || !p2.overlaps(&p1) {
		// fmt.Println(p1.String(), p2.String())
		t.Fail()
	}
}

func TestOverlap3(t *testing.T) {
	str := "" +
		"4  0,0 4,0 4,2 0,2\n" +
		"4  4,0 6,0 6,2 4,2\n" +
		"4  1,0 3,0 3,2 1,2\n" +
		"4  0,-1 4,-1 4,3 0,3\n" +
		"4  0,1 2,-1 4,1 2,3\n" +
		"4  1,-1 3,-1 3,3 1,3\n" +
		"4 3,0 3,-2 5,-2 5,0"

	var tokens = tokenize(str)
	var a, b, c, d, e, f, g polygon
	pos := a.parse(tokens)
	pos += b.parse(tokens[pos:])
	pos += c.parse(tokens[pos:])
	pos += d.parse(tokens[pos:])
	pos += e.parse(tokens[pos:])
	pos += f.parse(tokens[pos:])
	g.parse(tokens[pos:])

	// fmt.Println(a.String())
	// fmt.Println(b.String())
	// fmt.Println(c.String())
	// fmt.Println(d.String())

	if a.overlaps(&b) || b.overlaps(&a) {
		t.Fail()
	}
	if !a.overlaps(&c) || !c.overlaps(&a) {
		t.Fail()
	}
	if !d.overlaps(&e) || !e.overlaps(&d) {
		t.Fail()
	}
	if !a.overlaps(&f) || !f.overlaps(&a) {
		t.Fail()
	}
	if !d.overlaps(&f) || !f.overlaps(&d) {
		t.Fail()
	}
	if !d.overlaps(&c) || !c.overlaps(&d) {
		t.Fail()
	}
	if g.overlaps(&c) || c.overlaps(&g) {
		t.Fail()
	}
}

func TestOrthdot(t *testing.T) {
	var a, b vec
	a.set(4, 2)
	b.set(2, -1)

	var expected big.Rat
	expected.SetInt64(-8)

	if a.orthdot(&b).Cmp(&expected) != 0 {
		fmt.Println("orthdot:", a.orthdot(&b))
		t.Fail()
	}
}

func TestLinesIntersect(t *testing.T) {
	var a, b, c, d vec

	// test t-junction
	a.set(-2, -1)
	b.set(2, 1)
	c.set(0, 0)
	d.set(2, -1)

	if !linesIntersect(&b, &a, &c, &d) || !linesIntersect(&c, &d, &b, &a) {
		t.Fail()
	}

	if !linesIntersect(&b, &a, &d, &c) || !linesIntersect(&d, &c, &b, &a) {
		t.Fail()
	}

	if linesIntersect(&a, &b, &c, &d) || linesIntersect(&c, &d, &a, &b) {
		t.Fail()
	}

	if !linesIntersect(&b, &a, &c, &d) || !linesIntersect(&c, &d, &b, &a) {
		t.Fail()
	}

	// test overlap
	d.set(6, 3)
	if !linesIntersect(&a, &b, &c, &b) || !linesIntersect(&c, &b, &a, &b) {
		t.Fail()
	}

	if !linesIntersect(&a, &b, &c, &d) || !linesIntersect(&c, &d, &a, &b) {
		t.Fail()
	}

	if linesIntersect(&b, &a, &c, &d) || linesIntersect(&c, &d, &b, &a) {
		t.Fail()
	}

	if linesIntersect(&a, &c, &c, &d) || linesIntersect(&c, &d, &a, &c) {
		t.Fail()
	}

	// test parallel
	c.set(1, 0)
	d.set(3, 1)
	if linesIntersect(&a, &b, &c, &d) || linesIntersect(&c, &d, &a, &b) {
		t.Fail()
	}

	// test intersect
	if !linesIntersect(&a, &d, &c, &b) || !linesIntersect(&c, &b, &a, &d) {
		t.Fail()
	}
	if !linesIntersect(&d, &a, &c, &b) || !linesIntersect(&c, &b, &d, &a) {
		t.Fail()
	}

	// test separate
	d.set(3, 0)
	if linesIntersect(&a, &b, &c, &d) || linesIntersect(&c, &d, &a, &b) {
		t.Fail()
	}
	if linesIntersect(&b, &a, &c, &d) || linesIntersect(&c, &d, &b, &a) {
		t.Fail()
	}

	// test touching
	if linesIntersect(&a, &b, &b, &d) || linesIntersect(&b, &d, &a, &b) {
		t.Fail()
	}
	if linesIntersect(&b, &a, &b, &d) || linesIntersect(&b, &d, &b, &a) {
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

func TestReverseFacet(t *testing.T) {
	{
		facet := []int{23, 42, 2, 3, 4}
		reverseFacet(facet)

		expected := []int{4, 3, 2, 42, 23}

		for i := range facet {
			if facet[i] != expected[i] {
				fmt.Println("facet not reversed:", linestrip2str(facet))
				t.Fail()
			}
		}
	}

	{
		facet := []int{23, 42, 2, 3}
		reverseFacet(facet)

		expected := []int{3, 2, 42, 23}

		for i := range facet {
			if facet[i] != expected[i] {
				fmt.Println("facet not reversed:", linestrip2str(facet))
				t.Fail()
			}
		}
	}
}
