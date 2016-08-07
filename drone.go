package main

import (
	"bytes"
	"fmt"
	"math"
	"strings"
)

type drone struct {
	*problem

	srcfacets [][]int
	// srclines  []line
	srcpoints []vec
	dstpoints []vec
	polygons  []polygon

	indentstr string
	indent    int
}

func (dr *drone) solve(prob *problem) *solution {
	dr.problem = prob

	n := len(dr.problem.skeleton.lines)

	for i := 0; i < n; i++ {
		if s := dr.search(i, false, false); s != nil {
			return s
		}
		if s := dr.search(i, true, false); s != nil {
			return s
		}
		if s := dr.search(i, false, true); s != nil {
			return s
		}
		if s := dr.search(i, true, true); s != nil {
			return s
		}
	}
	return nil
}

func (dr *drone) reset() {
	dr.srcfacets = dr.srcfacets[:0]
	dr.srcpoints = dr.srcpoints[:0]
	dr.dstpoints = dr.dstpoints[:0]
	dr.polygons = dr.polygons[:0]
}

func (dr *drone) search(startline int, rotate, flip bool) *solution {

	dr.reset()

	dstline := dr.problem.skeleton.lines[startline]
	var v, d vec

	a := &dr.problem.skeleton.vertices[dstline.a]
	b := &dr.problem.skeleton.vertices[dstline.b]
	d.copy(b).sub(a)

	if rotate {
		v.copy(b)
		d.neg()
	} else {
		v.copy(a)
	}

	// var origin, v10 vec
	// v10.set(1, 0)

	// // transform startline to bottom square edge
	// v.neg()
	// prob.translate(&v)

	// fmt.Printf("d: %s\n", d.String())

	// calculate length of (a,b)
	dlen := d.dot(&d)
	{
		// weak spot of alogrithm
		// TODO: check inexact conversion
		f, _ := dlen.Float64()
		dlen.SetFloat64(math.Sqrt(f))
	}

	// add srcpoints and dstpoints for startline
	dr.dstpoints = make([]vec, 2)
	dr.dstpoints[0].copy(&v)
	dr.dstpoints[1].copy(&v).add(&d)
	if flip {
		dr.dstpoints[0], dr.dstpoints[1] = dr.dstpoints[1], dr.dstpoints[0]
	}

	dr.srcpoints = make([]vec, 2)
	dr.srcpoints[0].set(0, 0)
	dr.srcpoints[1].x.Set(dlen)
	dr.srcpoints[1].y.SetInt64(0)

	fmt.Println("problem:\n", dr.problem.String())
	fmt.Println("start:\n", dr.String())

	dr.indentstr = "    "
	if dr.addFacets(&line{0, 1}, flip) {
		return &solution{dr.srcpoints, dr.dstpoints, dr.srcfacets}
	}
	return nil
}

