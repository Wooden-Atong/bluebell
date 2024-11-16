package snowflake

import (
	_ "fmt"
	"time"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	snowflake.Epoch = st.UnixNano() / 1000000
	node, err = snowflake.NewNode(machineID)
	return
}
func GenID() int64 {
	return node.Generate().Int64()
}

/*
//🌟目前生成id是作为一个模块放入项目中，且没有真正涉及分布式。如果真的用的很频繁时，则可以写一个分布式服务嵌在内网中，用的时候就调一下
func main() {
	// 🌟起始时间设置为2020-07-01，可以用到往后的69年；machineID设置为1，因为实际上没有真正写一个分布式，随便传入一个便可以了
	if err := Init("2020-07-01", 1); err != nil {
		fmt.Printf("init failed, err:%v\n", err)
		return
	}
	id := GenID()
	fmt.Println(id)
}
*/