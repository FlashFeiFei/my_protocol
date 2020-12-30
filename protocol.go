package my_protocol

import (
	"bytes"
	"encoding/binary"
)

//tcp粘包，websocket不需要粘包(websocket已实现)
//自定义协议，实现封包和贴包
//为什么需要自定义协议 ？ 服务端每次读取的数据都是读取缓冲中的所有数据，但是缓存中的数据不一定完整
//比如数据传的是json。缓存中的数据是  {"name":张三}{"name":李四}{"name":阿强};   这个数据可能会读取不完整，在中间截断之类的

//1. 封包 (一个完整的数据) = 【请求头 + 总数据长度】 + 总数据
//2. 粘包:把每次读取的数据拼接成一个完整的数据包
const (
	ConstHeader         = "ly_protobuf" //数据包的头
	ConstHeaderLength   = 11
	ConstSaveDataLength = 4
)

//整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}

//封包 = (一个完整的数据) = 【请求头 + 总数据长度】 + 总数据
//返回一个完整的数据
func Packet(message []byte) []byte {
	//一个完整的数据 = 请求头 + 数据长度
	return append(append([]byte(ConstHeader), IntToBytes(len(message))...), message...)
}

//帖包
//buffer参数是，上一次截断的数据+下一次读取的数据
func Unpack(buffer []byte, readerChannel chan []byte) []byte {

	//获取数据的长度
	length := len(buffer)

	var i int

	for i = 0; i < length; i = i + 1 {

		if length < i+ConstHeaderLength+ConstSaveDataLength {
			//数据不完整的截断数据
			break
		}

		if string(buffer[i:i+ConstHeaderLength]) == ConstHeader {
			//如果是请求头开始的，说明这是一个新数据开始
			//获取新数据的数据总长度
			messageLength := BytesToInt(buffer[i+ConstHeaderLength : i+ConstHeaderLength+ConstSaveDataLength])
			if length < i+ConstHeaderLength+ConstSaveDataLength+messageLength {
				//本次buffer数据是一个不完整的数据，跳出循环,函数结束，返回截断数据的剩余内容
				break
			}

			//buffer数据中包含了整个完整的数据，直接数据数据
			data := buffer[i+ConstHeaderLength+ConstSaveDataLength : i+ConstHeaderLength+ConstSaveDataLength+messageLength]
			//一个完整的数据接受完成，通知下层
			readerChannel <- data

			//跳过这个完成的数据
			i += ConstHeaderLength + ConstSaveDataLength + messageLength - 1
		}
	}

	if i == length {
		return make([]byte, 0)
	}

	//返回不完整数据，等待下一次填充
	return buffer[i:]
}
