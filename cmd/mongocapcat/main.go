package main

import (
	"flag"
	"fmt"
	"os"

	"code.google.com/p/gopacket/pcap"
	"github.com/tmc/mongocaputils"
	"github.com/tmc/mongocaputils/mongoproto"
)

var (
	pcapFile      = flag.String("f", "-", "pcap file (or '-' for stdin)")
	packetBufSize = flag.Int("size", 1000, "size of packet buffer used for ordering within streams")
	verbose       = flag.Bool("v", false, "verbose output (to stderr)")
)

func main() {
	flag.Parse()

	pcap, err := pcap.OpenOffline(*pcapFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error opening pcap file:", err)
		os.Exit(1)
	}
	h := mongocaputils.NewPacketHandler(pcap)
	m := mongocaputils.NewMongoOpStream(*packetBufSize)

	ch := make(chan struct{})
	go func() {
		defer close(ch)
		for op := range m.Ops {
			if _, ok := op.(*mongoproto.OpUnknown); !ok {
				fmt.Println(op)
			}
		}
	}()

	if err := h.Handle(m, -1); err != nil {
		fmt.Fprintln(os.Stderr, "mongocapcat: error handling packet stream:", err)
	}
	<-ch
}
