/*
  Time-stamp: <[got.go] Elivoa @ Thursday, 2013-08-08 19:10:53>

  TODO:
    - Add Hooks: OnAppStart, AfterAppStart, ...
*/

package got

import (
	"fmt"
	"got/builtin"
	"got/config"
	"got/parser"
	"got/register"
	"got/route"
	"got/utils"
	"net/http"
)

// build phrase. only set config.
var Config *config.Configure

func init() {
	Config = config.Config
}

// BuildStart generates Start code and run server.
func BuildStart() {

	// register built-in module
	config.Config.RegisterModulePath(builtin.BuiltinModule.Path(), "BuiltinModule")

	printRegisteredModulePaths()

	// generate proton register sourcecode and compile and run.
	timer := utils.NewTimer()
	fmt.Println("> Generating startup codes...")
	app, err := parser.HackSource(Config.ModulePath)
	if err != nil {
		panic(fmt.Sprintf("build error: %v", err.Error()))
	}
	timer.Log("generating startup codes Done!")

	// start the server.
	app.Port = Config.Port
	app.Cmd().Run() // run and not return
}

func printRegisteredModulePaths() {
	// print registered modules.
	fmt.Println("> Registered Module paths:")
	for _, modulePath := range Config.ModulePath {
		fmt.Printf("    - module: %v.%v\n", modulePath.PackagePath, modulePath.Name)
	}
}

// called by generated server, start the server.
func Start() {
	welcome()

	// processing modules
	for key, module := range register.Modules.Map() {
		fmt.Println("> Register Module", key)
		if module.Register != nil {
			module.Register()
		}
	}

	// mapping static files
	for _, pair := range config.Config.StaticResources {
		http.Handle(pair[0],
			http.StripPrefix(pair[0], http.FileServer(http.Dir(pair[1]))),
		)
	}

	// got url matcher
	http.HandleFunc("/", route.RouteHandler)
	fmt.Println(">> got started...")
	http.ListenAndServe(fmt.Sprintf(":%v", Config.Port), nil)
}

func welcome() {
	fmt.Println("\n")
	fmt.Println("``````````````````````````````````````````````````")
	fmt.Println("`  GOT WebFramework     (EARLY BUILD 3)          `")
	fmt.Println("`                                                `")
	fmt.Println("``````````````````````````````````````````````````")
	fmt.Printf("Server Started, Listen localhost:%v\n\n", Config.Port)
	// PrintRegistry()
}

// ________________________________________________________________________________
// Print GOT Evnironment
//
func PrintRegistry() {
	register.Modules.PrintALL()

	fmt.Println("\n---- Pages ---------------------")
	register.Pages.PrintALL()

	fmt.Println("\n---- Components ---------------------")
	register.Components.PrintALL()

	fmt.Println("\n---- Mixins ---------------------")
	fmt.Println("... no mixins avaliable ...")

	fmt.Println("--------------------------------------------------------------------------------\n")
}
