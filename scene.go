package main

import (
	"bytes"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"

	"github.com/SkeleboyStudios/skeleIntro/assets"
)

const (
	CollisionGroupPlayaWall = 1 << iota
	CollisionGroupDoor
	CollisionGroupInterest
)

type SkeleScene struct {
	files []string
}

func (*SkeleScene) Type() string { return "Skele Scene" }

func (s *SkeleScene) Preload() {
	s.files = []string{
		"lobby/bg.png",
		"me/npc.png",
		"me/playa.png",
		"lobby/rsdoor.png",
		"lobby/mbdoor.png",
		"lobby/pdoor.png",
		"lobby/nanites.png",
		"lobby/mars.png",
		"lobby/sand.png",
		"title/bg.mp3",
		"title/cursor.png",
		"title/move.wav",
		"title/log.ttf",
		"title/log.png",
		"title/dots.png",
		"title/log.wav",
		"lab/bg.png",
		"lab/rsdoor.png",
		"lab/hood.png",
		"lab/len.png",
		"lab/lenSS.png",
		"president/bg.png",
		"president/doorSS.png",
		"president/diplomas.png",
		"president/discord.png",
		"president/desk.png",
		"president/crash.ogg",
		"president/diplomasSS.png",
		"president/donations.png",
		"president/engo.png",
		"president/safe.png",
		"president/safeSS.png",
		"space/bg.png",
		"space/doorSS.png",
		"space/moon.png",
		"space/tv.png",
		"space/tvSS.png",
		"space/sauce.png",
		"space/window.png",
		"space/windowSS.png",
	}

	for _, file := range s.files {
		data, err := assets.Asset(file)
		if err != nil {
			log.Fatalf("Unable to locate asset with URL: %v\n", file)
		}
		err = engo.Files.LoadReaderData(file, bytes.NewReader(data))
		if err != nil {
			log.Fatalf("Unable to load asset with URL: %v\n At %v", file, s.Type())
		}
	}

	engo.Input.RegisterButton("up", engo.KeyW, engo.KeyArrowUp)
	engo.Input.RegisterButton("down", engo.KeyS, engo.KeyArrowDown)
	engo.Input.RegisterButton("left", engo.KeyA, engo.KeyArrowLeft)
	engo.Input.RegisterButton("right", engo.KeyD, engo.KeyArrowRight)
	engo.Input.RegisterButton("A", engo.KeyJ, engo.KeyZ)
	engo.Input.RegisterButton("B", engo.KeyK, engo.KeyX)
	engo.Input.RegisterButton("X", engo.KeyL, engo.KeyC)
	engo.Input.RegisterButton("Y", engo.KeySemicolon, engo.KeyV)
	engo.Input.RegisterButton("FullScreen", engo.KeyFour, engo.KeyF4)
	engo.Input.RegisterButton("Exit", engo.KeyEscape)
}

