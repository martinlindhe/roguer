package rogue

// Player ...
type Player struct {
	Name  string
	Token string
	Spawn *Obj // points to a npc in game
}
