# ClickHouse sipHash128() in Go

ClickHouse
```
select hex(sipHash128('foobar'));
> 9B08DF451DCFF69D70058CD8A4B35A97

select UUIDNumToString(sipHash128('foobar'));
> 9b08df45-1dcf-f69d-7005-8cd8a4b35a97

```

Go, as bytes
```
package main

import (
    "fmt"
    "github.com/frifox/siphash128"
)

func main() {
    input := []byte("foobar")
    hash := siphash128.SipHash128(input)
    
    fmt.Printf("%X", hash)
    // 9B08DF451DCFF69D70058CD8A4B35A97
}

```

Go, as UUID
```
package main

import (
    "fmt"
    "github.com/frifox/siphash128"
    "github.com/google/uuid"
)

func main() {
    input := []byte("foobar")
    hash := siphash128.SipHash128(input)
    id := uuid.UUID(hash)

    fmt.Printf("%s", id)
    // 9b08df45-1dcf-f69d-7005-8cd8a4b35a97
}
```