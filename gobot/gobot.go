package main

import "github.com/bcspragu/Gobots/game"

//import "github.com/yongbin999/gobot/yboter"

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


//disconnect, still remain online on server
/*
Lost connection to server, trying to reconnect...
Failed to connect to server: dial tcp: lookup gobotgame.com: getaddrinfow: No such host is known.
exit status 1
*/

//switch statements still require return?