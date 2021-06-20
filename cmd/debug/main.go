package main

import (
	"github.com/pzierahn/project.go.omnetpp/simple"
	"log"
)

func fib(n float64) float64 {
	if n < 2 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

var result float64

func resultt() int {
	log.Println("resultt")
	return 88
}

func test() (res int) {

	defer func() {
		log.Println("defer")
	}()

	return resultt()
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	freeSlots := uint32(8)
	assign := simple.MathMinUint32(freeSlots, 134)
	freeSlots -= assign

	log.Println(assign, freeSlots, freeSlots-1)

	//test()

	//startTime := time.Now()

	//var r float64
	//
	//for inx := 0; inx < 100; inx++ {
	//	r = fib(30)
	//}
	//
	//result = r
	//
	//log.Printf("Duration: %v", time.Now().Sub(startTime))

	//log.Println(math.Ceil(1.0))
	//log.Printf("0x%x", (int64(0x12345678)<<32) | int64(0xa))

	//test := make(map[string]bool, 1)
	//test["1"] = true
	//test["2"] = true
	//test["3"] = true
	//test["4"] = true

	//ifaces, _ := net.Interfaces()
	//// handle err
	//for _, iface := range ifaces {
	//	addrs, _ := iface.Addrs()
	//
	//	log.Println(iface.Name)
	//	// handle err
	//	for _, addr := range addrs {
	//		//var ip net.IP
	//		//switch v := addr.(type) {
	//		//case *net.IPNet:
	//		//	ip = v.IP
	//		//case *net.IPAddr:
	//		//	ip = v.IP
	//		//}
	//		//// process IP address
	//		log.Printf("addr: %v", addr)
	//	}
	//}

	//localSID := rand.Uint32()
	//remoteSID := rand.Uint32()

	//log.Printf("localSID: %08x --> %04x", localSID, uint16(localSID))
	//log.Printf("remoteSID: %08x --> %04x", remoteSID, uint16(remoteSID))

	//pairTic := time.NewTicker(time.Second * 1)
	//pairTic.

	//var inx int
	//for range pairTic.C {
	//	inx++
	//	pairTic.Reset(time.Second * time.Duration(inx))
	//	log.Printf("tic")
	//}

	//test := make([]bool, 10)
	//
	//fmt.Printf("test: %v\n", test)
	//
	//copy(test[9:], []bool{true, true, true, true})
	//fmt.Printf("test: %v\n", test)

	//for inx := range test {
	//	if rand.Intn(2) == 0 {
	//		test[inx] = true
	//	}
	//
	//	//log.Printf("%v", test[inx])
	//}
	//
	//var buf bytes.Buffer
	//enc := gob.NewEncoder(&buf)
	//if err := enc.Encode(test); err != nil {
	//	panic(err)
	//}
	//
	//log.Printf("buf: %d", buf.Len())
	//
	//var zbuf bytes.Buffer
	//gw := gzip.NewWriter(&zbuf)
	//zenc := gob.NewEncoder(gw)
	//if err := zenc.Encode(test); err != nil {
	//	panic(err)
	//}
	//
	//log.Printf("zenc: %d", zbuf.Len())

	//test := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	//log.Println(test[3:4])

	//log.Println(simple.PrettyString(test))
	//
	//go func() {
	//	arith := new(Arith)
	//	err := rpc.Register(arith)
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//
	//	rpc.HandleHTTP()
	//	l, e := net.Listen("tcp", ":1234")
	//	if e != nil {
	//		log.Fatal("listen error:", e)
	//	}
	//
	//	//rpc.Accept(l)
	//
	//	go http.Serve(l, nil)
	//}()
	//
	//time.Sleep(time.Second * 1)
	//
	//client, err := rpc.DialHTTP("tcp", ":1234")
	//if err != nil {
	//	log.Fatal("dialing:", err)
	//}
	//
	//// Synchronous call
	//args := &Args{7, 8}
	//var reply int
	//err = client.Call("Arith.Multiply", args, &reply)
	//if err != nil {
	//	log.Fatal("arith error:", err)
	//}
	//log.Printf("Arith: %d*%d=%d", args.A, args.B, reply)
}
