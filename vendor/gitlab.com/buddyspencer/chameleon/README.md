# Chameleon

## How to use it?

### Download it
```
go get gitlab.com/buddyspencer/chameleon
```

### Use it
```
import (
	"fmt"
	. "gitlab.com/buddyspencer/chameleon"
)

func main() {
	fmt.Println(Red("test"))
}
```

### I need the Value
```
    test := Red("test")
    fmt.Println(test.Value())
```

## Add custom colors
```
AddCustomColor("pink", "38;5;199m")
fmt.Println(CustomColor("pink", "test"))
```
If you want to checkout all the available custom colors, you can run the chameleoncolors binary.
<img src="pictures/foregroundcolors.png" alt="foreground">
<img src="pictures/backgroundcolors.png" alt="background">
