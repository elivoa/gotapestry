/*
  Time-stamp: <[start.go] Elivoa @ Thursday, 2013-08-08 14:43:27>

  Application Entrance: The New World starts here.

  TODO: Remove this file. Start the application in another way.

*/
package syd

import (
	"fmt"
	"got"
	"got/config"
)

// Start collects all module information and call start the system.
// Note: only pass module location here.
// TODO make the system startup better.
func Start() {
	// Startup-1: register modules. (Do not do others)
	config.Config.RegisterModulePath(SYDModule.Path(), "SYDModule")

	// start got
	got.BuildStart()
}

func welcome() {
	fmt.Println("\n")
	fmt.Println("``````````````````````````````````````````````````")
	fmt.Println("`  SYD Sale Management System (ALPHA 1)          `")
	fmt.Println("`                                                `")
	fmt.Println("``````````````````````````````````````````````````")
	// TODO config the port.
	fmt.Printf("Server Started, Listen localhost:%v\n\n", got.Config.Port)
	// got.PrintRegistry()
}
