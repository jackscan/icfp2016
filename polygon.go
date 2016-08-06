package main

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
)

type vec struct {
	x, y big.Rat
}

type polygon struct {
	vertices []vec
}

type line struct {
	a, b int
}

type skeleton struct {
	vertices []vec
	lines    []line
}

func (v *vec) set(x, y int64) {
	v.x.SetInt64(x)
	v.y.SetInt64(y)
}

func (v *vec) copy(o *vec) *vec {
	v.x.Set(&o.x)
	v.y.Set(&o.y)
	return v
}

func (v *vec) add(o *vec) {
	v.x.Add(&v.x, &o.x)
	v.y.Add(&v.y, &o.y)
}

func (v *vec) sub(o *vec) {
	v.x.Sub(&v.x, &o.x)
	v.y.Sub(&v.y, &o.y)
}

func (v *vec) mul(s *big.Rat) {
	v.x.Mul(&v.x, s)
	v.y.Mul(&v.y, s)
}

func (v *vec) div(s *big.Rat) {
	v.x.Quo(&v.x, s)
	v.y.Quo(&v.y, s)
}

func (v *vec) neg() {
	v.x.Neg(&v.x)
	v.y.Neg(&v.y)
}

func (v *vec) dot(o *vec) *big.Rat {
	r := new(big.Rat)
	var tmp big.Rat
	r.Mul(&v.x, &o.x)
	tmp.Mul(&v.y, &o.y)
	return r.Add(r, &tmp)
}

func (v *vec) normalize() {
	len2 := v.dot(v)
	f, _ := len2.Float64()
	len2.SetFloat64(math.Sqrt(f))
	v.div(len2)
}

func (v *vec) isZero() bool {
	return v.x.Sign() == 0 && v.y.Sign() == 0
}

func (v *vec) equals(o *vec) bool {
	return v.x.Cmp(&o.x) == 0 && v.y.Cmp(&o.y) == 0
}

func clockwise(a, b, c *vec) bool {
	var ab, ac vec
	ab.copy(b).sub(a)
	ac.copy(c).sub(a)

	// ab.x * ac.y - ab.y * ac.x < 0
	ab.x.Mul(&ab.x, &ac.y)
	ab.y.Mul(&ab.y, &ac.x)
	ab.x.Sub(&ab.x, &ab.y)

	return ab.x.Sign() < 0
}

func counterclockwise(a, b, c *vec) bool {
	var ab, ac vec
	ab.copy(b).sub(a)
	ac.copy(c).sub(a)

	// ab.x * ac.y - ab.y * ac.x > 0
	ab.x.Mul(&ab.x, &ac.y)
	ab.y.Mul(&ab.y, &ac.x)
	ab.x.Sub(&ab.x, &ab.y)

	return ab.x.Sign() > 0
}

func (v *vec) mirror(o, a *vec) {
	a2 := a.dot(a)
	if a2.Num().Int64() == 0 {
		panic("invalid axis")
	}
	// fmt.Println("v, o, a: ", v.String(), o.String(), a.String())
	var d vec
	// fmt.Println("d: ", d.String())
	d.copy(v).sub(o)
	// fmt.Println("d-v: ", d.String())

	q := d.dot(a)
	q.Quo(q, a2)

	d.copy(a).mul(q)
	// fmt.Println("a*q: ", d.String())
	d.add(o)
	d.add(&d)
	d.sub(v)
	v.copy(&d)
	// fmt.Println("result: ", v.String())
}

func (l *line) equals(o *line) bool {
	return (l.a == o.a && l.b == o.b) || (l.b == o.a && l.a == o.b)
}

func (p *polygon) copy(o *polygon) {
	p.vertices = make([]vec, len(o.vertices))
	for i := range o.vertices {
		p.vertices[i].copy(&o.vertices[i])
	}
}

func (p *polygon) translate(offset *vec) {
	for i := range p.vertices {
		p.vertices[i].add(offset)
	}
}

