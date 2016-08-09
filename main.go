package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

func main() {

	var pname, sname string
	var probdir, soldir string
	var debug bool
	flag.StringVar(&pname, "p", "problem.txt", "problem file")
	flag.StringVar(&sname, "s", "solution.txt", "solution file to write")
	flag.BoolVar(&debug, "d", false, "print debug info")
	flag.StringVar(&probdir, "probdir", "", "solve problems in directory")
	flag.StringVar(&soldir, "soldir", "solutions", "folder to write solutions to")

	flag.Parse()

	if len(probdir) > 0 {
		dir, err := os.Open(probdir)
		defer dir.Close()
		if err != nil {
			log.Fatal(err)
		}

		names, err := dir.Readdirnames(-1)
		if err != nil {
			log.Fatal(err)
		}

		var wg sync.WaitGroup
		for _, filename := range names {
			ppath := probdir + "/" + filename
			spath := strings.TrimPrefix(filename, "problem-")
			spath = soldir + "/solution-" + spath

			log.Print("reading", ppath)
			prob := readProblem(ppath)

			wg.Add(1)
			go func(prob *problem, sname string) {
				defer wg.Done()
				if solve(prob, sname, false) {
					log.Print(ppath, "->", sname)
				} else {
					log.Print("failed to solve", ppath)
				}
			}(prob, spath)
		}

		wg.Wait()

		// fmt.Println(names)

	} else {

		prob := readProblem(pname)
		if solve(prob, sname, debug) {
			log.Println("solved", pname)
		} else {
			log.Fatal("failed to find solution")
		}
	}
}

func readProblem(pname string) *problem {
	pfile, err := os.Open(pname)
	if err != nil {
		log.Fatal(err)
	}
	defer pfile.Close()

	probstr, err := ioutil.ReadAll(pfile)
	if err != nil {
		log.Fatal(err)
	}

	var prob problem
	prob.parse(string(probstr))
	prob.skeleton.addIntersections()
	return &prob
}

func solve(prob *problem, sname string, debug bool) bool {

	var dr drone
	dr.debug = debug
	s := dr.solve(prob)
	if s == nil {
		return false
		// log.Fatal("failed to find solution")
	}

	if s.incomplete {
		sname += "-incomplete"
	}
	sfile, err := os.Create(sname)
	if err != nil {
		log.Fatal(err)
	}
	defer sfile.Close()
	sfile.WriteString(s.StdString())
	return true

}
