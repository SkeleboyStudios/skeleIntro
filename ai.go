package main

type AIComponent struct {
	Attacks          []Attack
	MinWait, MaxWait float32
}

type Attack struct {
	EffectFunc func(bad *Baddie, TargetPlayers []*Character, TargetBaddies []*Baddie)
	AttackTime float32
	MPCost     float32
}
