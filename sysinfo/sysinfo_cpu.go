package sysinfo

import (
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func getCPUUsageWindows() (usage float64) {

	// wmic cpu get loadpercentage
	cmd := exec.Command("wmic", "cpu", "get", "loadpercentage")
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}

	str := strings.TrimSpace(string(out))
	parts := strings.Split(str, "\n")

	percent, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(err)
	}

	usage = float64(percent)

	return
}

func getCPUUsageUnix() (usage float64) {

	cmd := exec.Command("ps", "-A", "-o", "%cpu")
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}

	str := strings.TrimSpace(string(out))
	parts := strings.Split(str, "\n")

	for _, part := range parts[1:] {
		us, err := strconv.ParseFloat(strings.TrimSpace(part), 32)
		if err != nil {
			panic(err)
		}

		usage += us
	}

	return
}

func GetCPUUsage() (usage float64) {
	if runtime.GOOS == "windows" {
		usage = getCPUUsageWindows()
	} else {
		usage = getCPUUsageUnix() / float64(runtime.NumCPU())
	}

	return
}
