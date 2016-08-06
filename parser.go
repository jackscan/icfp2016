package main

import (
	"fmt"
	"math/big"
	"regexp"
	"strconv"
)

func tokenize(str string) []string {
	sep := "[,\t\r\n ]+"
	sepre := regexp.MustCompile(sep)
	tokens := sepre.Split(str, -1)

	// drop leading empty tok
	if len(tokens[0]) == 0 {
		tokens = tokens[1:]
	}
	return tokens
}

func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func parseRat(r *big.Rat, s string) {
	_, err := fmt.Sscan(s, r)
	if err != nil {
		panic(err)
	}
}

func (p *polygon) parse(tokens []string) int {

	numverts := parseInt(tokens[0])
	p.vertices = make([]vec, numverts)

	pos := 1
	for iv := 0; iv < numverts; iv++ {
		parseRat(&p.vertices[iv].x, tokens[pos])
		pos++
		parseRat(&p.vertices[iv].y, tokens[pos])
		pos++
	}

	return pos
}

func (s *skeleton) parse(tokens []string) int {
	numlines := parseInt(tokens[0])
	s.lines = make([]line, numlines)

	pos := 1
	for il := 0; il < numlines; il++ {
		var a, b vec
		parseRat(&a.x, tokens[pos])
		parseRat(&a.y, tokens[pos+1])
		parseRat(&b.x, tokens[pos+2])
		parseRat(&b.y, tokens[pos+3])

		s.lines[il].a = s.addVertex(&a)
		s.lines[il].b = s.addVertex(&b)

		pos += 4
	}

	return pos
}

func (p *problem) parse(str string) {
	tokens := tokenize(str)

	numpolys := parseInt(tokens[0])

	p.polygons = make([]polygon, numpolys)

	pos := 1
	for ip := 0; ip < numpolys; ip++ {
		pos += p.polygons[ip].parse(tokens[pos:])
	}

	p.skeleton.parse(tokens[pos:])
}

func (s *solution) parseIncomplete(tokens []string) {
	numverts := parseInt(tokens[0])
	s.src = make([]vec, numverts)
	s.dst = make([]vec, numverts)

	pos := 1
	for iv := 0; iv < numverts; iv++ {
		parseRat(&s.src[iv].x, tokens[pos])
		pos++
		parseRat(&s.src[iv].y, tokens[pos])
		pos++
	}

	numfacets := parseInt(tokens[pos])
	s.facets = make([][]int, numfacets)
	pos++

	for f := 0; f < numfacets; f++ {
		numedges := parseInt(tokens[pos])
		s.facets[f] = make([]int, numedges)
		pos++

		for e := 0; e < numedges; e++ {
			s.facets[f][e] = parseInt(tokens[pos])
			pos++
		}
	}

	for iv := 0; iv < numverts; iv++ {
		s.dst[iv].copy(&s.src[iv])
	}
}
