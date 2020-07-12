package lineGo

import (
	talk "./talkservice"
	"github.com/apache/thrift/lib/go/thrift"
	"log"
	"net/http"
)

type LoginType int

const (
	authToken LoginType = iota
	qrCode
)

type LineLogin interface {
	Type() LoginType
	Value() string
}

type AuthTokenLogin string

func (a AuthTokenLogin) Type() LoginType {
	return authToken
}

func (a AuthTokenLogin) Value() string {
	return string(a)
}

func Token(authToken string) AuthTokenLogin {
	return AuthTokenLogin(authToken)
}

type QrCodeLogin struct{}

func (a QrCodeLogin) Type() LoginType {
	return qrCode
}

func (a QrCodeLogin) Value() string {
	return ""
}

var QrCode = QrCodeLogin{}

type OptionType int

const (
	withMid OptionType = iota
)

type LoginOption interface {
	Type() OptionType
	Value() string
}

type OptionWithMid string

func (o OptionWithMid) Type() OptionType {
	return withMid
}

func (o OptionWithMid) Value() string {
	return string(o)
}

func WithMid(mid string) OptionWithMid {
	return OptionWithMid(mid)
}

func createSession(authToken, path string, client *http.Client) *talk.TalkServiceClient {
	trans, err := thrift.NewTHttpClientWithOptions(path, thrift.THttpClientOptions{Client: client})
	if err != nil {
		log.Fatal(err)
	}
	httpTrans := trans.(*thrift.THttpClient)
	httpTrans.SetHeader("X-Line-Access", authToken)
	httpTrans.SetHeader("User-Agent", UserAgent)
	httpTrans.SetHeader("X-Line-Application", LineApp)
	prot := thrift.NewTCompactProtocol(trans)
	return talk.NewTalkServiceClient(thrift.NewTStandardClient(prot, prot))
}