func (p *polygon) mirror(origin, axis *vec) {
	for i := range p.vertices {
		p.vertices[i].mirror(origin, axis)
	}
}

func (p *polygon) pointIsInsideConvex(v *vec) bool {
	// assuming this polygon is defined counterclockwise
	n := len(p.vertices)
	if n < 2 {
		return false
	}
	i1 := n - 1
	for i2 := 0; i2 < n; i2++ {
		if !counterclockwise(&p.vertices[i2], v, &p.vertices[i1]) {
			return false
		}
		i1 = i2
	}
	return true
}

func (p *polygon) pointIsOutsideConvex(v *vec) bool {
	// assuming this polygon is defined counterclockwise
	n := len(p.vertices)
	if n < 2 {
		panic("invalid polygon")
	}
	i1 := n - 1
	for i2 := 0; i2 < n; i2++ {
		if counterclockwise(&p.vertices[i2], &p.vertices[i1], v) {
			return true
		}
		i1 = i2
	}
	return false
}

// func (p *polygon) isInsideConvex(o *polygon) bool {
// 	for _, v := range o.vertices {
// 		if !p.pointIsInsideConvex(&v) {
// 			return false
// 		}
// 	}
// 	return true
// }

func (p *polygon) overlaps(o *polygon) bool {
	// NOTE: assuming polygons are convex and defined counterclockwise
	// TODO: also assuming polygons do not overlap without containing vertices of other polygon

	anyOutside := false

	for i := range o.vertices {
		if p.pointIsInsideConvex(&o.vertices[i]) {
			return true
		} else if p.pointIsOutsideConvex(&o.vertices[i]) {
			anyOutside = true
		}
	}

	for i := range p.vertices {
		if o.pointIsInsideConvex(&p.vertices[i]) {
			return true
		} else if o.pointIsOutsideConvex(&p.vertices[i]) {
			anyOutside = true
		}
	}

	return !anyOutside
}

var zero = big.NewRat(0, 1)
var one = big.NewRat(1, 1)

func (p *polygon) inUnitSquare() bool {
	for _, v := range p.vertices {
		if v.x.Cmp(zero) < 0 || v.x.Cmp(one) > 0 ||
			v.y.Cmp(zero) < 0 || v.y.Cmp(one) > 0 {
			return false
		}
	}
	return true
}

func onUnitSquareEdge(a, b *vec) bool {
	var d vec
	d.copy(b).sub(a)
	return d.x.Sign() == 0 && (a.x.Cmp(one) == 0 || a.x.Sign() == 0) ||
		d.y.Sign() == 0 && (a.y.Cmp(one) == 0 || a.y.Sign() == 0)
}

func (s *skeleton) copy(o *skeleton) {
	s.lines = make([]line, len(o.lines))
	copy(s.lines, o.lines)

	for i := range o.vertices {
		s.vertices[i].copy(&o.vertices[i])
	}
}

func (s *skeleton) addVertex(v *vec) int {
	for i := range s.vertices {
		if v.equals(&s.vertices[i]) {
			return i
		}
	}
	r := len(s.vertices)
	s.vertices = append(s.vertices, *v)
	return r
}

func (s *skeleton) findVertex(v *vec) int {
	for i := range s.vertices {
		if v.equals(&s.vertices[i]) {
			return i
		}
	}
	return -1
}

func (s *skeleton) translate(offset *vec) {
	for i := range s.vertices {
		s.vertices[i].add(offset)
	}
}

func (s *skeleton) mirror(origin, axis *vec) {
	for i := range s.vertices {
		s.vertices[i].mirror(origin, axis)
	}
}

func (s *skeleton) findLine(l *line) int {
	return findLine(s.lines, l)
	// for i := range s.lines {
	// 	if l.equals(&s.lines[i]) {
	// 		return i
	// 	}
	// }
	// return -1
}

