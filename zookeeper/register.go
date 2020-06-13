package zookeeper

import "github.com/samuel/go-zookeeper/zk"

// Register will create an Ephemeral Node in GaiaRoot and sync local monitor data and and yurt status to zeus
type Register struct {
	conn *zk.Conn
	node ZKNode
}
