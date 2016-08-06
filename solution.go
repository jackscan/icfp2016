package main

import (
	"bytes"
	"fmt"
)

type solution struct {
	src    []vec
	dst    []vec
	facets [][]int
}

func (s *solution) mirror(pos, axis *vec) {
	for i := range s.dst {
		s.dst[i].mirror(pos, axis)
	}
}

func (s *solution) translate(offset *vec) {
	for i := range s.dst {
		s.dst[i].add(offset)
	}
}

func (s *solution) mirrorAt(a, b int, indices []int) {
	var pos, axis vec

	pos.copy(&s.dst[a])
	axis.copy(&s.dst[b]).sub(&pos)

	// fmt.Println(pos.String(), axis.String())

	for _, i := range indices {
		s.dst[i].mirror(&pos, &axis)
	}
}

func (s *solution) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%d\n", len(s.src)))
	for i := 0; i < len(s.src); i++ {
		// buf.WriteString(fmt.Sprintf("\t%s\n", s.src[i].String()))
		buf.WriteString(fmt.Sprintf("\t%s\t%d\n", s.src[i].String(), i))
	}
	buf.WriteString(fmt.Sprintf("%d\n", len(s.facets)))
	for f := 0; f < len(s.facets); f++ {
		buf.WriteString(fmt.Sprintf("\t%d", len(s.facets[f])))
		for _, i := range s.facets[f] {
			buf.WriteString(fmt.Sprintf(" %d", i))
		}
		buf.WriteString("\n")
	}
	buf.WriteString("\n")
	for i := 0; i < len(s.dst); i++ {
		// buf.WriteString(fmt.Sprintf("\t%s\n", s.dst[i].String()))
		buf.WriteString(fmt.Sprintf("\t%s\t%d\n", s.dst[i].String(), i))
	}

	return buf.String()
}

func (s *solution) StdString() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%d\n", len(s.src)))
	for i := 0; i < len(s.src); i++ {
		buf.WriteString(fmt.Sprintf("%s,%s\n", s.src[i].x.RatString(), s.src[i].y.RatString()))
	}
	buf.WriteString(fmt.Sprintf("%d\n", len(s.facets)))
	for f := 0; f < len(s.facets); f++ {
		buf.WriteString(fmt.Sprintf("%d", len(s.facets[f])))
		for _, i := range s.facets[f] {
			buf.WriteString(fmt.Sprintf(" %d", i))
		}
		buf.WriteString("\n")
	}
	for i := 0; i < len(s.dst); i++ {
		buf.WriteString(fmt.Sprintf("%s,%s\n", s.dst[i].x.RatString(), s.dst[i].y.RatString()))
	}

	return buf.String()

}
