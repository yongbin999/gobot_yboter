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
	switch {
		case ( bt.self_prevHP == nil):
				bt.self_prevHP = make(map[uint32]int)
		case (bt.self_prevHP[r.ID] == 0) :
				bt.self_prevHP[r.ID] = r.Health
	}
	switch {
		case ( bt.robot_positions == nil):
				bt.robot_positions = make(map[game.Loc]pos_stats)
	}
	switch {
		case ( bt.current_turn != b.Round):
				//copy everything from board to here
				bt.robot_positions =init_futureboard(b)
				bt.current_turn = b.Round
	}

	//print stats
	fmt.Printf("round:%2v bot:%2v loc:%3v H:%2v ", b.Round, r.ID, r.Loc,r.Health)
	

	//update oppoent
	update_targets(bt ,b , r)

	//defensive tatics
	action := game.Action{Kind: game.Wait}
	action = def_chain(bt,b,r)
	if action.Kind != game.Wait{
		print_action(&action)
		bt.self_prevHP[r.ID]  = r.Health
		return action
	}

	//offensive tatics
	action = off_chain(b,r)
	if action.Kind != game.Wait{
		print_action(&action)
		update_future_health(bt, r, &action) 
		bt.self_prevHP[r.ID]  = r.Health
		return action
	}

	//movement tatics
	action = move_to_target(b, r)
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


//---------------------------------------------------------------------------
//helpers:

