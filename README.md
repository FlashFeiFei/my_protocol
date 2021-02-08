# 自定义协议处理tcp粘包问题





- websocket不需要理论上不需要事先帖包，因为websocket传输数据的单元是msg，如果单个msg“过大”（具体多大尚不清楚）也需要自己处理粘包的问题


# 帖包，拆包的基本原理

- tpc流传输
- 发送数据之前，给这个数据添加上 “数据类型4个字节” + “数据总长度4个字节” ，共8个字节
- 解析数据的时候分两步读取,首先读取8个字节，如果能读取，说明是自定义协议，然后根据总长度第二次读取数据内容
- 根据类型将数据内容反序列化出来



[参看源码](https://github.com/aceld/zinx/blob/master/znet/datapack.go)



![图片](https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fi1.hdslb.com%2Fbfs%2Farchive%2Fd727df0385b9d6432d0a4e061dbc86e9f13898b5.jpg&refer=http%3A%2F%2Fi1.hdslb.com&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=jpeg?sec=1615340400&t=fac4f54646f6859076d1a1cd3f3a17ac "标题")


