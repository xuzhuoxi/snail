package impl

type Room struct {
	Name    string
	Private bool
	MaxUser int
}

type RoomGroup struct {
	RoomName []string
}
