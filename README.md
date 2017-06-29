# thrift_transport_pool

## 简介：

这是一个`thrift.TTransport`的Pool。提供Pool和Retry功能。它会自动`Open`这个`Transport`，并处理在建立链接过程中的重连以及传输数据失败后重连。

## 安装

```
go get -u github.com/xuxiangyang/thrift_transport_pool
```



## 使用

```
// 定义一个pool，共享使用
pool := thrift_transport_pool.NewPool(16, "127.0.0.1:3000", func(hostPort string) (thrift.TTransport, error) {
	socket, err := thrift.NewTSocket(hostPort)
    if err != nil {
		return nil, err
    }
    transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
    return transportFactory.GetTransport(socket)
})

// 假设在a函数里使用。假设是UserService

func a() {
  t, err := pool.Pop()
  if err != nil {
    panic(err)
  }
  defer pool.Push(t)
  user_client := user_thrift.NewUserServiceClientFactory(t, thrift.NewTBinaryProtocolFactoryDefault())
  user_client.Find("ID") //userService定义的函数
}
```

## 可能存在的问题

对于重试，目前重试的错误可能还不够完善
