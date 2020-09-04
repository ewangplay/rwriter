# RotateWriter

Rotating File Writer in Golang. Implements io.Writer.

## Usage

Get package:
```
go get github.com/ewangplay/rwriter
```

A sample:
```
package main

import (
	"io"
	"log"

	"github.com/ewangplay/rwriter"
)

func main() {
	var err error
	var w io.Writer

	cfg := &rwriter.Config{
		Module: "test",
		Path:   "/path/to/log/files",
	}
	w, err = rwriter.NewRotateWriter(cfg)
	if err != nil {
		log.Printf("Create rotate writer failed: %v\n", err)
		return
	}

	log.SetOutput(w)
	log.Println("Hello, rotate writer!")
}
```

The RototeWriter instance can be passed to any parameter that matches io.Writer interface.