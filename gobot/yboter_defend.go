package main

import "github.com/bcspragu/Gobots/game"

// start chain def tatics
func def_chain(bt *yboter,b *game.Board, r *game.Robot) game.Action {
	action :=  game.Action{Kind: game.Wait}
	action = def_lure(bt, b,r)
	if action.Kind != game.Wait{
		return action
	}
	action = def_guard(bt, b,r)
	if action.Kind != game.Wait{
		return action
	}

  return game.Action{Kind: game.Wait}
}

//guard when under 2+ attacks, check against prev health in stateboard
func def_guard(bt *yboter,b *game.Board, r *game.Robot) game.Action {
	nearby_count := count_enermies_adj(b,r)
	switch {
		case (bt.self_prevHP[r.ID] - r.Health > 15):
			if r.Health >= nearby_count*10 {
				return game.Action{Kind: game.Guard}
			}
	}
  return game.Action{Kind: game.Wait}
}

// flight or fight when overwhelmed by enermy
func def_lure(bt *yboter,b *game.Board, r *game.Robot) game.Action {
	//move back to lure enermy
	nearby_count := count_enermies_oct(b,r)
	nearby_friend_count := count_friend_oct(b,r)
	direction_backward :=  direction_back(b , r )
	direction_forward :=  direction_forward(b , r )

	enermyloc := game.Loc{}
	enermyloc = r.Loc
	enermyloc = enermyloc.Add(direction_forward)
	loc_type := b.LocType(enermyloc)
	if (loc_type ==game.Valid){
		opp_bot := b.At(enermyloc)

		loc := game.Loc{}
		loc = r.Loc
		loc = loc.Add(direction_backward)
		loc_type := b.LocType(loc)
		if (opp_bot !=nil){
			if (nearby_count >2 &&nearby_friend_count<=2  && loc_type== game.Valid){
					futurehealth := bt.robot_positions[opp_bot.Loc].future_health
					if (futurehealth>10 && !friendAt(b,loc)){
						return game.Action{
								Kind:      game.Move,
								Direction: direction_backward,
						}
				}
			}
		}
	}
  return game.Action{Kind: game.Wait}
}
