# Golang Soundtouch API

Base on [Nodejs version](https://github.com/CONNCTED/SoundTouch-NodeJS)

## Sample

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

  websocketCh, err := speaker.Listen()
  if err != nil {
    log.Fatal(err)
  }

  data, err := speaker.Volume()
  if err != nil {
    log.Fatal(err)
  }
  log.Printf("%v\n", data)
  log.Printf("%s\n", data.Raw)

  speaker.SetVolume(40)
  log.Printf("Set volume to 40")

  for message := range websocketCh {
    log.Printf(message)
  }
}

```

## License

MIT
