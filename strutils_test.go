package strutils_test

import (
	"fmt"
	"testing"

	"github.com/realjf/strutils"
)

func TestGetCode(t *testing.T) {
	code := "#aSNLRz2e#"
	fmt.Println(strutils.GetCode(code, 7))
	trial := "#TRIAL-jil23Cie#"
	fmt.Println(strutils.GetTrialCode(trial))
}

func TestCodeFormat(t *testing.T) {
	code := "#aSNLR2*/###z2e#"
	fmt.Println(strutils.CheckCodeFormat(code))
}

func TestCalcTokens(t *testing.T) {
	// str := "jifeonfe jogei jofj fe, joge, oj. ojweige 我是谁？"
	str := "魑魅魍魉"
	fmt.Println(strutils.CalcTokens(str))
}

func TestFilterRepeatedPunctuation(t *testing.T) {
	input := `嗯？！？！？！？！？！？！？！？！？！？！？！？！？！？！？！？！？！？！晚安，宝宝~！~！~！~！~！~！~！~！~！~！~！~！~！~！~！~！~！~！~！~！~！~！~！~！~！~， 乖~，mua~！[么么哒]`
	// fmt.Println(strings.Index(input, `~！`), strings.Index(input, `~`))
	fmt.Println(strutils.ReplaceRepeatingSubstrings(input, 3)) // Output: Hello! This. is. an example sentence.
}
