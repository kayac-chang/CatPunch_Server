package protocol

import (
	"net/http"

	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/foundation"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/myhttp"
)

// Request ...
type GameRequest struct {
	Token      string
	BetIndex   int64
	GameTypeID string
	PlayerID   int64
}

// InitData ...
func (c *GameRequest) InitData(r *http.Request) {
	c.Token = r.Header.Get("Authorization")
	postData := myhttp.PostData(r)
	c.BetIndex = foundation.InterfaceToInt64(postData["bet"])
	c.GameTypeID = foundation.InterfaceToString(postData["gametypeid"])
	c.PlayerID = foundation.InterfaceToInt64(postData["playerid"])
}

// // Respon ...
// type Respon struct {

// }

// // InitData ...
// func (c *Respon) InitData(r *http.Request) {
// 	postData := myhttp.PostData(r)
// 	c.Token = foundation.InterfaceToString(postData["token"])
// 	c.BetIndex = foundation.InterfaceToInt64(postData["bet"])
// 	c.GameTypeID = foundation.InterfaceToString(postData["gametypeid"])
// 	c.PlayerID = foundation.InterfaceToInt64(postData["playerid"])
// }
