package main

import (
    "fmt"
    "os"
    "flag"
    "net/url"
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

    for i := range args {
        u, err := url.Parse(args[i])

        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            continue
        }

        println(u.RequestURI())
        println(u.Scheme)
        println(u.Host)
        println(u.Path)
        println(u.RawQuery)
        println(u.Fragment)
    }

    if len(args) == 0 {
        usage()
    }
}