func (s *skeleton) findFacet(start int, prefix []int) []int {
	e := prefix[len(prefix)-1]

	// fmt.Println("start:", start, "prefix:", linestrip2str(prefix))

	if e == prefix[0] {
		return prefix
	}

linesloop:
	for i := start; i < len(s.lines); i++ {
		var next []int
		if s.lines[i].a == e {
			next = append(prefix, s.lines[i].b)
		} else if s.lines[i].b == e {
			next = append(prefix, s.lines[i].a)
		}

		if next != nil {
			n := len(next)

			// fmt.Println("check:", next[n-1])

			// check self intersect
			for j := 1; j < len(prefix); j++ {
				if prefix[j] == next[n-1] {
					continue linesloop
				}
			}

			// check counterclockwise
			a := s.vertices[next[n-3]]
			b := s.vertices[next[n-2]]
			c := s.vertices[next[n-1]]
			if !counterclockwise(&a, &b, &c) {
				// fmt.Println("clockwise:", next[n-3], next[n-2], next[n-1])
				continue
			}

			next = s.findFacet(0, next)
			if next != nil {
				return next
			}
		}
	}

	return nil
}

func (s *skeleton) findNextFacet(prev []int) []int {

	for n := len(prev); n > 2; n-- {
		start := s.findLine(&line{prev[n-2], prev[n-1]}) + 1

		// fmt.Println("start: ", start)

		next := s.findFacet(start, prev[:n-1])
		if next != nil {
			return next
		}
	}

	return nil
}

func (s *skeleton) getPolygon(facet []int) *polygon {
	var p polygon
	p.vertices = make([]vec, len(facet)-1)
	for i := 0; i < len(facet)-1; i++ {
		p.vertices[i].copy(&s.vertices[facet[i]])
	}
	return &p
}

func facetsEqual(a, b []int) bool {
	na, nb := len(a), len(b)
	if na != nb {
		return false
	}

	shift := -1
	for i := range b {
		if a[0] == b[i] {
			shift = i
			break
		}
	}

	if shift < 0 {
		return false
	}

	for i := range a {
		if a[i] != b[(i+shift)%nb] {
			return false
		}
	}

	return true
}

func facetsContain(facets [][]int, f []int) bool {
	for _, a := range facets {
		if facetsEqual(a, f) {
			return true
		}
	}
	return false
}

func reverseFacet(facet []int) {
	n := len(facet)
	for i := 0; i < n/2; i++ {
		j := n - 1 - i
		facet[i], facet[j] = facet[j], facet[i]
	}
}

func findLine(lines []line, l *line) int {
	for i := range lines {
		if l.equals(&lines[i]) {
			return i
		}
	}
	return -1
}

func linestrip2str(strip []int) string {
	var buf bytes.Buffer
	for i := 0; i < len(strip); i++ {
		buf.WriteString(fmt.Sprintf("%d ", strip[i]))
	}
	return buf.String()
}

func (v *vec) String() string {
	return fmt.Sprintf("%s, %s", v.x.RatString(), v.y.RatString())
}

func (p *polygon) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%d:\n", len(p.vertices)))
	for i := 0; i < len(p.vertices); i++ {
		buf.WriteString(fmt.Sprintf("\t%s, %s\n", p.vertices[i].x.RatString(), p.vertices[i].y.RatString()))
	}

	return buf.String()
}

func (s *skeleton) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("verts %d:\n", len(s.vertices)))
	for i := 0; i < len(s.vertices); i++ {
		buf.WriteString(fmt.Sprintf("\t%s\n", s.vertices[i].String()))
	}
	buf.WriteString(fmt.Sprintf("lines %d:\n", len(s.lines)))
	for i := 0; i < len(s.lines); i++ {
		buf.WriteString(fmt.Sprintf("\t%d - %d\n", s.lines[i].a, s.lines[i].b))
	}

	return buf.String()
}

func facetString(facet []int) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%d", len(facet)))
	for _, i := range facet {
		buf.WriteString(fmt.Sprintf(" %d", i))
	}
	return buf.String()
}
