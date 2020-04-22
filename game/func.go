package game

import (
	"encoding/json"

	"github.com/YWJSonic/ServerUtility/foundation"
	"github.com/YWJSonic/ServerUtility/foundation/fileload"
	"github.com/YWJSonic/ServerUtility/iserver"
	"github.com/YWJSonic/ServerUtility/restfult"
	"github.com/YWJSonic/ServerUtility/socket"
	"gitlab.fbk168.com/gamedevjp/cat/server/game/cache"
	"gitlab.fbk168.com/gamedevjp/cat/server/game/gamerule"
)

// NewGameServer ...
func NewGameServer(jsStr string) {

	config := foundation.StringToJSON(jsStr)
	baseSetting := iserver.NewSetting()
	baseSetting.SetData(config)

	gamejsStr := fileload.Load("./file/gameconfig.json")
	var gameRule = &gamerule.Rule{
		GameTypeID: config["GameTypeID"].(string),
	}
	if err := json.Unmarshal([]byte(gamejsStr), &gameRule); err != nil {
		panic(err)
	}

	cacheRedis := cache.Setting{
		URL: baseSetting.RedisURL,
	}

	var gameserver = iserver.NewService()
	var game = &Game{
		IGameRule: gameRule,
		Server:    gameserver,
		Cache:     cache.NewCache(cacheRedis),
	}
	gameserver.Restfult = restfult.NewRestfultService()
	gameserver.Socket = socket.NewSocket()
	gameserver.IGame = game

	// start Server
	gameserver.Launch(baseSetting)

	// start DB service
	setting := gameserver.Setting.DBSetting()
	gameserver.LaunchDB("gamedb", setting)
	gameserver.LaunchDB("logdb", setting)

	// start restful service
	go gameserver.LaunchRestfult(game.RESTfulURLs())
	// go gameserver.LaunchSocket(game.SocketURLs())

	<-gameserver.ShotDown
}
