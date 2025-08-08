package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"k8s.io/klog/v2"
)

func main() {
	log.SetFlags(0)

	flagSet := flag.NewFlagSet("test", flag.ExitOnError)
	klog.InitFlags(flagSet)
	_ = flagSet.Parse([]string{"--v", "0"})
	cmd := newRootCmd()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
