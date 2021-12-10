package main

import (
	"github.com/EngoEngine/engo"
)

type SaveData struct {
	HasNaniteKey          bool
	NaniteKeyInSafe       bool
	HasHoodKey            bool
	HoodKeyInSafe         bool
	HasPPE                bool
	HasDeskKey            bool
	DeskKeyInSafe         bool
	IsDrawerBroken        bool
	NaniteBoxChecks       int
	MarsChecks            int
	HasSpookyBoard        bool
	HasSpookyBoardPointer bool
	HasSpaceKey           bool
	SpaceKeyInSafe        bool
	KeyCount              int
	IsSafeOpen            bool
}

var CurrentSave = &SaveData{}

func main() {
	opts := engo.RunOptions{
		Title:         "Skeleboy Studios",
		Width:         640,
		Height:        360,
		ScaleOnResize: true,
	}
	CurrentSave.HasSpaceKey = true
	CurrentSave.HasHoodKey = true
	CurrentSave.HasNaniteKey = true
	CurrentSave.HasDeskKey = true
	engo.Run(opts, &SkeleScene{PlayerLocation: engo.Point{X: 300, Y: 125}})
}
