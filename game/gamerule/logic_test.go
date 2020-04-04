package gamerule

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/YWJSonic/ServerUtility/foundation/fileload"
)

func TestNew(t *testing.T) {
	att := CatAttach{
		FreeCount:    0,
		IsLockBet:    0,
		LockBetIndex: 1,
	}
	for i := 0; i < 200; i++ {

		gamejsStr := fileload.Load("../../file/gameconfig.json")
		var gameRule = &Rule{}
		if err := json.Unmarshal([]byte(gamejsStr), &gameRule); err != nil {
			panic(err)
		}

		result := gameRule.newlogicResult(0, att)
		fmt.Println(result)
		att.FreeCount = result.Otherdata["freecount"].(int64)
		// result.Otherdata["freecount"]
	}
}
