package main

import (
    "unicode"
    "io"
    "strings"
)

type parse struct {
    s string
    sep string
    pos int
    end int
    n int
}

func initParse(s string) *parse {
    return &parse{
        s: s, 
        pos: 0, 
        end: len(s),
    }
}

func (p *parse) next() (string, error) {
    i := p.pos

Loop:

    for ; i + 4 <= p.end; i++ {
        // find href
        if p.s[i] == 'h' && p.s[i:i+4] == "href" {
            i += 4

            // find beginning of link
            for ; i < p.end; i++ {
                x := p.s[i]
                if unicode.IsSpace(int32(x)) || x == '=' {
                    continue
                }

                if x == '"' {
                    i++
                    break
                }

                goto Loop
            }

            start := i

            // find the end of link
            for ; i < p.end; i++ {
                if p.s[i] == '"' {
                    break
                }
            }

            p.pos = i
            return strings.TrimSpace(p.s[start:p.pos]), nil
        }
    }

    return "", io.EOF
}

func (p *parse) all() ([]string) {
    var buf []string

    for s, err := p.next(); err == nil; {
        buf = append(buf, s)
        s, err = p.next()
    }

    return buf
}