/*
	0 = wait
	1 = Collision
	2 = Attack
	3 = Destruct
	4 = Self
	Direction_north Direction = 0
    Direction_south Direction = 1
    Direction_east  Direction = 2
    Direction_west  Direction = 3
*/
func print_action(action *game.Action) {
	switch{
	case action.Kind ==0:
		fmt.Printf("wait: %v\n", print_direction(action))
	case action.Kind ==1:
		fmt.Printf("move: %v\n", print_direction(action))
	case action.Kind ==2:
		fmt.Printf("attack: %v\n",print_direction(action))
	case action.Kind ==3:
		fmt.Printf("Destruct: %v\n",print_direction(action))
	case action.Kind ==4:
		fmt.Printf("defend: %v\n",print_direction(action)) /// ??
	default:
		fmt.Printf("unknow action: %v\n", action)
	}
}
func print_direction(action *game.Action) string {
	switch{
	case action.Direction ==0:
		return "north"
	case action.Direction ==1:
		return "south"
	case action.Direction ==2:
		return "east"
	case action.Direction ==3:
		return "west"
	default:
		return "unknown"
	}
}
func nearestOpponent(b *game.Board, r *game.Robot) *game.Robot {
	// Probably faster ways of doing this.. traversing outward
	bots := b.Bots(game.OpponentFaction)

	var closest *game.Robot
	var closestDist int
	for _, bot := range bots {
			d := game.Distance(r.Loc, bot.Loc)
			if closest == nil || d < closestDist {
				closest, closestDist = bot, d
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

		loc := game.Loc{}
		loc = r.Loc
		loc = loc.Add(d)
		if opponentAt(b, loc)  {
			counter +=1
		}

		if (d == game.North || d == game.South ){
			loc = loc.Add(game.East)
			//loc = r.Loc.Add(game.East)
			if opponentAt(b, loc) {
				counter +=1
			}
			loc = loc.Add(game.West)
			//loc = r.Loc.Add(game.West)
			if opponentAt(b, loc) {
				counter +=1
			}
		}
	}
	return counter
}
func count_enermies_oct_weak(b *game.Board, r *game.Robot) int {
	counter :=0
	ds := []game.Direction{
		game.North,
		game.South,
		game.East,
		game.West,
	}
	for _, d := range ds {

		loc := game.Loc{}
		loc = r.Loc
		loc = loc.Add(d)
		if opponentAt(b, loc) && b.At(loc).Health <=15 {
			counter +=1
		}

		if (d == game.North || d == game.South ){
			loc = loc.Add(game.East)
			//loc = r.Loc.Add(game.East)
			if opponentAt(b, loc) && b.At(loc).Health <=15  {
				counter +=1
			}
			loc = loc.Add(game.West)
			//loc = r.Loc.Add(game.West)
			if opponentAt(b, loc) && b.At(loc).Health <=15  {
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
		loc := game.Loc{}
		loc = r.Loc
		loc = loc.Add(d)
		if opponentAt(b, loc) {
			counter +=1
		}
	}
	return counter
}
func count_friend_oct(b *game.Board, r *game.Robot) int {
	counter :=0
	ds := []game.Direction{
		game.North,
		game.South,
		game.East,
		game.West,
	}
	for _, d := range ds {

		loc := game.Loc{}
		loc = r.Loc
		loc = loc.Add(d)
		if friendAt(b, loc) {
			counter +=1
		}

		if (d == game.North || d == game.South ){
			loc = loc.Add(game.East)
			//loc = r.Loc.Add(game.East)
			if friendAt(b, loc) {
				counter +=1
			}
			loc = loc.Add(game.West)
			//loc = r.Loc.Add(game.West)
			if friendAt(b, loc) {
				counter +=1
			}
		}
	}
	return counter
}
func count_friend_adj(b *game.Board, r *game.Robot) int {
	counter :=0
	ds := []game.Direction{
		game.North,
		game.South,
		game.East,
		game.West,
	}
	for _, d := range ds {
		loc := game.Loc{}
		loc = r.Loc
		loc = loc.Add(d)
		if friendAt(b, loc) {
			counter +=1
		}
	}
	return counter
}
func count_friend_adj_loc(b *game.Board, loca game.Loc) int {
	counter :=0
	ds := []game.Direction{
		game.North,
		game.South,
		game.East,
		game.West,
	}
	for _, d := range ds {
		loc := game.Loc{}
		loc = loca
		loc = loc.Add(d)
		if friendAt(b, loc) {
			counter +=1
		}
	}
	return counter
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
func direction_back(b *game.Board, r *game.Robot) game.Direction{
	if r.Faction ==1{
		return  game.West
	}
	if r.Faction ==2{
		return  game.East
	}
	return game.West
}
func direction_forward(b *game.Board, r *game.Robot) game.Direction{
	if r.Faction ==2{
		return  game.West
	}
	if r.Faction ==1{
		return  game.East
	}
	return game.West
}
func direction_enermy(opp *game.Robot, r *game.Robot) game.Direction{
	//clockwise direction, less conflict when back and forth
	switch {
	case opp.Loc.X < r.Loc.X :
		return game.West
	case opp.Loc.Y < r.Loc.Y :
		return game.North
	case opp.Loc.X > r.Loc.X :
		return game.East
	case opp.Loc.Y > r.Loc.Y :
		return game.South
	// should not reach default
	default: 
		return game.North
	}
}

func init_futureboard(b *game.Board) map[game.Loc]pos_stats{
	cells := make(map[game.Loc]pos_stats)
	friendbots :=b.Bots(game.MyFaction)
	enermybots :=b.Bots(game.OpponentFaction)

	for _, bot := range friendbots{
		posstats := pos_stats{}
		posstats.future_whats_here = "friend"
		posstats.future_health = bot.Health
		cells[bot.Loc] = posstats
	}
	for _, bot := range enermybots{
		posstats := pos_stats{}
		posstats.future_whats_here = "enermy"
		posstats.future_health = bot.Health
		cells[bot.Loc] = posstats
	}
	return cells
}
func update_future_position(bt *yboter, r *game.Robot, action *game.Action) {
	loc := game.Loc{}
	loc = r.Loc
	loc = loc.Add(action.Direction)
	bt.robot_positions[loc] = bt.robot_positions[r.Loc]
	bt.robot_positions[r.Loc] = pos_stats{}
}
func update_future_health(bt *yboter, r *game.Robot, action *game.Action) {
	loc := game.Loc{}
	loc = r.Loc
	loc = loc.Add(action.Direction)
	bt.robot_positions[loc] = bt.robot_positions[r.Loc]
	bt.robot_positions[r.Loc] = pos_stats{}
}



//---------------------------------------------------------------------
//actions:

func def_chain(bt *yboter,b *game.Board, r *game.Robot) game.Action {
	action :=  game.Action{Kind: game.Wait}
	action = def_lure(b,r)
	if action.Kind != game.Wait{
		return action
	}
	action = def_guard(bt, b,r)
	if action.Kind != game.Wait{
		return action
	}

  return game.Action{Kind: game.Wait}
}

func def_guard(bt *yboter,b *game.Board, r *game.Robot) game.Action {
	nearby_count := count_enermies_adj(b,r)
	//save the last health
	switch {
		case (bt.self_prevHP[r.ID] - r.Health > 15):
			if r.Health >= nearby_count*10 {
				return game.Action{Kind: game.Guard}
			}
	}
  return game.Action{Kind: game.Wait}
}

func def_lure(b *game.Board, r *game.Robot) game.Action {
	//move back to lure enermy
	if(r.Health>25){
	nearby_count := count_enermies_oct(b,r)
	nearby_friend_count := count_friend_oct(b,r)
	direction :=  direction_back(b , r )
	loc := game.Loc{}
	loc = r.Loc
	loc = loc.Add(direction)
	loc_type := b.LocType(loc)
	if (nearby_count >=2 &&nearby_friend_count<=2  && loc_type== game.Valid){
		bot_atloc := b.At(loc)
		if (bot_atloc == nil && !friendAt(b,loc)){
				return game.Action{
						Kind:      game.Move,
						Direction: direction,
				}
			}
		}
	}
  return game.Action{Kind: game.Wait}
}

func off_chain(b *game.Board, r *game.Robot) game.Action {
	action := game.Action{Kind: game.Wait}
	action = off_selfdestruct(b,r)
	if action.Kind != game.Wait{
		return action
	}
	//lure when right condition
	
	//apply attack
	action = off_attack(b,r)
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

func off_attack(b *game.Board, r *game.Robot) game.Action {
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
			return game.Action{
				Kind:      game.Attack,
				Direction: d,
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



func move_to_target(b *game.Board, r *game.Robot) game.Action {
	opp := nearestOpponent(b, r)
		if opp == nil {
			return game.Action{Kind: game.Wait}
	}

	direction_opp := direction_enermy(opp,r)
	direction_forward := direction_forward(b,r)

	switch {
		
	// move forward when possible 
	case direction_opp == game.West:
		loc := game.Loc{}
		loc = r.Loc
		loc = loc.Add(game.West)
		if (!friendAt(b, loc) ){
			if(count_friend_adj_loc(b, loc)<3){
				return game.Action{
					Kind:      game.Move,
					Direction: game.West,
				}
			}
			//check if move sideways possible
			loc = r.Loc
			loc = loc.Add(game.North)
			if (r.Loc.Y <=b.Center().Y && !friendAt(b, loc)){
				return game.Action{
					Kind:      game.Move,
					Direction: game.North,
				}
			}
			loc = r.Loc
			loc = loc.Add(game.South)
			if (r.Loc.Y >b.Center().Y && !friendAt(b, loc)){
				return game.Action{
					Kind:      game.Move,
					Direction: game.South,
				}
			}
		}

		loc = r.Loc
		loc = loc.Add(game.North)
		if (r.Loc.Y <=b.Center().Y && !friendAt(b, loc)){
			return game.Action{
				Kind:      game.Move,
				Direction: game.North,
			}
		}
		loc = r.Loc
		loc = loc.Add(game.South)
		if (r.Loc.Y >b.Center().Y && !friendAt(b, loc)){
			return game.Action{
				Kind:      game.Move,
				Direction: game.South,
			}
		}


	case direction_opp == game.East:
		loc := game.Loc{}
		loc = r.Loc
		loc = loc.Add(game.East)
		if (!friendAt(b, loc) ){
			if(count_friend_adj_loc(b, loc)<3){
				return game.Action{
					Kind:      game.Move,
					Direction: game.East,
				}
			}
			//check if move sideways possible
			loc = r.Loc
			loc = loc.Add(game.North)
			if (r.Loc.Y <=b.Center().Y && !friendAt(b, loc)){
				return game.Action{
					Kind:      game.Move,
					Direction: game.North,
				}
			}
			loc = r.Loc
			loc = loc.Add(game.South)
			if (r.Loc.Y >b.Center().Y && !friendAt(b, loc)){
				return game.Action{
					Kind:      game.Move,
					Direction: game.South,
				}
			}
		}

		loc = r.Loc
		loc = loc.Add(game.North)
		if (r.Loc.Y <=b.Center().Y && !friendAt(b, loc)){
			return game.Action{
				Kind:      game.Move,
				Direction: game.North,
			}
		}
		loc = r.Loc
		loc = loc.Add(game.South)
		if (r.Loc.Y >b.Center().Y && !friendAt(b, loc)){
			return game.Action{
				Kind:      game.Move,
				Direction: game.South,
			}
		}


	case direction_opp == game.North:
		loc := game.Loc{}
		loc = r.Loc
		loc = loc.Add(direction_opp)
		if (!friendAt(b, loc)){
			return game.Action{
				Kind:      game.Move,
				Direction: game.North,
			}
		}
		loc = r.Loc
		loc = loc.Add(direction_forward)
		if (!friendAt(b, loc)){
			return game.Action{
				Kind:      game.Move,
				Direction: direction_forward,
			}
		}


	case direction_opp == game.South:
		loc := game.Loc{}
		loc = r.Loc
		loc = loc.Add(direction_opp)
		if (!friendAt(b, loc)){
			return game.Action{
				Kind:      game.Move,
				Direction: game.South,
			}
		}
		loc = r.Loc
		loc = loc.Add(direction_forward)
		if (!friendAt(b, loc)){
			return game.Action{
				Kind:      game.Move,
				Direction: direction_forward,
			}
		}
	}

  return game.Action{Kind: game.Wait}
}