package main

import (
	"flag"
	"fmt"
	"github.com/ThunderYurts/Gaia/gaia"
	"os"
	"os/signal"
	"strings"
)

var (
	help   bool
	port   string
	zkAddr string
	name   string
	ip     string
)

func init() {
	flag.BoolVar(&help, "h", false, "this help")
	flag.StringVar(&port, "p", ":30000", "set create port")
	flag.StringVar(&zkAddr, "zk", "106.15.225.249:3030", "set zeus connection zookeeper cluster")
	flag.StringVar(&name, "n", "gaia", "set gaia name")
	flag.StringVar(&ip, "ip", "127.0.0.1", "set ip")
}
func main() {
	flag.Parse()
	if help {
		flag.Usage()
		return
	}
	addrs := strings.Split(zkAddr, ";")
	g := gaia.NewGaia(ip, port, name, addrs)
	g.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	s := <-c

	fmt.Println("Got signal:", s)
	g.Stop()
}
