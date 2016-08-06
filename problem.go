package main

import (
	"bytes"
	"fmt"
)

type problem struct {
	polygons []polygon
	skeleton
}

func (p *problem) copy(o *problem) {
	p.polygons = make([]polygon, len(o.polygons))
	for i := range o.polygons {
		p.polygons[i].copy(&o.polygons[i])
	}

	p.skeleton.copy(&o.skeleton)
}

func (p *problem) translate(offset *vec) {
	for i := range p.polygons {
		p.polygons[i].translate(offset)
	}
	p.skeleton.translate(offset)
}

func (p *problem) mirror(origin, axis *vec) {
	for i := range p.polygons {
		p.polygons[i].mirror(origin, axis)
	}
	p.skeleton.mirror(origin, axis)
}

func (p *problem) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("polygons %d:\n", len(p.polygons)))
	for i := 0; i < len(p.polygons); i++ {
		buf.WriteString(p.polygons[i].String())
	}

	buf.WriteString(fmt.Sprintf("skeleton:\n"))
	buf.WriteString(p.skeleton.String())

	return buf.String()
}
