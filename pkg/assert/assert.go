package assert

import "log"

func OK(err error) {
	if err != nil {
		log.Fatalf("Error: %s\n", err.Error())
	}
}
