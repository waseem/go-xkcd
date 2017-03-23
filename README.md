# go-xkcd
xkcd downloads and prints json of an xkcd comic given its number.

Usage
-----

```go
package main

import (
    "fmt"
    "log"
    "os"
    "github.com/waseem/go-xkcd"
    )

func main() {
  comic, err := xkcd.GetComic(os.Args[1])
    if err != nil {
      log.Fatal(err)
    }

  fmt.Printf("Comic: %d\n", comic.Num)
    fmt.Printf("+----------+\n\n")
    fmt.Printf("Title: %s\n", comic.Title)
    fmt.Printf("Safe Title: %s\n", comic.SafeTitle)
    fmt.Printf("Transcript: %s\n", comic.Transcript)
    fmt.Printf("Alt: %s\n", comic.Alt)
    fmt.Printf("Image: %s\n", comic.Img)
    fmt.Printf("Link: %s\n", comic.Link)
    fmt.Printf("Date: %s-%s-%s\n", comic.Year, comic.Month, comic.Day)
    fmt.Printf("News: %s\n", comic.News)
}

```

`xkcd.GetComic(0)` will return the latest published comic.
