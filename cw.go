package main

import (
    "flag"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "os"
    "strings"
)

type fetchState struct {
    id    string
    state bool
}

type urlResult struct {
    id   string
    data []byte
    url  *url.URL
}

var (
    logger = log.New(os.Stdout, "", 0)
)

func (u *urlResult) storageWriter(status chan<- fetchState) {
    //logger.Printf("%s\n", u.data)
    status <- fetchState{"", true}
}

func (u *urlResult) getUrl(status, result chan<- fetchState) {
    urlstr := u.url.String()
    sep := strings.Index(urlstr, "?")
    fs := fetchState{urlstr, true}

    if sep > 0 {
        urlstr = urlstr[:sep]
    }

    res, err := http.Get(urlstr)

    if err != nil {
        fs.state = false
        status <- fs
        return
    }

    debugResponse(res)

    u.data, err = ioutil.ReadAll(res.Body)

    if err != nil {
        fs.state = false
        status <- fs
        return
    }

    res.Body.Close() /* close fd */
    go u.storageWriter(result)
    status <- fs
}

func getAndListenUrls(urls []*url.URL) {
    netlen := len(urls)
    fslen := netlen
    netState := make(chan fetchState)
    fsState := make(chan fetchState)

    for _, u := range urls {
        ur := &urlResult{url: u}
        go ur.getUrl(netState, fsState)
    }

    for {
        select {
        case fs := <-netState:
            debugFetchState(&fs)
            netlen--
        case fs := <-fsState:
            debugFetchState(&fs)
            fslen--
        }

        if netlen == 0 && fslen == 0 {
            break
        }
    }
}

func debugFetchState(fs *fetchState) {
    logger.Printf("id: %s", fs.id)

    if fs.state {
        logger.Printf("state: OK\n")
    } else {
        logger.Printf("state: ERR\n")
    }
}

func debugUrl(u *url.URL) {
    logger.Println("Host:", u.Host)
    logger.Println("Path:", u.Path)
    logger.Println("Request URI:", u.RequestURI())
    logger.Println("Scheme:", u.Scheme)
    logger.Println("Query:", u.RawQuery)
    logger.Println("Fragment:", u.Fragment)
}

func debugResponse(r *http.Response) {
    logger.Println("Status:", r.Status)
    logger.Println("StatusCode:", r.StatusCode)
    logger.Println("Proto:", r.Proto)
    logger.Println("Header:")

    for k, v := range r.Header {
        logger.Println("\t", k, ":", v)
    }
}

func usage() {
    fmt.Fprintf(os.Stdout, "usage: cw url [url ...]\n")
    flag.PrintDefaults()
    os.Exit(0)
}

func main() {
    flag.Usage = usage
    flag.Parse()
    args := flag.Args()

    urls := make([]*url.URL, 0, len(args))

    for i := range args {
        u, err := url.Parse(args[i])

        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            continue
        }

        urls = append(urls, u)
    }

    getAndListenUrls(urls)

    if len(args) == 0 {
        usage()
    }
}
