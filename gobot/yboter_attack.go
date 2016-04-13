package main

import "github.com/bcspragu/Gobots/game"


func off_chain(bt *yboter, b *game.Board, r *game.Robot) game.Action {
	action := game.Action{Kind: game.Wait}
	action = off_selfdestruct(b,r)
	if action.Kind != game.Wait{
		return action
	}
	//lure when right condition
	
	//apply attack
	action = off_attack(bt, b,r)
	if action.Kind != game.Wait{
		return action
	}

	action = off_preattack(b,r)
	if action.Kind != game.Wait{
		return action
	}

  return game.Action{Kind: game.Wait}
}

func off_selfdestruct(b *game.Board, r *game.Robot) game.Action {

	//should  find death toll before do it, hp <15
	nearby_count := count_enermies_oct_weak(b,r)
	friend_nearby_count := count_friend_oct(b,r)
	//nearby_count_oct := count_enermies_otc(b,r)
	//if nearby * avg dmage > your health destory
	if (nearby_count > friend_nearby_count && nearby_count*10 > r.Health  && b.Round <92){
		return game.Action{
	    	Kind: game.SelfDestruct,
		}
	}

  return game.Action{Kind: game.Wait}
}

func off_attack(bt *yboter, b *game.Board, r *game.Robot) game.Action {
	ds := []game.Direction{
		game.North,
		game.East,
		game.South,
		game.West,
	}
	//collective attack not exceed their base not implemented yet

	//attack the lowest instead of follow direction
	for _, d := range ds {
		loc := game.Loc{}
		loc = r.Loc
		loc = loc.Add(d)
		if opponentAt(b, loc) {
			opp_bot := b.At(loc)
			futurehealth := bt.robot_positions[opp_bot.Loc].future_health
			if futurehealth >0{
				return game.Action{
					Kind:      game.Attack,
					Direction: d,
				}
			}
		}
	}
  return game.Action{Kind: game.Wait}
}

func off_preattack(b *game.Board, r *game.Robot) game.Action {
	opp := nearestOpponent(b, r)
		if opp == nil {
			return game.Action{Kind: game.Wait}
	}

	direction_opp := direction_enermy(opp,r)
		//if enermy is marching toward you and attack
	if game.Distance(r.Loc, opp.Loc) == 1 && count_friend_adj(b,opp) == 0{
		return game.Action{
				Kind:      game.Attack,
				Direction: direction_opp,
			}
		}

	return game.Action{Kind: game.Wait}
}