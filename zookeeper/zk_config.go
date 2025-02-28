package zookeeper

// ZKNode is a pure data structure synced into zk to show node status
type ZKNode struct {
	Memory     float64
	CPU        float64
	CreateAddr string
}
