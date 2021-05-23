package main

import (
	"errors"
	"fmt"
	"log"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//log.Println(math.Ceil(1.0))
	//log.Printf("0x%x", (int64(0x12345678)<<32) | int64(0xa))

	//test := make(map[string]bool, 1)
	//test["1"] = true
	//test["2"] = true
	//test["3"] = true
	//test["4"] = true

	for inx := 0; inx < 10; inx++ {
		fmt.Println(inx)
	}

	fmt.Println("---------------")

	for inx := 10 - 1; inx >= 0; inx-- {
		fmt.Println(inx)
	}

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
