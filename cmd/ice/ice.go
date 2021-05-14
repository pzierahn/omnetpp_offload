package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pion/randutil"
	"net/http"
	"time"

	"github.com/pion/ice/v2"
)

type ICEInfo struct {
	LocalUfrag    string
	LocalPwd      string
	Candidates    []string
	IsControlling bool
}

type Controlling struct {
	IsControlling bool
}

var (
	isControlling bool
	iceAgent      *ice.Agent
)

func main() {
	var err error

	flag.BoolVar(&isControlling, "controlling", false, "is ICE Agent controlling")
	flag.Parse()

	//if isControlling {
	//	fmt.Println("Local Agent is controlling")
	//} else {
	//	fmt.Println("Local Agent is controlled")
	//}
	//fmt.Print("Press 'Enter' when both processes have started")
	//if _, err = bufio.NewReader(os.Stdin).ReadBytes('\n'); err != nil {
	//	panic(err)
	//}

	iceAgent, err = ice.NewAgent(&ice.AgentConfig{
		PortMax: 51088,
		PortMin: 51088,
		NetworkTypes: []ice.NetworkType{
			//ice.NetworkTypeTCP4,
			//ice.NetworkTypeTCP6,
			ice.NetworkTypeUDP4,
			ice.NetworkTypeUDP6,
		},
		NAT1To1IPs: []string{
			"31.18.129.212",
			//"2a02:8108:3cbf:f718:9dee:dd20:e44:58e3",
		},
	})
	if err != nil {
		panic(err)
	}

	var candidates ICEInfo
	candidates.IsControlling = isControlling

	// When we have gathered a new ICE Candidate send it to the remote peer
	err = iceAgent.OnCandidate(func(c ice.Candidate) {
		if c == nil {
			return
		}

		fmt.Printf("######## OnCandidate: '%s' --> '%s'\n", c.String(), c.Marshal())

		candidates.Candidates = append(candidates.Candidates, c.Marshal())
	})

	if err != nil {
		panic(err)
	}

	// When ICE Connection state has change print to stdout
	err = iceAgent.OnConnectionStateChange(func(c ice.ConnectionState) {
		fmt.Printf("ICE Connection State has changed: %s\n", c.String())
	})

	if err != nil {
		panic(err)
	}

	// Get the local auth details and send to remote peer
	localUfrag, localPwd, err := iceAgent.GetLocalUserCredentials()
	if err != nil {
		panic(err)
	}

	fmt.Println("localUfrag:", localUfrag)
	fmt.Println("localPwd:", localPwd)

	candidates.LocalUfrag = localUfrag
	candidates.LocalPwd = localPwd

	if err = iceAgent.GatherCandidates(); err != nil {
		panic(err)
	}

	fmt.Println("GatherCandidates done")

	time.Sleep(time.Second * 3)

	jbyt, _ := json.MarshalIndent(candidates, "", "  ")
	fmt.Println(string(jbyt))

	_, err = http.Post(
		"https://8ca70b82a4b0.ngrok.io/candidate",
		"application/json",
		bytes.NewReader(jbyt))

	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 4)

	other, _ := json.MarshalIndent(Controlling{IsControlling: isControlling}, "", "  ")
	resp, err := http.Post(
		"https://8ca70b82a4b0.ngrok.io/exchange",
		"application/json",
		bytes.NewReader(other))
	if err != nil {
		panic(err)
	}

	var remote ICEInfo
	err = json.NewDecoder(resp.Body).Decode(&remote)
	if err != nil {
		panic(err)
	}

	fmt.Println("exchange: ", remote)

	//remoteUfrag := <-remoteAuthChannel
	//remotePwd := <-remoteAuthChannel

	for _, itm := range remote.Candidates {
		can, err := ice.UnmarshalCandidate(itm)
		if err != nil {
			panic(err)
		}

		err = iceAgent.AddRemoteCandidate(can)
		if err != nil {
			panic(err)
		}
	}

	var conn *ice.Conn

	// Start the ICE Agent. One side must be controlled, and the other must be controlling
	if isControlling {
		conn, err = iceAgent.Dial(context.Background(), remote.LocalUfrag, remote.LocalPwd)
	} else {
		conn, err = iceAgent.Accept(context.Background(), remote.LocalUfrag, remote.LocalPwd)
	}
	if err != nil {
		panic(err)
	}

	fmt.Println("RemoteAddr:", conn.RemoteAddr())

	//server := grpc.NewServer()
	//pb.RegisterBrokerServer(server, &brk)
	//pb.RegisterStorageServer(server, &storage.Server{})
	//err = server.Serve(conn)

	// Send messages in a loop to the remote peer
	go func() {
		for {
			time.Sleep(time.Second * 3)

			val, err := randutil.GenerateCryptoRandomString(15, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
			if err != nil {
				panic(err)
			}
			if _, err = conn.Write([]byte(val)); err != nil {
				panic(err)
			}

			fmt.Printf("Sent: '%s'\n", val)
		}
	}()

	// Receive messages in a loop from the remote peer
	buf := make([]byte, 1500)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Received: '%s'\n", string(buf[:n]))
	}
}
