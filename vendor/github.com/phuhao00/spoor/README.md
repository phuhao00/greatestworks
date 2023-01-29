# spoor
make logger switchable and adapted project


## 💡 Usage

You can import `spoor` using:

```go
import (
    "github.com/phuhao00/spoor"
)
```

Then use one of the helpers below:

## fileWriter
```go
fileWriter := spoor.NewFileWriter("log", 0, 0, 0)
l := spoor.NewSpoor(spoor.DEBUG, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile, spoor.WithFileWriter(fileWriter))
l.DebugF("hhhh")
select {}

```

## consoleWriter

```` go
l := NewSpoor(DEBUG, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile, WithNormalWriter(os.Stdout))
l.DebugF("hhhh")
````
## elasticWriter

````go


````
## clickHouseWriter

````go


````
## logbusWriter

````go


````