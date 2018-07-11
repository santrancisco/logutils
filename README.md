# San modified version

This modified version of logutils includes fatih/color libary. I used a simple math operation to determine the color we want to attach to each log line.

Assuming the highest severity log is always RED, we pre-determine the color for the last 4 highest severity from Blue, Green, Yellow and Red. Anything fall outside of this will have default terminal color.

Considering most of the time, my filter ares: `[]logutils.LogLevel{"DEBUG","INFO", "WARN", "ERROR"},`, this will color all 4 log severity levels. 

# logutils

logutils is a Go package that augments the standard library "log" package
to make logging a bit more modern, without fragmenting the Go ecosystem
with new logging packages.

## The simplest thing that could possibly work

Presumably your application already uses the default `log` package. To switch, you'll want your code to look like the following:

```go
package main

import (
	"log"
	"os"

	"github.com/hashicorp/logutils"
)

func main() {
	filter := &logutils.LevelFilter{
		Levels: []logutils.LogLevel{"DEBUG", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("WARN"),
		Writer: os.Stderr,
	}
	log.SetOutput(filter)

	log.Print("[DEBUG] Debugging") // this will not print
	log.Print("[WARN] Warning") // this will
	log.Print("[ERROR] Erring") // and so will this
	log.Print("Message I haven't updated") // and so will this
}
```

This logs to standard error exactly like go's standard logger. Any log messages you haven't converted to have a level will continue to print as before.
