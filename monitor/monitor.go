package monitor

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/ThunderYurts/Gaia/gconst"
	"github.com/ThunderYurts/Gaia/zookeeper"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/samuel/go-zookeeper/zk"
)

// SimpleMonitor will fetch all container and get them brief info
type SimpleMonitor struct {
	ctx        context.Context
	client     *client.Client
	conn       *zk.Conn
	name       string
	createAddr string
}

// NewSimpleMonitor is a help function
func NewSimpleMonitor(ctx context.Context, client *client.Client, conn *zk.Conn, name string, createAddr string) SimpleMonitor {
	return SimpleMonitor{
		ctx:        ctx,
		client:     client,
		conn:       conn,
		name:       name,
		createAddr: createAddr,
	}
}

// SyncNodeStat will return list of ContainerStat
func (sm *SimpleMonitor) SyncNodeStat() error {
	// create or get the version of zknode
	fmt.Printf("check zknode %s\n", gconst.GaiaRoot+"/"+sm.name)
	exist, stat, err := sm.conn.Exists(gconst.GaiaRoot + "/" + sm.name)
	if err != nil {
		return err
	}
	if !exist {
		_, err := sm.conn.Create(gconst.GaiaRoot+"/"+sm.name, []byte{}, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
		if err != nil {
			return err
		}
		_, stat, err = sm.conn.Get(gconst.GaiaRoot + "/" + sm.name)
		if err != nil {
			return err
		}
	}
	fmt.Printf("has create znode %v\n", gconst.GaiaRoot+"/"+sm.name)
	for {
		select {
		case <-sm.ctx.Done():
			{
				return nil
			}
		default:
			{
				arg, err := filters.ParseFlag(gconst.YurtFilter, filters.NewArgs())
				if err != nil {
					panic(err)
				}
				containers, err := sm.client.ContainerList(sm.ctx, types.ContainerListOptions{
					Filters: arg,
				})
				//fmt.Printf("containers in check: %v\n", containers)

				if err != nil {
					fmt.Printf("72 %s\n", err.Error())
					return err
				}
				Memory := float64(0)
				CPU := float64(0)
				for _, container := range containers {
					id := container.ID
					stat, err := sm.client.ContainerStats(sm.ctx, id, false)
					if err != nil {
						fmt.Printf("81 %s\n", err.Error())
						continue
					}
					var dat map[string]interface{}
					buf := new(bytes.Buffer)
					_, err = buf.ReadFrom(stat.Body)
					if err != nil {
						fmt.Println(err.Error())
						continue
					}
					str := buf.String()
					err = json.Unmarshal([]byte(str), &dat)
					if err != nil {
						fmt.Println(err.Error())
						continue
					}
					cpuStats := dat["cpu_stats"]
					mcpuStats := cpuStats.(map[string]interface{})
					totalUsage := mcpuStats["cpu_usage"].(map[string]interface{})
					memoryStats := dat["memory_stats"].(map[string]interface{})
					if memoryStats["usage"] != nil {
						Memory = Memory + memoryStats["usage"].(float64)
					}
					if totalUsage["total_usage"] != nil {
						CPU = CPU + totalUsage["total_usage"].(float64)
					}

					// fmt.Printf("memory usage %v\n", memoryStats["usage"].(float64))
					// fmt.Printf("cpu usage %v\n", totalUsage["total_usage"].(float64))
				}
				nodeInfo := zookeeper.ZKNode{CPU: CPU, Memory: Memory, CreateAddr: sm.createAddr}
				buf := new(bytes.Buffer)
				enc := gob.NewEncoder(buf)
				err = enc.Encode(nodeInfo)
				if err != nil {
					return err
				}
				stat, err = sm.conn.Set(gconst.GaiaRoot+"/"+sm.name, buf.Bytes(), stat.Version)
				// TODO sync data to zookeeper
				//fmt.Printf("update value %v\n", nodeInfo)
				if err != nil {
					return err
				}
			}
		}
	}
}
