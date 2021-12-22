package main

import (
	"fmt"
	"flag"
)

func main(){
	id := flag.String("id", "shield", "shield identity")
	group := flag.String("group", "shield", "shield group name")
	listen := flag.String("listen", "0.0.0.0:0", "listen address include IP and port, disabled by default")
	brokers := flag.String("brokers", "", "list of broker address include IP and port that shield will attached to")
	engine := flag.String("engine", "nftables", "firewall engine. current support:\n- nftables \n- cilium")

	flag.Parse()

    fmt.Println("id:", *id)
    fmt.Println("group:", *group)
    fmt.Println("listen:", *listen)
    fmt.Println("brokers:", *brokers)
    fmt.Println("engine:", *engine)
	fmt.Println("args:", flag.Args())

}
//https://gobyexample.com/command-line-flags
// shield -id jesus -group religion -listen 127.0.0.1:12345 -brokers 10.10.2.12:9022,10.10.1.11:9022 -engine nftables/cilium
