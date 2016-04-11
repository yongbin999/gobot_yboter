package main

import "github.com/bcspragu/Gobots/game"

// Bot moves to the center and does nothing else
type bot struct{}

func (bot) Act(b *game.Board, r *game.Robot) game.Action {
  return game.Action{
    Kind:      game.Move,
    Direction: game.Towards(r.Loc, b.Center()),
  }
}

func main() {
    game.StartServerForFactory("yboter", "IlxVTAKDBsIisHXuGesBcBuds", game.ToFactory(bot{}))
}

/*
# github.com/bcspragu/Gobots/samplebots
..\..\bcspragu\Gobots\samplebots\aggro.go:7: undefined: game.Turn
..\..\bcspragu\Gobots\samplebots\aggro.go:17: undefined: game.Turn
..\..\bcspragu\Gobots\samplebots\aggro.go:23: undefined: game.Turn
..\..\bcspragu\Gobots\samplebots\main.go:57: not enough arguments in call to c.RegisterAI
..\..\bcspragu\Gobots\samplebots\pathfinder.go:9: undefined: game.Turn
..\..\bcspragu\Gobots\samplebots\pathfinder.go:20: undefined: game.Turn
..\..\bcspragu\Gobots\samplebots\pathfinder.go:41: undefined: game.Turn
..\..\bcspragu\Gobots\samplebots\pathfinder.go:52: undefined: game.Turn
..\..\bcspragu\Gobots\samplebots\pathfinder.go:57: undefined: game.Turn
..\..\bcspragu\Gobots\samplebots\pathfinder.go:62: undefined: game.Turn
..\..\bcspragu\Gobots\samplebots\pathfinder.go:62: too many errors
# github.com/bcspragu/Gobots/cloud
..\..\bcspragu\Gobots\cloud\cloud.go:55: undefined: cloudlaunch in cloudlaunch.Config

*/