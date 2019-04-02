package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/axw/gocov/gocovutil"
)

func main() {
	p, err := gocovutil.ReadPackages([]string{"-"})
	if err != nil {
		panic(err)
	}

	var cClasses, tClasses, cLines, tLines, cMethods, tMethods, cStmt, tStmt int

	for _, pkg := range p {
		tTypes := make(map[string]struct{})
		cTypes := make(map[string]struct{})
		for _, fn := range pkg.Functions {
			fnCovered := false
			for _, stmt := range fn.Statements {
				tStmt++
				if stmt.Reached > 0 {
					cLines += stmt.End - stmt.Start
					cStmt++
					fnCovered = true
				}
			}

			if parts := strings.Split(fn.Name, "."); len(parts) > 1 {
				tTypes[parts[0]] = struct{}{}
				if fnCovered {
					cTypes[parts[0]] = struct{}{}
				}
			}

			tLines += fn.End - fn.Start
			tMethods++
			if fnCovered {
				cMethods++
			}
		}
		tClasses += len(tTypes)
		cClasses += len(cTypes)
	}

	w := os.Stdout
	writeServiceMessage(w, "CodeCoverageAbsCCovered", cClasses)
	writeServiceMessage(w, "CodeCoverageAbsCTotal", tClasses)
	writeServiceMessage(w, "CodeCoverageAbsLCovered", cLines)
	writeServiceMessage(w, "CodeCoverageAbsLTotal", tLines)
	writeServiceMessage(w, "CodeCoverageAbsMCovered", cMethods)
	writeServiceMessage(w, "CodeCoverageAbsMTotal", tMethods)
	writeServiceMessage(w, "CodeCoverageAbsSCovered", cStmt)
	writeServiceMessage(w, "CodeCoverageAbsSTotal", tStmt)
}

func writeServiceMessage(w io.Writer, key string, value interface{}) {
	fmt.Fprintf(w, "##teamcity[buildStatisticValue key='%s' value='%v']\n", key, value)
}
