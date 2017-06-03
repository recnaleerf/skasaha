package skasaha

import (
	"github.com/bwmarrin/snowflake"
)

var (
	node *snowflake.Node
)

type Snowflake int64

func (sn Snowflake) String() string {
	return snowflake.ID(sn).String()
}

func initSnowflake() {
	var (
		err error
	)

	node, err = snowflake.NewNode(0)
	if err != nil {
		panic(err)
	}
}

func NewSnowflake() Snowflake {
	id := node.Generate()

	return Snowflake(id)
}
