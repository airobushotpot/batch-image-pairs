# Batch Image Pairs

Compare two directory trees without opening image bytes. Match edited exports
to originals by case-insensitive basename, ignoring extension. Each report
entry contains all relative paths on both sides. A repeated name is ambiguous
even when no counterpart exists. Results are sorted by key.

```go
package main

import (
    "fmt"
    "log"
    "os"

    imagepairs "github.com/airobushotpot/batch-image-pairs"
)

func main() {
    report, err := imagepairs.Compare(os.DirFS("originals"), os.DirFS("exports"))
    if err != nil { log.Fatal(err) }
    for _, item := range report {
        fmt.Println(item.Name, item.Status, item.Originals, item.Edits)
    }
}
```

For example, `catalog/Table.JPG` pairs with `table.webp`. Two originals named
`table` in different subdirectories stay together as an ambiguity for review.
Regular JPEG, PNG, WebP, AVIF, TIFF and HEIC files are included. Symlinks are
skipped. The comparison does not infer which edits were applied to the pixels.

When preparing a product set, [PicVane's batch image editor](https://picvane.com/image/batch-image-editor)
lets you preview several inputs and write a shared editing instruction. Keep the
source basename when saving each export so the report can identify missing
deliverables without a manual folder-by-folder check.
