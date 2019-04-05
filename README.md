# logstash-go
## 日志转发工具

##1、编译
````
//windows 编译 linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

//Mac下编译Linux, Windows平台的64位可执行程序：
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build test.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build test.go

//Linux下编译Mac, Windows平台的64位可执行程序：
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build test.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build test.go
````
## 2、部署 logstash

将编译好的 logstash 放到服务器的任意目录
给执行权限
````
$ chmod +x logstash
$ ./logstash
````

## 3、打开客户端链接对应的端口
client.html 

`ws://example.com:9191`

## 6、配置服务器变量

```
HTTP_LOG_SERVER=127.0.0.1
HTTP_LOG_PORT=9192

UDP_LOG_SERVER=127.0.0.1
UDP_LOG_PORT=9093
```

## 5、调用 php
```
<?php
class Test extends Command
{
    use UdpLogTrait;
    
    public function index(){
    
        $data = [];
        $this->write($data);
    }    
}
?>

<?php
class Test extends Command
{
    use HttpLogTrait;
    
    public function index(){
    
        $data = [];
        $this->write($data);
    }    
}
?>
```


## 6、服务器应该打开对应的端口配置供应商的防火墙规则