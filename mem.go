package main

import (
    "runtime"
    "fmt"
)

func memstat() {
    s := new(runtime.MemStats)
    runtime.ReadMemStats(s)

    fmt.Println("\nTotal:")
    fmt.Printf("    Alloc:        %d(%.3fMB)\n", s.Alloc, float32(s.Alloc)/1024/1024)
    fmt.Printf("    TotalAlloc:   %d(%.3fMB)\n", s.TotalAlloc, float32(s.TotalAlloc)/1024/1024)
    fmt.Printf("    Sys:          %d(%.3fMB)\n", s.Sys, float32(s.Sys)/1024/1024)
    fmt.Printf("    Lookups:      %d\n", s.Lookups)
    fmt.Printf("    Mallocs:      %d\n", s.Mallocs)
    fmt.Printf("    Frees:        %d\n", s.Frees)

    fmt.Println("\nHeap:")
    fmt.Printf("    HeapAlloc:    %d(%.3fMB)\n", s.HeapAlloc, float32(s.HeapAlloc)/1024/1024)
    fmt.Printf("    HeapSys:      %d(%.3fMB)\n", s.HeapSys, float32(s.HeapSys)/1024/1024)
    fmt.Printf("    HeapIdle:     %d(%.3fMB)\n", s.HeapIdle, float32(s.HeapIdle)/1024/1024)
    fmt.Printf("    HeapInuse:    %d(%.3fMB)\n", s.HeapInuse, float32(s.HeapInuse)/1024/1024)
    fmt.Printf("    HeapReleased: %d(%.3fMB)\n", s.HeapReleased, float32(s.HeapReleased)/1024/1024)
    fmt.Printf("    HeapObjects:  %d\n", s.HeapObjects)

    //fmt.Println("\nPer-Size alloc:")

    //for i, _ := range s.BySize {
    //    fmt.Printf("%d %d %d\n", s.BySize[i].Size, s.BySize[i].Mallocs, s.BySize[i].Frees)
    //}
}
