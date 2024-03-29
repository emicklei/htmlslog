## html slog

```
package main

import (
	"fmt"
	"log/slog"

	"github.com/emicklei/htmlslog"
)

func main() {
	handler := htmlslog.New(htmlslog.Options{Level: slog.LevelInfo, Title: "Test"})
	mylog := slog.New(handler)
	mylog.Info("Hello, world!")
	fmt.Println(handler.Close())
}
```

See https://goplay.tools/snippet/CkxYaQODkMt