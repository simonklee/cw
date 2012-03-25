package main

import (
    "bytes"
    "unicode"
)

// Stateless link parser. Finds the next link in s, starting at s[i:].
// Returns raw link and last index in s or an empty string and -1
func linkParse(s []byte, i int) ([]byte, int) {
    n := len(s)
    sep := []byte("href")
Loop:

    for ; i+4 <= n; i++ {
        // find href
        if s[i] == 'h' && bytes.Equal(s[i:i+4], sep) {
            i += 4

            // find beginning of link
            for ; i < n; i++ {
                v := s[i]
                if unicode.IsSpace(int32(v)) || v == '=' {
                    continue
                }

                if v == '"' {
                    i++
                    break
                }

                goto Loop
            }

            start := i

            // find the end of link
            for ; i < n; i++ {
                if s[i] == '"' {
                    break
                }
            }

            return bytes.TrimSpace(s[start:i]), i
        }
    }

    return nil, -1
}

func linkParseAll(data []byte) [][]byte {
    var buf [][]byte
    link, i := linkParse(data, 0)

    for i != -1 {
        buf = append(buf, link)
        link, i = linkParse(data, i)
    }

    return buf
}
