# sugar
sugar is a comprehensive, efficient, and reusable util function library of go


## 💡 Usage

You can import `sugar` using:

```go
import (
    "github.com/phuhao00/sugar"
)
```

Then use one of the helpers below:

```go
clamp := sugar.Clamp(2, 3, 5)
fmt.Printf("clamp:%v", clamp)
```