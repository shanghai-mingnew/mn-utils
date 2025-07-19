### 日志打印说明

```
logger.Debugln("test debug")
logger.Infoln("test info")
logger.Warnln("test warn")
logger.Errorln("test error")
logger.DebugStackln("test debug") 打印堆栈并输出
logger.InfoStackln("test info") 打印堆栈并输出
logger.WarnStackln("test warn") 打印堆栈并输出
logger.ErrorStackln("test error") 打印堆栈并输出
logger.Fatalln("test fatal")  打印堆栈后退出程序
logger.Panicln("test panic")  打印堆栈后触发Panic
```

不推荐的写法

```
logger.Infoln(fmt.Sprintf("%s 同步数据集中断", suffix))
logger.Errorln("解析数据集的数据输入失败，表名：" + dataset.TableName + "错误信息：" + err.Error())
```

推荐写法

```
logger.Infoln(suffix, "同步数据集中断")
logger.Errorln("解析数据集的数据输入失败，表名：", dataset.TableName, "错误信息：", err)
```

默认info等级，可通过设置环境变量LOG_LEVEL=debug修改日志等级
或者通过SetLevel("debug")方式修改
高于等于error的日志，会输出到stderr
低于等于warn的日志，会输出到stdout

**各日志等级说明**

* debug  仅调试
* info   辅助日志信息
* warn   意料中的错误或警告
* error  需要被关注的错误
* fatal  打印堆栈后退出程序,退出code 1
* panic  打印堆栈后触发Panic,退出code 2

**日志打印**

1. 调试日志仅使用Debug方法
2. 打包环境日志等级必须大于等于info
3. 对于意料中的错误，使用warn方法
4. 对于系统错误，使用error方式；如err:=os.Remove("xxx");logger.Errorln(err)

**设置输出到文件**

```
  区分日志等级的写法：
	logger.SetLogFile("./info.log", "./err.log", logger.LogFileConf{})
  不区分日志等级的写法：
	logger.SetLogFile("./test.log", "", logger.LogFileConf{})
```

**GIN**

提供GinLoggerMid()方法处理gin路由日志
http状态为400以下则打印debug日志
http状态为400及以上则打印error日志
使用方法

```
  gin.SetMode(gin.ReleaseMode)
  r := gin.New()
  r.Use(gin.Recovery())
  r.Use(logger.GinLoggerMid())

```

**GORM**

提供NewGormLog(time)方法处理gorm日志
参数：慢语句时长
例：logger.NewGormLog(10 * time.Second)
说明：执行超过10s的语句会打印warn日志
使用方法

```
  sqliteConn, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{Logger: logger.NewGormLog(10 * time.Second)})

```
