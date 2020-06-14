package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	fmt.Println("18")
	res, err := cli.ImagePull(ctx, "registry.cn-hangzhou.aliyuncs.com/brynelee/httpbin", types.ImagePullOptions{})
	fmt.Println("20")
	if err != nil {
		panic(err)
	}
	defer res.Close()
	io.Copy(os.Stdout, res)
	// containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	// if err != nil {
	// 	panic(err)
	// }
	// for _, container := range containers {
	// 	id := container.ID
	// 	stat, err := cli.ContainerStats(ctx, id, false)
	// 	if err != nil {
	// 		fmt.Println(err.Error())
	// 		continue
	// 	}
	// 	var dat map[string]interface{}
	// 	buf := new(bytes.Buffer)
	// 	buf.ReadFrom(stat.Body)
	// 	str := buf.String()
	// 	err = json.Unmarshal([]byte(str), &dat)
	// 	if err != nil {
	// 		fmt.Println(err.Error())
	// 		continue
	// 	}
	// 	cpuStats := dat["cpu_stats"]
	// 	mcpuStats := cpuStats.(map[string]interface{})
	// 	totalUsage := mcpuStats["cpu_usage"].(map[string]interface{})
	// 	memoryStats := dat["memory_stats"].(map[string]interface{})
	// 	fmt.Printf("memory usage %v\n", memoryStats["usage"].(float64))
	// 	fmt.Printf("cpu usage %v\n", totalUsage["total_usage"].(float64))
	// }
}
