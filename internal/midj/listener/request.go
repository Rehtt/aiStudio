package listener

var request func(r ReqCb)

type reqFunc func(r ReqCb)

func Request(r reqFunc) {
	request = r
}
