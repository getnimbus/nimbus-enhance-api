package setting

const (
	HttpPort     string = ":8081"
	QueryTimeout        = 20 // for sql context timeout

	StatusNotReady   = 1
	StatusProcessing = 2
	StatusDone       = 3
	StatusFail       = 4
)
