package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/types"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/loader"

	"github.com/knzm/go-fixname/lint"
	"github.com/knzm/go-fixname/rename"
)

func containsHardErrors(errors []error) bool {
	for _, err := range errors {
		if err, ok := err.(types.Error); ok && err.Soft {
			continue
		}
		return true
	}
	return false
}

func loadProgram(pkgs map[string]bool, verbose bool) (*loader.Program, error) {
	conf := loader.Config{
		Build:      &build.Default,
		ParserMode: parser.ParseComments,

		// TODO(adonovan): enable this.  Requires making a lot of code more robust!
		AllowErrors: false,
	}
	// Optimization: don't type-check the bodies of functions in our
	// dependencies, since we only need exported package members.
	conf.TypeCheckFuncBodies = func(p string) bool {
		return pkgs[p] || pkgs[strings.TrimSuffix(p, "_test")]
	}

	if verbose {
		conf.AfterTypeCheck = func(info *loader.PackageInfo, files []*ast.File) {
			log.Println("Checked:", info)
		}
	}

	for pkg := range pkgs {
		conf.ImportWithTests(pkg)
	}

	// Ideally we would just return conf.Load() here, but go/types
	// reports certain "soft" errors that gc does not (Go issue 14596).
	// As a workaround, we set AllowErrors=true and then duplicate
	// the loader's error checking but allow soft errors.
	// It would be nice if the loader API permitted "AllowErrors: soft".
	conf.AllowErrors = true
	prog, err := conf.Load()
	if err != nil {
		return nil, err
	}

	var errpkgs []string
	// Report hard errors in indirectly imported packages.
	for _, info := range prog.AllPackages {
		if containsHardErrors(info.Errors) {
			errpkgs = append(errpkgs, info.Pkg.Path())
		}
	}
	if errpkgs != nil {
		var more string
		if len(errpkgs) > 3 {
			more = fmt.Sprintf(" and %d more", len(errpkgs)-3)
			errpkgs = errpkgs[:3]
		}
		return nil, fmt.Errorf("couldn't load packages due to errors: %s%s",
			strings.Join(errpkgs, ", "), more)
	}
	return prog, nil
}

type Option struct {
	inplace bool
	check   bool
	verbose bool
	filter  Filter
	args    []string
}

func Main(option *Option) error {
	pkgs := make(map[string]bool)
	for _, arg := range option.args {
		pkgs[arg] = true
	}

	iprog, err := loadProgram(pkgs, option.verbose)
	if err != nil {
		return err
	}

	renamer := rename.New(iprog)

	writeFunc := rename.Diff
	verbose := option.verbose
	quiet := true
	if option.check {
		writeFunc = func(filename string, content []byte) error {
			return nil
		}
		quiet = false
	} else if option.inplace {
		writeFunc = rename.WriteFile
		quiet = false
	}
	renamer.SetWriteFunc(writeFunc)
	renamer.SetVerbose(verbose)
	renamer.SetQuiet(quiet)

	infolist := make([]*loader.PackageInfo, 0, len(iprog.Imported))
	for _, info := range iprog.Imported {
		infolist = append(infolist, info)
	}
	for _, info := range iprog.Created {
		infolist = append(infolist, info)
	}
	for _, info := range infolist {
		for _, f := range info.Files {
			lint.WalkNames(iprog.Fset, f, func(id *ast.Ident, thing interface{}) {
				if !option.filter.byThing(thing) {
					return
				}
				if obj := info.Info.Defs[id]; obj != nil {
					if spec := lint.Check(id); spec != nil {
						if !option.filter.byCategory(spec.Category) {
							return
						}
						if option.check {
							pos := iprog.Fset.Position(id.Pos())
							filename := filepath.Base(pos.String())
							path := path.Join(info.Pkg.Path(), filename)
							fmt.Fprintf(os.Stderr, "%s: %s %s should be %s\n", path, thing, id, spec.To)
						}
						renamer.Rename(obj, *spec)
					}
				}
			})
		}
	}

	return renamer.Update()
}

var (
	flagInplace = flag.Bool("inplace", false, "edit in-place")
	flagCheck   = flag.Bool("check", false, "perform lint check")
	flagVerbose = flag.Bool("verbose", false, "show verbose messages")
	flagFilter  = flag.String("filter", "", "specify filter conditions by comma-separated string")
)

func init() {
	flag.Usage = func() {
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "A valid filter item is an element of the following sets:\n")
		fmt.Fprintf(os.Stderr, "  {caps, underscore}\n")
		fmt.Fprintf(os.Stderr, "  {const, var, type, struct field, func}\n")
	}
}

func parseFilter(str string) (*Filter, error) {
	var filter Filter
	for _, e := range strings.Split(str, ",") {
		e = strings.TrimSpace(e)
		if e == "" {
			continue
		}
		switch strings.ToLower(e) {
		case "caps":
			filter.category |= Caps
		case "underscore":
			filter.category |= Underscore
		case "const":
			filter.thing |= Const
		case "var":
			filter.thing |= Var
		case "type":
			filter.thing |= Type
		case "struct field":
			filter.thing |= StructField
		case "func":
			filter.thing |= Func
		default:
			return nil, fmt.Errorf("Unknown filter: %s", e)
		}
	}
	return &filter, nil
}

func ParseOption() *Option {
	flag.Parse()

	filter, err := parseFilter(*flagFilter)
	if err != nil {
		log.Fatal(err)
	}

	return &Option{
		inplace: *flagInplace,
		check:   *flagCheck,
		verbose: *flagVerbose,
		filter:  *filter,
		args:    flag.Args(),
	}
}

func main() {
	option := ParseOption()
	if err := Main(option); err != nil {
		log.Fatal(err)
	}
}
