package main

import (
	"fmt"
	"github.com/patrickz98/project.go.omnetpp/simple"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

type rollingAverage struct {
	points []float64
	avg    float64
	index  int
}

func (rol *rollingAverage) push(elem float64) (avg float64) {

	del := rol.points[rol.index]
	rol.points[rol.index] = elem

	rol.avg = (rol.avg*float64(len(rol.points)) - del + elem) / float64(len(rol.points))

	rol.index = (rol.index + 1) % len(rol.points)

	return rol.avg
}

func newAverage(size int) (avg rollingAverage) {

	avg.points = make([]float64, size)

	return
}

func main() {

	memo, err := mem.VirtualMemory()
	if err != nil {
		panic(err)
	}
	fmt.Println("mem", simple.PrettyString(memo))

	stats, err := cpu.Info()
	if err != nil {
		panic(err)
	}

	fmt.Println("stats:", simple.PrettyString(stats))

	hostInfo, err := host.Info()
	if err != nil {
		panic(err)
	}
	fmt.Println("hostInfo:", simple.PrettyString(hostInfo))

	avg, err := load.Avg()
	if err != nil {
		panic(err)
	}
	fmt.Println("avg:", simple.PrettyString(avg))

	//log.Println("OS:", runtime.GOOS)
	//log.Println("ARCH:", runtime.GOARCH)
	//log.Println("NumCPU:", runtime.NumCPU())
	//
	//avg := newAverage(30)
	//
	//for inx := 0; inx < 1000000; inx++ {
	//
	//	usage := sysinfo.GetCPUUsage()
	//	log.Printf("usage=%7.3f idle=%7.3f avg=%7.3f", usage, 100-usage, avg.push(100-usage))
	//
	//	time.Sleep(time.Millisecond * 500)
	//	//break
	//}
}
