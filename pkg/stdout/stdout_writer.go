package stdout

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"
)

type outputWriter struct {
}

func SetOutputWriter() {
	log.SetFlags(0)
	log.SetOutput(new(outputWriter))
}

func (writer outputWriter) Write(bytes []byte) (int, error) {

	// 获取调用者信息（跳过2层调用栈）
	_, filePath, line, ok := runtime.Caller(3)

	if !ok {
		filePath = "unknown"
		line = 0
	}

	lineLength := len(filePath)
	lineMaxLength := 36

	if lineLength > lineMaxLength {
		filePath = "....." + string(filePath[lineLength-lineMaxLength+4:lineLength])
	} else if lineLength < lineMaxLength {
		filePath = PadLeft(filePath, " ", lineMaxLength+1)
	}

	s := fmt.Sprintf("{%s:%d}", filePath, line)

	return fmt.Print(time.Now().Format(TimeStampFormat) + " SOUT - " + s + " > " + string(bytes))
}

func PadLeft(s string, padStr string, maxLen int) string {
	var padCountInt = 1 + ((maxLen - len(padStr)) / len(padStr))
	var retStr = strings.Repeat(padStr, padCountInt) + s

	return retStr[(len(retStr) - maxLen):]
}
