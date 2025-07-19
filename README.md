# mn-utils

基础组件库



### logger 日志打印

[详细示例](https://github.com/shanghai-mingnew/mn-utils/blob/main/logger/readme.md)
```
import "github.com/shanghai-mingnew/mn-utils/logger" 
logger 日志打印
```

### mnnet 网络相关

```
import "github.com/shanghai-mingnew/mn-utils/mnnet" 

mnnet.Listen("tcp4", ":8080") // 独占监听端口，防止端口复用被其他程序占用
```