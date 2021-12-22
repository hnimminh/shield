package main

import (
	"os"
	"fmt"
	"text/template"
	"strings"
	"strconv"
)



type Rule struct{
	SrcIP []string
    DstIP []string
    SrcPort []int
	DstPort []int
	Transport string
	Action string
    Comment string
}

type Security struct{
    Builtin bool
    Rules []Rule
}



func main(){

    NFTEMPLATE := `
flush ruleset
table inet SHIELD {

    set WhiteList {
    }

    set BlackList {
    }

    set TemporaryBlocks {
    }

    chain Inbound {
        type filter hook input priority 0; policy drop;
            iifname lo accept
            #---------------------------------------------------------------------
            ip saddr @WhiteList accept
            ip saddr @BlackList drop
            ip saddr @TemporaryBlocks drop
            #---------------------------------------------------------------------
            {{- if .Thread}}
            ip frag-off & 0x1fff != 0 counter drop comment "IP FRAGMENTS"
            tcp flags != syn ct state new drop comment "FIRST MEET BUT NOT SYN"
            tcp flags & (fin|syn) == (fin|syn) drop comment "NEW BUT FIN"
            tcp flags & (syn|rst) == (syn|rst) drop comment "NEVER MET BUT RESET"
            tcp flags & (fin|syn|rst|psh|ack|urg) < (fin) drop comment "ATTACK"
            tcp flags & (fin|syn|rst|psh|ack|urg) == (fin|psh|urg) drop comment "XMAS ATTACK"
            tcp flags & (fin|syn|rst|psh|ack|urg) == 0x0 counter drop comment "NULL"
            tcp flags syn tcp option maxseg size 1-536 counter drop comment "TCPSEGSIZE"
            ct state invalid counter drop comment "INVALID STATE"
            ct state {established, related} counter accept
            #---------------------------------------------------------------------
            {{- end}}
			{{- range $Rule := .Rules }}
            {{$Rule.Transport}} ip saddr { {{StrsJoin $Rule.SrcIP}} } sport { {{IntsJoin $Rule.SrcPort}} } ip daddr { {{StrsJoin $Rule.DstIP}} } dport { {{IntsJoin $Rule.DstPort}} } counter {{ $Rule.Action }} comment "{{ $Rule.Comment }}"
			{{- end}}
            # count and drop any other traffic
            counter drop
    }
    chain Outbound {
        type filter hook output priority 0;
    }
    chain Forward {
        type filter hook forward priority 0; policy drop;
    }
}
`

    sec := Security{true, []Rule{ Rule{ []string{"1.1.1.1", "3.3.3.3"},
										[]string{"2.2.2.2", "4.4.4.4"},
										[]int{10000, 10001}, []int{443, 80},
										"tcp",
										"drop",
										"Testing Rule TCP"},
									Rule{[]string{"5.5.5.5"},
										 []string{"6.6.6.6"},
										 []int{10200}, []int{5060},
										 "udp",
										 "accept",
										 "Testing Rule UDP"}}}

    tmpfw, _ := template.New("firewall").Funcs(template.FuncMap{"StrsJoin": strsjoin, "IntsJoin": intsjoin}).Parse(NFTEMPLATE)
    err := tmpfw.Execute(os.Stdout, sec)
	if err != nil {
		fmt.Println(err)
	}
}


func strsjoin(input []string) string {
	return strings.Join(input, ",")
}


func intsjoin(input []int) string{
	strs := []string{}
	for _, number := range input {
		strs = append(strs, strconv.Itoa(number))
	}
	return strsjoin(strs)
}

