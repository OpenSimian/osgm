package main

import (
	"encoding/json"
	"flag"
	"net"
	"syscall"
	"time"

	"code.google.com/p/gcfg"

	"github.com/M-O-S-E-S/mgm/core"
	"github.com/M-O-S-E-S/mgm/mgm"
	"github.com/jcelliott/lumber"
	pscpu "github.com/shirou/gopsutil/cpu"
	psmem "github.com/shirou/gopsutil/mem"
	psnet "github.com/shirou/gopsutil/net"
)

type nodeConfig struct {
	Node struct {
		OpensimBinDir string
		RegionDir     string
		MGMAddress    string
	}

	Opensim struct {
		MinRegionPort   int
		MaxRegionPort   int
		MinConsolePort  int
		MaxConsolePort  int
		ExternalAddress string
	}
}

type mgmNode struct {
	logger core.Logger
}

func main() {
	n := mgmNode{lumber.NewConsoleLogger(lumber.DEBUG)}

	cfgPtr := flag.String("config", "/opt/mgm/node.gcfg", "path to config file")
	flag.Parse()

	//read configuration file
	config := nodeConfig{}
	err := gcfg.ReadFileInto(&config, *cfgPtr)
	if err != nil {
		n.logger.Fatal("Error reading config file: ", err)
		return
	}

	n.logger.Info("config loaded successfully")

	hStats := make(chan mgm.HostStat, 8)
	mgmCommands := make(chan []byte, 32)
	socketClosed := make(chan bool)

	go n.collectHostStatistics(hStats)

	rMgr := newRegionManager(config.Node.OpensimBinDir, config.Node.RegionDir)
	rMgr.P

	for {
		conn, err := net.Dial("tcp", config.Node.MGMAddress)
		if err != nil {
			n.logger.Fatal("Cannot connect to MGM")
			time.Sleep(10 * time.Second)
			continue
		}
		n.logger.Info("MGM Node connected to MGM")
		go n.readConnection(conn, mgmCommands, socketClosed)

	ProcessingPackets:
		for {
			select {
			case <-socketClosed:
				n.logger.Error("Disconnected from MGM")
				break ProcessingPackets
			case msg := <-mgmCommands:
				n.logger.Info("recieved message from MGM: ", string(msg))
			case stats := <-hStats:
				nmsg := core.NetworkMessage{}
				nmsg.MessageType = "host_stats"
				nmsg.HStats = stats
				data, _ := json.Marshal(nmsg)

				_, err = conn.Write(data)
				if err == syscall.EPIPE {
					break
				}
				if err != nil {
					n.logger.Error("Error sending data: ", err)
				}
			}
		}
	}

}

func (node mgmNode) readConnection(conn net.Conn, readBytes chan<- []byte, closing chan bool) {
	d := json.NewDecoder(conn)

	for {
		nmsg := core.NetworkMessage{}
		err := d.Decode(&nmsg)
		if err != nil {
			if err.Error() == "EOF" {
				conn.Close()
				return
			}
			node.logger.Error("Error decoding mgm message: ", err)
		}

		switch nmsg.MessageType {
		default:
			node.logger.Info("Received invalid message from an MGM node: ", nmsg.MessageType)
		}
	}
}

func (node mgmNode) writeConnection(conn net.Conn, writeBytes <-chan core.UserObject) {

}

func (node mgmNode) collectHostStatistics(out chan mgm.HostStat) {
	for {
		//start calculating network sent
		fInet, err := psnet.NetIOCounters(false)
		if err != nil {
			node.logger.Error("Error reading networking", err)
		}

		s := mgm.HostStat{}
		c, err := pscpu.CPUPercent(time.Second, true)
		if err != nil {
			node.logger.Error("Error readin CPU: ", err)
		}
		s.CPUPercent = c

		v, err := psmem.VirtualMemory()
		if err != nil {
			node.logger.Error("Error reading Memory", err)
		}
		s.MEMTotal = v.Total / 1000
		s.MEMUsed = (v.Total - v.Available) / 1000
		s.MEMPercent = v.UsedPercent

		lInet, err := psnet.NetIOCounters(false)
		if err != nil {
			node.logger.Error("Error reading networking", err)
		}
		s.NetSent = (lInet[0].BytesSent - fInet[0].BytesSent)
		s.NetRecv = (lInet[0].BytesRecv - fInet[0].BytesRecv)

		out <- s
	}
}
