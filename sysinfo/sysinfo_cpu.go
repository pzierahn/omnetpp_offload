package sysinfo

import (
	"context"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func getCPUUsageWindows(ctx context.Context) (usage float64) {

	// wmic cpu get loadpercentage
	cmd := exec.CommandContext(ctx, "wmic", "cpu", "get", "loadpercentage")
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}

	str := strings.TrimSpace(string(out))
	parts := strings.Split(str, "\n")

	if len(parts) < 2 {
		return
	}

	percent, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(err)
	}

	usage = float64(percent)

	return
}

func getCPUUsageUnix(ctx context.Context) (usage float64) {

	cmd := exec.CommandContext(ctx, "ps", "-A", "-o", "%cpu")
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

	usage /= float64(runtime.NumCPU())

	return
}

func GetCPUUsage(ctx context.Context) (usage float64) {
	if runtime.GOOS == "windows" {
		usage = getCPUUsageWindows(ctx)
	} else {
		usage = getCPUUsageUnix(ctx)
	}

	return
}
