package gkafkav2

type Msg struct {
	Key       string
	Value     string
	Partition int32
	Offset    int64
}
