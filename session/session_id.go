package main

import (
    "math/rand"
    "time"
    "strconv"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
    letterIdxBits = 6                    // 6 bits to represent a letter index
    letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
)

func randStringBytesMask(n int) string {
    rand.Seed(time.Now().UnixNano())
    b := make([]byte, n)
    for i := 0; i < n; {
        if idx := int(rand.Int63() & letterIdxMask); idx < len(letterBytes) {
            b[i] = letterBytes[idx]
            i++
        }
    }
    return string(b)
}

func SessionID() string {
    return randStringBytesMask(10) + strconv.FormatInt(time.Now().Unix(), 10)
}