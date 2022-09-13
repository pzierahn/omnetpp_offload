package scenario

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
	"path/filepath"
)

func connectSSH(addr string) (client *ssh.Client, err error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	key, err := os.ReadFile(filepath.Join(home, "/.ssh/id_rsa"))
	if err != nil {
		err = fmt.Errorf("unable to read private key: %v", err)
		return
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		err = fmt.Errorf("unable to parse private key: %v", err)
		return
	}

	config := &ssh.ClientConfig{
		User: "fioo",
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
			//ssh.Password("PASSWORD"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return ssh.Dial("tcp", addr, config)
}
