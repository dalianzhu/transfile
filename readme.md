## 启动agent
找个服务器，启动agent
```
transfile -a 127.0.0.1:9886 agent
```
收发文件需要先建立服务端


## 发送文件：
-c 为管道。不同的管道传送不同的文件。

put的时候，put端会等待get端到来。
```
transfile -a 127.0.0.1:9886 -c 1 put hello.tgz 
# 配置环境变量后
transfile put hello.tgz
```

## 接收文件：
```
transfile -a 127.0.0.1:9886 -c 1 get hello.tgz 
transfile -a 127.0.0.1:9886 -c 1 get // 使用默认的名字
# 配置环境变量后
transfile get // 使用默认的名字
```

## 环境变量（可选）
```
export f_address=127.0.0.1:9886
export f_code=1
```