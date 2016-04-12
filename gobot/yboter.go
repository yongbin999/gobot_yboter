package main

import "github.com/bcspragu/Gobots/game"

type yboter struct{
	targets map[uint32]uint32
	aggression uint32
	prevHP map[uint32]int
	enermy_HP map[uint32]int
}

func (bt *yboter) Act(b *game.Board, r *game.Robot) game.Action {

	nearby_count := count_enermies_adj(b,r)
	//save the last health
	switch {
		case ( bt.prevHP == nil):
				bt.prevHP = make(map[uint32]int)
		case (bt.prevHP[r.ID] == 0) :
				bt.prevHP[r.ID] = r.Health
				//bt.boardside = r.Faction
		case (bt.prevHP[r.ID] - r.Health > 15):
			if r.Health >= nearby_count*10 {
				return game.Action{Kind: game.Guard}
			}
	}

	bt.prevHP[r.ID]  = r.Health

	//update oppoent
	update_opp(bt ,b , r)

	//if nearby * avg dmage > your health destory
	action := game.Action{Kind: game.Wait}
	action = off_chain(b,r)
	if action.Kind != game.Wait{
		return action
	}

	//move to target
	action = move_to_target(b, r)
	if action.Kind != game.Wait{
		return action
	}

  	return game.Action{Kind: game.Wait}
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

func friendAt(b *game.Board, loc game.Loc) bool {
	if !b.IsInside(loc) {
		return false
	}
	r := b.At(loc)
	if r == nil {
		return false
	}
	return r.Faction == game.MyFaction
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
func count_enermies_adj(b *game.Board, r *game.Robot) int {
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
	}
	return counter
}

func update_opp(bt *yboter,b *game.Board, r *game.Robot) {
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
    	if opp != nil {
      		bt.targets[r.ID] = opp.ID
    	}
	}
}
func direction_back(b *game.Board, r *game.Robot) game.Direction{
	if r.Faction ==1{
		return  game.West
	}
	if r.Faction ==2{
		return  game.East
	}
	return game.West
}


//---------------------------------------------------------------------
//actions:



func off_chain(b *game.Board, r *game.Robot) game.Action {
	action := game.Action{Kind: game.Wait}
	action = off_selfdestruct(b,r)
	if action.Kind != game.Wait{
		return action
	}
	//lure when right condition
	action = off_lure(b,r)
	if action.Kind != game.Wait{
		return action
	}
	//apply attack
	action = off_attack(b,r)
	if action.Kind != game.Wait{
		return action
	}
  return game.Action{Kind: game.Wait}
}

func off_selfdestruct(b *game.Board, r *game.Robot) game.Action {
	nearby_count := count_enermies_oct(b,r)
	//nearby_count_oct := count_enermies_otc(b,r)
	//if nearby * avg dmage > your health destory
	if (nearby_count*5 > r.Health ){
		return game.Action{
	    	Kind: game.SelfDestruct,
		}
	}

  return game.Action{Kind: game.Wait}
}

func off_lure(b *game.Board, r *game.Robot) game.Action {
	//move back to lure enermy
	nearby_count := count_enermies_oct(b,r)
	direction :=  direction_back(b , r )
	location := r.Loc.Add(direction)
	loc_type := b.LocType(location)
	if (nearby_count >=3 && loc_type== game.Valid){
		bot_atloc := b.At(location)
		if (bot_atloc == nil && !friendAt(b,location)){
				return game.Action{
						Kind:      game.Move,
						Direction: direction,
				}
			}
		}
  return game.Action{Kind: game.Wait}
}

func off_attack(b *game.Board, r *game.Robot) game.Action {
	ds := []game.Direction{
		game.North,
		game.South,
		game.East,
		game.West,
	}
	//collective attack not exceed their base not implemented yet
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



func move_to_target(b *game.Board, r *game.Robot) game.Action {
	opp := nearestOpponent(b, r.Loc)
		if opp == nil {
			return game.Action{Kind: game.Wait}
	}
	// Move to target.
	// Don't worry about collisions, since we already shot at all neighbors.
	// TODO: but what about friends?
	// TODO: and why don't you compute the vector angle?
	switch {
	case game.Distance(r.Loc, opp.Loc) == 1:
		return game.Action{Kind: game.Wait}

	case opp.Loc.X < r.Loc.X && !friendAt(b,r.Loc.Add(game.West)):
		return game.Action{
			Kind:      game.Move,
			Direction: game.West,
		}
	case opp.Loc.X > r.Loc.X && !friendAt(b,r.Loc.Add(game.East)):
		return game.Action{
			Kind:      game.Move,
			Direction: game.East,
		}
	case opp.Loc.Y < r.Loc.Y && !friendAt(b,r.Loc.Add(game.North)):
		return game.Action{
			Kind:      game.Move,
			Direction: game.North,
		}
	case opp.Loc.Y > r.Loc.Y && !friendAt(b,r.Loc.Add(game.South)):
		return game.Action{
			Kind:      game.Move,
			Direction: game.South,
		}
	}
	// TODO: impossibru?
	//return game.Action{Kind: game.Wait}
  return game.Action{Kind: game.Wait}
}