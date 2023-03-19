# kpm
KCL Package Manager

```
package main

import (
    "log"
    "os"

    "github.com/otiai10/copy"
)

func main() {
    // 复制目录及其内容到新位置
    err := copy.Copy("oldDir", "newDir")
    if err != nil {
        log.Fatal(err)
    }

    // 删除旧目录
    err = os.RemoveAll("oldDir")
    if err != nil {
        log.Fatal(err)
    }
}
```