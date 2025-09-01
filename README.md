# Core-kit

Go utilities for microservices and bots.

## Features

- Lightweight logger with optional metadata integration
- Environment and system metadata discovery (hostname, pod, node, app info)
- Ready-to-use utilities for microservices and bots
- Unit tests and benchmarks included

## Installation

```bash
go get github.com/gdnd-community/core-kit
```


## Usage 


###### Default & Fast Usage
```go

import "github.com/gdnd-community/core-kit/pkg/log"

func main() {
    log.Init("debug")
    log.Info("Service started", map[string]interface{}{
        "module": "main",
    })
}
```

###### With Meta ( This option is suitable for microservices. )
```go

import "github.com/gdnd-community/core-kit/pkg/log"
import "github.com/gdnd-community/core-kit/pkg/meta"

func main() {
    meta := meta.Discover("user-service", "1.0.0", "dev")
    log.Init("debug", logger.WithMetadata(meta))

    log.Warn("Run boi")
}
```



