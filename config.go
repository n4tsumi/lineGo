package lineGo

import "context"

const (
	LineHost = "https://legy-jp-addr.line.naver.jp"
	Register = "/api/v4/TalkService.do"
	Normal   = LineHost + "/S4"
	Polling  = LineHost + "/P4"

	SystemName = "lineGo"

	UserAgent = "Line/5.24.1"
	LineApp   = "DESKTOPMAC\t5.24.1\tOS X\t10.15.1"
)

var ctx = context.Background()
