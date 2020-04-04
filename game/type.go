package game

import (
	"errors"

	"gitlab.fbk168.com/gamedevjp/cat/server/game/cache"
	"gitlab.fbk168.com/gamedevjp/cat/server/game/catattach"
	"github.com/golang/protobuf/ptypes"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/igame"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/iserver"
	_ "gitlab.fbk168.com/gamedevjp/backend-utility/utility/mysql"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/playerinfo"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/restfult"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/socket"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/thirdparty/transaction/protoc"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/user"
)

// Game ...
type Game struct {
	Server    *iserver.Service
	IGameRule igame.ISlotRule
	Cache     *cache.GameCache
	// ProtocolMap map[string]func(r *http.Request) protocol.IProtocol
}

// RESTfulURLs ...
func (g *Game) RESTfulURLs() []restfult.Setting {
	return []restfult.Setting{
		restfult.Setting{
			RequestType: "POST",
			URL:         "game/init",
			Fun:         g.gameinit,
			ConnType:    restfult.Client,
		},
		restfult.Setting{
			RequestType: "POST",
			URL:         "game/result",
			Fun:         g.gameresult,
			ConnType:    restfult.Client,
		},
	}
}

// SocketURLs ...
func (g *Game) SocketURLs() []socket.Setting {
	return []socket.Setting{
		socket.Setting{
			URL: "lobby/createNewSocket",
			Fun: g.createNewSocket,
		},
	}
}

// NewUser *Not Use
func (g *Game) NewUser(token, gameAccount string) *user.Info {
	return &user.Info{}
}

// GetUser ...
func (g *Game) GetUser(userToken string) (*user.Info, *protoc.Error, error) {
	if g.Server.Setting.ServerMod == "dev" {
		return &user.Info{
			UserServerInfo: &playerinfo.AccountInfo{},
			UserGameInfo: &playerinfo.Info{
				IDStr:  "devtest",
				Money:  10000000,
				MoneyU: 10000000,
			},
			IAttach: catattach.NewUserAttach(g.Cache, 0),
		}, nil, nil
	}

	userProto, errorProto, err := g.Server.Transfer.AuthUser(userToken)
	if err != nil {
		if errorProto != nil {
			return nil, errorProto, err
		}
		return nil, nil, err
	}

	return &user.Info{
		UserServerInfo: &playerinfo.AccountInfo{},
		UserGameInfo: &playerinfo.Info{
			IDStr:  userProto.GetUserId(),
			Money:  int64(userProto.GetBalance()),
			MoneyU: userProto.GetBalance(),
		},
		IAttach: catattach.NewUserAttach(g.Cache, 0),
	}, nil, nil
}

// NewOrder ...
func (g *Game) NewOrder(token, userIDStr string, betMoney int64) (*protoc.Order, *protoc.Error, error) {
	if g.Server.Setting.ServerMod == "dev" {
		return &protoc.Order{
			UserId:  userIDStr,
			GameId:  g.IGameRule.GetGameTypeID(),
			Bet:     uint64(betMoney),
			OrderId: "testOrder",
		}, nil, nil
	}
	orderProto, errorProto, err := g.Server.Transfer.NewOrder(token, &protoc.Order{
		UserId: userIDStr,
		GameId: g.IGameRule.GetGameTypeID(),
		Bet:    uint64(betMoney),
	})
	if err != nil {
		if errorProto != nil {
			return nil, errorProto, err
		}
		return nil, nil, err
	}
	return orderProto, nil, nil

}

// EndOrder ...
func (g *Game) EndOrder(token string, orderProto *protoc.Order) (*protoc.Order, *protoc.Error, error) {
	orderProto.CompletedAt = ptypes.TimestampNow()
	if g.Server.Setting.ServerMod == "dev" {
		return orderProto, nil, nil
	}
	return g.Server.Transfer.EndOrder(token, orderProto)
}

// GetUserByGameID ...
func (g *Game) GetUserByGameID(token string, userID int64) (*user.Info, error) {
	return &user.Info{}, nil
}

// CheckGameType *Not Use
func (g *Game) CheckGameType(clientGameTypeID string) bool {
	return true
}

// CheckToken *Not Use
func (g *Game) CheckToken(token string) error {
	if serverToken, err := g.getToken(); err != nil {
		return errors.New("getToken error: ")
	} else if serverToken != token {
		return errors.New("token not equal: ")
	}
	return nil
}

func (g *Game) getToken() (string, error) {
	token, err := g.Server.HTTPConn.HTTPPostRawRequest("", nil)
	if err != nil {
		return "", err
	}
	return string(token), nil
}
