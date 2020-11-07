package auth

import (
	"fmt"
	"math/rand"
	"time"
)

// Total name perms: 992
var (
	// 31
	colors = [...]string{
		"white", "silver", "gray", "black", "blue", "green", "cyan", "yellow",
		"gold", "orange", "brown", "red", "violet", "pink", "magenta", "purple",
		"maroon", "crimson", "plum", "fuchsia", "lavender", "slate", "navy",
		"azure", "aqua", "olive", "teal", "lime", "beige", "tan", "sienna",
	}
	// 32
	animals = [...]string{"ant", "bear", "bird", "cat", "chicken", "cow", "deer",
		"dog", "donkey", "duck", "fish", "fox", "frog", "horse", "kangaroo", "koala",
		"lemur", "lion", "lizard", "monkey", "octopus", "pig", "shark", "sheep",
		"sloth", "spider", "squirrel", "tiger", "toad", "weasel", "whale", "wolf",
	}
)

func GenerateRandomName() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%s-%s", colors[rand.Intn(len(colors))],
		animals[rand.Intn(len(animals))])
}