func (dr *drone) addFacets(srcline *line, flip bool) bool {

	// NOTE: srcfacet does not contain index for start and end twice
	// while dstfacet does

	dr.println("{")
	dr.indent++
	defer func() {
		dr.indent--
		dr.println("}")
	}()

	oldFacets := dr.srcfacets
	oldSrc := dr.srcpoints
	oldDst := dr.dstpoints
	oldPolygons := dr.polygons
	// oldLines := dr.srclines

	if flip {
		// TODO: check if we need to reverse the src line here
		srcline.a, srcline.b = srcline.b, srcline.a
	}

	dr.println("line", srcline.a, srcline.b)

	dstline := line{
		dr.problem.skeleton.findVertex(&dr.dstpoints[srcline.a]),
		dr.problem.skeleton.findVertex(&dr.dstpoints[srcline.b]),
	}

	dstfacet := make([]int, 2, len(dr.problem.skeleton.vertices))
	dstfacet[0] = dstline.a
	dstfacet[1] = dstline.b

	// destination facet
	dstfacet = dr.problem.skeleton.findFacet(0, dstfacet)
	for dstfacet != nil {

		dr.println("dstfacet:", linestrip2str(dstfacet))

		// add facet into src facet
		srcfacet := dr.addFacet(srcline, dstfacet, flip)
		if srcfacet != nil {

			dr.println("added facet:", linestrip2str(srcfacet))
			dr.print()

			n := len(srcfacet)
			i1, i2 := n-1, 0
			for ; i2 < n; i2++ {
				sline := line{srcfacet[i2], srcfacet[i1]}
				i1 = i2
				// check if line is already in other facets
				if facetsHaveLine(dr.srcfacets[:len(dr.srcfacets)-1], &sline) {
					dr.println("line", sline.a, sline.b, "already processed")
					continue
				}

				// check if line lies on square edge, except for initial line
				if onUnitSquareEdge(&dr.srcpoints[sline.a], &dr.srcpoints[sline.b]) {
					dr.println("line", sline.a, sline.b, "on square edge")
					continue
				}

				// dr.println("adding line", sline.a, sline.b)

				if !dr.addFacets(&sline, !flip) {
					dr.println("fail")
					break
				}
				// dr.srclines = append(dr.srclines, sline)
			}

			dr.println(i2, "/", n)

			if i2 == n {
				dr.println("succeeded")
				return true
			}

			// dr.println(dr.problem.skeleton.getPolygon(facet).String())
		}

		//reset drone
		dr.srcfacets = oldFacets
		dr.srcpoints = oldSrc
		dr.dstpoints = oldDst
		// dr.srclines = oldLines
		dr.polygons = oldPolygons

		dr.println("failed to add facet", linestrip2str(dstfacet))
		// dr.println(dr.String())

		// continue search with next dst facet
		dstfacet = dr.skeleton.findNextFacet(dstfacet)
	}

	dr.println("no more facets")

	return false
}

func (dr *drone) addFacet(srcline *line, dstfacet []int, flip bool) []int {

	// TODO: check if dst polygon overlaps reversed polygons for holes

	if flip {
		rev := make([]int, len(dstfacet))
		copy(rev, dstfacet)
		reverseFacet(rev)
		dstfacet = rev
		dr.println("reversed dstfacet:", linestrip2str(dstfacet))
	}

	p := dr.skeleton.getPolygon(dstfacet[:len(dstfacet)-1])

	dr.println("adding polygon")
	dr.println(p.StdString())

	var trans, axis, srca, srcd vec
	srca.copy(&dr.srcpoints[srcline.a])
	trans.copy(&srca).sub(&dr.dstpoints[srcline.a])
	dr.println("translate", trans.String())
	p.translate(&trans)

	srcd.copy(&dr.srcpoints[srcline.b]).sub(&srca)
	axis.copy(&srcd).add(&dr.dstpoints[srcline.b])
	axis.sub(&dr.dstpoints[srcline.a])

	if axis.isZero() {
		axis.x.Neg(&srcd.y)
		axis.y.Set(&srcd.x)
	}

	dr.println("transform", srca.String(), axis.String())
	p.mirror(&srca, &axis)

	if !flip {
		dr.println("t2", srca.String(), srcd.String())
		p.mirror(&srca, &srcd)
	} else {
		dr.println("flip")
	}

	dr.println("transformed")
	dr.println(p.StdString())

	// check polygon is inside square
	if !p.inUnitSquare() {
		dr.println("polygon is outside")
		return nil
	}

	// check overlap
	for i := range dr.polygons {
		if p.overlaps(&dr.polygons[i]) {
			dr.println("polygon overlaps")
			dr.println(dr.polygons[i].StdString())
			return nil
		}
	}

	srcfacet := make([]int, len(p.vertices))
	for i := range p.vertices {
		var added bool
		v := &p.vertices[i]
		srcfacet[i], added = dr.addPoint(v)
		if added {
			var dstv vec
			dstv.copy(v)
			// reverse transformation
			// dr.println("reverse transform", dstv.String())
			if !flip {
				dstv.mirror(&srca, &srcd)
				// dr.println("rev t2", dstv.String())
			}
			dstv.mirror(&srca, &axis)
			// dr.println("rev", dstv.String())
			dstv.sub(&trans)
			dr.println("adding dstpoint", dstv.String())
			dr.dstpoints = append(dr.dstpoints, dstv)
		}
	}

	// since polygon.overlaps is not correct check also if facet is new
	if facetsContain(dr.srcfacets, srcfacet) {
		dr.println("facet already known")
		return nil
	}

	// add facet
	dr.srcfacets = append(dr.srcfacets, srcfacet)
	dr.polygons = append(dr.polygons, *p)

	return srcfacet
}

