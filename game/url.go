package game

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/code"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/foundation"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/httprouter"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/igame"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/messagehandle"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/socket"
	"gitlab.fbk168.com/gamedevjp/cat/server/game/constants"
	"gitlab.fbk168.com/gamedevjp/cat/server/game/db"
	"gitlab.fbk168.com/gamedevjp/cat/server/game/protocol"
)

func (g *Game) createNewSocket(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")

	if err := g.CheckToken(token); err != nil {
		log.Fatal("createNewSocket: not this token\n")
		return
	}

	c, err := g.Server.Socket.Upgrade(w, r, r.Header)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	g.Server.Socket.AddNewConn("f", c, func(msg socket.Message) error {
		fmt.Println("#-- socket --#", msg)
		return nil
	})
	// g.Server.Socket.AddNewConn(user.GetGameInfo().GameAccount, c, g.SocketMessageHandle)

	time.Sleep(time.Second * 3)
	g.Server.Socket.ConnMap["f"].Send(websocket.CloseMessage, []byte{})
}

// SocketMessageHandle ...
func (g *Game) SocketMessageHandle(msg socket.Message) error {
	fmt.Println("#-- socket --#", msg)
	return nil
}

func (g *Game) gameinit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var result = make(map[string]interface{})
	var proto protocol.InitRequest
	proto.InitData(r)

	// get user
	user, _, err := g.GetUser(proto.Token)
	if err != nil {
		err := messagehandle.New()
		err.ErrorCode = code.NoThisPlayer
		g.Server.HTTPResponse(w, "", err)
		return
	}

	result["player"] = map[string]interface{}{
		"gameaccount": g.IGameRule.GetGameTypeID(),
		"id":          user.UserGameInfo.IDStr,
		"money":       user.UserGameInfo.Money,
	}
	result["reel"] = g.IGameRule.GetReel()
	result["betrate"] = g.IGameRule.GetBetSetting()

	user.IAttach.Save()
	g.Server.HTTPResponse(w, result, messagehandle.New())
}

func (g *Game) gameresult(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var proto protocol.GameRequest
	proto.InitData(r)

	// get user
	// user, err := g.GetUserByGameID(proto.Token, proto.PlayerID)
	user, _, err := g.GetUser(proto.Token)
	if err != nil {
		err := messagehandle.New()
		err.Msg = "GameTypeError"
		err.ErrorCode = code.GameTypeError
		// messagehandle.ErrorLogPrintln("GetPlayerInfoByPlayerID-2", err, token, betIndex, betMoney)
		g.Server.HTTPResponse(w, "", err)
		return
	}
	if user.UserGameInfo.Money < g.IGameRule.GetBetMoney(proto.BetIndex) {
		err := messagehandle.New()
		err.Msg = "NoMoneyToBet"
		err.ErrorCode = code.NoMoneyToBet
		g.Server.HTTPResponse(w, "", err)
		return
	}

	order, errproto, err := g.NewOrder(proto.Token, user.UserGameInfo.IDStr, g.IGameRule.GetBetMoney(proto.BetIndex))
	if errproto != nil {
		err := messagehandle.New()
		err.Msg = errproto.String()
		err.ErrorCode = code.ULGInfoFormatError
		g.Server.HTTPResponse(w, "", err)
		return
	}
	if err != nil {
		err := messagehandle.New()
		err.Msg = "ULGInfoFormatError"
		err.ErrorCode = code.ULGInfoFormatError
		g.Server.HTTPResponse(w, "", err)
		return
	}

	fmt.Println(order)
	if err != nil {
		errMsg := messagehandle.New()
		errMsg.Msg = err.Error()
		errMsg.ErrorCode = code.NoMoneyToBet
		g.Server.HTTPResponse(w, "", errMsg)
		return
	}
	// get attach
	user.LoadAttach()

	oldMoney := user.UserGameInfo.Money
	// get game result
	RuleRequest := &igame.RuleRequest{
		BetIndex: proto.BetIndex,
		Attach:   &user.IAttach,
	}
	result := g.IGameRule.GameRequest(RuleRequest)
	user.UserGameInfo.SumMoney(result.Totalwinscore - result.BetMoney)

	for _, newAtt := range result.Attach {
		user.IAttach.SetAttach(newAtt)
	}
	user.IAttach.Save()
	resultMap := make(map[string]interface{})
	resultMap["totalwinscore"] = result.Totalwinscore
	resultMap["playermoney"] = user.UserGameInfo.GetMoney()
	resultMap["normalresult"] = result.GameResult["normalresult"]
	resultMap["attach"] = result.Attach

	respin, ok := result.OtherData["isrespin"]
	if ok && respin == 1 {
		resultMap["isrespin"] = 1
		resultMap["respin"] = result.GameResult["respin"]
	} else {
		resultMap["isrespin"] = 0
		resultMap["respin"] = []interface{}{}
	}

	msg := foundation.JSONToString(result.GameResult)
	msg = strings.ReplaceAll(msg, "\"", "\\\"")
	errMsg := db.SetLog(
		g.Server.DBConn("logdb"),
		user.UserGameInfo.IDStr,
		0,
		time.Now().Unix(),
		constants.ActionGameResult,
		oldMoney,
		user.UserGameInfo.Money,
		result.Totalwinscore,
		"",
		"",
		"",
		msg,
	)
	if errMsg.ErrorCode != code.OK {
		g.Server.HTTPResponse(w, resultMap, errMsg)
		return
	}

	g.EndOrder(proto.Token, order)
	g.Server.HTTPResponse(w, resultMap, messagehandle.New())
}
