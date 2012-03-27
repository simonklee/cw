package main

import (
    "io/ioutil"
    "net/url"
    "os"
    "testing"
)

var data = []byte(`<!doctype html>
<!--[if lt IE 9]><html class="ie"><![endif]-->
<!--[if gte IE 9]><!--><html><!--<![endif]-->
<head> 
<meta charset="utf-8"/>
<meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1"/>
<meta name="viewport" content="width=device-width, initial-scale=1"/>
<title>Redis Protocol ¿ simon.klee</title> 
<link rel="icon" href="/static/favicon.ico" type="image/vnd.microsoft.icon">
<link href='http://fonts.googleapis.com/css?family=Contrail+One' rel='stylesheet' type='text/css'>
<link href="/static/styles.css" rel="stylesheet"> 
<!--[if lt IE 9]>
<script src="//html5shim.googlecode.com/svn/trunk/html5.js"></script>
<![endif]-->
</head>
<body>
<div id="wrapper"><div id="wrapper_inner">
<nav> 
<div class="container"> 
<h1 class="brand"><a href="/">simon.klee</a></h1>
<ul class="navlist"> 
<li><a href="http://github.com/simonz05/godis">go/godis</a></li> 
<li><a href="http://github.com/simonz05/redoc">go/redoc</a></li> 
<li><a href="http://github.com/simonz05/odis">python/odis</a></li> 
</ul> 
</div>
</nav>`)

func TestUrlResolver(t *testing.T) {
    baseurl, _ := url.Parse("http://simonklee.org")
    urls := urlResolver(&Entry{Data: data, URL: baseurl})

    if len(urls) != 6 {
        error_(t, 6, len(urls), nil)
    }
}

func TestUrlResolverLong(t *testing.T) {
    f, err := os.OpenFile("test.html", os.O_RDONLY, 0644)

    if err != nil {
        error_(t, nil, nil, err)
    }

    data, err = ioutil.ReadAll(f)
    f.Close()

    baseurl, _ := url.Parse("http://github.com")
    urls := urlResolver(&Entry{Data: data, URL: baseurl})

    if len(urls) != 243 {
        error_(t, 243, len(urls), nil)
    }
}

func BenchmarkResolveUrl(b *testing.B) {
    f, _ := os.OpenFile("test.html", os.O_RDONLY, 0644)
    data, _ = ioutil.ReadAll(f)
    f.Close()

    baseurl, _ := url.Parse("http://github.com")

    for i := 0; i < b.N; i++ {
        urlResolver(&Entry{Data: data, URL: baseurl})
    }
}
