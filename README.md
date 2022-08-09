# evbus
This package is a simple asynchronous event bus.

## How to get
```
go get github.com/junyaU/evbus
```
## Usage
```go
package main
import "github.com/junyaU/evbus"

func main() {
    bus := evbus.New()
    
    fn := func(v string) {
        fmt.Printf("%d\n", v)
    }
    
    bus.Subscribe("topic", fn)
    
    bus.Publish("topic", "hello")
    bus.Publish("topic", "world")
    
    bus.Wait()
}
```

## License
The MIT License (MIT) - see LICENSE for more details

