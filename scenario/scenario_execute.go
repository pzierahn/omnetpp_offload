package scenario

import (
	"context"
	"fmt"
	"github.com/pzierahn/project.go.omnetpp/consumer"
	"github.com/pzierahn/project.go.omnetpp/gconfig"
	"github.com/pzierahn/project.go.omnetpp/simple"
	"github.com/pzierahn/project.go.omnetpp/stargate"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (runner RunnerRemote) RunScenario(scenario, connect string, trail int) (duration time.Duration, err error) {

	log.Printf("Staring scenario scenarioId=%s connect=%s trailId=%d", scenario, connect, trail)

	// session.Setenv() doesn't work.
	// see https://vic.demuzere.be/articles/environment-variables-setenv-ssh-golang/
	envVars := []string{
		fmt.Sprintf("SCENARIOID=%s", scenario),
		fmt.Sprintf("CONNECT=%s", connect),
		fmt.Sprintf("TRAILID=%d", trail),
	}
	envPrefix := strings.Join(envVars, " ")

	cmd := []string{
		// Source GO paths.
		fmt.Sprintf("source ~/.profile"),
		// Switch to project dir
		fmt.Sprintf("cd %s", runner.sim.OppEdge),
	}

	// Delete deprecated data.
	deleteDirs := []string{
		filepath.Join(runner.sim.Simulation, "opp-edge-results"),
		filepath.Join(runner.sim.Simulation, "results"),
		filepath.Join(runner.sim.Simulation, "out"),
	}

	for _, dir := range deleteDirs {
		cmd = append(cmd, "rm -rf "+dir)
	}

	cmd = append(cmd,
		fmt.Sprintf(
			"%s go run cmd/consumer/opp_edge_run.go -broker %s -path %s -config %s",
			envPrefix,
			runner.sim.Broker,
			runner.sim.Simulation,
			filepath.Join(runner.sim.Simulation, "opp-edge-config.json"),
		),
	)

	log.Printf("bashscript:\n%s\n", strings.Join(cmd, "; \n"))

	session, err := runner.sshClient.NewSession()
	if err != nil {
		err = fmt.Errorf("unable to create a new session: %s", err)
		return
	}
	defer func() {
		_ = session.Close()
	}()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	start := time.Now()
	if err = session.Run(strings.Join(cmd, "; ")); err != nil {
		err = fmt.Errorf("unable to run scenario: %s", err)
		return
	}

	duration = time.Now().Sub(start)

	return
}

func (runner RunnerLocal) RunScenario(scenario, connect string, trail int) (duration time.Duration, err error) {

	log.Printf("Staring scenario scenarioId=%s connect=%s trailId=%d", scenario, connect, trail)

	envVars := map[string]string{
		"SCENARIOID": scenario,
		"CONNECT":    connect,
		"TRAILID":    fmt.Sprint(trail),
	}
	for key, val := range envVars {
		err = os.Setenv(key, val)
		if err != nil {
			return
		}
	}

	// Delete deprecated data.
	deleteDirs := []string{
		filepath.Join(runner.sim.Simulation, "opp-edge-results"),
		filepath.Join(runner.sim.Simulation, "results"),
		filepath.Join(runner.sim.Simulation, "out"),
	}

	for _, dir := range deleteDirs {
		_ = os.RemoveAll(dir)
	}

	var config *consumer.Config
	err = simple.UnmarshallFile(filepath.Join(runner.sim.Simulation, "opp-edge-config.json"), &config)
	if err != nil {
		return
	}

	config.Path = runner.sim.Simulation

	brokerConfig := gconfig.Broker{
		Address:      runner.sim.Broker,
		BrokerPort:   gconfig.DefaultBrokerPort,
		StargatePort: stargate.DefaultPort,
	}

	start := time.Now()
	ctx := context.Background()
	consumer.OffloadSimulation(ctx, brokerConfig, config)
	duration = time.Now().Sub(start)

	return
}
