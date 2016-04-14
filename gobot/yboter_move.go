package main

import "github.com/bcspragu/Gobots/game"
//import "fmt"

//start movement tatics chain
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

  return move_forward(bt,b,r)
}

//switch case on 4 direction when moving towards closest enermy
func move_to_target(bt *yboter, b *game.Board, r *game.Robot) game.Action {
	opp := nearestOpponent(b, r)
		if opp == nil {
			return game.Action{Kind: game.Wait}
	}

	direction_opp := direction_enermy(opp,r)
	loc :=game.Loc{}
	action := game.Action{}

	//order matters, going counter clockwise
	switch {	
	case direction_opp == game.West:
		action = move_to_direction(bt,b,r,direction_opp)
		if (action.Kind !=game.Wait){
			loc = r.Loc
			loc = loc.Add(direction_opp)
			if(count_friend_adj_loc(b, loc)<3 ){ // 
				return action
			}
		}

	case direction_opp == game.South:
		action = move_to_direction(bt,b,r,direction_opp)
		if (action.Kind !=game.Wait){
			return action
		}

	case direction_opp == game.East:
		action = move_to_direction(bt,b,r,direction_opp)
		if (action.Kind !=game.Wait){
			loc = r.Loc
			loc = loc.Add(direction_opp)
			if(count_friend_adj_loc(b, loc)<3 ){ // || game.Distance(r.Loc, loc)>2
				return action
			}
		}

	case direction_opp == game.North:
		action = move_to_direction(bt,b,r,direction_opp)
		if (action.Kind !=game.Wait){
			return action
		}
	}
	
	// move sideways if falls out from the case conditions
	return move_fanout(bt,b,r)
}

//fan out to the north and south when its too crowded
func move_fanout_lure(bt *yboter, b *game.Board, r *game.Robot) game.Action {
		opp := nearestOpponent(b, r)
			if opp == nil {
				return game.Action{Kind: game.Wait}
		}
		//if enermy is marching toward you and attack
		if game.Distance(r.Loc, opp.Loc) == 3 && count_friend_adj(b,opp) == 0 {
			if ( count_friend_oct(b,r)==1 || count_friend_oct(b,r)>3){
				return move_fanout(bt,b,r)
			}
		}

		return game.Action{Kind: game.Wait}
}

//fan out to the north and south
func move_fanout(bt *yboter, b *game.Board, r *game.Robot) game.Action {
		loc := game.Loc{}
		fut_pos := pos_stats{}
		loc = r.Loc
		loc = loc.Add(game.North)
		fut_pos = bt.robot_positions[loc]

		if (r.Loc.Y <b.Center().Y &&  fut_pos ==pos_stats{} && b.LocType(loc)==game.Valid){
			return game.Action{
				Kind:      game.Move,
				Direction: game.North,
			}
		}
		loc = r.Loc
		loc = loc.Add(game.South)
		fut_pos = bt.robot_positions[loc]

		if (r.Loc.Y >=b.Center().Y &&  fut_pos ==pos_stats{} && b.LocType(loc)==game.Valid){
			return game.Action{
				Kind:      game.Move,
				Direction: game.South,
			}
		}

  		return game.Action{Kind: game.Wait}
}

// go forward when possible
func move_forward(bt *yboter, b *game.Board, r *game.Robot) game.Action {
		direction_forward := direction_forward(b,r)
		loc := game.Loc{}
		fut_pos := pos_stats{}
		loc = r.Loc
		loc = loc.Add(direction_forward)
		fut_pos = bt.robot_positions[loc]
		if (fut_pos ==pos_stats{} && b.LocType(loc)==game.Valid){
			return game.Action{
				Kind:      game.Move,
				Direction: direction_forward,
			}
		}
  		return game.Action{Kind: game.Wait}
}
// move towards a direction
func move_to_direction(bt *yboter, b *game.Board, r *game.Robot, dir game.Direction) game.Action {
		loc := game.Loc{}
		fut_pos := pos_stats{}
		loc = r.Loc
		loc = loc.Add(dir)
		fut_pos = bt.robot_positions[loc]
		//fmt.Printf("stats:%v, clear:%v ", fut_pos,fut_pos ==pos_stats{})

		if (fut_pos ==pos_stats{} && b.LocType(loc)==game.Valid){
			return game.Action{
				Kind:      game.Move,
				Direction: dir,
			}
		}
  		return game.Action{Kind: game.Wait}
}
