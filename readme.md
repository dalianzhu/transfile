## 启动agent
找个服务器，启动agent

./transfile -m agent

## put文件
设置address为服务器的地址

export address=127.0.0.1:9886

设置一个code，比如1，get文件的时候需要提供这个code

export code=1

./transfile put -f ./file.tgz

## get文件
export address=127.0.0.1:9886

code与put的一致

export code=1

./transfile get -f ./filecp.tgz
