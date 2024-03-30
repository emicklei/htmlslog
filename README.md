## html slog

```
package main

import (
	"fmt"
	"log/slog"

	"github.com/emicklei/htmlslog"
)

func main() {
	handler := htmlslog.New(os.Stdout, htmlslog.Options{Level: slog.LevelInfo, Title: "Test"})
	mylog := slog.New(handler)
	mylog.Info("Hello, world!")
	handler.Close()
}
```