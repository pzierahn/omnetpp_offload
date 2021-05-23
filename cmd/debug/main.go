package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"
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

	go func() {
		arith := new(Arith)
		err := rpc.Register(arith)
		if err != nil {
			log.Fatalln(err)
		}

		rpc.HandleHTTP()
		l, e := net.Listen("tcp", ":1234")
		if e != nil {
			log.Fatal("listen error:", e)
		}

		//rpc.Accept(l)

		go http.Serve(l, nil)
	}()

	time.Sleep(time.Second * 1)

	client, err := rpc.DialHTTP("tcp", ":1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	// Synchronous call
	args := &Args{7, 8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	log.Printf("Arith: %d*%d=%d", args.A, args.B, reply)
}
