package main

import (
	"fmt"
	"time"
)

func main() {
	// Test how Go parses MST
	t, err := time.Parse(time.RFC1123, "Mon, 02 Jan 2006 15:04:05 MST")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Parsed: %v\n", t)
	fmt.Printf("UTC: %v\n", t.UTC())
	fmt.Printf("Zone: %v, Offset: %v\n", t.Zone())

	// What the test expects
	shanghai := time.FixedZone("Asia/Shanghai", 8*3600)
	expected := time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("MST", -7*3600)).In(shanghai)
	fmt.Printf("Expected: %v\n", expected)

	// What we get
	actual := t.In(shanghai)
	fmt.Printf("Actual: %v\n", actual)
}
