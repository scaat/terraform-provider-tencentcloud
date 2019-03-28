package tencentcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"
	"sync/atomic"
	"time"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

const FILED_SP = "#"

var firstLogTime = ""
var logAtomaticId int64 = 0

//init log config
func InitLogConfig(saveLocalFile bool) {
	if firstLogTime == "" {
		firstLogTime = fmt.Sprintf("%x", time.Now().Unix())
	}

	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)

	if saveLocalFile {
		logOut := &lumberjack.Logger{
			Filename:   "./tencentcloud.log",
			MaxSize:    2 * 1024, //2G
			MaxBackups: 10,
			Compress:   true,
			LocalTime:  true,
		}
		log.SetOutput(logOut)
	}
}

//get logid  for trace, return a new logid if ctx is nil
func GetLogId(ctx context.Context) string {

	if ctx != nil {
		logId, ok := ctx.Value("logId").(string)
		if ok {
			return logId
		}
	}
	return fmt.Sprintf("%s-%d", firstLogTime, atomic.AddInt64(&logAtomaticId, 1))
}

//write data to file
func writeToFile(filePath string, data interface{}) error {

	if strings.HasPrefix(filePath, "~") {
		usr, err := user.Current()
		if err != nil {
			return fmt.Errorf("Get current user fail,reason %s", err.Error())
		}
		if usr.HomeDir != "" {
			filePath = strings.Replace(filePath, "~", usr.HomeDir, 1)
		}
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("stat old file error,reason %s", err.Error())
	}
	if !os.IsNotExist(err) {
		if fileInfo.IsDir() {
			return fmt.Errorf("old filepath is a dir,can not delete")
		}
		if err := os.Remove(filePath); err != nil {
			return fmt.Errorf("delete old file error,reason %s", err.Error())
		}
	}
	jsonStr, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return fmt.Errorf("json decode error,reason %s", err.Error())
	}
	return ioutil.WriteFile(filePath, []byte(jsonStr), 422)
}
