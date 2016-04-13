package main

import "github.com/bcspragu/Gobots/game"


func move_chain(bt *yboter,b *game.Board, r *game.Robot) game.Action {
	action :=  game.Action{Kind: game.Wait}
	action = move_fanout_lure(bt, b,r)
	if action.Kind != game.Wait{
		return action
	}
	action = move_to_target(bt, b,r)
	if action.Kind != game.Wait{
		return action
	}

  return game.Action{Kind: game.Wait}
}

func move_to_target(bt *yboter, b *game.Board, r *game.Robot) game.Action {
	opp := nearestOpponent(b, r)
		if opp == nil {
			return game.Action{Kind: game.Wait}
	}

	direction_opp := direction_enermy(opp,r)
	loc :=game.Loc{}
	fut_loc := game.Loc{}
	fut_loc = r.Loc
	fut_loc = fut_loc.Add(direction_opp)
	fut_pos := bt.robot_positions[fut_loc]

	switch {	
	case direction_opp == game.West:
		if (fut_pos !=pos_stats{}){
			loc = r.Loc
			loc = loc.Add(game.West)
			if(count_friend_adj_loc(b, loc)<3){
				return game.Action{
					Kind:      game.Move,
					Direction: game.West,
				}
			}
		}
		//check if move sideways possible
		return move_fanout(bt,b,r)


	case direction_opp == game.East:
		if (fut_pos !=pos_stats{}){
			loc = r.Loc
			loc = loc.Add(game.West)
			if(count_friend_adj_loc(b, loc)<3){
				return game.Action{
					Kind:      game.Move,
					Direction: game.East,
				}
			}
		}
		//check if move sideways possible
		return move_fanout(bt,b,r)


	case direction_opp == game.North:
		loc := game.Loc{}
		loc = r.Loc
		loc = loc.Add(direction_opp)
		fut_pos = bt.robot_positions[loc]
		if (fut_pos !=pos_stats{}){
			return game.Action{
				Kind:      game.Move,
				Direction: game.North,
			}
		}
		//if sideway not possible try forward
		return move_forward(bt,b,r)


	case direction_opp == game.South:
		loc := game.Loc{}
		loc = r.Loc
		loc = loc.Add(direction_opp)
		fut_pos = bt.robot_positions[loc]
		if (fut_pos != pos_stats{} ){
			return game.Action{
				Kind:      game.Move,
				Direction: game.South,
			}
		}
		//if sideway not possible try forward
		return move_forward(bt,b,r)
	}

  return game.Action{Kind: game.Wait}
}

func move_fanout_lure(bt *yboter, b *game.Board, r *game.Robot) game.Action {
		opp := nearestOpponent(b, r)
			if opp == nil {
				return game.Action{Kind: game.Wait}
		}
		//if enermy is marching toward you and attack
		if game.Distance(r.Loc, opp.Loc) == 2 && count_friend_adj(b,opp) == 0{
			if ( count_friend_oct(b,r)==1 || count_friend_oct(b,r)>=3){
				return move_fanout(bt,b,r)
			}
		}

		return game.Action{Kind: game.Wait}
}

func move_fanout(bt *yboter, b *game.Board, r *game.Robot) game.Action {
		loc := game.Loc{}
		fut_pos := pos_stats{}
		loc = r.Loc
		loc = loc.Add(game.North)
		fut_pos = bt.robot_positions[loc]
		if (r.Loc.Y <=b.Center().Y &&  fut_pos !=pos_stats{}){
			return game.Action{
				Kind:      game.Move,
				Direction: game.North,
			}
		}
		loc = r.Loc
		loc = loc.Add(game.South)
		fut_pos = bt.robot_positions[loc]
		if (r.Loc.Y >b.Center().Y &&  fut_pos !=pos_stats{}){
			return game.Action{
				Kind:      game.Move,
				Direction: game.South,
			}
		}

  		return game.Action{Kind: game.Wait}
}

func move_forward(bt *yboter, b *game.Board, r *game.Robot) game.Action {
		direction_forward := direction_forward(b,r)
		loc := game.Loc{}
		fut_pos := pos_stats{}
		loc = r.Loc
		loc = loc.Add(direction_forward)
		fut_pos = bt.robot_positions[loc]
		if (fut_pos !=pos_stats{} ){
			return game.Action{
				Kind:      game.Move,
				Direction: direction_forward,
			}
		}
  		return game.Action{Kind: game.Wait}
}
