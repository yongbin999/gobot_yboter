
### installation
* go get -u github.com/bcspragu/Gobots/
* go install -u github.com/bcspragu/Gobots/

### compile and run
* cd Go/workspace/src/github.com/yongbin999/gobot
* go install ./
* go run *.go

### running dual terminal bots to test in matches
* go install ./pathfinder/
* go run ./pathfinder/*.go

go to the web client and create a match
http://gobotgame.com/bots

<hr>
### about this bot
all the main functions are split into files
* current turn state board
* defense module
* attack module
* movement module 