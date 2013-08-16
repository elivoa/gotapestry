/*

TODO:
   - Make an util to generate file and run it.

*/

package parser

import (
	"bytes"
	"fmt"
	"github.com/robfig/revel"
	"go/build"
	"got/config"
	"got/core"
	"got/debug"
	"got/utils"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"text/template"
)

var importErrorPattern = regexp.MustCompile("cannot find package \"([^\"]+)\"")

// Build the app:
// 1. Generate the the main.go file.
// 2. Run the appropriate "go build" command.
// Requires that revel.Init has been called previously.
// Returns the path to the built binary, and an error if there was a problem building it.
func HackSource(modulePaths []*config.ModulePath) (app *App, compileError *Error) {
	if modulePaths == nil || len(modulePaths) == 0 {
		panic("Generating Error: No modules found!!!")
	}

	// First, clear the generated files (to avoid them messing with ProcessSource).
	cleanSource("generated")

	sourceInfo, compileError := ParseSource(modulePaths, true) // find only
	if compileError != nil {
		return nil, compileError
	}

	// set it to cache

	// // Add the db.import to the import paths.
	// if dbImportPath, found := revel.Config.String("db.import"); found {
	// 	sourceInfo.InitImportPaths = append(sourceInfo.InitImportPaths, dbImportPath)
	// }
	importPaths := make(map[string]string)
	// pageSpecs := []*StructInfo{}
	typeArrays := [][]*StructInfo{sourceInfo.Structs}
	for _, specs := range typeArrays {
		for _, spec := range specs {
			switch spec.ProtonKind {
			case core.PAGE, core.COMPONENT, core.MIXIN:
				addAlias(importPaths, spec.ImportPath, spec.PackageName)
			}
		}
	}

	modules := [][]string{}
	for _, p := range modulePaths {
		importPath := utils.PackagePath(p.PackagePath)
		packageName := filepath.Base(importPath)
		addAlias(importPaths, importPath, packageName)
		modules = append(modules, []string{packageName, p.Name})
	}

	// Generate two source files.
	data := map[string]interface{}{
		// "Controllers":    sourceInfo.ControllerSpecs(), // empty, leave it there
		// "ValidationKeys": sourceInfo.ValidationKeys,    // empty, levae it there
		"ImportPaths": importPaths,
		"Structs":     sourceInfo.Structs,
		"ModulePaths": modulePaths,
		"Modules":     modules,
		"ProtonKindLabel": map[core.Kind]string{
			core.PAGE:      "core.PAGE",
			core.COMPONENT: "core.COMPONENT",
			core.MIXIN:     "core.MIXIN",
		},
		// "TestSuites":     sourceInfo.TestSuites(),
	}
	genSource("generated", "main.go", MAIN, data)

	// // Read build config.
	// buildTags := revel.Config.StringDefault("build.tags", "")

	// Build the user program (all code under app).
	// It relies on the user having "go" installed.
	goPath, err := exec.LookPath("go")
	if err != nil {
		debug.Log("Go executable not found in PATH.")
	}

	pkg, err := build.Default.Import("got", "", build.FindOnly)
	if err != nil {
		debug.Log("Failure importing %v", revel.ImportPath)
	}
	binName := path.Join(pkg.BinDir, path.Base(config.Config.BasePath))
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	gotten := make(map[string]struct{})
	for {
		buildCmd := exec.Command(goPath, "build",
			// "-tags", buildTags,
			"-o", binName, path.Join("generated"))
		// debug.Log("Exec: %v", buildCmd.Args)
		output, err := buildCmd.CombinedOutput()

		// If the build succeeded, we're done.
		if err == nil {
			fmt.Println("bulid success, return")
			return NewApp(binName), nil
		}
		panic(string(output))

		// See if it was an import error that we can go get.
		matches := importErrorPattern.FindStringSubmatch(string(output))
		if matches == nil {
			return nil, newCompileError(output)
		}

		// Ensure we haven't already tried to go get it.
		pkgName := matches[1]
		if _, alreadyTried := gotten[pkgName]; alreadyTried {
			return nil, newCompileError(output)
		}
		gotten[pkgName] = struct{}{}

		// Execute "go get <pkg>"
		getCmd := exec.Command(goPath, "get", pkgName)
		// debug.Log("Exec: ", getCmd.Args)
		getOutput, err := getCmd.CombinedOutput()
		if err != nil {
			panic(string(getOutput))
			// revel.ERROR.Println(string(getOutput))
			return nil, newCompileError(output)
		}

		// Success getting the import, attempt to build again.
	}

	panic("Not reachable")
	return nil, nil
}

