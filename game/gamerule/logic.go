package gamerule

import (
	"fmt"

	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/foundation"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/gamesystem"
)

type result struct {
	Normalresult   map[string]interface{}
	Otherdata      map[string]interface{}
	Normaltotalwin int64
	Respinresult   []interface{}
	Respintotalwin int64
}

// Result att 0: freecount
func (r *Rule) newlogicResult(betMoney int64, att CatAttach) result {
	var totalWin int64

	normalresult, otherdata, normaltotalwin := r.outputGame(betMoney, att.FreeCount)
	// result = foundation.AppendMap(result, otherdata)
	// result["normalresult"] = normalresult
	// result["islockbet"] = 0
	totalWin += normaltotalwin
	fmt.Println("----normalresult----", normalresult)
	fmt.Println("----otherdata----", otherdata)
	fmt.Println("----normaltotalwin----", normaltotalwin)

	if otherdata["isrespin"].(int) == 1 {
		respinresult, respintotalwin := r.outRespin(totalWin)
		totalWin = respintotalwin
		// result["respin"] = respinresult
		// result["isrespin"] = 1
		fmt.Println("----respinresult----", respinresult)
		fmt.Println("----respintotalwin----", respintotalwin)
	}

	if otherdata["isfreegame"].(int) == 1 {
		freeresult, freetotalwin := r.outputFreeSpin(betMoney)
		totalWin += freetotalwin
		// result["freegame"] = freeresult
		// result["isfreegame"] = 1
		fmt.Println("----freeresult----", freeresult)
		fmt.Println("----freetotalwin----", freetotalwin)
	}

	// result["totalwinscore"] = totalWin
	return result{
		Normalresult:   normalresult,
		Otherdata:      otherdata,
		Normaltotalwin: normaltotalwin,
	}

}

func (r *Rule) outputGame(betMoney int64, freecount int64) (map[string]interface{}, map[string]interface{}, int64) {
	var totalScores int64
	var result map[string]interface{}
	otherdata := make(map[string]interface{})
	islink := false

	ScrollIndex, plate := gamesystem.NewPlate(r.ScrollSize, r.normalReel())

	// count++
	// plate = TestPlate(count % 4)
	gameresult := r.winresultArray(plate)

	otherdata["isfreegame"] = 0
	otherdata["freecount"] = freecount % r.FreeGameTrigger
	otherdata["isrespin"] = 0

	if r.isFreeGameCount(plate) {
		freecount++
		if freecount >= r.FreeGameTrigger {
			otherdata["isfreegame"] = 1
			otherdata["freecount"] = freecount
		} else {
			otherdata["freecount"] = freecount
		}
	}

	if r.isRespin(plate) {
		otherdata["isrespin"] = 1
	}

	if len(gameresult) > 0 {
		islink = true
		totalScores = betMoney * int64(gameresult[0][3])
	}

	result = gamesystem.ResultMap(ScrollIndex, plate, totalScores, islink)
	return result, otherdata, totalScores
}

func (r *Rule) outputFreeSpin(betMoney int64) ([]interface{}, int64) {
	var result []interface{}
	var totalScores int64
	var freewinScore int64
	var freeresult map[string]interface{}
	islink := false

	for i, max := 0, len(r.FreeGameWinRate); i < max; i++ {
		ScrollIndex, plate := gamesystem.NewPlate(r.ScrollSize, r.FreeReel)
		gameresult := r.winresultArray(plate)
		freewinScore = 0
		islink = false
		if len(gameresult) > 0 {
			islink = true
			freewinScore = betMoney * int64(gameresult[0][3]) * int64(r.FreeGameWinRate[i])
		}

		totalScores += freewinScore
		freeresult = gamesystem.ResultMap(ScrollIndex, plate, freewinScore, islink)
		result = append(result, freeresult)
	}
	return result, totalScores
}

func (r *Rule) outRespin(normalScore int64) (map[string]interface{}, int64) {
	var totalScores int64
	islink := false

	ScrollIndex, plate := gamesystem.NewPlate([]int{1}, [][]int{r.respuinScroll()})
	gameresult := r.respinResult(plate)

	if len(gameresult) > 0 {
		islink = true
		totalScores = normalScore * int64(gameresult[0][1])
	}

	result := gamesystem.ResultMap(ScrollIndex, plate, totalScores, islink)
	return result, totalScores
}

// winresultArray ...
func (r *Rule) winresultArray(plate []int) [][]int {
	var result [][]int
	var dynamicresult []int

	for _, ItemResult := range r.ItemResults {
		if r.isWin(plate, ItemResult) {

			if r.isDynamicResult(ItemResult) {
				dynamicresult = r.dynamicScore(plate, ItemResult)
				result = append(result, dynamicresult)
				break
			} else {
				result = append(result, ItemResult)
				break
			}
		}
	}

	return result
}

// RespinResult result 0: icon index, 1: win rate
func (r *Rule) respinResult(plate []int) [][]int {
	var result [][]int

	switch plate[0] {
	case 2:
		result = append(result, []int{2, 5})
	case 3:
		result = append(result, []int{3, 7})
	case 4:
		result = append(result, []int{4, 10})
	}

	return result
}

// IsFreeGameCount ...
func (r *Rule) isFreeGameCount(plate []int) bool {
	if plate[1] == 1 {
		return true
	}
	return false

}

// IsRespin ...
func (r *Rule) isRespin(plate []int) bool {
	if plate[0] != 10 && plate[1] == 1 && plate[2] == 0 {
		return true
	}
	return false

}

func (r *Rule) isWin(plates []int, result []int) bool {
	IsWin := false

	if r.isBounsGame(result) {
		if r.isRespin(plates) {
			return true
		}

		return false
	}

	for i, plate := range plates {
		IsWin = false

		if plate == r.Space {
			return false
		}

		if plate == r.wild1() || plate == r.wild2() {
			IsWin = true
		} else {

			switch result[i] {
			case plate:
				IsWin = true
			case -1000:
				IsWin = true
			case -1001: // any 7
				if foundation.IsInclude(plate, []int{5, 6}) {
					IsWin = true
				}
			case -1002: // any bar
				if foundation.IsInclude(plate, []int{6, 7, 8, 9}) {
					IsWin = true
				}
			}
		}
		if !IsWin {
			return IsWin
		}
	}

	return IsWin
}

// isBounsGame bouns game reul: itemresult score < 0
func (r *Rule) isBounsGame(plates []int) bool {
	if plates[len(plates)-1] < 0 {
		return true
	}
	return false
}

func (r *Rule) dynamicScore(plant, currendResult []int) []int {
	dynamicresult := make([]int, len(currendResult))
	copy(dynamicresult, currendResult)

	switch currendResult[3] {
	case -100:
		for _, result := range r.ItemResults {
			if result[0] == plant[0] {
				dynamicresult[3] = result[3]
				break
			}
		}
	}

	return dynamicresult
}

func (r *Rule) isDynamicResult(result []int) bool {
	if result[3] < 0 {
		return true
	}
	return false
}
