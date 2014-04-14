/*
  Time-stamp: <[start.go] Elivoa @ Sunday, 2014-04-13 00:08:03>

  Application Entrance: The New World starts here.

  TODO: Remove this file. Start the application in another way.

*/
package syd

import (
	"fmt"
	"github.com/elivoa/got/config"
	"got"
)

// Start collects all module information and call start the system.
// Note: only pass module location here.
// TODO make the system startup better.
func Start() {

	// Startup-1: register modules. (Do not do others)

	// config.Config.RegisterModulePath
	// fmt.Println("009")
	// fmt.Println(SYDModule.Name)
	// fmt.Println(SYDModule.BasePath)
	// fmt.Println(SYDModule.PackagePath)
	// fmt.Println(SYDModule.Description)
	// fmt.Println(SYDModule.IsStartupModule)
	// fmt.Println(SYDModule.Register)
	// fmt.Println("009")

	config.Config.RegisterModule(SYDModule)
	// config.Config.RegisterModulePath(SYDModule.Path(), "SYDModule", SYDModule.IsStartupModule)

	// start got
	got.BuildStart() // build and start
}

func welcome() {
	fmt.Print("\n\n")
	fmt.Println("``````````````````````````````````````````````````")
	fmt.Println("`  SYD Sale Management System (ALPHA 1)          `")
	fmt.Println("`                                                `")
	fmt.Println("``````````````````````````````````````````````````")
	// TODO config the port.
	fmt.Printf("Server Started, Listen localhost:%v\n\n", got.Config.Port)
	// got.PrintRegistry()
}
