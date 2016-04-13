

### installation
go get -u github.com/bcspragu/Gobots/
go install -u github.com/bcspragu/Gobots/



### compile and run
cd Go/workspace/src/github.com/yongbin999/gobot
go install ./
go run *.go



## running daul bots to test
go install ./pathfinder/
go run ./pathfinder/*.go


## about this bot
all the main functions are split into files
* current turn state board
* defense module
* attack module
* movement module 