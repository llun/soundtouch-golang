# Golang Soundtouch API

Base on [Nodejs version](https://github.com/CONNCTED/SoundTouch-NodeJS)

## Example

```golang
package main

import (
  "log"

  "github.com/llun/soundtouch-golang"
)

func main() {
  speakerCh := make(chan *soundtouch.Speaker, 1)
  soundtouch.Lookup(speakerCh)
  speaker := <-speakerCh
  data, err := speaker.Volume()
  if err != nil {
    log.Fatal(err)
  }
  log.Printf("%v\n", data)
  log.Printf("%s\n", data.Raw)
}
```

## License

MIT
