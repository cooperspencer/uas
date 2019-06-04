package chameleon

import (
	"fmt"
	"runtime"
)

var (
	cfcolor = make(map[string]string)
)

type Chameleon struct {
	value interface{}
	color string
	end   string
}

func (c Chameleon) String() string {
	if c.color != "" {
		return fmt.Sprintf("%s%v%s", c.color, c.value, c.end)
	}
	return fmt.Sprint("%v", c.value)
}

func (c Chameleon) Value() interface{} {
	return c.value
}

func ret(color string, s interface{}) Chameleon {
	if runtime.GOOS == "windows" {
		return Chameleon{value: s}
	}
	switch s.(type) {
	case Chameleon:
		return s.(Chameleon).addColor(color)
	default:
		return Chameleon{value: s, color: color, end: "\x1b[0m"}
	}
}

func Bold(s interface{}) Chameleon {
	return ret("\x1b[1m", s)
}

func Underlined(s interface{}) Chameleon {
	return ret("\x1b[4m", s)
}

func Dim(s interface{}) Chameleon {
	return ret("\x1b[2m", s)
}

func Inverted(s interface{}) Chameleon {
	return ret("\x1b[7m", s)
}

func Hidden(s interface{}) Chameleon {
	return ret("\x1b[8m", s)
}

// Foregroundcolors

func Black(s interface{}) Chameleon {
	return ret("\x1b[30m", s)
}

func Red(s interface{}) Chameleon {
	return ret("\x1b[31m", s)
}

func Green(s interface{}) Chameleon {
	return ret("\x1b[32m", s)
}

func Yellow(s interface{}) Chameleon {
	return ret("\x1b[33m", s)
}

func Blue(s interface{}) Chameleon {
	return ret("\x1b[34m", s)
}

func Magenta(s interface{}) Chameleon {
	return ret("\x1b[35m", s)
}

func Cyan(s interface{}) Chameleon {
	return ret("\x1b[36m", s)
}

func Lightgray(s interface{}) Chameleon {
	return ret("\x1b[37m", s)
}

func Darkgray(s interface{}) Chameleon {
	return ret("\x1b[90m", s)
}

func Lightred(s interface{}) Chameleon {
	return ret("\x1b[91m", s)
}

func Lightgreen(s interface{}) Chameleon {
	return ret("\x1b[92m", s)
}

func Lightyellow(s interface{}) Chameleon {
	return ret("\x1b[93m", s)
}

func Lightblue(s interface{}) Chameleon {
	return ret("\x1b[94m", s)
}

func Lightmagenta(s interface{}) Chameleon {
	return ret("\x1b[95m", s)
}

func Lightcyan(s interface{}) Chameleon {
	return ret("\x1b[96m", s)
}

func White(s interface{}) Chameleon {
	return ret("\x1b[97m", s)
}

// Backgroundcolors

func BBlack(s interface{}) Chameleon {
	return ret("\x1b[40m", s)
}

func BRed(s interface{}) Chameleon {
	return ret("\x1b[41m", s)
}

func BGreen(s interface{}) Chameleon {
	return ret("\x1b[42m", s)
}

func BYellow(s interface{}) Chameleon {
	return ret("\x1b[43m", s)
}

func BBlue(s interface{}) Chameleon {
	return ret("\x1b[44m", s)
}

func BMagenta(s interface{}) Chameleon {
	return ret("\x1b[45m", s)
}

func BCyan(s interface{}) Chameleon {
	return ret("\x1b[46m", s)
}

func BLightgray(s interface{}) Chameleon {
	return ret("\x1b[47m", s)
}

func BDarkgray(s interface{}) Chameleon {
	return ret("\x1b[100m", s)
}

func BLightred(s interface{}) Chameleon {
	return ret("\x1b[101m", s)
}

func BLightgreen(s interface{}) Chameleon {
	return ret("\x1b[102m", s)
}

func BLightyellow(s interface{}) Chameleon {
	return ret("\x1b[103m", s)
}

func BLightblue(s interface{}) Chameleon {
	return ret("\x1b[104m", s)
}

func BLightmagenta(s interface{}) Chameleon {
	return ret("\x1b[105m", s)
}

func BLightcyan(s interface{}) Chameleon {
	return ret("\x1b[106m", s)
}

func BWhite(s interface{}) Chameleon {
	return ret("\x1b[107m", s)
}

