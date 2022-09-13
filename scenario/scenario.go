package scenario

import (
	"fmt"
	"github.com/docker/docker/client"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"time"
)

type Simulation struct {
	Broker         string
	ClientSSHAddr  string
	OppEdgePath    string
	SimulationPath string
}

type Runner interface {
	RunScenario(scenario, connect string, trail int) (duration time.Duration, err error)
}

type RunnerRemote struct {
	sim       Simulation
	docker    *client.Client
	sshClient *ssh.Client
}

type RunnerLocal struct {
	sim Simulation
}

func (runner RunnerRemote) UpdateRepo() {

	log.Printf("Updating repository %s", runner.sim.OppEdgePath)

	session, err := runner.sshClient.NewSession()
	if err != nil {
		log.Fatalf("unable to create a new session: %s", err)
	}
	defer func() {
		_ = session.Close()
	}()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	if err := session.Run(fmt.Sprintf("cd %s; git pull", runner.sim.OppEdgePath)); err != nil {
		log.Fatalf("unable to update project: %s", err)
	}
}

func NewScenarioRemote(sim Simulation) Runner {

	var runner RunnerRemote
	runner.sim = sim

	var err error
	runner.sshClient, err = connectSSH(sim.ClientSSHAddr)
	if err != nil {
		log.Fatalf("unable to connect ssh client: %s", err)
	}

	runner.UpdateRepo()

	return runner
}

func NewScenario(sim Simulation) Runner {
	var runner RunnerLocal
	runner.sim = sim

	return runner
}
