package main

import (
	"github.com/EngoEngine/engo"
)

type SaveData struct {
	HasNaniteKey    bool
	HasHoodKey      bool
	HasPPE          bool
	NaniteBoxChecks int
	MarsChecks      int
}

var CurrentSave = &SaveData{}

func main() {
	opts := engo.RunOptions{
		Title:         "Skeleboy Studios",
		Width:         640,
		Height:        360,
		ScaleOnResize: true,
	}
	engo.Run(opts, &SkeleScene{})
}
