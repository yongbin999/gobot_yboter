package main

import "github.com/bcspragu/Gobots/game"

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

func count_enermies_oct(b *game.Board, r *game.Robot) int {
	counter :=0
	ds := []game.Direction{
		game.North,
		game.South,
		game.East,
		game.West,
	}
	for _, d := range ds {
		loc := r.Loc.Add(d)
		if opponentAt(b, loc) {
			counter +=1
		}

		if (d == game.North || d == game.South ){
			loc = r.Loc.Add(game.East)
			if opponentAt(b, loc) {
				counter +=1
			}
			loc = r.Loc.Add(game.West)
			if opponentAt(b, loc) {
				counter +=1
			}
		}
	}
	return counter
}

func act_attack(b *game.Board, r *game.Robot) game.Action {
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

	return game.Action{Kind: game.Wait}
}


func find_target(bt *yboter,b *game.Board, r *game.Robot) game.Action {
	tgt, ok := bt.targets[r.ID]
	var opp *game.Robot
	if ok {
		opp = b.Find(func(q *game.Robot) bool {
			return q.ID == tgt
		})
	}
	if !ok || opp == nil {
		if bt.targets == nil {
			bt.targets = make(map[uint32]uint32)
		}
		opp = nearestOpponent(b, r.Loc)
		if opp == nil {
			return game.Action{Kind: game.Wait}
		}
		bt.targets[r.ID] = opp.ID
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