// Custom Foregroundcolor
func AddCustomColor(colorname, color string) {
	cfcolor[colorname] = fmt.Sprintf("\033[" + color)
}

func CustomColor(colorname string, s interface{}) Chameleon {
	return ret(cfcolor[colorname], s)
}

func (c Chameleon) addColor(color string) Chameleon {
	c.color += color
	return c
}

func (c Chameleon) Bold() Chameleon {
	return c.addColor("\x1b[1m")
}

func (c Chameleon) Underlined() Chameleon {
	return c.addColor("\x1b[4m")
}

func (c Chameleon) Dim() Chameleon {
	return c.addColor("\x1b[2m")
}

func (c Chameleon) Inverted() Chameleon {
	return c.addColor("\x1b[7m")
}

func (c Chameleon) Hidden() Chameleon {
	return c.addColor("\x1b[8m")
}

// Foregroundcolors

func (c Chameleon) Black() Chameleon {
	return c.addColor("\x1b[30m")
}

func (c Chameleon) Red() Chameleon {
	return c.addColor("\x1b[31m")
}

func (c Chameleon) Green() Chameleon {
	return c.addColor("\x1b[32m")
}

func (c Chameleon) Yellow() Chameleon {
	return c.addColor("\x1b[33m")
}

func (c Chameleon) Blue() Chameleon {
	return c.addColor("\x1b[34m")
}

func (c Chameleon) Magenta() Chameleon {
	return c.addColor("\x1b[35m")
}

func (c Chameleon) Cyan() Chameleon {
	return c.addColor("\x1b[36m")
}

func (c Chameleon) Lightgray() Chameleon {
	return c.addColor("\x1b[37m")
}

func (c Chameleon) Darkgray() Chameleon {
	return c.addColor("\x1b[90m")
}

func (c Chameleon) Lightred() Chameleon {
	return c.addColor("\x1b[91m")
}

func (c Chameleon) Lightgreen() Chameleon {
	return c.addColor("\x1b[92m")
}

func (c Chameleon) Lightyellow() Chameleon {
	return c.addColor("\x1b[93m")
}

func (c Chameleon) Lightblue() Chameleon {
	return c.addColor("\x1b[94m")
}

func (c Chameleon) Lightmagenta() Chameleon {
	return c.addColor("\x1b[95m")
}

func (c Chameleon) Lightcyan() Chameleon {
	return c.addColor("\x1b[96m")
}

func (c Chameleon) White() Chameleon {
	return c.addColor("\x1b[97m")
}

// Backgroundcolors

func (c Chameleon) BBlack() Chameleon {
	return c.addColor("\x1b[40m")
}

func (c Chameleon) BRed() Chameleon {
	return c.addColor("\x1b[41m")
}

func (c Chameleon) BGreen() Chameleon {
	return c.addColor("\x1b[42m")
}

func (c Chameleon) BYellow() Chameleon {
	return c.addColor("\x1b[43m")
}

func (c Chameleon) BBlue() Chameleon {
	return c.addColor("\x1b[44m")
}

func (c Chameleon) BMagenta() Chameleon {
	return c.addColor("\x1b[45m")
}

func (c Chameleon) BCyan() Chameleon {
	return c.addColor("\x1b[46m")
}

func (c Chameleon) BLightgray() Chameleon {
	return c.addColor("\x1b[47m")
}

func (c Chameleon) BDarkgray() Chameleon {
	return c.addColor("\x1b[100m")
}

func (c Chameleon) BLightred() Chameleon {
	return c.addColor("\x1b[101m")
}

func (c Chameleon) BLightgreen() Chameleon {
	return c.addColor("\x1b[102m")
}

func (c Chameleon) BLightyellow() Chameleon {
	return c.addColor("\x1b[103m")
}

func (c Chameleon) BLightblue() Chameleon {
	return c.addColor("\x1b[104m")
}

func (c Chameleon) BLightmagenta() Chameleon {
	return c.addColor("\x1b[105m")
}

func (c Chameleon) BLightcyan() Chameleon {
	return c.addColor("\x1b[106m")
}

func (c Chameleon) BWhite() Chameleon {
	return c.addColor("\x1b[107m")
}

// Custom Foregroundcolor
func (c Chameleon) AddCustomColor(colorname, color string) {
	cfcolor[colorname] = fmt.Sprintf("\033[" + color)
}

func (c Chameleon) CustomColor(colorname string, ) Chameleon {
	return c.addColor(cfcolor[colorname])
}