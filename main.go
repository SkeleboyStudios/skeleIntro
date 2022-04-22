package main

import (
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
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
	RecruitedLen          bool
	RecruitedMe           bool
	PlayerLocation        engo.Point
	DrinkCount            int
	CookieCount           int
	BandageCount          int
	HasMedKit             bool
}

var CurrentSave = &SaveData{
	PlayerLocation: engo.Point{X: 300, Y: 125},
}

func main() {
	common.AddShader(fightShader)
	skeleScene := &SkeleScene{}
	engo.RegisterScene(skeleScene)
	engo.RegisterScene(&GhostFightScene{})
	opts := engo.RunOptions{
		Title:         "Skeleboy Studios",
		Width:         640,
		Height:        360,
		ScaleOnResize: true,
	}
	engo.Run(opts, &GhostFightScene{})
}
