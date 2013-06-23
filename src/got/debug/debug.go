package debug

import (
	"fmt"
	"log"
	"runtime/debug"
)

var (
	DebugLog   = true
	PrintStack = true

	FLAG_print_form_submit_details = true
)

// ________________________________________________________________________________
// Logging things
//
func Log(format string, params ...interface{}) {
	if DebugLog {
		log.Printf(format, params...)
	}
}

// ________________________________________________________________________________
// Error Handling
//
func Error(err error) error {
	log.Printf("~~~~~~ Error Occured ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~" +
		"~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	log.Printf(err.Error())
	log.Printf("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~" +
		"~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	if PrintStack {
		fmt.Println("StackTrace >>")
		debug.PrintStack()
		fmt.Println()
	}
	return err
}

func PrintError(err error) {
	log.Printf("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	log.Printf(err.Error())
	log.Printf("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
}

func PrintFormMap(name string, m map[string][]string) {
	fmt.Printf("\n---  [DEBUG PRINT]    (%v) --------\n", name)
	i := 1
	for k, v := range m {
		fmt.Printf("%3d: %-30v  -->  %v\n", i, k, v)
		i++
	}
}

// ________________________________________________________________________________
/*
   TODO: Finish this
*/
func PrintAllFieldTags(target interface{}, fieldName string) {
	fmt.Println("~~~~~ [GOT DEBUG TOOL] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	fmt.Printf("~ :: Reflect Tags of field '%v'\n", fieldName)
	fmt.Printf("~ \t key : value")
	fmt.Printf("~ \t key : value")
	fmt.Println("~~~~~ END ~~~~~...............")
}

func test() {
	fmt.Printf("+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+\n")
	fmt.Printf("+ %v\n", nil)
}
