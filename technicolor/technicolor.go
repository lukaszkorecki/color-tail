package technicolor

import (
	"math/rand"
)

var (
	resetMark  = "\033[39m"
	colorTable = map[string]string{
		"Red":           "\033[31m",
		"Green":         "\033[32m",
		"Yellow":        "\033[33m",
		"Blue":          "\033[34m",
		"Magenta":       "\033[35m",
		"Cyan":          "\033[36m",
		"Light gray":    "\033[37m",
		"Dark gray":     "\033[90m",
		"Light red":     "\033[91m",
		"Light green":   "\033[92m",
		"Light yellow":  "\033[93m",
		"Light blue":    "\033[94m",
		"Light magenta": "\033[95m",
		"Light cyan":    "\033[96m",
	}

	Colors = (func() []string {
		c := make([]string, len(colorTable)-1)
		for key, _ := range colorTable {
			c = append(c, key)
		}
		return c
	})()

	ctLen = int32(len(Colors))
)

// Colorizes given string with specified color
// if the color doesn't exist we fallback to a random one
func Paint(str, color string) string {

	c, ok := colorTable[color]
	s := ""
	if !ok {
		s = RandPaint(str)
	} else {
		s = c + str + resetMark
	}

	return s
}

// Colorizes the string with randomly picked color
func RandPaint(str string) string {
	idx := randIndex()
	key := Colors[idx]
	return Paint(str, key)
}

func RandColorName() string {
	return Colors[randIndex()]
}
func randIndex() int {

	return int(rand.Int31n(ctLen))
}
