# evbus
This package is a simple asynchronous event bus.

## How to get
```
go get github.com/junyaU/evbus
```
## Usage
```
import github.com/junyaU/evbus

bus := evbus.New()

fn := func(v int) {
    fmt.Printf("%d\n", a + b)
}

bus.Subscribe("topic", fn)

bus.Publish("topic", 10)
bus.Publish("topic", 15)

bus.Wait()
```

## License
The MIT License (MIT) - see LICENSE for more details