/// add by elivoa
func calcImports(src *SourceInfo) map[string]string {
	aliases := make(map[string]string)
	typeArrays := [][]*StructInfo{src.Structs /*, src.TestSuites()*/}
	for _, specs := range typeArrays {
		for _, spec := range specs {
			// fmt.Println("  > i:", spec.ImportPath, " o:", spec.PackageName)
			switch spec.ProtonKind {
			case core.PAGE, core.COMPONENT, core.MIXIN:
				addAlias(aliases, spec.ImportPath, spec.PackageName)
			}

			//# method imports, don't import this.
			//#
			// for _, methSpec := range spec.MethodSpecs {
			// 	for _, methArg := range methSpec.Args {
			// 		if methArg.ImportPath == "" {
			// 			continue
			// 		}

			// 		addAlias(aliases, methArg.ImportPath, methArg.TypeExpr.PkgName)
			// 	}
			// }
		}
	}

	// Add the "InitImportPaths", with alias "_"
	// for _, importPath := range src.InitImportPaths {
	// 	if _, ok := aliases[importPath]; !ok {
	// 		aliases[importPath] = "_"
	// 	}
	// }

	return aliases
}

////@ cleared
func cleanSource(dirs ...string) {
	for _, dir := range dirs {
		tmpPath := path.Join(config.Config.SrcPath, dir)
		err := os.RemoveAll(tmpPath)
		if err != nil {
			revel.ERROR.Println("Failed to remove dir:", err)
		}
	}
}

//// @cleaned
// genSource renders the given template to produce source code, which it writes
// to the given directory and file.
func genSource(dir, filename, templateSource string, args map[string]interface{}) {
	// generate source
	sourceCode := ExecuteTemplate(
		template.Must(template.New("").Parse(templateSource)),
		args)

	// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	// fmt.Println(sourceCode)
	// fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")

	// Create a fresh dir.
	tmpPath := path.Join(config.Config.SrcPath, dir)

	err := os.RemoveAll(tmpPath)
	if err != nil {
		revel.ERROR.Println("Failed to remove dir:", err)
	}
	err = os.Mkdir(tmpPath, 0777)
	if err != nil {
		revel.ERROR.Fatalf("Failed to make tmp directory: %v", err)
	}

	// Create the file
	file, err := os.Create(path.Join(tmpPath, filename))
	defer file.Close()
	if err != nil {
		revel.ERROR.Fatalf("Failed to create file: %v", err)
	}
	_, err = file.WriteString(sourceCode)
	if err != nil {
		revel.ERROR.Fatalf("Failed to write to file: %v", err)
	}
}

// Execute a template and returns the result as a string.
func ExecuteTemplate(tmpl revel.ExecutableTemplate, data interface{}) string {
	var b bytes.Buffer
	if err := tmpl.Execute(&b, data); err != nil {
		panic(err.Error())
	}
	return b.String()
}

// Looks through all the method args and returns a set of unique import paths
// that cover all the method arg types.
// Additionally, assign package aliases when necessary to resolve ambiguity.
// func calcImportAliases(src *SourceInfo) map[string]string {
// 	aliases := make(map[string]string)
// 	typeArrays := [][]*StructInfo{src.ControllerSpecs() /*, src.TestSuites()*/}
// 	for _, specs := range typeArrays {
// 		for _, spec := range specs {
// 			addAlias(aliases, spec.ImportPath, spec.PackageName)

