package main

import (
	"bytes"
	"fmt"
	"math"
)

type drone struct {
	*problem

	srcfacets [][]int
	srclines  []line
	srcpoints []vec
	dstpoints []vec
	polygons  []polygon
}

func (dr *drone) search(prob *problem, startline int, rotate, flip bool) bool {
	dr.problem = prob
	dstline := prob.skeleton.lines[startline]
	var v, d vec

	a := &prob.skeleton.vertices[dstline.a]
	b := &prob.skeleton.vertices[dstline.b]
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
		// TODO: check inexact conversion
		f, _ := dlen.Float64()
		dlen.SetFloat64(math.Sqrt(f))
	}

	// d.div(dlen)
	// d.add(&v10)
	//
	// if d.isZero() {
	// 	d.set(0, 1)
	// }
	//
	// // fmt.Printf("%s, %s\n", v.String(), d.String())
	//
	// // d.set(1, 1)
	// prob.mirror(&origin, &d)
	//
	// if flip {
	// 	var two big.Rat
	// 	two.SetInt64(2)
	// 	origin.copy(a).add(b)
	// 	origin.div(&two)
	// 	prob.mirror(&origin, &v10)
	// }

	// // add srcpoints and dstpoints for startline
	// dr.dstpoints = make([]vec, 2)
	// dr.dstpoints[0].set(0, 0)
	// dr.dstpoints[1].x.Set(dlen)
	// dr.dstpoints[1].y.SetInt64(0)
	// dr.srcpoints = make([]vec, 2)
	// dr.srcpoints[0].copy(&dr.dstpoints[0])
	// dr.srcpoints[1].copy(&dr.dstpoints[1])

	// add srcpoints and dstpoints for startline
	dr.dstpoints = make([]vec, 2)
	dr.dstpoints[0].copy(&v)
	dr.dstpoints[1].copy(&v).add(&d)
	dr.srcpoints = make([]vec, 2)
	dr.srcpoints[0].set(0, 0)
	dr.srcpoints[1].x.Set(dlen)
	dr.srcpoints[1].y.SetInt64(0)

	fmt.Println("problem:\n", dr.problem.String())
	fmt.Println("start:\n", dr.String())

	// TODO: revert transformation on found solution
	return dr.addFacets(&line{0, 1}, false)
}

func (dr *drone) addFacets(srcline *line, flip bool) bool {

	// NOTE: srcfacet does not contain index for start and end twice
	// while dstfacet does

	oldFacets := dr.srcfacets
	oldSrc := dr.srcpoints
	oldDst := dr.dstpoints
	oldPolygons := dr.polygons

	if flip {
		// TODO: check if we need to reverse the src line here
		srcline.a, srcline.b = srcline.b, srcline.a
	}

	fmt.Println("line", srcline.a, srcline.b)

	// check if line is already in lines
	if findLine(dr.srclines, srcline) >= 0 {
		fmt.Println("line already processed")
		return true
	}

	// check if line lies on square edge, except for initial line
	if onUnitSquareEdge(&dr.srcpoints[srcline.a], &dr.srcpoints[srcline.b]) &&
		len(dr.srcfacets) > 0 {
		fmt.Println("line on square edge")
		return true
	}

	dr.srclines = append(dr.srclines, *srcline)
	oldLines := dr.srclines

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

		fmt.Println("dstfacet:", linestrip2str(dstfacet))

		// add facet into src facet
		srcfacet := dr.addFacet(srcline, dstfacet, flip)
		if srcfacet != nil {

			fmt.Println("added facet:", facetString(srcfacet))
			fmt.Println("new:", dr.String())

			n := len(srcfacet)
			i1 := n - 1
			for i2 := 0; i2 < n; i2++ {
				sline := line{srcfacet[i2], srcfacet[i1]}
				if !dr.addFacets(&sline, !flip) {
					break
				}
				i1 = i2
			}

			if i1 == n-1 {
				return true
			}

			// fmt.Println(dr.problem.skeleton.getPolygon(facet).String())
		}

		//reset drone
		dr.srcfacets = oldFacets
		dr.srcpoints = oldSrc
		dr.dstpoints = oldDst
		dr.srclines = oldLines
		dr.polygons = oldPolygons
		// continue search with next dst facet
		dstfacet = dr.skeleton.findNextFacet(dstfacet)
	}

	return false
}

func (dr *drone) addFacet(srcline *line, dstfacet []int, flip bool) []int {

	// if !flip {
	// 	rev := make([]int, len(dstfacet))
	// 	copy(rev, dstfacet)
	// 	reverseFacet(rev)
	// 	dstfacet = rev
	// 	fmt.Println("reversed dstfacet:", linestrip2str(dstfacet))
	// }

	p := dr.skeleton.getPolygon(dstfacet[:len(dstfacet)-1])

	fmt.Println("adding polygon\n", p.String())

	var trans, axis, srca, srcd vec
	srca.copy(&dr.srcpoints[srcline.a])
	trans.copy(&srca).sub(&dr.dstpoints[srcline.a])
	p.translate(&trans)

	srcd.copy(&dr.srcpoints[srcline.b]).sub(&srca)
	axis.copy(&srcd).add(&dr.dstpoints[srcline.b])
	axis.sub(&dr.dstpoints[srcline.a])

	if axis.isZero() {
		axis.x.Neg(&srcd.y)
		axis.y.Set(&srcd.x)
	}

	p.mirror(&srca, &axis)

	if !flip {
		p.mirror(&srca, &srcd)
	} else {
		fmt.Println("flip")
	}

	fmt.Println("transformed\n", p.String())

	// check polygon is inside square
	if !p.inUnitSquare() {
		fmt.Println("polygon is outside")
		return nil
	}

	// check overlap
	for i := range dr.polygons {
		if p.overlaps(&dr.polygons[i]) {
			fmt.Println("polygon overlaps", dr.polygons[i].String())
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
			if !flip {
				dstv.mirror(&srca, &srcd)
			}
			dstv.mirror(&srca, &axis)
			fmt.Println("adding dstpoint", dstv.String())
			dr.dstpoints = append(dr.dstpoints, dstv)
		}
	}

	// since polygon.overlaps is not correct check also if facet is new
	if facetsContain(dr.srcfacets, srcfacet) {
		fmt.Println("facet already known")
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
	fmt.Println("adding srcpoint", p.String())
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
		buf.WriteString(fmt.Sprintf("\t%d", len(dr.srcfacets[f])))
		for _, i := range dr.srcfacets[f] {
			buf.WriteString(fmt.Sprintf(" %d", i))
		}
		buf.WriteString("\n")
	}
	buf.WriteString("\ndst:\n")
	for i := 0; i < len(dr.dstpoints); i++ {
		buf.WriteString(fmt.Sprintf("\t%s\t%d\n", dr.dstpoints[i].String(), i))
	}

	return buf.String()
}