func (s *SkeleScene) Setup(u engo.Updater) {
	w := u.(*ecs.World)

	rand.Seed(time.Now().UnixNano())

	var renderable *common.Renderable
	var notrenderable *common.NotRenderable
	w.AddSystemInterface(&common.RenderSystem{}, renderable, notrenderable)

	var animatable *common.Animationable
	var notanimatable *common.NotAnimationable
	var animSys = &common.AnimationSystem{}
	w.AddSystemInterface(animSys, animatable, notanimatable)

	var collisionable *common.Collisionable
	var notcollisionable *common.NotCollisionable
	w.AddSystemInterface(&common.CollisionSystem{Solids: CollisionGroupPlayaWall}, collisionable, notcollisionable)

	var audioable *common.Audioable
	var notaudioable *common.NotAudioable
	var audioSys = &common.AudioSystem{}
	w.AddSystemInterface(audioSys, audioable, notaudioable)

	// w.AddSystem(&systems.FullScreenSystem{})
	// w.AddSystem(&systems.ExitSystem{})

	var moveable *Moveable
	w.AddSystemInterface(&MoveSystem{}, moveable, nil)

	var doorable *Doorable
	w.AddSystemInterface(&DoorSystem{}, doorable, nil)

	var interestable *InterestAble
	w.AddSystemInterface(&InterestSystem{}, interestable, nil)

	var cursorable *CursorAble
	var notcursorable *NotCursorAble
	var curSys CursorSystem
	curSys.ClickSoundURL = "title/move.wav"
	curSys.CursorURL = "title/cursor.png"
	w.AddSystemInterface(&curSys, cursorable, notcursorable)

	w.AddSystem(&CombatLogSystem{
		BackgroundURL: "title/log.png",
		DotURL:        "title/dots.png",
		FontURL:       "title/log.ttf",
		LineDelay:     0.3,
		LetterDelay:   0.1,
	})

	var phaseable *common.BasicFace
	w.AddSystemInterface(&PhaseSystem{}, phaseable, nil)

	selFont := &common.Font{
		Size: 48,
		FG:   color.RGBA{R: 0xb7, G: 0xf7, B: 0xff, A: 0xff},
		URL:  "title/log.ttf",
	}
	selFont.CreatePreloaded()

	w.AddSystem(&AcceptSystem{Fnt: selFont, BackgroundURL: "title/log.png"})

	bgm := audio{BasicEntity: ecs.NewBasic()}
	bgmPlayer, _ := common.LoadedPlayer("title/bg.mp3")
	bgm.AudioComponent = common.AudioComponent{Player: bgmPlayer}
	bgmPlayer.Repeat = true
	bgmPlayer.Play()
	w.AddEntity(&bgm)

	logSnd := audio{BasicEntity: ecs.NewBasic()}
	logPlayer, _ := common.LoadedPlayer("title/log.wav")
	logSnd.AudioComponent = common.AudioComponent{Player: logPlayer}
	logSnd.AudioComponent.Player.SetVolume(0.15)
	w.AddEntity(&logSnd)

	crashSnd := audio{BasicEntity: ecs.NewBasic()}
	crashPlayer, _ := common.LoadedPlayer("president/crash.ogg")
	crashSnd.AudioComponent = common.AudioComponent{Player: crashPlayer}
	w.AddEntity(&crashSnd)

	playaSS := common.NewSpritesheetWithBorderFromFile("me/playa.png", 23, 45, 1, 1)
	playa := playa{BasicEntity: ecs.NewBasic()}
	playa.Drawable = playaSS.Drawable(0)
	playa.SetZIndex(2)
	playa.Position = CurrentSave.PlayerLocation
	playa.Scale = engo.Point{X: 2, Y: 2}
	playa.Height = 2 * playa.Drawable.Height()
	playa.Width = 2 * playa.Drawable.Width()
	playa.AnimationComponent = common.NewAnimationComponent(playaSS.Drawables(), 0.2)
	playa.AddAnimations([]*common.Animation{
		&common.Animation{
			Name:   "upstop",
			Frames: []int{9},
		},
		&common.Animation{
			Name:   "downstop",
			Frames: []int{0},
		},
		&common.Animation{
			Name:   "leftstop",
			Frames: []int{3},
		},
		&common.Animation{
			Name:   "rightstop",
			Frames: []int{6},
		},
		&common.Animation{
			Name:   "up",
			Frames: []int{9, 10},
			Loop:   true,
		},
		&common.Animation{
			Name:   "down",
			Frames: []int{1, 2},
			Loop:   true,
		},
		&common.Animation{
			Name:   "left",
			Frames: []int{4, 5},
			Loop:   true,
		},
		&common.Animation{
			Name:   "right",
			Frames: []int{7, 8},
			Loop:   true,
		},
	})
	playa.SelectAnimationByName("downstop")
	playa.CollisionComponent = common.CollisionComponent{Main: CollisionGroupPlayaWall, Group: CollisionGroupDoor | CollisionGroupInterest}
	playa.AddShape(common.Shape{
		Lines: []engo.Line{
			engo.Line{P1: engo.Point{X: 6, Y: 82}, P2: engo.Point{X: 42, Y: 82}},
			engo.Line{P1: engo.Point{X: 42, Y: 82}, P2: engo.Point{X: 42, Y: 88}},
			engo.Line{P1: engo.Point{X: 42, Y: 88}, P2: engo.Point{X: 6, Y: 88}},
			engo.Line{P1: engo.Point{X: 6, Y: 88}, P2: engo.Point{X: 6, Y: 82}},
		},
	})
	playa.Speed = 145.0
	playa.PlayerCharacter = true
	w.AddEntity(&playa)
	w.AddSystem(&common.EntityScroller{SpaceComponent: &playa.SpaceComponent, TrackingBounds: engo.AABB{Min: engo.Point{X: -1000, Y: -1000}, Max: engo.Point{X: 1000, Y: 15000}}})

	newRoom(w, engo.Point{X: 0, Y: 0}, "lobby/bg.png", []wallInfo{
		wallInfo{
			Position: engo.Point{X: 112, Y: 0},
			Width:    378,
			Height:   144,
		},
		wallInfo{
			Position: engo.Point{X: 0, Y: 250},
			Width:    600,
			Height:   150,
		},
		wallInfo{
			Position: engo.Point{X: 0, Y: 144},
			Shapes: []common.Shape{
				common.Shape{Lines: []engo.Line{
					engo.Line{
						P1: engo.Point{X: 0, Y: 0},
						P2: engo.Point{X: 112, Y: 0},
					},
					engo.Line{
						P1: engo.Point{X: 112, Y: 0},
						P2: engo.Point{X: 0, Y: 112},
					},
					engo.Line{
						P1: engo.Point{X: 0, Y: 112},
						P2: engo.Point{X: 0, Y: 0},
					},
				}},
			},
		},
		wallInfo{
			Position: engo.Point{X: 490, Y: 144},
			Shapes: []common.Shape{
				common.Shape{Lines: []engo.Line{
					engo.Line{
						P1: engo.Point{X: 0, Y: 0},
						P2: engo.Point{X: 112, Y: 0},
					},
					engo.Line{
						P1: engo.Point{X: 112, Y: 0},
						P2: engo.Point{X: 112, Y: 112},
					},
					engo.Line{
						P1: engo.Point{X: 112, Y: 112},
						P2: engo.Point{X: 0, Y: 0},
					},
				}},
			},
		},
		wallInfo{
			Position: engo.Point{X: 65, Y: 140},
			Shapes: []common.Shape{
				common.Shape{Lines: []engo.Line{
					engo.Line{
						P1: engo.Point{X: 0, Y: 54},
						P2: engo.Point{X: 0, Y: 38},
					},
					engo.Line{
						P1: engo.Point{X: 0, Y: 38},
						P2: engo.Point{X: 16, Y: 22},
					},
					engo.Line{
						P1: engo.Point{X: 16, Y: 22},
						P2: engo.Point{X: 72, Y: 22},
					},
					engo.Line{
						P1: engo.Point{X: 72, Y: 22},
						P2: engo.Point{X: 72, Y: 38},
					},
					engo.Line{
						P1: engo.Point{X: 72, Y: 38},
						P2: engo.Point{X: 56, Y: 54},
					},
					engo.Line{
						P1: engo.Point{X: 56, Y: 54},
						P2: engo.Point{X: 0, Y: 54},
					},
				}},
			},
		},
		wallInfo{
			Position: engo.Point{X: 450, Y: 100},
			Shapes: []common.Shape{
				common.Shape{
					Ellipse: common.Ellipse{Rx: 32, Ry: 32, Cx: 32, Cy: 32},
				},
			},
		},
		wallInfo{
			Position: engo.Point{X: 355, Y: 175},
			Shapes: []common.Shape{
				common.Shape{
					Lines: []engo.Line{
						engo.Line{P1: engo.Point{X: 46, Y: 24}, P2: engo.Point{X: 170, Y: 24}},
						engo.Line{P1: engo.Point{X: 170, Y: 24}, P2: engo.Point{X: 170, Y: 36}},
						engo.Line{P1: engo.Point{X: 170, Y: 36}, P2: engo.Point{X: 134, Y: 72}},
						engo.Line{P1: engo.Point{X: 134, Y: 72}, P2: engo.Point{X: 10, Y: 72}},
						engo.Line{P1: engo.Point{X: 10, Y: 72}, P2: engo.Point{X: 10, Y: 60}},
						engo.Line{P1: engo.Point{X: 10, Y: 60}, P2: engo.Point{X: 46, Y: 24}},
					},
				},
			},
		},
	}, []doorInfo{
		doorInfo{
			URL:          "lobby/rsdoor.png",
			CellWidth:    32,
			CellHeight:   70,
			BorderWidth:  1,
			BorderHeight: 1,
			Position:     engo.Point{X: 8, Y: 114},
			Shapes: []common.Shape{
				common.Shape{Lines: []engo.Line{
					engo.Line{P1: engo.Point{X: 0, Y: 134}, P2: engo.Point{X: 64, Y: 70}},
					engo.Line{P1: engo.Point{X: 64, Y: 70}, P2: engo.Point{X: 64, Y: 84}},
					engo.Line{P1: engo.Point{X: 64, Y: 84}, P2: engo.Point{X: 0, Y: 142}},
					engo.Line{P1: engo.Point{X: 0, Y: 142}, P2: engo.Point{X: 0, Y: 134}},
				}},
			},
			OpenFrames:  []int{0, 1, 2, 3},
			CloseFrames: []int{3, 2, 1, 0},
			Button:      "left",
			TeleportTo:  engo.Point{X: 476, Y: 626},
		},
		doorInfo{
			URL:          "lobby/mbdoor.png",
			CellWidth:    40,
			CellHeight:   72,
			BorderWidth:  1,
			BorderHeight: 1,
			Position:     engo.Point{X: 384, Y: 28},
			Shapes: []common.Shape{
				common.Shape{Lines: []engo.Line{
					engo.Line{P1: engo.Point{X: 16, Y: 116}, P2: engo.Point{X: 64, Y: 116}},
					engo.Line{P1: engo.Point{X: 64, Y: 116}, P2: engo.Point{X: 64, Y: 128}},
					engo.Line{P1: engo.Point{X: 64, Y: 128}, P2: engo.Point{X: 16, Y: 128}},
					engo.Line{P1: engo.Point{X: 16, Y: 128}, P2: engo.Point{X: 16, Y: 116}},
				}},
			},
			OpenFrames:  []int{0, 1, 2, 3, 4, 5},
			CloseFrames: []int{5, 4, 3, 2, 1, 0},
			Button:      "up",
			TeleportTo:  engo.Point{X: 394, Y: 1652},
		},
		doorInfo{
			URL:          "lobby/pdoor.png",
			CellWidth:    32,
			CellHeight:   54,
			BorderWidth:  1,
			BorderHeight: 1,
			Position:     engo.Point{X: 134, Y: 58},
			Shapes: []common.Shape{
				common.Shape{Lines: []engo.Line{
					engo.Line{P1: engo.Point{X: 20, Y: 86}, P2: engo.Point{X: 60, Y: 86}},
					engo.Line{P1: engo.Point{X: 60, Y: 86}, P2: engo.Point{X: 60, Y: 96}},
					engo.Line{P1: engo.Point{X: 60, Y: 96}, P2: engo.Point{X: 60, Y: 96}},
					engo.Line{P1: engo.Point{X: 60, Y: 96}, P2: engo.Point{X: 20, Y: 86}},
				}},
			},
			OpenFrames:  []int{0, 1, 2, 3, 4},
			CloseFrames: []int{4, 3, 2, 1, 0},
			Button:      "up",
			TeleportTo:  engo.Point{X: 128, Y: 1170},
		},
	}, []interestInfo{
		interestInfo{
			URL:      "lobby/mars.png",
			Position: engo.Point{X: 450, Y: 100},
			Shapes: []common.Shape{
				common.Shape{
					Lines: []engo.Line{
						engo.Line{P1: engo.Point{X: 0, Y: 32}, P2: engo.Point{X: 64, Y: 32}},
						engo.Line{P1: engo.Point{X: 64, Y: 32}, P2: engo.Point{X: 64, Y: 64}},
						engo.Line{P1: engo.Point{X: 64, Y: 64}, P2: engo.Point{X: 0, Y: 64}},
						engo.Line{P1: engo.Point{X: 0, Y: 64}, P2: engo.Point{X: 0, Y: 32}},
					},
				},
			},
			Func: func() {
				CurrentSave.MarsChecks++
				msgs := []string{}
				if CurrentSave.MarsChecks < 2 {
					msgs = append(msgs,
						"It's Mars!",
						"The grand prize for the winner of Marsbound",
						"Should I mount it on a trophy?",
						"Would that look too tacky for space?",
					)
				} else if CurrentSave.MarsChecks < 3 {
					msgs = append(msgs,
						"It's still Mars!",
						"I didn't want it hung up so it wouldn't",
						"accidentally fall and break.",
					)
				} else if CurrentSave.MarsChecks < 4 {
					msgs = append(msgs,
						"One little poke couldn't hurt",
						"...",
						"A piece fell off.",
						"Oops.",
					)
				} else {
					msgs = append(msgs,
						"Not gonna touch it again.",
						"Planets are actually very expensive.",
						"Can't have pieces falling off all willy-nilly.",
					)
				}
				for _, msg := range msgs {
					engo.Mailbox.Dispatch(CombatLogMessage{
						Msg:  msg,
						Fnt:  selFont,
						Clip: logPlayer,
					})
				}
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: ListenPhase,
				})
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: LogClearPhase,
				})
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: WalkPhase,
				})
				engo.Mailbox.Dispatch(PhaseDequeuMessage{})
			},
		},
		interestInfo{
			URL:      "lobby/nanites.png",
			Position: engo.Point{X: 65, Y: 140},
			Shapes: []common.Shape{
				common.Shape{
					Lines: []engo.Line{
						engo.Line{P1: engo.Point{X: 0, Y: 58}, P2: engo.Point{X: 0, Y: 14}},
						engo.Line{P1: engo.Point{X: 0, Y: 14}, P2: engo.Point{X: 76, Y: 14}},
						engo.Line{P1: engo.Point{X: 76, Y: 14}, P2: engo.Point{X: 52, Y: 58}},
						engo.Line{P1: engo.Point{X: 52, Y: 58}, P2: engo.Point{X: 0, Y: 58}},
					},
				},
			},
			Func: func() {
				CurrentSave.NaniteBoxChecks++
				msgs := []string{}
				if CurrentSave.NaniteBoxChecks < 2 {
					msgs = append(msgs,
						"It's a box of nanites and mods!",
						"These little guys buff up and help out",
						"Rogue Scientists!",
					)
				} else if CurrentSave.NaniteBoxChecks < 10 {
					msgs = append(msgs, "It's")
					r := rand.Intn(11)
					switch r {
					case 0:
						msgs = append(msgs,
							"An Absorbant Module featuring Crumplezones!",
							"Wow!",
							"It adds several layers of defense!",
						)
					case 1:
						msgs = append(msgs,
							"Len!",
							"He's the starter nanite!",
							"He gives you laser-based abilities!",
						)
					case 2:
						msgs = append(msgs,
							"Kelvin!",
							"He's the cool nanite.",
							"Gives you ice-based abilities!",
						)
					case 3:
						msgs = append(msgs,
							"Prometheus!",
							"Such a hot-head!",
							"Gives fire-based abilities!",
						)
					case 4:
						msgs = append(msgs,
							"Gauss",
							"He's got a magnetic personality!",
							"Movement and speed based abilities",
						)
					case 5:
						msgs = append(msgs,
							"Faraday",
							"A shocking guy",
							"Lightning-based abilities!",
						)
					case 6:
						msgs = append(msgs,
							"a slime!",
							"After accidentally feeding it after midnight",
							"This guy grew until it nearly destroyed the city!",
						)
					case 7:
						msgs = append(msgs,
							"Parts to an auto-turret.",
							"These cuddly guys were made by Dr. Shockley",
							"To comfort his friends!",
							"(and shoot his enemies)",
						)
					case 8:
						msgs = append(msgs,
							"a Repo-tron 40k!",
							"This sophiscated 3D printer can print anything",
							"a rogue scientist might need!",
						)
					case 9, 10:
						if !CurrentSave.HasPPE {
							msgs = append(msgs,
								"It's a pair of nitrile gloves and goggles!",
								"Added PPE to your inventory!",
							)
							CurrentSave.HasPPE = true
						} else {
							msgs = append(msgs,
								"There's some ISO-certified PPE here!",
								"But you've already got some!",
							)
						}
					default:
						msgs = append(msgs,
							"Nothing...",
							"Uh oh",
							"There's gotta be more in here!",
						)
					}
				} else if CurrentSave.NaniteBoxChecks < 20 {
					//no need to keep digging
					msgs = append(msgs,
						"I've already looked through this box enough",
						"There couldn't possibly be anything left!",
					)
				} else if CurrentSave.NaniteBoxChecks < 21 {
					msgs = append(msgs,
						"Okay. Fine. I'll look through again.",
						"See? Nothing left.",
						"Except...wait a minute...",
						"It's a toad out on patrol!",
						"You exchange glances.",
						"It blushes before running into its toad-hole.",
					)
				} else {
					msgs = append(msgs,
						"The hole just sits there.",
						"Your friend is not coming back.",
						"Unless you take drastic measures",
						"Like refresh the page!",
					)
				}
				if !CurrentSave.HasNaniteKey {
					msgs = append(msgs, "Hey, it looks like there's a key", "at the bottom of the box!", "Would you like to take it?")
				}
				for _, msg := range msgs {
					engo.Mailbox.Dispatch(CombatLogMessage{
						Msg:  msg,
						Fnt:  selFont,
						Clip: logPlayer,
					})
				}
				if !CurrentSave.HasNaniteKey {
					engo.Mailbox.Dispatch(AcceptSetMessage{
						AcceptFunc: func() {
							CurrentSave.HasNaniteKey = true
						},
					})
					engo.Mailbox.Dispatch(PhaseSetMessage{
						Phase: AcceptPhase,
					})
				} else {
					engo.Mailbox.Dispatch(PhaseSetMessage{
						Phase: ListenPhase,
					})
				}
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: LogClearPhase,
				})
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: WalkPhase,
				})
				engo.Mailbox.Dispatch(PhaseDequeuMessage{})
			},
		},
		interestInfo{
			URL:      "lobby/sand.png",
			Position: engo.Point{X: 355, Y: 175},
			Shapes: []common.Shape{
				common.Shape{
					Lines: []engo.Line{
						engo.Line{P1: engo.Point{X: 0, Y: 0}, P2: engo.Point{X: 184, Y: 0}},
						engo.Line{P1: engo.Point{X: 184, Y: 0}, P2: engo.Point{X: 184, Y: 72}},
						engo.Line{P1: engo.Point{X: 184, Y: 72}, P2: engo.Point{X: 0, Y: 72}},
						engo.Line{P1: engo.Point{X: 0, Y: 72}, P2: engo.Point{X: 0, Y: 0}},
					},
				},
			},
			Func: func() {
				msgs := []string{
					"Let's Save Summer!",
					"I'm currently still plotting...",
					"er... planning out this one!",
					"would you like to visit the site?",
				}
				for _, msg := range msgs {
					engo.Mailbox.Dispatch(CombatLogMessage{
						Msg:  msg,
						Fnt:  selFont,
						Clip: logPlayer,
					})
				}
				engo.Mailbox.Dispatch(AcceptSetMessage{
					AcceptFunc: func() {
						navigateToPageImpl("https://www.letssavesummer.com")
						audioSys.Pause()
						engo.Mailbox.Dispatch(PhaseDequeuMessage{})
						engo.Mailbox.Dispatch(CombatLogMessage{
							Msg:  "Oh? You're still here?",
							Fnt:  selFont,
							Clip: logPlayer,
						})
						engo.Mailbox.Dispatch(AcceptSetMessage{
							AcceptFunc: func() {
								audioSys.Restart()
							},
						})
						engo.Mailbox.Dispatch(PhaseSetMessage{
							Phase: AcceptPhase,
						})
						engo.Mailbox.Dispatch(PhaseSetMessage{
							Phase: LogClearPhase,
						})
						engo.Mailbox.Dispatch(PhaseSetMessage{
							Phase: WalkPhase,
						})
						engo.Mailbox.Dispatch(PhaseDequeuMessage{})
					},
				})
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: AcceptPhase,
				})
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: LogClearPhase,
				})
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: WalkPhase,
				})
				engo.Mailbox.Dispatch(PhaseDequeuMessage{})
			},
		},
	})

	lab := newRoom(w, engo.Point{X: 0, Y: 500}, "lab/bg.png", []wallInfo{
		wallInfo{
			Position: engo.Point{X: 352, Y: 88},
			Width:    138,
			Height:   56,
		},
		wallInfo{
			Position: engo.Point{X: 154, Y: 120},
			Width:    180,
			Height:   48,
		},
		wallInfo{
			Position: engo.Point{X: 46, Y: 166},
			Width:    62,
			Height:   48,
		},
		wallInfo{
			Position: engo.Point{X: 0, Y: 0},
			Width:    0,
			Height:   0,
			Shapes: []common.Shape{common.Shape{Lines: []engo.Line{
				engo.Line{P1: engo.Point{X: 108, Y: 166}, P2: engo.Point{X: 152, Y: 122}},
				engo.Line{P1: engo.Point{X: 152, Y: 122}, P2: engo.Point{X: 152, Y: 166}},
				engo.Line{P1: engo.Point{X: 152, Y: 166}, P2: engo.Point{X: 108, Y: 210}},
				engo.Line{P1: engo.Point{X: 108, Y: 210}, P2: engo.Point{X: 108, Y: 166}},
			}}},
		},
		wallInfo{
			Position: engo.Point{X: 0, Y: 214},
			Width:    0,
			Height:   0,
			Shapes: []common.Shape{common.Shape{Lines: []engo.Line{
				engo.Line{P1: engo.Point{X: 0, Y: 0}, P2: engo.Point{X: 44, Y: 0}},
				engo.Line{P1: engo.Point{X: 44, Y: 0}, P2: engo.Point{X: 0, Y: 44}},
				engo.Line{P1: engo.Point{X: 0, Y: 44}, P2: engo.Point{X: 0, Y: 0}},
			}}},
		},
		wallInfo{
			Position: engo.Point{X: 0, Y: 254},
			Width:    600,
			Height:   20,
		},
		wallInfo{
			Position: engo.Point{X: 332, Y: 144},
			Width:    0,
			Height:   0,
			Shapes: []common.Shape{common.Shape{Lines: []engo.Line{
				engo.Line{P1: engo.Point{X: 0, Y: 0}, P2: engo.Point{X: 20, Y: 0}},
				engo.Line{P1: engo.Point{X: 20, Y: 0}, P2: engo.Point{X: 0, Y: 20}},
				engo.Line{P1: engo.Point{X: 0, Y: 20}, P2: engo.Point{X: 0, Y: 0}},
			}}},
		},
		wallInfo{
			Position: engo.Point{X: 490, Y: 144},
			Width:    0,
			Height:   0,
			Shapes: []common.Shape{common.Shape{Lines: []engo.Line{
				engo.Line{P1: engo.Point{X: 0, Y: 0}, P2: engo.Point{X: 114, Y: 0}},
				engo.Line{P1: engo.Point{X: 114, Y: 0}, P2: engo.Point{X: 114, Y: 114}},
				engo.Line{P1: engo.Point{X: 114, Y: 114}, P2: engo.Point{X: 0, Y: 0}},
			}}},
		},
	}, []doorInfo{
		doorInfo{
			URL:          "lab/rsdoor.png",
			Position:     engo.Point{X: 525, Y: 110},
			TeleportTo:   engo.Point{X: 65, Y: 140},
			CellWidth:    32,
			CellHeight:   70,
			BorderWidth:  1,
			BorderHeight: 1,
			Shapes: []common.Shape{common.Shape{Lines: []engo.Line{
				engo.Line{P1: engo.Point{X: 0, Y: 70}, P2: engo.Point{X: 0, Y: 82}},
				engo.Line{P1: engo.Point{X: 0, Y: 82}, P2: engo.Point{X: 64, Y: 146}},
				engo.Line{P1: engo.Point{X: 64, Y: 146}, P2: engo.Point{X: 64, Y: 134}},
				engo.Line{P1: engo.Point{X: 64, Y: 134}, P2: engo.Point{X: 0, Y: 70}},
			}}},
			OpenFrames:  []int{0, 1, 2, 3},
			CloseFrames: []int{3, 2, 1, 0},
			Button:      "right",
		},
	}, []interestInfo{
		interestInfo{
			URL:      "lab/hood.png",
			Position: engo.Point{X: 46, Y: 0},
			Shapes: []common.Shape{common.Shape{Lines: []engo.Line{
				engo.Line{P1: engo.Point{X: 110, Y: 166}, P2: engo.Point{X: 142, Y: 166}},
				engo.Line{P1: engo.Point{X: 142, Y: 166}, P2: engo.Point{X: 96, Y: 212}},
				engo.Line{P1: engo.Point{X: 96, Y: 212}, P2: engo.Point{X: 64, Y: 212}},
				engo.Line{P1: engo.Point{X: 64, Y: 212}, P2: engo.Point{X: 110, Y: 166}},
			}}},
			Func: func() {
				msgs := []string{"The hood is packed with dangerous chemicals!"}
				if CurrentSave.HasPPE {
					msgs = append(msgs, "But you have PPE!", "Would you like to put it on and look inside?")
					for _, msg := range msgs {
						engo.Mailbox.Dispatch(CombatLogMessage{
							Msg:  msg,
							Fnt:  selFont,
							Clip: logPlayer,
						})
					}
					engo.Mailbox.Dispatch(AcceptSetMessage{
						AcceptFunc: func() {
							msgs2 := []string{
								"Inside the hood is a key shaped mold.",
								"You dust the mold off. Now it's just a key!",
								"You obtained THE LAB KEY",
							}
							for _, msg := range msgs2 {
								engo.Mailbox.Dispatch(CombatLogMessage{
									Msg:  msg,
									Fnt:  selFont,
									Clip: logPlayer,
								})
							}
							engo.Mailbox.Dispatch(PhaseDequeuMessage{})
							engo.Mailbox.Dispatch(PhaseSetMessage{
								Phase: ListenPhase,
							})
							CurrentSave.HasHoodKey = true
							engo.Mailbox.Dispatch(PhaseSetMessage{
								Phase: LogClearPhase,
							})
							engo.Mailbox.Dispatch(PhaseSetMessage{
								Phase: WalkPhase,
							})
							engo.Mailbox.Dispatch(PhaseDequeuMessage{})
						},
					})
					engo.Mailbox.Dispatch(PhaseSetMessage{
						Phase: AcceptPhase,
					})
					engo.Mailbox.Dispatch(PhaseSetMessage{
						Phase: LogClearPhase,
					})
					engo.Mailbox.Dispatch(PhaseSetMessage{
						Phase: WalkPhase,
					})
					engo.Mailbox.Dispatch(PhaseDequeuMessage{})
				} else {
					msgs = append(msgs, "It would be dangerous to open it without PPE.")
					for _, msg := range msgs {
						engo.Mailbox.Dispatch(CombatLogMessage{
							Msg:  msg,
							Fnt:  selFont,
							Clip: logPlayer,
						})
					}
					engo.Mailbox.Dispatch(PhaseSetMessage{
						Phase: ListenPhase,
					})
					engo.Mailbox.Dispatch(PhaseSetMessage{
						Phase: LogClearPhase,
					})
					engo.Mailbox.Dispatch(PhaseSetMessage{
						Phase: WalkPhase,
					})
					engo.Mailbox.Dispatch(PhaseDequeuMessage{})
				}
			},
		},
		interestInfo{
			URL:      "lab/len.png",
			Position: engo.Point{X: 400, Y: 25},
			Shapes: []common.Shape{common.Shape{Lines: []engo.Line{
				engo.Line{P1: engo.Point{X: 0, Y: 64}, P2: engo.Point{X: 64, Y: 64}},
				engo.Line{P1: engo.Point{X: 64, Y: 64}, P2: engo.Point{X: 64, Y: 128}},
				engo.Line{P1: engo.Point{X: 64, Y: 128}, P2: engo.Point{X: 0, Y: 128}},
				engo.Line{P1: engo.Point{X: 0, Y: 128}, P2: engo.Point{X: 0, Y: 64}},
			}}},
			Func: func() {
				msgs := []string{
					"Hello!",
					"I am Len!",
					"A nanite that grants Rogue Scientists",
					"science-based powers!",
					"Can't wait to help you in-game!",
				}
				for _, msg := range msgs {
					engo.Mailbox.Dispatch(CombatLogMessage{
						Msg:  msg,
						Fnt:  selFont,
						Clip: logPlayer,
					})
				}
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: ListenPhase,
				})
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: LogClearPhase,
				})
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: WalkPhase,
				})
				engo.Mailbox.Dispatch(PhaseDequeuMessage{})
			},
		},
	})

	//len Animation
	lenSS := common.NewSpritesheetWithBorderFromFile("lab/lenSS.png", 32, 64, 1, 1)
	lenAnim := animation{}
	lenAnim.AnimationComponent = common.NewAnimationComponent(lenSS.Drawables(), 0.3)
	lenAnim.AddDefaultAnimation(&common.Animation{
		Name:   "float",
		Frames: []int{0, 1, 2, 1},
		Loop:   true,
	})
	lenAnim.AnimationComponent.SelectAnimationByName("float")
	lab.interests[1].Drawable = lenSS.Drawable(0)
	lab.interests[1].Scale = engo.Point{X: 2, Y: 2}
	lab.interests[1].SetZIndex(5)
	animSys.Add(lab.interests[1].GetBasicEntity(), lenAnim.GetAnimationComponent(), lab.interests[1].GetRenderComponent())

	//diploma animation
	dipSS := common.NewSpritesheetWithBorderFromFile("president/diplomasSS.png", 64, 64, 1, 1)
	dipAnim := animation{}
	dipAnim.AnimationComponent = common.NewAnimationComponent(dipSS.Drawables(), 0.3)
	dipAnim.AddDefaultAnimation(&common.Animation{
		Name:   "undisturbed",
		Frames: []int{0},
	})
	dipAnim.AddAnimation(&common.Animation{
		Name:   "sparkle",
		Frames: []int{1, 2, 3, 2},
		Loop:   true,
	})
	dipAnim.AddAnimation(&common.Animation{
		Name:   "empty",
		Frames: []int{1},
		Loop:   true,
	})

	// safe animation
	safeSS := common.NewSpritesheetWithBorderFromFile("president/safeSS.png", 20, 20, 1, 1)
	safeAnim := animation{}
	safeAnim.AnimationComponent = common.NewAnimationComponent(safeSS.Drawables(), 0.3)
	safeAnim.AddAnimation(&common.Animation{
		Name:   "0",
		Frames: []int{0},
		Loop:   true,
	})
	safeAnim.AddAnimation(&common.Animation{
		Name:   "1",
		Frames: []int{1},
		Loop:   true,
	})
	safeAnim.AddAnimation(&common.Animation{
		Name:   "2",
		Frames: []int{2},
		Loop:   true,
	})
	safeAnim.AddAnimation(&common.Animation{
		Name:   "3",
		Frames: []int{3},
		Loop:   true,
	})
	safeAnim.AddAnimation(&common.Animation{
		Name:   "4",
		Frames: []int{4},
	})
	safeAnim.AddAnimation(&common.Animation{
		Name:   "open",
		Frames: []int{4, 5, 6, 7, 8},
	})
	pres := newRoom(w, engo.Point{X: 0, Y: 1000}, "president/bg.png", []wallInfo{
		wallInfo{
			Width:    380,
			Height:   50,
			Position: engo.Point{X: 112, Y: 96},
		},
		wallInfo{
			Width:    600,
			Height:   100,
			Position: engo.Point{X: 0, Y: 256},
		},
		wallInfo{
			Shapes: []common.Shape{common.Shape{Lines: []engo.Line{
				engo.Line{P1: engo.Point{X: 0, Y: 144}, P2: engo.Point{X: 112, Y: 144}},
				engo.Line{P1: engo.Point{X: 112, Y: 144}, P2: engo.Point{X: 0, Y: 256}},
				engo.Line{P1: engo.Point{X: 0, Y: 256}, P2: engo.Point{X: 0, Y: 144}},
			}}},
		},
		wallInfo{
			Shapes: []common.Shape{common.Shape{Lines: []engo.Line{
				engo.Line{P1: engo.Point{X: 490, Y: 144}, P2: engo.Point{X: 600, Y: 144}},
				engo.Line{P1: engo.Point{X: 600, Y: 144}, P2: engo.Point{X: 600, Y: 256}},
				engo.Line{P1: engo.Point{X: 600, Y: 256}, P2: engo.Point{X: 490, Y: 144}},
			}}},
		},
		wallInfo{
			Shapes: []common.Shape{common.Shape{Lines: []engo.Line{
				engo.Line{P1: engo.Point{X: 492, Y: 190}, P2: engo.Point{X: 574, Y: 190}},
				engo.Line{P1: engo.Point{X: 574, Y: 190}, P2: engo.Point{X: 574, Y: 230}},
				engo.Line{P1: engo.Point{X: 574, Y: 230}, P2: engo.Point{X: 516, Y: 230}},
				engo.Line{P1: engo.Point{X: 516, Y: 230}, P2: engo.Point{X: 492, Y: 190}},
			}}},
		},
		wallInfo{
			Position: engo.Point{X: 152, Y: 155},
			Width:    116,
			Height:   30,
		},
		wallInfo{
			Shapes: []common.Shape{common.Shape{Lines: []engo.Line{
				engo.Line{P1: engo.Point{X: 362, Y: 142}, P2: engo.Point{X: 448, Y: 144}},
				engo.Line{P1: engo.Point{X: 448, Y: 144}, P2: engo.Point{X: 440, Y: 160}},
				engo.Line{P1: engo.Point{X: 440, Y: 160}, P2: engo.Point{X: 362, Y: 160}},
				engo.Line{P1: engo.Point{X: 362, Y: 160}, P2: engo.Point{X: 362, Y: 142}},
			}}},
		},
		wallInfo{
			Shapes: []common.Shape{common.Shape{Lines: []engo.Line{
				engo.Line{P1: engo.Point{X: 102, Y: 144}, P2: engo.Point{X: 142, Y: 144}},
				engo.Line{P1: engo.Point{X: 142, Y: 144}, P2: engo.Point{X: 132, Y: 154}},
				engo.Line{P1: engo.Point{X: 132, Y: 154}, P2: engo.Point{X: 102, Y: 154}},
				engo.Line{P1: engo.Point{X: 102, Y: 154}, P2: engo.Point{X: 102, Y: 144}},
			}}},
		},
		wallInfo{
			Position: engo.Point{X: 376, Y: 186},
			Width:    54,
			Height:   40,
		},
		wallInfo{
			Position: engo.Point{X: 512, Y: 186},
			Width:    62,
			Height:   44,
		},
	}, []doorInfo{
		doorInfo{
			URL:          "president/doorSS.png",
			Position:     engo.Point{X: 120, Y: 236},
			TeleportTo:   engo.Point{X: 130, Y: 85},
			CellWidth:    40,
			CellHeight:   10,
			BorderWidth:  1,
			BorderHeight: 1,
			Shapes: []common.Shape{common.Shape{Lines: []engo.Line{
				engo.Line{P1: engo.Point{X: 10, Y: 0}, P2: engo.Point{X: 10, Y: 20}},
				engo.Line{P1: engo.Point{X: 10, Y: 20}, P2: engo.Point{X: 70, Y: 20}},
				engo.Line{P1: engo.Point{X: 70, Y: 20}, P2: engo.Point{X: 70, Y: 0}},
				engo.Line{P1: engo.Point{X: 70, Y: 0}, P2: engo.Point{X: 10, Y: 0}},
			}}},
			OpenFrames:  []int{0, 1, 2},
			CloseFrames: []int{2, 1, 0},
			Button:      "down",
		},
	}, []interestInfo{
		interestInfo{
			URL:      "president/diplomas.png",
			Position: engo.Point{X: 228, Y: 75},
			Shapes: []common.Shape{common.Shape{Lines: []engo.Line{
				engo.Line{P1: engo.Point{X: 0, Y: 68}, P2: engo.Point{X: 62, Y: 68}},
				engo.Line{P1: engo.Point{X: 62, Y: 68}, P2: engo.Point{X: 62, Y: 88}},
				engo.Line{P1: engo.Point{X: 62, Y: 88}, P2: engo.Point{X: 0, Y: 88}},
				engo.Line{P1: engo.Point{X: 0, Y: 88}, P2: engo.Point{X: 0, Y: 68}},
			}}},
			Func: func() {
				msgs := []string{}
				if CurrentSave.IsDrawerBroken {
					if CurrentSave.HasSpookyBoardPointer {
						msgs = []string{
							"There's nothing else inside.",
							"The empty hole in the wall serves as a",
							"reminder of your brute strength.",
						}
					} else {
						msgs = []string{
							"The diplomas were disturbed when you",
							"flung the drawer handle through the wall.",
							"Geeze. This is a disaster.",
							"Wait a second...",
							"In the hole there.",
							"Look inside?",
						}
					}
				} else {
					msgs = []string{
						"There's a bunch of diplomas on the wall",
						"just gathering dust.",
						"A PhD in WHAT?",
						"No WAY is that a thing.",
					}
				}
				for _, msg := range msgs {
					engo.Mailbox.Dispatch(CombatLogMessage{
						Msg:  msg,
						Fnt:  selFont,
						Clip: logPlayer,
					})
				}
				if CurrentSave.IsDrawerBroken {
					if CurrentSave.HasSpookyBoardPointer {
						engo.Mailbox.Dispatch(PhaseSetMessage{
							Phase: ListenPhase,
						})
					} else {
						engo.Mailbox.Dispatch(AcceptSetMessage{
							AcceptFunc: func() {
								msgs2 := []string{
									"The light was glinting off of",
									"The pointer of a spooky board.",
									"You know.",
									"For talking to the dead.",
									"OoooooOOOOOoooo",
									"Found the spooky board pointer!",
								}
								CurrentSave.HasSpookyBoardPointer = true
								dipAnim.SelectAnimationByName("empty")
								engo.Mailbox.Dispatch(PhaseDequeuMessage{})
								for _, msg := range msgs2 {
									engo.Mailbox.Dispatch(CombatLogMessage{
										Msg:  msg,
										Fnt:  selFont,
										Clip: logPlayer,
									})
								}
								engo.Mailbox.Dispatch(PhaseSetMessage{
									Phase: ListenPhase,
								})
								engo.Mailbox.Dispatch(PhaseSetMessage{
									Phase: LogClearPhase,
								})
								engo.Mailbox.Dispatch(PhaseSetMessage{
									Phase: WalkPhase,
								})
								engo.Mailbox.Dispatch(PhaseDequeuMessage{})
							},
						})
						engo.Mailbox.Dispatch(PhaseSetMessage{
							Phase: AcceptPhase,
						})
					}
				} else {
					engo.Mailbox.Dispatch(PhaseSetMessage{
						Phase: ListenPhase,
					})
				}
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: LogClearPhase,
				})
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: WalkPhase,
				})
				engo.Mailbox.Dispatch(PhaseDequeuMessage{})
			},
		},
		interestInfo{
			URL:      "president/discord.png",
			Position: engo.Point{X: 500, Y: 154},
			Shapes: []common.Shape{common.Shape{Lines: []engo.Line{
				engo.Line{P1: engo.Point{X: -10, Y: 52}, P2: engo.Point{X: 64, Y: 52}},
				engo.Line{P1: engo.Point{X: 64, Y: 52}, P2: engo.Point{X: 64, Y: 80}},
				engo.Line{P1: engo.Point{X: 64, Y: 80}, P2: engo.Point{X: -10, Y: 80}},
				engo.Line{P1: engo.Point{X: -10, Y: 80}, P2: engo.Point{X: -10, Y: 52}},
			}}},
			Func: func() {
				msgs := []string{
					"It's some sort of top-secret",
					"communication device.",
					"Would you like to turn it on?",
				}
				for _, msg := range msgs {
					engo.Mailbox.Dispatch(CombatLogMessage{
						Msg:  msg,
						Fnt:  selFont,
						Clip: logPlayer,
					})
				}
				engo.Mailbox.Dispatch(AcceptSetMessage{
					AcceptFunc: func() {
						navigateToPageImpl("https://discord.gg/QpyyrUY6JR")
						audioSys.Pause()
						engo.Mailbox.Dispatch(PhaseDequeuMessage{})
						engo.Mailbox.Dispatch(CombatLogMessage{
							Msg:  "Done listening?",
							Fnt:  selFont,
							Clip: logPlayer,
						})
						engo.Mailbox.Dispatch(AcceptSetMessage{
							AcceptFunc: func() {
								audioSys.Restart()
							},
						})
						engo.Mailbox.Dispatch(PhaseSetMessage{
							Phase: AcceptPhase,
						})
						engo.Mailbox.Dispatch(PhaseSetMessage{
							Phase: LogClearPhase,
						})
						engo.Mailbox.Dispatch(PhaseSetMessage{
							Phase: WalkPhase,
						})
						engo.Mailbox.Dispatch(PhaseDequeuMessage{})
					},
				})
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: AcceptPhase,
				})
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: LogClearPhase,
				})
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: WalkPhase,
				})
				engo.Mailbox.Dispatch(PhaseDequeuMessage{})
			},
		},
		interestInfo{
			URL:      "president/desk.png",
			Position: engo.Point{X: 152, Y: 120},
			Shapes: []common.Shape{common.Shape{Lines: []engo.Line{
				engo.Line{P1: engo.Point{X: 6, Y: 50}, P2: engo.Point{X: 118, Y: 50}},
				engo.Line{P1: engo.Point{X: 118, Y: 50}, P2: engo.Point{X: 118, Y: 94}},
				engo.Line{P1: engo.Point{X: 118, Y: 94}, P2: engo.Point{X: 6, Y: 94}},
				engo.Line{P1: engo.Point{X: 6, Y: 94}, P2: engo.Point{X: 6, Y: 50}},
			}}},
			Func: func() {
				msgs := []string{"It's an old oak desk."}
				r := rand.Intn(11)
				accept := false
				acceptFunc := func() {}
				switch r {
				case 0:
					msgs = append(msgs, "There's a headset on the desk.")
					msgs = append(msgs, "wanna put it on?")
					accept = true
					acceptFunc = func() {
						navigateToPageImpl("https://open.spotify.com/playlist/3sFTfG9vBVX1NidgBizVZ7?si=08a7ecd1b7af4338")
						audioSys.Pause()
						engo.Mailbox.Dispatch(PhaseDequeuMessage{})
						engo.Mailbox.Dispatch(CombatLogMessage{
							Msg:  "Done listening?",
							Fnt:  selFont,
							Clip: logPlayer,
						})
						engo.Mailbox.Dispatch(AcceptSetMessage{
							AcceptFunc: func() {
								audioSys.Restart()
							},
						})
						engo.Mailbox.Dispatch(PhaseSetMessage{
							Phase: AcceptPhase,
						})
						engo.Mailbox.Dispatch(PhaseSetMessage{
							Phase: LogClearPhase,
						})
						engo.Mailbox.Dispatch(PhaseSetMessage{
							Phase: WalkPhase,
						})
						engo.Mailbox.Dispatch(PhaseDequeuMessage{})
					}
				case 1, 2, 3:
					if CurrentSave.IsDrawerBroken {
						msgs = append(msgs, "The drawer here is completely obliterated.")
						msgs = append(msgs, "Guess I don't know my own strength!")
					} else {
						msgs = append(msgs, "There's a key still in one of the drawers.")
						msgs = append(msgs, "Want to try to open it?")
						accept = true
						acceptFunc = func() {
							engo.Mailbox.Dispatch(PhaseDequeuMessage{})
							//gotta roll a 8 or higher!
							roll := rand.Intn(20)
							messages := []string{}
							if roll >= 7 {
								messages = append(messages, "You gently tug at the drawer handle")
								messages = append(messages, "...")
								messages = append(messages, "oops.")
								CurrentSave.IsDrawerBroken = true
								crashPlayer.Play()
								dipAnim.SelectAnimationByName("sparkle")
							} else {
								messages = append(messages, "You yank on the drawer")
								messages = append(messages, "with everything you can muster!")
								messages = append(messages, "... !!!")
								messages = append(messages, "... !!! ??? !!!")
								messages = append(messages, "... it won't budge!")
							}
							for _, msg := range messages {
								engo.Mailbox.Dispatch(CombatLogMessage{
									Msg:  msg,
									Fnt:  selFont,
									Clip: logPlayer,
								})
							}
							engo.Mailbox.Dispatch(PhaseSetMessage{
								Phase: ListenPhase,
							})
							engo.Mailbox.Dispatch(PhaseSetMessage{
								Phase: LogClearPhase,
							})
							engo.Mailbox.Dispatch(PhaseSetMessage{
								Phase: WalkPhase,
							})
							engo.Mailbox.Dispatch(PhaseDequeuMessage{})
						}
					}
				case 4:
					msgs = append(msgs, "There's no work being done on the laptop.")
					msgs = append(msgs, "Only a ton of unanswered emails, a ")
					msgs = append(msgs, "realllly long to-do list, ")
					msgs = append(msgs, "and a lot of weird puppet-based websites open.")
				case 7, 8, 9:
					msgs = append(msgs, "There's a floppy disc on the desk labeled")
					msgs = append(msgs, "...haunted?")
					msgs = append(msgs, "Put it in the computer and try it?")
					accept = true
					acceptFunc = func() {
						engo.Mailbox.Dispatch(PhaseDequeuMessage{})
						CurrentSave.PlayerLocation = playa.Position
						engo.SetSceneByName("Ghost Fight!!!", true)
					}
				case 5, 6, 10:
					if CurrentSave.IsDrawerBroken {
						msgs = append(msgs, "Looks like when the drawer broke")
						msgs = append(msgs, "It knocked a bunch of the papers away.")
						msgs = append(msgs, "Underneath them was a key!")
						msgs = append(msgs, "Obtained the Desk Key!")
						CurrentSave.HasDeskKey = true
					} else {
						msgs = append(msgs, "There's a bunch of papers, floppy discs,")
						msgs = append(msgs, "half-eaten food containers, and other")
						msgs = append(msgs, "debris strewn around the desk.")
					}
				}
				for _, msg := range msgs {
					engo.Mailbox.Dispatch(CombatLogMessage{
						Msg:  msg,
						Fnt:  selFont,
						Clip: logPlayer,
					})
				}
				if accept {
					engo.Mailbox.Dispatch(AcceptSetMessage{
						AcceptFunc: acceptFunc,
					})
					engo.Mailbox.Dispatch(PhaseSetMessage{
						Phase: AcceptPhase,
					})
				} else {
					engo.Mailbox.Dispatch(PhaseSetMessage{
						Phase: ListenPhase,
					})
				}
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: LogClearPhase,
				})
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: WalkPhase,
				})
				engo.Mailbox.Dispatch(PhaseDequeuMessage{})
			},
		},
		interestInfo{
			URL:      "president/donations.png",
			Position: engo.Point{X: 360, Y: 112},
			Shapes: []common.Shape{common.Shape{Lines: []engo.Line{
				engo.Line{P1: engo.Point{X: -2, Y: 34}, P2: engo.Point{X: 100, Y: 34}},
				engo.Line{P1: engo.Point{X: 100, Y: 34}, P2: engo.Point{X: 78, Y: 56}},
				engo.Line{P1: engo.Point{X: 78, Y: 56}, P2: engo.Point{X: -24, Y: 56}},
				engo.Line{P1: engo.Point{X: -24, Y: 56}, P2: engo.Point{X: -2, Y: 34}},
			}}},
			Func: func() {
				msgs := []string{
					"Raising money for a good cause!",
					"Buying shiny new machines!",
					"Would you like to drop in a few coins?",
				}
				for _, msg := range msgs {
					engo.Mailbox.Dispatch(CombatLogMessage{
						Msg:  msg,
						Fnt:  selFont,
						Clip: logPlayer,
					})
				}
				engo.Mailbox.Dispatch(AcceptSetMessage{
					AcceptFunc: func() {
						navigateToPageImpl("https://www.buymeacoffee.com/Letssavesummer")
						audioSys.Pause()
						engo.Mailbox.Dispatch(PhaseDequeuMessage{})
						engo.Mailbox.Dispatch(CombatLogMessage{
							Msg:  "Thank you!!!",
							Fnt:  selFont,
							Clip: logPlayer,
						})
						engo.Mailbox.Dispatch(AcceptSetMessage{
							AcceptFunc: func() {
								audioSys.Restart()
							},
						})
						engo.Mailbox.Dispatch(PhaseSetMessage{
							Phase: AcceptPhase,
						})
						engo.Mailbox.Dispatch(PhaseSetMessage{
							Phase: LogClearPhase,
						})
						engo.Mailbox.Dispatch(PhaseSetMessage{
							Phase: WalkPhase,
						})
						engo.Mailbox.Dispatch(PhaseDequeuMessage{})
					},
				})
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: AcceptPhase,
				})
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: LogClearPhase,
				})
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: WalkPhase,
				})
				engo.Mailbox.Dispatch(PhaseDequeuMessage{})
			},
		},
		interestInfo{
			URL:      "president/engo.png",
			Position: engo.Point{X: 376, Y: 186},
			Shapes: []common.Shape{common.Shape{Lines: []engo.Line{
				engo.Line{P1: engo.Point{X: 24, Y: 0}, P2: engo.Point{X: 100, Y: 0}},
				engo.Line{P1: engo.Point{X: 100, Y: 0}, P2: engo.Point{X: 46, Y: 56}},
				engo.Line{P1: engo.Point{X: 46, Y: 56}, P2: engo.Point{X: -32, Y: 56}},
				engo.Line{P1: engo.Point{X: -32, Y: 56}, P2: engo.Point{X: 24, Y: 0}},
			}}},
			Func: func() {
				msgs := []string{
					"I goof off for hours with this thing ",
					"instead of working.",
					"I mean... ehm.",
					"It's the game engine all this is built on.",
					"Fun to use. Open source.",
					"Wanna check out the website?",
				}
				for _, msg := range msgs {
					engo.Mailbox.Dispatch(CombatLogMessage{
						Msg:  msg,
						Fnt:  selFont,
						Clip: logPlayer,
					})
				}
				engo.Mailbox.Dispatch(AcceptSetMessage{
					AcceptFunc: func() {
						navigateToPageImpl("https://engoengine.github.io")
						audioSys.Pause()
						engo.Mailbox.Dispatch(PhaseDequeuMessage{})
						engo.Mailbox.Dispatch(CombatLogMessage{
							Msg:  "Wasn't that a blast?",
							Fnt:  selFont,
							Clip: logPlayer,
						})
						engo.Mailbox.Dispatch(AcceptSetMessage{
							AcceptFunc: func() {
								audioSys.Restart()
							},
						})
						engo.Mailbox.Dispatch(PhaseSetMessage{
							Phase: AcceptPhase,
						})
						engo.Mailbox.Dispatch(PhaseSetMessage{
							Phase: LogClearPhase,
						})
						engo.Mailbox.Dispatch(PhaseSetMessage{
							Phase: WalkPhase,
						})
						engo.Mailbox.Dispatch(PhaseDequeuMessage{})
					},
				})
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: AcceptPhase,
				})
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: LogClearPhase,
				})
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: WalkPhase,
				})
				engo.Mailbox.Dispatch(PhaseDequeuMessage{})
			},
		},
		interestInfo{
			URL:      "president/safe.png",
			Position: engo.Point{X: 104, Y: 116},
			Shapes: []common.Shape{common.Shape{Lines: []engo.Line{
				engo.Line{P1: engo.Point{X: -6, Y: 26}, P2: engo.Point{X: 52, Y: 26}},
				engo.Line{P1: engo.Point{X: 52, Y: 26}, P2: engo.Point{X: 30, Y: 48}},
				engo.Line{P1: engo.Point{X: 30, Y: 48}, P2: engo.Point{X: -28, Y: 48}},
				engo.Line{P1: engo.Point{X: -28, Y: 48}, P2: engo.Point{X: -6, Y: 26}},
			}}},
			Func: func() {
				msgs := []string{"It's a top-secret safe!"}
				var acceptFunc func()
				if CurrentSave.IsSafeOpen && CurrentSave.HasSpookyBoard {
					msgs = append(msgs, "...that's already open!")
				} else {
					if CurrentSave.HasNaniteKey && !CurrentSave.NaniteKeyInSafe {
						msgs = append(msgs, "This key slot glows with the power of nanites!")
						msgs = append(msgs, "Would you like to put the nanite key in the slot?")
						acceptFunc = func() {
							CurrentSave.KeyCount++
							checkKeyCount(&safeAnim, selFont, logPlayer)
							CurrentSave.NaniteKeyInSafe = true
							engo.Mailbox.Dispatch(PhaseDequeuMessage{})
							msgs2 := []string{
								"You put the nanite key in the safe.",
								"The safe hums with nanite energy.",
							}
							for _, msg := range msgs2 {
								engo.Mailbox.Dispatch(CombatLogMessage{
									Msg:  msg,
									Fnt:  selFont,
									Clip: logPlayer,
								})
							}
							engo.Mailbox.Dispatch(PhaseSetMessage{
								Phase: ListenPhase,
							})
							engo.Mailbox.Dispatch(PhaseSetMessage{
								Phase: LogClearPhase,
							})
							engo.Mailbox.Dispatch(PhaseSetMessage{
								Phase: WalkPhase,
							})
							engo.Mailbox.Dispatch(PhaseDequeuMessage{})
						}
					} else if CurrentSave.HasDeskKey && !CurrentSave.DeskKeyInSafe {
						msgs = append(msgs, "This key slot is oaken.")
						msgs = append(msgs, "Pretty strange for an electronic safe.")
						msgs = append(msgs, "Would you like to put the desk key in the slot?")
						acceptFunc = func() {
							CurrentSave.KeyCount++
							checkKeyCount(&safeAnim, selFont, logPlayer)
							CurrentSave.DeskKeyInSafe = true
							engo.Mailbox.Dispatch(PhaseDequeuMessage{})
							msgs2 := []string{
								"You put the desk key in the oaken slot.",
								"The safe begins to photosynthesize.",
							}
							for _, msg := range msgs2 {
								engo.Mailbox.Dispatch(CombatLogMessage{
									Msg:  msg,
									Fnt:  selFont,
									Clip: logPlayer,
								})
							}
							engo.Mailbox.Dispatch(PhaseSetMessage{
								Phase: ListenPhase,
							})
							engo.Mailbox.Dispatch(PhaseSetMessage{
								Phase: LogClearPhase,
							})
							engo.Mailbox.Dispatch(PhaseSetMessage{
								Phase: WalkPhase,
							})
							engo.Mailbox.Dispatch(PhaseDequeuMessage{})
						}
					} else if CurrentSave.HasHoodKey && !CurrentSave.HoodKeyInSafe {
						msgs = append(msgs, "This key slot looks lab grown.")
						msgs = append(msgs, "Would you like to put the lab key in the slot?")
						acceptFunc = func() {
							CurrentSave.KeyCount++
							checkKeyCount(&safeAnim, selFont, logPlayer)
							CurrentSave.HoodKeyInSafe = true
							engo.Mailbox.Dispatch(PhaseDequeuMessage{})
							msgs2 := []string{
								"You put the lab key in the safe.",
								"The safe begins to fizz and pop.",
								"Hope the chemicals on that key didn't ",
								"hurt anything.",
							}
							for _, msg := range msgs2 {
								engo.Mailbox.Dispatch(CombatLogMessage{
									Msg:  msg,
									Fnt:  selFont,
									Clip: logPlayer,
								})
							}
							engo.Mailbox.Dispatch(PhaseSetMessage{
								Phase: ListenPhase,
							})
							engo.Mailbox.Dispatch(PhaseSetMessage{
								Phase: LogClearPhase,
							})
							engo.Mailbox.Dispatch(PhaseSetMessage{
								Phase: WalkPhase,
							})
							engo.Mailbox.Dispatch(PhaseDequeuMessage{})
						}
					} else if CurrentSave.HasSpaceKey && !CurrentSave.SpaceKeyInSafe {
						msgs = append(msgs, "This key slot is floating!!")
						msgs = append(msgs, "Would you like to put the space key in the slot?")
						acceptFunc = func() {
							engo.Mailbox.Dispatch(PhaseDequeuMessage{})
							msgs2 := []string{
								"You put the space key in the safe.",
								"The safe appears much lighter.",
							}
							for _, msg := range msgs2 {
								engo.Mailbox.Dispatch(CombatLogMessage{
									Msg:  msg,
									Fnt:  selFont,
									Clip: logPlayer,
								})
							}
							engo.Mailbox.Dispatch(PhaseSetMessage{
								Phase: ListenPhase,
							})
							engo.Mailbox.Dispatch(PhaseSetMessage{
								Phase: LogClearPhase,
							})
							engo.Mailbox.Dispatch(PhaseSetMessage{
								Phase: WalkPhase,
							})
							engo.Mailbox.Dispatch(PhaseDequeuMessage{})
							checkKeyCount(&safeAnim, selFont, logPlayer)
							CurrentSave.KeyCount++
							CurrentSave.SpaceKeyInSafe = true
						}
					} else if CurrentSave.DeskKeyInSafe && CurrentSave.NaniteKeyInSafe &&
						CurrentSave.HoodKeyInSafe && CurrentSave.SpaceKeyInSafe &&
						!CurrentSave.IsSafeOpen {
						msgs = append(msgs,
							"Oh, wow! Looks like you have collected",
							"All 4 keys!",
							"Great job!",
							"Open the safe?",
						)
						acceptFunc = func() {
							engo.Mailbox.Dispatch(PhaseDequeuMessage{})
							safeAnim.SelectAnimationByName("open")
							msgs2 := []string{
								"Inside the safe is...",
								"A board game?",
								"Looks like one of those boards for",
								"talking to spirits.",
								"Obtained the spooky board!",
							}
							CurrentSave.HasSpookyBoard = true
							CurrentSave.IsSafeOpen = true
							for _, msg := range msgs2 {
								engo.Mailbox.Dispatch(CombatLogMessage{
									Msg:  msg,
									Fnt:  selFont,
									Clip: logPlayer,
								})
							}
							engo.Mailbox.Dispatch(PhaseSetMessage{
								Phase: ListenPhase,
							})
							engo.Mailbox.Dispatch(PhaseSetMessage{
								Phase: LogClearPhase,
							})
							engo.Mailbox.Dispatch(PhaseSetMessage{
								Phase: WalkPhase,
							})
							engo.Mailbox.Dispatch(PhaseDequeuMessage{})
						}
					} else {
						msgs = append(msgs, "You don't have any more keys.")
						msgs = append(msgs, "Look around for more!")
					}
				}
				for _, msg := range msgs {
					engo.Mailbox.Dispatch(CombatLogMessage{
						Msg:  msg,
						Fnt:  selFont,
						Clip: logPlayer,
					})
				}
				if acceptFunc != nil {
					engo.Mailbox.Dispatch(AcceptSetMessage{
						AcceptFunc: acceptFunc,
					})
					engo.Mailbox.Dispatch(PhaseSetMessage{
						Phase: AcceptPhase,
					})
				} else {
					engo.Mailbox.Dispatch(PhaseSetMessage{
						Phase: ListenPhase,
					})
				}
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: LogClearPhase,
				})
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: WalkPhase,
				})
				engo.Mailbox.Dispatch(PhaseDequeuMessage{})
			},
		},
	})

	animSys.Add(pres.interests[0].GetBasicEntity(), dipAnim.GetAnimationComponent(), pres.interests[0].GetRenderComponent())
	pres.interests[3].GetRenderComponent().Scale = engo.Point{X: 2, Y: 2}
	pres.interests[5].GetRenderComponent().Scale = engo.Point{X: 2, Y: 2}
	animSys.Add(pres.interests[5].GetBasicEntity(), safeAnim.GetAnimationComponent(), pres.interests[5].GetRenderComponent())

	space := newRoom(w, engo.Point{X: 0, Y: 1500}, "space/bg.png", []wallInfo{
		wallInfo{
			Shapes: []common.Shape{common.Shape{Lines: []engo.Line{
				engo.Line{P1: engo.Point{X: 54, Y: 256}, P2: engo.Point{X: 142, Y: 168}},
				engo.Line{P1: engo.Point{X: 142, Y: 168}, P2: engo.Point{X: 140, Y: 130}},
				engo.Line{P1: engo.Point{X: 140, Y: 130}, P2: engo.Point{X: 50, Y: 220}},
				engo.Line{P1: engo.Point{X: 50, Y: 220}, P2: engo.Point{X: 54, Y: 256}},
			}}},
		},
		wallInfo{
			Shapes: []common.Shape{common.Shape{Lines: []engo.Line{
				engo.Line{P1: engo.Point{X: 462, Y: 130}, P2: engo.Point{X: 552, Y: 218}},
				engo.Line{P1: engo.Point{X: 552, Y: 218}, P2: engo.Point{X: 552, Y: 256}},
				engo.Line{P1: engo.Point{X: 552, Y: 256}, P2: engo.Point{X: 460, Y: 168}},
				engo.Line{P1: engo.Point{X: 460, Y: 168}, P2: engo.Point{X: 462, Y: 130}},
			}}},
		},
		wallInfo{
			Position: engo.Point{X: 138, Y: 128},
			Width:    320,
			Height:   42,
		},
		wallInfo{
			Position: engo.Point{X: 82, Y: 190},
			Width:    158,
			Height:   48,
		},
		wallInfo{
			Position: engo.Point{X: 52, Y: 254},
			Width:    500,
			Height:   4,
		},
	}, []doorInfo{
		doorInfo{
			URL:          "space/doorSS.png",
			Position:     engo.Point{X: 388, Y: 236},
			TeleportTo:   engo.Point{X: 382, Y: 88},
			CellWidth:    40,
			CellHeight:   10,
			BorderWidth:  1,
			BorderHeight: 1,
			Shapes: []common.Shape{common.Shape{Lines: []engo.Line{
				engo.Line{P1: engo.Point{X: -10, Y: 0}, P2: engo.Point{X: -10, Y: 20}},
				engo.Line{P1: engo.Point{X: 10, Y: 20}, P2: engo.Point{X: 70, Y: 20}},
				engo.Line{P1: engo.Point{X: 70, Y: 20}, P2: engo.Point{X: 70, Y: 0}},
				engo.Line{P1: engo.Point{X: 70, Y: 0}, P2: engo.Point{X: -10, Y: 0}},
			}}},
			OpenFrames:  []int{0, 1, 2},
			CloseFrames: []int{2, 1, 0},
			Button:      "down",
		},
	}, []interestInfo{
		interestInfo{
			URL:      "space/moon.png",
			Position: engo.Point{X: 102, Y: 144},
			Shapes: []common.Shape{common.Shape{Lines: []engo.Line{
				engo.Line{P1: engo.Point{X: -20, Y: 0}, P2: engo.Point{X: 136, Y: 0}},
				engo.Line{P1: engo.Point{X: 136, Y: 0}, P2: engo.Point{X: 136, Y: 100}},
				engo.Line{P1: engo.Point{X: 136, Y: 100}, P2: engo.Point{X: -20, Y: 100}},
				engo.Line{P1: engo.Point{X: -20, Y: 100}, P2: engo.Point{X: -20, Y: 0}},
			}}},
			Func: func() {
				msgs := []string{
					"Oh noooooo...",
					"the moon...",
					"it's broken!",
					"I knew I should've found the",
					"STUD FINDER",
					"before hanging it!",
				}
				for _, msg := range msgs {
					engo.Mailbox.Dispatch(CombatLogMessage{
						Msg:  msg,
						Fnt:  selFont,
						Clip: logPlayer,
					})
				}
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: ListenPhase,
				})
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: LogClearPhase,
				})
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: WalkPhase,
				})
				engo.Mailbox.Dispatch(PhaseDequeuMessage{})
			},
		},
		interestInfo{
			URL:      "space/tv.png",
			Position: engo.Point{X: 402, Y: 66},
			Shapes: []common.Shape{common.Shape{Lines: []engo.Line{
				engo.Line{P1: engo.Point{X: 0, Y: 75}, P2: engo.Point{X: 60, Y: 75}},
				engo.Line{P1: engo.Point{X: 60, Y: 75}, P2: engo.Point{X: 60, Y: 105}},
				engo.Line{P1: engo.Point{X: 60, Y: 105}, P2: engo.Point{X: 0, Y: 105}},
				engo.Line{P1: engo.Point{X: 0, Y: 105}, P2: engo.Point{X: 0, Y: 75}},
			}}},
			Func: func() {
				msgs := []string{
					"Mars is all over the news!",
					"Something big must be going on!",
					"Want to check it out?",
				}
				for _, msg := range msgs {
					engo.Mailbox.Dispatch(CombatLogMessage{
						Msg:  msg,
						Fnt:  selFont,
						Clip: logPlayer,
					})
				}
				engo.Mailbox.Dispatch(AcceptSetMessage{
					AcceptFunc: func() {
						navigateToPageImpl("https://www.marsbound.space")
						audioSys.Pause()
						engo.Mailbox.Dispatch(PhaseDequeuMessage{})
						engo.Mailbox.Dispatch(CombatLogMessage{
							Msg:  "How was your trip?",
							Fnt:  selFont,
							Clip: logPlayer,
						})
						engo.Mailbox.Dispatch(AcceptSetMessage{
							AcceptFunc: func() {
								audioSys.Restart()
							},
						})
						engo.Mailbox.Dispatch(PhaseSetMessage{
							Phase: AcceptPhase,
						})
						engo.Mailbox.Dispatch(PhaseSetMessage{
							Phase: LogClearPhase,
						})
						engo.Mailbox.Dispatch(PhaseSetMessage{
							Phase: WalkPhase,
						})
						engo.Mailbox.Dispatch(PhaseDequeuMessage{})
					},
				})
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: AcceptPhase,
				})
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: LogClearPhase,
				})
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: WalkPhase,
				})
				engo.Mailbox.Dispatch(PhaseDequeuMessage{})
			},
		},
		interestInfo{
			URL:      "space/sauce.png",
			Position: engo.Point{X: 528, Y: 128},
			Shapes: []common.Shape{common.Shape{Lines: []engo.Line{
				engo.Line{P1: engo.Point{X: -56, Y: 56}, P2: engo.Point{X: -6, Y: 56}},
				engo.Line{P1: engo.Point{X: -6, Y: 56}, P2: engo.Point{X: -6, Y: 116}},
				engo.Line{P1: engo.Point{X: -6, Y: 116}, P2: engo.Point{X: -56, Y: 116}},
				engo.Line{P1: engo.Point{X: -56, Y: 116}, P2: engo.Point{X: -56, Y: 56}},
			}}},
			Func: func() {
				msgs := []string{
					"Yum! Spicy!",
					"Wanna check out the SAUCE?",
				}
				for _, msg := range msgs {
					engo.Mailbox.Dispatch(CombatLogMessage{
						Msg:  msg,
						Fnt:  selFont,
						Clip: logPlayer,
					})
				}
				engo.Mailbox.Dispatch(AcceptSetMessage{
					AcceptFunc: func() {
						navigateToPageImpl("hhttps://github.com/Noofbiz/MarsBound")
						audioSys.Pause()
						engo.Mailbox.Dispatch(PhaseDequeuMessage{})
						engo.Mailbox.Dispatch(CombatLogMessage{
							Msg:  "Wasn't that delicious?",
							Fnt:  selFont,
							Clip: logPlayer,
						})
						engo.Mailbox.Dispatch(AcceptSetMessage{
							AcceptFunc: func() {
								audioSys.Restart()
							},
						})
						engo.Mailbox.Dispatch(PhaseSetMessage{
							Phase: AcceptPhase,
						})
						engo.Mailbox.Dispatch(PhaseSetMessage{
							Phase: LogClearPhase,
						})
						engo.Mailbox.Dispatch(PhaseSetMessage{
							Phase: WalkPhase,
						})
						engo.Mailbox.Dispatch(PhaseDequeuMessage{})
					},
				})
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: AcceptPhase,
				})
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: LogClearPhase,
				})
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: WalkPhase,
				})
				engo.Mailbox.Dispatch(PhaseDequeuMessage{})
			},
		},
		interestInfo{
			URL:      "space/window.png",
			Position: engo.Point{X: 25, Y: -5},
			Shapes: []common.Shape{common.Shape{Lines: []engo.Line{
				engo.Line{P1: engo.Point{X: 204, Y: 156}, P2: engo.Point{X: 354, Y: 156}},
				engo.Line{P1: engo.Point{X: 354, Y: 156}, P2: engo.Point{X: 354, Y: 186}},
				engo.Line{P1: engo.Point{X: 354, Y: 186}, P2: engo.Point{X: 204, Y: 186}},
				engo.Line{P1: engo.Point{X: 204, Y: 186}, P2: engo.Point{X: 204, Y: 156}},
			}}},
			Func: func() {
				var acceptFunc func()
				msgs := []string{
					"Outside you see glittering stars.",
					"Space is calling you!",
					"Apply for an internship at Blue Origin today!",
				}
				if !CurrentSave.HasSpaceKey {
					// roll perception
					// gotta get a 9 or higher!
					roll := rand.Intn(20)
					if roll >= 8 {
						msgs = append(msgs,
							"Wait a second.",
							"In the window there!",
							"It's a keyboard!",
							"Show off your",
							"SICK TYPING SKILLS?",
						)
						acceptFunc = func() {
							engo.Mailbox.Dispatch(PhaseDequeuMessage{})
							rollAgain := rand.Intn(20)
							var msgs2 []string
							if rollAgain < 7 {
								msgs2 = append(msgs2,
									"You bash on the keyboard",
									"with all your might!",
									"That was fun!",
									"You never were good at typing.",
									"Oops.",
									"Looks like something broke.",
									"You just slip it in your pocket",
									"If they can't find the key,",
									"they won't know it's broken!",
									"You obtained the",
									"SPACE KEY",
								)
								CurrentSave.HasSpaceKey = true
							} else if roll < 10 {
								msgs2 = append(msgs2,
									"You tap away at the keyboard.",
									"Nothing really special about it.",
									"Kinda boring.",
								)
							} else {
								msgs2 = append(msgs2,
									"You begin tapping away at the keyboard",
									"On the screen behind you, an intense",
									"game starts up. You get really into it.",
									"You lose track of time.",
									"After playing for what feels like days",
									"The keyboard gives out.",
									"The space key finally pops right out!",
									"You keep it as a momento of that epic game.",
									"You obtained the",
									"SPACE KEY",
								)
								CurrentSave.HasSpaceKey = true
							}
							for _, msg := range msgs2 {
								engo.Mailbox.Dispatch(CombatLogMessage{
									Msg:  msg,
									Fnt:  selFont,
									Clip: logPlayer,
								})
							}
							engo.Mailbox.Dispatch(PhaseSetMessage{
								Phase: ListenPhase,
							})
							engo.Mailbox.Dispatch(PhaseSetMessage{
								Phase: LogClearPhase,
							})
							engo.Mailbox.Dispatch(PhaseSetMessage{
								Phase: WalkPhase,
							})
							engo.Mailbox.Dispatch(PhaseDequeuMessage{})
						}
					}
				}
				for _, msg := range msgs {
					engo.Mailbox.Dispatch(CombatLogMessage{
						Msg:  msg,
						Fnt:  selFont,
						Clip: logPlayer,
					})
				}
				if acceptFunc != nil {
					engo.Mailbox.Dispatch(AcceptSetMessage{
						AcceptFunc: acceptFunc,
					})
					engo.Mailbox.Dispatch(PhaseSetMessage{
						Phase: AcceptPhase,
					})
				} else {
					engo.Mailbox.Dispatch(PhaseSetMessage{
						Phase: ListenPhase,
					})
				}
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: LogClearPhase,
				})
				engo.Mailbox.Dispatch(PhaseSetMessage{
					Phase: WalkPhase,
				})
				engo.Mailbox.Dispatch(PhaseDequeuMessage{})
			},
		},
	})

	//tv Animation
	tvSS := common.NewSpritesheetWithBorderFromFile("space/tvSS.png", 35, 30, 1, 1)
	tvAnim := animation{}
	tvAnim.AnimationComponent = common.NewAnimationComponent(tvSS.Drawables(), 0.3)
	tvAnim.AddDefaultAnimation(&common.Animation{
		Name:   "scroll",
		Frames: []int{0, 1, 2, 3, 4, 5, 6, 7},
		Loop:   true,
	})
	tvAnim.AnimationComponent.SelectAnimationByName("scroll")
	space.interests[1].Drawable = tvSS.Drawable(0)
	space.interests[1].Scale = engo.Point{X: 2, Y: 2}
	space.interests[1].SetZIndex(6)
	animSys.Add(space.interests[1].GetBasicEntity(), tvAnim.GetAnimationComponent(), space.interests[1].GetRenderComponent())

	//window Animation
	windowSS := common.NewSpritesheetWithBorderFromFile("space/windowSS.png", 275, 80, 1, 1)
	windowAnim := animation{}
	windowAnim.AnimationComponent = common.NewAnimationComponent(windowSS.Drawables(), 0.3)
	windowAnim.AddDefaultAnimation(&common.Animation{
		Name:   "twinkle",
		Frames: []int{0, 1, 2, 1},
		Loop:   true,
	})
	windowAnim.AnimationComponent.SelectAnimationByName("twinkle")
	space.interests[3].Drawable = lenSS.Drawable(0)
	space.interests[3].Scale = engo.Point{X: 2, Y: 2}
	space.interests[3].SetZIndex(5)
	animSys.Add(space.interests[3].GetBasicEntity(), windowAnim.GetAnimationComponent(), space.interests[3].GetRenderComponent())

	msgs := []string{
		"Where am I?",
		"Oh well...",
		"Welcome to Skeleboy Studios!",
		"My name is Jerry!",
		"I make games!",
		"Currently I'm working on Marsbound",
		"An adventure to mars!",
		"Look around to see what else is afoot!",
	}
	for _, msg := range msgs {
		engo.Mailbox.Dispatch(CombatLogMessage{
			Msg:  msg,
			Fnt:  selFont,
			Clip: logPlayer,
		})
	}
	engo.Mailbox.Dispatch(PhaseSetMessage{
		Phase: ListenPhase,
	})
	engo.Mailbox.Dispatch(PhaseSetMessage{
		Phase: LogClearPhase,
	})
	engo.Mailbox.Dispatch(PhaseSetMessage{
		Phase: WalkPhase,
	})
}

func checkKeyCount(a *animation, fnt *common.Font, clip *common.Player) {
	switch CurrentSave.KeyCount {
	case 0:
		a.SelectAnimationByName("0")
	case 1:
		a.SelectAnimationByName("1")
	case 2:
		a.SelectAnimationByName("2")
	case 3:
		a.SelectAnimationByName("3")
	default:
		a.SelectAnimationByName("4")
	}
}
