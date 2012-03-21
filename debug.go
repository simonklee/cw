package main

import (
    "net/http"
    "net/url"
)

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
