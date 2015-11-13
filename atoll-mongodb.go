
package main

import (
  "os"
  "log"
  "fmt"
  "github.com/codegangsta/cli"
)

func fatalError(err error) {
  log.Fatalf("Error: %v", err)
}

func main() {
  app := cli.NewApp()
  app.Name = "atoll-mongodb"
  app.Usage = "MongoDB monitoring plugin for Atoll"
  app.Version = "0.1.1"
  app.Flags = []cli.Flag {
    cli.StringFlag{
      Name: "host",
      Value: "localhost",
      Usage: "MongoDB admin host",
    },
    cli.IntFlag{
      Name: "port",
      Value: 27017,
      Usage: "MongoDB port",
    },
  }
  app.Action = func(c *cli.Context) {
    mongodb := MongoDB {
      c.String("host"),
      uint16(c.Int("port")),
    };
    data, err := mongodb.Monitor();
    if err != nil {
      fatalError(err);
    } else {
      fmt.Printf("%s", data);
    }
  }

  app.Run(os.Args)
}
