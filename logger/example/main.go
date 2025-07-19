package main

import "mn-utils/logger"

func main() {
	logger.SetLevel("debug")
	logger.SetLogFile("./test.log", "./err.log", logger.LogFileConf{})
	logger.DebugStackln("test debug", "xxxx")
	logger.Infoln("test info")
	logger.Warnln("test warn")
	logger.Errorln("test error")
	// logger.Fatalln("test fatal")
	logger.Panicln("test panic")
}
