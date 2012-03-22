package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "time"
)

var (
    logger = log.New(os.Stdout, "", 0)
    store  *Store
)

func usage() {
    fmt.Fprintf(os.Stdout, "usage: cw url [url ...]\n")
    flag.PrintDefaults()
    os.Exit(0)
}

func main() {
    flag.Usage = usage
    flag.Parse()
    args := flag.Args()

    if len(args) == 0 {
        usage()
    } else {
        c := newContext()

        for _, u := range args {
            c.Add(u)
        }

        time.Sleep(1e9)
    }
}
