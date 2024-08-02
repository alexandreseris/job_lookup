package main

import "os"

func IsTest() bool {
	v, _ := os.LookupEnv("JOB_LOOKUP_TEST")
	return v == "TRUE"
}
