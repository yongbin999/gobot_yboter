package main

import "github.com/bcspragu/Gobots/game"
import "fmt"

//---------------------------------------------------------------------------
//helpers: print infos

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
