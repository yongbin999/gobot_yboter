package main

import "github.com/bcspragu/Gobots/game"
import "fmt"

//---------------------------------------------------------------------------
//helpers:

func init_yboter_states(bt *yboter,b *game.Board,r *game.Robot) {
	switch {
		case ( bt.self_prevHP == nil):
				bt.self_prevHP = make(map[uint32]int)
		case (bt.self_prevHP[r.ID] == 0) :
				bt.self_prevHP[r.ID] = r.Health
	}
	switch {
		case ( bt.robot_positions == nil):
				init_futureboard(bt, b)
				//fmt.Printf("stats : %v\n", bt.robot_positions)
				bt.current_turn = b.Round
		case ( bt.current_turn != b.Round):
				//copy everything from board to here
				init_futureboard(bt, b)
				fmt.Printf("botsmap:%v \n",bt.robot_positions)
				bt.current_turn = b.Round
	}
}

func update_targets(bt *yboter,b *game.Board, r *game.Robot) {

	//need also update thier health
	if bt.targets == nil {
		bt.targets = make(map[uint32]uint32)
	}
	opp := nearestOpponent(b, r)
    if opp != nil {
      	bt.targets[r.ID] = opp.ID
	}
}

func init_futureboard(bt *yboter,b *game.Board) {
	bt.robot_positions = make(map[game.Loc]pos_stats)
	friendbots :=b.Bots(game.MyFaction)
	enermybots :=b.Bots(game.OpponentFaction)

	for _, bot := range friendbots{
		posstats := pos_stats{}
		posstats.future_whats_here = "friend"
		posstats.future_health = bot.Health
		bt.robot_positions[bot.Loc] = posstats
	}
	for _, bot := range enermybots{
		posstats := pos_stats{}
		posstats.future_whats_here = "enermy"
		posstats.future_health = bot.Health
		bt.robot_positions[bot.Loc] = posstats
	}
}
func update_future_position(bt *yboter, r *game.Robot, action *game.Action) {
	loc := game.Loc{}
	loc = r.Loc
	loc = loc.Add(action.Direction)
	bt.robot_positions[loc] = bt.robot_positions[r.Loc]
	bt.robot_positions[r.Loc] = pos_stats{}
}
func update_future_health(bt *yboter, b *game.Board ,r *game.Robot, action *game.Action) {
	switch{
		case action.Kind ==game.Attack: {	
			loc := game.Loc{}
			loc = r.Loc
			loc = loc.Add(action.Direction)
			currentstats := bt.robot_positions[loc]
			currentstats.future_health += -10
			bt.robot_positions[loc] = currentstats
		}	
		case action.Kind ==game.SelfDestruct: {
			locations := b.LocsAround(r.Loc)
			for _, loc := range locations{
				currentstats := bt.robot_positions[loc]
				currentstats.future_health += -15
				bt.robot_positions[loc] = currentstats
			}
		}
	}

}
