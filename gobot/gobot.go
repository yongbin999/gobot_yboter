package main

import "github.com/bcspragu/Gobots/game"

//import "github.com/yongbin999/gobot/yboter"

type yboter struct{
	targets map[uint32]uint32
	aggression uint32
}

func (bt *yboter) Act(b *game.Board, r *game.Robot) game.Action {

	//if enermy adjecent, attack
	act_attack(b,r)
	//else find targets
	find_target(bt,b,r)

	return game.Action{
	    Kind:      game.Move,
	    Direction: game.Towards(r.Loc, b.Center()),
	}
}

//if distance with closes enermy == 1 and 0 or have adjcent friend neighbor, attack
//if distance with closes enermy == 1 and 0 or fork adjcent friend neighbor, march

//if  self.hp - adjucent enemy count * damage <0 , suicide
//if surrounding boxes enermycount * avg damage > self.hp, suicide. 

//if theres wall go towards it. and for a slope 

//if triangle trap formed and head to head shift up or down as a team?

//if emermy not direct align, then leave a space, else can align one after another 



func main() {
    game.StartServerForFactory("yboter", "IlxVTAKDBsIisHXuGesBcBuds", game.ToFactory(&yboter{}))
}
