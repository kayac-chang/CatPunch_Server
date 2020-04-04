package gamerule

import (
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/attach"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/igame"
)

// CatAttach ...
type CatAttach struct {
	FreeCount    int64 `json:"freecount"`
	IsLockBet    int64 `json:"islockbet"`
	LockBetIndex int64 `json:"lockbetindex"`
}

// Rule ...
type Rule struct {
	BetRate             []int64 `json:"BetRate"`
	BetRateDefaultIndex int     `json:"BetRateDefaultIndex"`
	BetRateLinkIndex    []int64 `json:"BetRateLinkIndex"`
	FreeGameTrigger     int64   `json:"FreeGameTrigger"`
	FreeGameWinRate     []int   `json:"FreeGameWinRate"`
	FreeReel            [][]int `json:"FreeReel"`
	GameIndex           int     `json:"GameIndex"`
	GameTypeID          string  `json:"GameTypeID"`
	ItemResults         [][]int `json:"ItemResults"`
	Items               []int   `json:"Items"`
	NormalReel          [][]int `json:"NormalReel"`
	RespinReel          [][]int `json:"RespinReel"`
	RTPSetting          int     `json:"RTPSetting"`
	ScrollSize          []int   `json:"ScrollSize"`
	Space               int     `json:"Space"`
	Version             string  `json:"Version"`
	Wild                []int   `json:"Wild"`
	WinBetRateLimit     int64   `json:"WinBetRateLimit"`
	WinScoreLimit       int64   `json:"WinScoreLimit"`
}

// GetGameTypeID ...
func (r *Rule) GetGameTypeID() string {
	return r.GameTypeID
}

// GetBetMoney ...
func (r *Rule) GetBetMoney(index int64) int64 {
	return r.BetRate[index]
}

// GetReel ...
func (r *Rule) GetReel() map[string][][]int {
	scrollmap := map[string][][]int{
		"normalreel": r.normalReel(),
		"respinreel": {r.respuinScroll()},
	}
	return scrollmap
}

// GetBetSetting ...
func (r *Rule) GetBetSetting() map[string]interface{} {
	tmp := make(map[string]interface{})
	tmp["betrate"] = r.BetRate                         //betRate
	tmp["betratelinkindex"] = r.BetRateLinkIndex       //betRateLinkIndex
	tmp["betratedefaultindex"] = r.BetRateDefaultIndex //betRateDefaultIndex
	return tmp
}

// CheckGameType ...
func (r *Rule) CheckGameType(gameTypeID string) bool {
	if r.GameTypeID != gameTypeID {
		return false
	}
	return true
}

func (r *Rule) normalReel() [][]int {
	return r.NormalReel
}

func (r *Rule) wild1() int {
	return r.Wild[0]
}

func (r *Rule) wild2() int {
	return r.Wild[1]
}

// RespuinScroll ...
func (r *Rule) respuinScroll() []int {
	return r.RespinReel[r.RTPSetting]
}

// GameRequest ...
func (r *Rule) GameRequest(config *igame.RuleRequest) *igame.RuleRespond {
	result := make(map[string]interface{})
	otherData := make(map[string]interface{})
	var totalWin int64

	catAttach := r.GetAttach(*config.Attach)
	if catAttach.FreeCount >= int64(r.FreeGameTrigger) {
		catAttach.FreeCount %= int64(r.FreeGameTrigger)
		catAttach = r.newAttach()
	}
	if catAttach.IsLockBet > 0 {
		config.BetIndex = catAttach.LockBetIndex
	}

	betMoney := r.GetBetMoney(config.BetIndex)
	gameResult := r.newlogicResult(betMoney, catAttach)

	result["normalresult"] = gameResult.Normalresult
	result["isrespin"] = 0
	totalWin += gameResult.Normaltotalwin
	catAttach.FreeCount = gameResult.Otherdata["freecount"].(int64)
	if catAttach.FreeCount > 0 {
		catAttach.LockBetIndex = config.BetIndex
		catAttach.IsLockBet = 1
	}

	if gameResult.Respinresult != nil {
		result["respin"] = gameResult.Respinresult
		result["isrespin"] = 1
		totalWin += gameResult.Respintotalwin
	}

	result["totalwinscore"] = totalWin

	return &igame.RuleRespond{
		Attach:        r.outPutAttach(catAttach),
		BetMoney:      betMoney,
		Totalwinscore: totalWin,
		GameResult:    result,
		OtherData:     otherData,
	}
}

// GetAttach 0:free game count
func (r *Rule) GetAttach(att attach.IAttach) CatAttach {
	var info CatAttach
	count := att.Get(int64(r.GameIndex), 0)
	info.FreeCount = count.GetIValue()
	lockindex := att.Get(int64(r.GameIndex), 1)
	info.LockBetIndex = lockindex.GetIValue()
	if info.FreeCount > 0 {
		info.IsLockBet = 1
	}
	return info
}
func (r *Rule) outPutAttach(catAtt CatAttach) []*attach.Info {
	resAtt := make([]*attach.Info, 0, 2)
	freeCountAtt := attach.NewInfo(int64(r.GameIndex), 0, false)
	freeCountAtt.SetIValue(catAtt.FreeCount)
	lockAtt := attach.NewInfo(int64(r.GameIndex), 1, false)
	lockAtt.SetIValue(catAtt.LockBetIndex)

	resAtt = append(resAtt, freeCountAtt)
	resAtt = append(resAtt, lockAtt)
	return resAtt
}

func (r *Rule) newAttach() CatAttach {
	return CatAttach{
		FreeCount:    0,
		IsLockBet:    0,
		LockBetIndex: 0,
	}
}
