package main

import "github.com/bcspragu/Gobots/game"
import "fmt"

type yboter struct{
	targets map[uint32]uint32
	aggression uint32 						//unused
	self_prevHP map[uint32]int 				//decides  guard if under heavy attack
	robot_positions map[game.Loc]pos_stats	//sync all bot actions for each round avoid collisions
	current_turn int 						//reini this board every turn
}
type pos_stats struct{
	future_whats_here string
	future_health int
}


func (bt *yboter) Act(b *game.Board, r *game.Robot) game.Action {
	//init variables
	init_yboter_states(bt,b,r)
	//update oppoent
	update_targets(bt ,b , r)

	//print stats into console
	fmt.Printf("round:%2v bot:%2v loc:%3v H:%2v ", b.Round, r.ID, r.Loc,r.Health)
	
	//choose action by running the tatic chains
	return ai_action(bt,b,r)
}


func ai_action(bt *yboter, b *game.Board, r *game.Robot) game.Action {
	//defensive tatics
	action := game.Action{Kind: game.Wait}
	action = def_chain(bt,b,r)
	if action.Kind != game.Wait{
		print_action(&action)
		bt.self_prevHP[r.ID]  = r.Health
		return action
	}

	//offensive tatics
	action = off_chain(bt, b,r)
	if action.Kind != game.Wait{
		print_action(&action)
		update_future_health(bt, b, r, &action) 
		bt.self_prevHP[r.ID]  = r.Health
		return action
	}

	//movement tatics
	action = move_to_target(bt, b, r)
	if action.Kind != game.Wait{
		print_action(&action)
		update_future_position(bt, r, &action) 
		bt.self_prevHP[r.ID]  = r.Health
		return action
	}

	//save health to next one
	print_action(&action)
	bt.self_prevHP[r.ID]  = r.Health
  	return game.Action{Kind: game.Wait}
}


//---------------------------------------------------------------------
//actions:
//in other files

