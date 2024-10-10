package wret

//go:generate msgp

type Error struct {
	Message string `msg:"message"`
}
