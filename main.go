package main

import (
	"flag"
	"kpt-look-alike/functions"
	"log"
	"os"
)

func main() {

	flag.Parse()
	args := flag.Args()

	if !(args[0] == "apply" && args[1] == "fn") {
		os.Exit(-1)
	}

	functionName := args[2]
	if functionName == "hydrate" {
		err := functions.Hydrate(args[3], args[4:])
		if err != nil {
			log.Printf("Error occured %s", err)
			os.Exit(-1)
		}
	}
}

/*

kla apply  -f ./main.yaml image="nginx:1" replicas=2
kla apply fn hydrate -f ./samples/main.yaml image="nginx:1" replicas=2
*/
