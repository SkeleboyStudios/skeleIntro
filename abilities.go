package main

import (
	"math/rand"
	"strconv"

	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type Ability struct {
	Title        string
	Shorthand    string
	Description  string
	MPCost       float32
	TargetType   Target
	CastTimeFunc func(You *Character)
	EffectFunc   func(You *Character, TargetPlayers []*Character, TargetBaddies []*Baddie)
}

var RegularAttackAbility = Ability{
	Title:       "Regular boring old attack.",
	Shorthand:   "Fight",
	Description: "Normal Series: Normal punch.",
	MPCost:      0,
	TargetType:  TargetTypeSingleEnemy,
	CastTimeFunc: func(You *Character) {
		You.totalCastTime = 1.5 - (You.Dex / 100)
		if You.totalCastTime < 0.2 {
			You.totalCastTime = 0.2
		}
	},
	EffectFunc: func(You *Character, TargetPlayers []*Character, TargetBaddies []*Baddie) {
		//Regular attack GO!

	},
}

var DefendAbility = Ability{
	Title:        "Defend yourself!",
	Shorthand:    "Defend",
	Description:  "Protect yourself this turn to take less damage!",
	TargetType:   TargetTypeNone,
	CastTimeFunc: func(You *Character) {},
	EffectFunc:   func(You *Character, TargetPlayers []*Character, TargetBaddies []*Baddie) {},
}

var LookAroundAbility = Ability{
	Title:        "Look Around!",
	Shorthand:    "Look",
	Description:  "Look around the fight area for \nclues! Maybe something useful \nwill turn up!",
	MPCost:       0,
	TargetType:   TargetTypeNone,
	CastTimeFunc: func(You *Character) {},
	EffectFunc: func(You *Character, TargetPlayers []*Character, TargetBaddies []*Baddie) {
		//Roll perception
		msgs := []string{"You look around the room..."}
		if rng := rand.Intn(20) + 1; rng <= 5 {
			//Didn't see Anything
			msgs = append(msgs, "But don't see anything of note.")
		} else if rng <= 19 {
			//See something
			if !CurrentSave.HasSpookyBoard {
				msgs = append(msgs, "That safe over there looks SUSPICIOUS!")
				You.AddAbility(SafeSearchAbility)
			} else if !CurrentSave.HasSpookyBoardPointer {
				msgs = append(msgs,
					"Is that wall...glinting?",
					"Maybe you should check it out!",
				)
				You.AddAbility(GlintSearchAbility)
			} else if !CurrentSave.HasMedKit {
				msgs = append(msgs,
					"There's a safety kit on the back wall here",
					"Maybe there's a medkit in there",
					"That could help anyone who gets injured!",
				)
				You.AddAbility(MedKitSearchAbility)
			} else if !CurrentSave.HasSalt {
				msgs = append(msgs,
					"There's a Himylian Salt Lamp on the desk here",
					"Salt circles can help fight ghosts!",
					"Or so I've heard...",
				)
				You.AddAbility(ScratchSaltAbility)
			} else {
				msgs = append(msgs,
					"You've already seen everything in the room.",
				)
			}
		} else {
			//See EVERYTHING
			if !CurrentSave.HasSpookyBoard {
				msgs = append(msgs, "That safe over there looks SUSPICIOUS!")
				You.AddAbility(SafeSearchAbility)
			}
			if !CurrentSave.HasSpookyBoardPointer {
				msgs = append(msgs,
					"Is that wall...glinting?",
					"Maybe you should check it out!",
				)
				You.AddAbility(GlintSearchAbility)
			}
			if !CurrentSave.HasMedKit {
				msgs = append(msgs,
					"There's a safety kit on the back wall here",
					"Maybe there's a medkit in there",
					"That could help anyone who gets injured!",
				)
				You.AddAbility(MedKitSearchAbility)
			}
			if !CurrentSave.HasSalt {
				msgs = append(msgs,
					"There's a Himylian Salt Lamp on the desk here",
					"Salt circles can help fight ghosts!",
					"Or so I've heard...",
				)
				You.AddAbility(ScratchSaltAbility)
			}
			if CurrentSave.HasSpookyBoard && CurrentSave.HasSpookyBoardPointer &&
				CurrentSave.HasMedKit && CurrentSave.HasSalt {
				msgs = append(msgs,
					"You've already seen everything in the room.",
				)
			}
		}
		for _, msg := range msgs {
			engo.Mailbox.Dispatch(CombatLogMessage{
				Msg:  msg,
				Fnt:  You.Font,
				Clip: You.Clip,
			})
		}
	},
}

var SafeSearchAbility = Ability{
	Title:        "Search the Safe!",
	Shorthand:    "SSafe",
	Description:  "Search that safe! Myabe it has something useful inside!",
	MPCost:       0,
	TargetType:   TargetTypeNone,
	CastTimeFunc: func(You *Character) {},
	EffectFunc: func(You *Character, TargetPlayers []*Character, TargetBaddies []*Baddie) {
		msgs := []string{
			"Looks like the safe has 4 key holes",
			"If you have the keys in your INVENTORY",
			"You should be able to easily open it.",
			"If not though,",
			"Maybe one of the ghost's ENERGY BLASTS can break through.",
		}
		You.AddAbility(DistractAndDodgeAbility)
		msgs = append(msgs,
			"There's also a pin-pad that *might* open it if you guess it right.",
			"I can never remember it, but you might be able to ask ME what the pin is",
			"If you've recruited me...",
		)
		You.AddAbility(GuessPinAbility)
		if CurrentSave.RecruitedMe {
			You.AddAbility(AskPinAbility)
		}
		for _, msg := range msgs {
			engo.Mailbox.Dispatch(CombatLogMessage{
				Msg:  msg,
				Fnt:  You.Font,
				Clip: You.Clip,
			})
		}
	},
}

var AskPinAbility = Ability{
	Title:        "Ask about the Pin!",
	Shorthand:    "Ask",
	Description:  "Ask the President what his super secret safe pin is",
	MPCost:       0,
	TargetType:   TargetTypeSingleFriend,
	CastTimeFunc: func(You *Character) {},
	EffectFunc: func(You *Character, TargetPlayers []*Character, TargetBaddies []*Baddie) {
		tar := TargetPlayers[0]
		msgs := []string{
			"You turn to " + tar.Name + " and ask about the secret safe",
		}
		if tar.Name == "Me" {
			msgs = append(msgs,
				"My pin? I can never remember it",
				"So I never set it.",
				"It's the factory-default: 1234!",
			)
			You.RemoveAbility("Guess the Pin!")
			You.RemoveAbility("Ask about the Pin!")
			You.AddAbility(InputPinAbility)
		} else {
			if You.Name == "Me" {
				msgs = append(msgs,
					"Oh yeah!",
					"I AM THE PRESIDENT!",
				)
				You.RemoveAbility("Guess the Pin!")
				You.RemoveAbility("Ask about the Pin!")
				You.AddAbility(InputPinAbility)
			} else {
				msgs = append(msgs,
					"They look confused.",
					"Maybe that wasn't the president?",
				)
			}
		}
		for _, msg := range msgs {
			engo.Mailbox.Dispatch(CombatLogMessage{
				Msg:  msg,
				Fnt:  You.Font,
				Clip: You.Clip,
			})
		}
	},
}

var GuessPinAbility = Ability{
	Title:        "Guess the Pin!",
	Shorthand:    "Guess",
	Description:  "Try a random pin on the safe!",
	MPCost:       0,
	TargetType:   TargetTypeNone,
	CastTimeFunc: func(You *Character) {},
	EffectFunc: func(You *Character, TargetPlayers []*Character, TargetBaddies []*Baddie) {
		msgs := []string{
			"It's FOUR DIGITS!",
			"That's literally a one in ten thousand chance!",
			"But, despite the odds, you try to guess anyway...",
			"You punch in 4 random numbers and hit enter.",
		}
		if guess := rand.Intn(10000); guess == 1234 {
			cash, _ := common.LoadedPlayer("fight/cash.wav")
			if cash.IsPlaying() {
				cash.Pause()
				cash.Rewind()
			}
			cash.Play()
			msgs = append(msgs,
				"Wow. You actually guessed it!",
				"Great job!",
				"You opened the safe!",
			)
			You.RemoveAbility("Guess the Pin!")
			You.RemoveAbility("Ask about the Pin!")
			You.RemoveAbility("Input the Pin!")
			You.RemoveAbility("Distract and Dodge")
			You.AddAbility(GrabItemInSafeAbility)
		} else {
			//wrong sound
			msgs = append(msgs,
				"Guess that wasn't it.",
				"You can always try again!",
			)
		}
		for _, msg := range msgs {
			engo.Mailbox.Dispatch(CombatLogMessage{
				Msg:  msg,
				Fnt:  You.Font,
				Clip: You.Clip,
			})
		}
	},
}

var InputPinAbility = Ability{
	Title:        "Input the Pin!",
	Shorthand:    "Input",
	Description:  "Input the pin and open the safe!",
	MPCost:       0,
	TargetType:   TargetTypeNone,
	CastTimeFunc: func(You *Character) {},
	EffectFunc: func(You *Character, TargetPlayers []*Character, TargetBaddies []*Baddie) {
		msgs := []string{
			"You carefully input the pin.",
		}
		for _, msg := range msgs {
			engo.Mailbox.Dispatch(CombatLogMessage{
				Msg:  msg,
				Fnt:  You.Font,
				Clip: You.Clip,
			})
			You.RemoveAbility("Guess the Pin!")
			You.RemoveAbility("Ask about the Pin!")
			You.RemoveAbility("Input the Pin!")
			You.RemoveAbility("Distract and Dodge")
			You.AddAbility(GrabItemInSafeAbility)
		}
	},
}

var DistractAndDodgeAbility = Ability{
	Title:        "Distract and Dodge",
	Shorthand:    "DnD",
	Description:  "Distract the ghost to get its attention, then dodge its energy blast! The safe will take the brunt of the blast!",
	MPCost:       0,
	TargetType:   TargetTypeSingleEnemy,
	CastTimeFunc: func(You *Character) {},
	EffectFunc: func(You *Character, TargetPlayers []*Character, TargetBaddies []*Baddie) {
		if CurrentSave.IsSafeOpen {

		} else {

		}
		msgs := []string{
			"The Blood Mouthed Ghost begins charging up his laser!",
			"HEY! NOT-SPOOKY-AT-ALL! You couldn't hit me with that blast",
			"Just like you couldn't scare my little sister's imaginary friend!",
			"You shout.",
			"The Spookster turns to you.",
			"The blast grows bigger than it ever has!",
		}
		//Ghost Blast Animation!
		//Dodge! Quickly!
		if rand.Intn(100)+1+20+int(You.Dex) > int(TargetBaddies[0].Dex)+rand.Intn(100)+1 {
			msgs = append(msgs,
				"You narrowly dodge the blast",
				"It hits the safe dead-on!",
				"Looks like it broke the door!",
				"Go check it out!",
			)
			You.RemoveAbility("Guess the Pin!")
			You.RemoveAbility("Ask about the Pin!")
			You.RemoveAbility("Input the Pin!")
			You.RemoveAbility("Distract and Dodge")
			You.Abilities = append(You.Abilities, GrabItemInSafeAbility)
		} else {
			dmg := rand.Intn(30) + 15 + int(TargetBaddies[0].Str+TargetBaddies[0].Int)
			msgs = append(msgs,
				"The ghost fires the blast right into your FACE!",
				"Ooof. That's gotta hurt",
				"The ghost deals "+strconv.Itoa(dmg)+" damage to you!",
			)
			//big hit sound and animation!
			// if snd, err := common.LoadedPlayer("ghostblast.wav"); err == nil {
			// 	snd.Play()
			// }
			//engo.Mailbox.Dispatch(ScreenShakeMessage{})
			You.HP -= float32(dmg)
		}
		for _, msg := range msgs {
			engo.Mailbox.Dispatch(CombatLogMessage{
				Msg:  msg,
				Fnt:  You.Font,
				Clip: You.Clip,
			})
		}
	},
}

var GrabItemInSafeAbility = Ability{
	Title:        "Grab whatever's in that safe!",
	Shorthand:    "SGrab",
	Description:  "Grab what's inside the now open safe and add it to your inventory",
	MPCost:       0,
	CastTimeFunc: func(You *Character) {},
	EffectFunc: func(You *Character, TargetPlayers []*Character, TargetBaddies []*Baddie) {
		msgs := []string{}
		if CurrentSave.HasSpookyBoard {
			msgs = []string{
				"You try to grab the item from the safe",
				"But there's nothing there!",
				"Someone must've gotten here before you!",
			}
		} else {
			msgs = []string{
				"Inside the safe is a SPOOKYBOARD!",
				"The SPOOKYBOARD was added to your inventory",
			}
			CurrentSave.HasSpookyBoard = true
			CurrentSave.IsSafeOpen = true
		}
		You.RemoveAbility("Grab whatever's in that safe!")
		for _, msg := range msgs {
			engo.Mailbox.Dispatch(CombatLogMessage{
				Msg:  msg,
				Fnt:  You.Font,
				Clip: You.Clip,
			})
		}
	},
}

var GlintSearchAbility = Ability{
	Title:        "Look closer at the wall!",
	Shorthand:    "SWall",
	Description:  "Take a closer look at the gleam in the wall.",
	MPCost:       0,
	CastTimeFunc: func(You *Character) {},
	EffectFunc: func(You *Character, TargetPlayers []*Character, TargetBaddies []*Baddie) {
		msgs := []string{
			"Looks like something shiny",
			"Is just behind the wall here!",
			"The hole it gleams through is too small to get it out of.",
			"Maybe if we were to damage it?",
		}
		//Add baddie "Wall" to be knocked down
		You.RemoveAbility("Look closer at the wall!")
		for _, msg := range msgs {
			engo.Mailbox.Dispatch(CombatLogMessage{
				Msg:  msg,
				Fnt:  You.Font,
				Clip: You.Clip,
			})
		}
	},
}

var MedKitSearchAbility = Ability{
	Title:        "Search the med kit!",
	Shorthand:    "SKit",
	Description:  "Scavenge the med kit for band-aids and medicine!",
	MPCost:       5,
	CastTimeFunc: func(You *Character) {},
	EffectFunc: func(You *Character, TargetPlayers []*Character, TargetBaddies []*Baddie) {
		msgs := []string{
			"You approach the medkit",
		}
		if rand.Intn(100)+1+int(You.Int) > 35 {
			msgs = append(msgs, "And open it!", "Inside you find")
			bandaidcount := rand.Intn(6) - 2
			watercount := rand.Intn(6) - 3
			if bandaidcount > 0 {
				msgs = append(msgs, strconv.Itoa(bandaidcount)+" bandages")
				CurrentSave.BandageCount += bandaidcount
				if watercount > 0 {
					msgs = append(msgs, "and")
				}
			}
			if watercount > 0 {
				msgs = append(msgs, strconv.Itoa(watercount)+" sports drinks")
				CurrentSave.DrinkCount += watercount
			} else if bandaidcount <= 0 {
				msgs = append(msgs, "nothing.")
			}
			msgs = append(msgs,
				"You didn't have enough time to search the whole bag",
				"There's plenty more stuff inside!",
			)
		} else {
			dmg := rand.Intn(20) + 5
			msgs = append(msgs,
				"You reach into the bag",
				"But that wasn't a zipper! They're teeth!",
				"That's not a med-kit! It's a Mimic!",
				"And it takes a chomp at your arm!",
				"Ouchie! That looks like "+strconv.Itoa(dmg)+" points of damage!",
			)
			You.HP -= float32(dmg)
			// A mimic appears!
		}
		for _, msg := range msgs {
			engo.Mailbox.Dispatch(CombatLogMessage{
				Msg:  msg,
				Fnt:  You.Font,
				Clip: You.Clip,
			})
		}
	},
}

var ReadSpookyBoardAbility = Ability{
	Title:        "Ready Spooky Board!",
	Shorthand:    "RSpookB",
	Description:  "Take a reading from the spooky board!",
	MPCost:       0,
	CastTimeFunc: func(You *Character) {},
	EffectFunc: func(You *Character, TargetPlayers []*Character, TargetBaddies []*Baddie) {

	},
}

var ScratchSaltAbility = Ability{}

var ShieldsUpAbility = Ability{}

var HeatBeamEyesAbility = Ability{}

var HealBeamEyesAbility = Ability{}

var CoverAbility = Ability{}

var YouthfulSaltSplashAbility = Ability{}
