package main

import "github.com/bcspragu/Gobots/game"

//import "github.com/yongbin999/gobot/yboter"

type yboter struct{
	targets map[uint32]uint32
	aggression uint32
}

func (bt *yboter) Act(b *game.Board, r *game.Robot) game.Action {

	//if enermy adjecent, attack
	ds := []game.Direction{
		game.North,
		game.South,
		game.East,
		game.West,
	}
	for _, d := range ds {
		loc := r.Loc.Add(d)
		if opponentAt(b, loc) {
			return game.Action{
				Kind:      game.Attack,
				Direction: d,
			}
		}
	}

	var opp *game.Robot
	opp = nearestOpponent(b, r.Loc)
		if opp == nil {
			return game.Action{Kind: game.Wait}
		}

	// Move to target.
	// Don't worry about collisions, since we already shot at all neighbors.
	// TODO: but what about friends?
	// TODO: and why don't you compute the vector angle?
	switch {
	case opp.Loc.X < r.Loc.X:
		return game.Action{
			Kind:      game.Move,
			Direction: game.West,
		}
	case opp.Loc.X > r.Loc.X:
		return game.Action{
			Kind:      game.Move,
			Direction: game.East,
		}
	case opp.Loc.Y < r.Loc.Y:
		return game.Action{
			Kind:      game.Move,
			Direction: game.North,
		}
	case opp.Loc.Y > r.Loc.Y:
		return game.Action{
			Kind:      game.Move,
			Direction: game.South,
		}
	}
	// TODO: impossibru?
	//return game.Action{Kind: game.Wait}

	//else move to center
	return game.Action{
	    Kind:      game.Move,
	    Direction: game.Towards(r.Loc, b.Center()),
	}
}

func nearestOpponent(b *game.Board, loc game.Loc) *game.Robot {
	// Probably faster ways of doing this.. traversing outward
	var closest *game.Robot
	var closestDist int
	for y, row := range b.Cells {
		for x, r := range row {
			curr := game.Loc{x, y}
			if r == nil || r.Faction != game.OpponentFaction {
				continue
			}
			d := game.Distance(loc, curr)
			if closest == nil || d < closestDist {
				closest, closestDist = r, d
			}
		}
	}
	return closest
}

func opponentAt(b *game.Board, loc game.Loc) bool {
	if !b.IsInside(loc) {
		return false
	}
	r := b.At(loc)
	if r == nil {
		return false
	}
	return r.Faction == game.OpponentFaction
}

//if distance with closes enermy == 1 and 0 or have adjcent friend neighbor, attack
//if distance with closes enermy == 1 and 0 or fork adjcent friend neighbor, march

//if  self.hp - adjucent enemy count * damage <0 , suicide
//if surrounding boxes enermycount * avg damage > self.hp, suicide. 

//if theres wall go towards it. and for a slope 

//if triagle formed




func main() {
    game.StartServerForFactory("testing", "IlxVTAKDBsIisHXuGesBcBuds", game.ToFactory(&yboter{}))
}


/*
when i did ...
# github.com/bcspragu/Gobots/samplebots
..\..\bcspragu\Gobots\samplebots\aggro.go:7: undefined: game.Action
..\..\bcspragu\Gobots\samplebots\aggro.go:17: undefined: game.Action
..\..\bcspragu\Gobots\samplebots\aggro.go:23: undefined: game.Action
..\..\bcspragu\Gobots\samplebots\main.go:57: not enough arguments in call to c.RegisterAI
..\..\bcspragu\Gobots\samplebots\pathfinder.go:9: undefined: game.Action
..\..\bcspragu\Gobots\samplebots\pathfinder.go:20: undefined: game.Action
..\..\bcspragu\Gobots\samplebots\pathfinder.go:41: undefined: game.Action
..\..\bcspragu\Gobots\samplebots\pathfinder.go:52: undefined: game.Action
..\..\bcspragu\Gobots\samplebots\pathfinder.go:57: undefined: game.Action
..\..\bcspragu\Gobots\samplebots\pathfinder.go:62: undefined: game.Action
..\..\bcspragu\Gobots\samplebots\pathfinder.go:62: too many errors
# github.com/bcspragu/Gobots/cloud
..\..\bcspragu\Gobots\cloud\cloud.go:55: undefined: cloudlaunch in cloudlaunch.Config

*/