// 			for _, methSpec := range spec.MethodSpecs {
// 				for _, methArg := range methSpec.Args {
// 					if methArg.ImportPath == "" {
// 						continue
// 					}

// 					addAlias(aliases, methArg.ImportPath, methArg.TypeExpr.PkgName)
// 				}
// 			}
// 		}
// 	}

// 	// Add the "InitImportPaths", with alias "_"
// 	for _, importPath := range src.InitImportPaths {
// 		if _, ok := aliases[importPath]; !ok {
// 			aliases[importPath] = "_"
// 		}
// 	}

// 	return aliases
// }

func addAlias(aliases map[string]string, importPath, pkgName string) {
	alias, ok := aliases[importPath]
	if ok {
		return
	}
	alias = makePackageAlias(aliases, pkgName)
	aliases[importPath] = alias
}

func makePackageAlias(aliases map[string]string, pkgName string) string {
	i := 0
	alias := pkgName
	for containsValue(aliases, alias) {
		alias = fmt.Sprintf("%s%d", pkgName, i)
		i++
	}
	return alias
}

func containsValue(m map[string]string, val string) bool {
	for _, v := range m {
		if v == val {
			return true
		}
	}
	return false
}

// Parse the output of the "go build" command.
// Return a detailed Error.
func newCompileError(output []byte) *Error {
	errorMatch := regexp.MustCompile(`(?m)^([^:#]+):(\d+):(\d+:)? (.*)$`).
		FindSubmatch(output)
	if errorMatch == nil {
		revel.ERROR.Println("Failed to parse build errors:\n", string(output))
		return &Error{
			SourceType:  "Go code",
			Title:       "Go Compilation Error",
			Description: "See console for build error.",
		}
	}

	// Read the source for the offending file.
	var (
		relFilename    = string(errorMatch[1]) // e.g. "src/revel/sample/app/controllers/app.go"
		absFilename, _ = filepath.Abs(relFilename)
		line, _        = strconv.Atoi(string(errorMatch[2]))
		description    = string(errorMatch[4])
		compileError   = &Error{
			SourceType:  "Go code",
			Title:       "Go Compilation Error",
			Path:        relFilename,
			Description: description,
			Line:        line,
		}
	)

	fileStr, err := revel.ReadLines(absFilename)
	if err != nil {
		compileError.MetaError = absFilename + ": " + err.Error()
		revel.ERROR.Println(compileError.MetaError)
		return compileError
	}

	compileError.SourceLines = fileStr
	return compileError
}

const MAIN = `// DO NOT EDIT THIS FILE -- GENERATED CODE
package main

import (
    "fmt"
    _got "got"
    "got/config"
    "got/register"
    "got/cache"
    "got/parser"
    "got/route"
	{{range $k, $v := $.ImportPaths}}
    {{$v}} "{{$k}}"{{end}}
)

func main() {
    fmt.Println("=============== STARTING ================================================")

    // restore config.ModulePath
    {{range .ModulePaths}}
    config.Config.RegisterModulePath("{{.PackagePath}}", "{{.Name}}"){{end}}

    // parse source again.
    sourceInfo, compileError := parser.ParseSource(config.Config.ModulePath, false) // deep parse
    if compileError != nil {
        panic(compileError.Error())
    }

	// cache source info into cache.
    cache.SourceCache = sourceInfo

    // register real module
    register.RegisterModule(
    {{range .Modules}}
      {{index . 0}}.{{index . 1}},{{end}}
    )

    // register pages & components
    {{range .Structs}}{{if .IsProton}}
    route.RegisterProton("{{.ImportPath}}", "{{.StructName}}", "{{.ModulePackage}}", &{{index $.ImportPaths .ImportPath}}.{{.StructName}}{}){{end}}{{end}}

    // start the server
    _got.Start()
}
`

// // cache StructCache
// if !findOnly {
// 	fmt.Println("********************************************************************************")
// 	for _, si := range srcInfo.Structs {
// 		// cache.StructCache.GetCreate(si.IsProton)
// 		fmt.Println(si)
// 	}
// }
