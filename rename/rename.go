package rename

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/types"
	"log"
	"sort"

	"golang.org/x/tools/go/loader"

	"github.com/knzm/go-fixname/lint"
)

type Renamer struct {
	iprog        *loader.Program
	objsToUpdate map[types.Object]lint.Spec
	verbose      bool
	quiet        bool
	writeFunc    func(filename string, content []byte) error
}

func New(iprog *loader.Program) *Renamer {
	return &Renamer{
		iprog:        iprog,
		objsToUpdate: make(map[types.Object]lint.Spec),
	}
}

func (r *Renamer) Rename(obj types.Object, spec lint.Spec) {
	r.objsToUpdate[obj] = spec
}

func (r *Renamer) Update() error {
	writeFunc := r.writeFunc
	if writeFunc == nil {
		writeFunc = WriteFile
	}

	infolist := make([]*loader.PackageInfo, 0, len(r.iprog.Imported))
	for _, info := range r.iprog.Imported {
		infolist = append(infolist, info)
	}
	for _, info := range r.iprog.Created {
		infolist = append(infolist, info)
	}
	sort.Slice(infolist, func(i, j int) bool {
		return infolist[i].Pkg.Path() < infolist[j].Pkg.Path()
	})

	var nidents int
	filesToUpdate := make(map[string]*ast.File)
	for _, info := range infolist {
		astfileMap := make(map[string]*ast.File)
		for _, f := range info.Files {
			filename := r.iprog.Fset.File(f.Pos()).Name()
			astfileMap[filename] = f
		}

		processObjects := func(m map[*ast.Ident]types.Object) {
			for id, obj := range m {
				if spec, ok := r.objsToUpdate[obj]; ok {
					filename := r.iprog.Fset.File(id.Pos()).Name()
					filesToUpdate[filename] = astfileMap[filename]
					id.Name = spec.To
					nidents++
				}
			}
		}

		processObjects(info.Info.Defs)
		processObjects(info.Info.Uses)
	}

	var nerrs, npkgs int
	for _, info := range infolist {
		files := make([]*ast.File, len(info.Files))
		for i, f := range info.Files {
			files[i] = f
		}
		sort.Slice(files, func(i, j int) bool {
			return files[i].Pos() < files[j].Pos()
		})

		first := true
		for _, f := range files {
			filename := r.iprog.Fset.File(f.Pos()).Name()
			if filesToUpdate[filename] == nil {
				continue
			}

			if first {
				npkgs++
				first = false
				if r.verbose {
					log.Printf("Updating package %s", info.Pkg.Path())
				}
			}

			var buf bytes.Buffer
			if err := format.Node(&buf, r.iprog.Fset, f); err != nil {
				log.Printf("failed to pretty-print syntax tree: %v", err)
				nerrs++
				continue
			}

			if err := writeFunc(filename, buf.Bytes()); err != nil {
				log.Print(err)
				nerrs++
				continue
			}
		}
	}
	if !r.quiet {
		log.Printf("Renamed %s in %s in %s.",
			plural(nidents, "occurrence", "occurrences"),
			plural(len(filesToUpdate), "file", "files"),
			plural(npkgs, "package", "packages"))
	}

	if nerrs > 0 {
		return fmt.Errorf("failed to rewrite %s", plural(nerrs, "file", "files"))
	}

	return nil
}

func (r *Renamer) SetQuiet(quiet bool) {
	r.quiet = quiet
}

func (r *Renamer) SetVerbose(verbose bool) {
	r.verbose = verbose
}

func (r *Renamer) SetWriteFunc(writeFunc func(filename string, content []byte) error) {
	r.writeFunc = writeFunc
}
