# 自定义协议处理tcp粘包问题





- websocket不需要理论上不需要事先帖包，因为websocket传输数据的单元是msg，如果单个msg“过大”（具体多大尚不清楚）也需要自己处理粘包的问题
