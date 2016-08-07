package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
)

func main() {

	var pname, sname string
	flag.StringVar(&pname, "p", "problem.txt", "problem file")
	flag.StringVar(&sname, "s", "solution.txt", "solution file to write")
	flag.Parse()

	pfile, err := os.Open(pname)
	if err != nil {
		log.Fatal(err)
	}

	probstr, err := ioutil.ReadAll(pfile)
	if err != nil {
		log.Fatal(err)
	}

	var prob problem
	prob.parse(string(probstr))

	// fmt.Printf("%s\n", string(probstr))

	var dr drone
	dr.debug = true
	s := dr.solve(&prob)
	if s == nil {
		log.Fatal("failed to find solution")
	} else {
		if s.inexact {
			sname += "-inexact"
		}
		sfile, err := os.Create(sname)
		if err != nil {
			log.Fatal(err)
		}
		defer sfile.Close()
		sfile.WriteString(s.StdString())
	}

	// fmt.Println("finished:", dr.String())

}
