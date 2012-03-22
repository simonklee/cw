package main

import (
    "testing"
)

func error_(t *testing.T, expected, got interface{}, err error) {
    t.Errorf("expected `%v` got `%v`, err(%v)", expected, got, err)
}

type parseTest struct {
    in, out string
}

var parseTests = []parseTest{
    {"href=\"http://foo.com\"", "http://foo.com"},
    {" href = \" http://foo.com/\" ", "http://foo.com/"},
    {" href = \" http://foo.com/?foo=bar \" ", "http://foo.com/?foo=bar"},
}

var multistring = `<!doctype html>
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
</nav>`

func TestParserOnePass(t *testing.T) {
    for _, test := range parseTests {
        p := initParse(test.in)

        if s, err := p.next(); s != test.out {
            error_(t, test.out, s, err)
        }

        if s, err := p.next(); err == nil {
            error_(t, s, nil, err)
        }

        t.Log(test.in, test.out)
    }
}

func TestParseMultiPass(t *testing.T) {
    p := initParse(multistring)

    if o := p.all(); len(o) != 6 {
        error_(t, 6, len(o), nil)
    }
}

func BenchmarkParseOnePass(b *testing.B) {
    for i := 0; i < b.N; i++ {
        p := initParse(multistring)
        _, err := p.next()

        for err == nil {
            _, err = p.next()
        }
    }
}

func BenchmarkParseMultiPass(b *testing.B) {
    p := initParse(multistring)

    for i := 0; i < b.N; i++ {
        p = initParse(multistring)
        p.all()
    }
}

func BenchmarkParseMultiPassLong(b *testing.B) {
    s := ""

    for i := 0; i < 100; i++ {
        s += multistring
    }

    p := initParse(s)

    for i := 0; i < b.N; i++ {
        p.pos = 0
        p.all()
    }
}
