package main

import (
    "bytes"
    "unicode"
)

// Stateless link parser. Finds the next link in s, starting at s[i:].
// Returns raw link and last index in s or an empty string and -1
func linkParse(s []byte, off int) ([]byte, int) {
    n := len(s)
    sep := []byte("href")
    seplen := len(sep)
Loop:

    for ; off+seplen <= n; off++ {
        // find href
        if s[off] == sep[0] && bytes.Equal(s[off:off+seplen], sep) {
            off += seplen

            // find beginning of link
            for ; off < n; off++ {
                v := s[off]
                if unicode.IsSpace(int32(v)) || v == '=' {
                    continue
                }

                if v == '"' {
                    off++
                    break
                }

                goto Loop
            }

            start := off

            // find the end of link
            for ; off < n; off++ {
                if s[off] == '"' {
                    break
                }
            }

            return bytes.TrimSpace(s[start:off]), off
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
