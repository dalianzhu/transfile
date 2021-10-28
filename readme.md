## 启动agent
找个服务器，启动agent
```
transfile -a 127.0.0.1:9886 agent
```
收发文件需要先建立服务端


## 发送文件：
```
transfile -a 127.0.0.1:9886 -c 1 put hello.tgz 
```

## 接收文件：
```
transfile -a 127.0.0.1:9886 -c 1 get hello.tgz 
transfile -a 127.0.0.1:9886 -c 1 get // 使用默认的名字
```

## 环境变量（可选）
```
export f_address=127.0.0.1:9886
export f_code=1
```