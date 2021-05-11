package main

import (
	"github.com/patrickz98/project.go.omnetpp/sysinfo"
	"log"
	"time"
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

	// wmic cpu get loadpercentage
	// ps aux
	// ps -A -o user,%cpu,command

	//var avg float64
	avg := newAverage(30)

	for inx := 0; inx < 30*100; inx++ {

		usage := sysinfo.GetCPUUsage()
		log.Printf("usage=%7.3f avg=%7.3f", usage, avg.push(usage))

		time.Sleep(time.Second * 2)
		//break
	}
}
