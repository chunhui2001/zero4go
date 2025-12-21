package gkafkav2

type Msg struct {
	Key       []byte
	Headers   map[string][]byte
	Value     []byte
	Partition int32
	Offset    int64
}