func (dr *drone) addPoint(p *vec) (int, bool) {
	for i := range dr.srcpoints {
		if p.equals(&dr.srcpoints[i]) {
			return i, false
		}
	}
	dr.println("adding srcpoint", p.String())
	r := len(dr.srcpoints)
	dr.srcpoints = append(dr.srcpoints, *p)
	return r, true
}

func (dr *drone) String() string {
	var buf bytes.Buffer
	// buf.WriteString(fmt.Sprintf("problem:\n%s\n", dr.problem.String()))
	buf.WriteString(fmt.Sprintf("src:\n%d\n", len(dr.srcpoints)))
	for i := 0; i < len(dr.srcpoints); i++ {
		buf.WriteString(fmt.Sprintf("\t%s\t%d\n", dr.srcpoints[i].String(), i))
	}
	buf.WriteString(fmt.Sprintf("%d\n", len(dr.srcfacets)))
	for f := 0; f < len(dr.srcfacets); f++ {
		// buf.WriteString(fmt.Sprintf("\t%d:", len(dr.srcfacets[f])))
		for _, i := range dr.srcfacets[f] {
			buf.WriteString(fmt.Sprintf(" %d", i))
		}
		buf.WriteString("\n")
	}
	buf.WriteString("\ndst:\n")
	for i := 0; i < len(dr.dstpoints); i++ {
		buf.WriteString(fmt.Sprintf("\t%s\t%d\n", dr.dstpoints[i].String(), i))
	}

	// buf.WriteString("\nlines:\n")
	// for _, l := range dr.srclines {
	// 	buf.WriteString(fmt.Sprintf("\t%d - %d\n", l.a, l.b))
	// }

	return buf.String()
}

func (dr *drone) print() {
	// var buf bytes.Buffer
	// buf.WriteString(fmt.Sprintf("problem:\n%s\n", dr.problem.String()))
	dr.println(fmt.Sprintf("src %d:", len(dr.srcpoints)))
	for i := 0; i < len(dr.srcpoints); i++ {
		dr.println(fmt.Sprintf("~ %s\t%d", dr.srcpoints[i].String(), i))
	}
	dr.println(fmt.Sprintf("%d", len(dr.srcfacets)))
	for f := 0; f < len(dr.srcfacets); f++ {
		// buf.WriteString(fmt.Sprintf("\t%d:", len(dr.srcfacets[f])))
		dr.printIndent()
		fmt.Print("~")
		for _, i := range dr.srcfacets[f] {
			fmt.Printf(" %d", i)
		}

		fmt.Println("\t\t", dr.polygons[f].StdString())
	}
	// dr.println("polys:")
	// for p := range dr.polygons {
	// 	dr.println(fmt.Sprintf("~ %s", dr.polygons[i].String(), i))
	// }
	dr.println("dst:")
	for i := 0; i < len(dr.dstpoints); i++ {
		dr.println(fmt.Sprintf("~ %s\t%d", dr.dstpoints[i].String(), i))
	}

	// dr.println("lines:")
	// for _, l := range dr.srclines {
	// 	dr.println(fmt.Sprintf("~ %d - %d", l.a, l.b))
	// }
}

func (dr *drone) printIndent() {
	fmt.Print(strings.Repeat(dr.indentstr, dr.indent))
}

func (dr *drone) println(str ...interface{}) {
	dr.printIndent()
	fmt.Println(str...)
}
