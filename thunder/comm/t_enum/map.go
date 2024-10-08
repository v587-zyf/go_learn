package enum

// 2 3 4
// 1 0 5
// 8 7 6
type MAP_TOWARD int

var MAP_TOWARD_ALL = []MAP_TOWARD{MAP_NOW, MAP_LEFT, MAP_RIGHT, MAP_TOP, MAP_BOTTOM, MAP_LEFT_TOP, MAP_RIGHT_TOP, MAP_RIGHT_BOTTOM, MAP_LEFT_BOTTOM}

const (
	MAP_NOW MAP_TOWARD = iota
	MAP_LEFT
	MAP_LEFT_TOP
	MAP_TOP
	MAP_RIGHT_TOP
	MAP_RIGHT
	MAP_RIGHT_BOTTOM
	MAP_BOTTOM
	MAP_LEFT_BOTTOM
)
