package tian

import "errors"

const (
	C_dujitang        = "dujitang"   // 毒鸡汤
	C_mingyan         = "mingyan"    // 名人名言
	C_godreply        = "godreply"   // 神回复
	C_wanan           = "wanan"      // 晚安心语
	C_saylove         = "saylove"    // 土味情话
	C_caipu           = "caipu"      // 菜谱
	C_englishSentence = "ensentence" // 英语一句话
	C_lizhiguyan      = "lzmy"       // 励志古言
)

var (
	ErrNotfoundCaiPu = errors.New("not found cookbook")
)
