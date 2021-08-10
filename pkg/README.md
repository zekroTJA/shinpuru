# `pkg`

Here you can find all publicly available packages, where *most** packages are also used by shinpuru.

I try to document these packages as well as possible and keep the API as consistent as possible so that you can rely on these packages.

In some packages, you might find a `v2`, `v3` (, ...) directory. When I need to change the API of a package so that it would break something, I create a new version which contains the updated version so that the original version maintains consistency.

You can simply `go get` these packages.
```
go get github.com/zekroTJA/shinpuru/pkg/<package_name>
```

Then, you can import it in your own application.
```go
package main

import (
    "fmt"
	  "github.com/zekroTJA/shinpuru/pkg/bytecount"
)

func main() {
	  fmt.Println(bytecount.Format(123456789))
    // -> "117.738 MiB"
}
```

---
**Some packages (or pakcage versions) are actually not anymore used in shinpuru because of new features and updates. For example, package [lctimer](lctimer) was replaced with [robfig/cron/v3](https://github.com/robfig/cron), so the package is not used anymore but will be provided further to maintain package consistency.*