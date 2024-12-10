package commands

import (
	"fmt"
	"net"
	"sync"

	"github.com/KirantiLoh/ssh-automator/internal/model"
	"github.com/KirantiLoh/ssh-automator/internal/utils"
	"golang.org/x/crypto/ssh"
)

func RunCommandsSSH(server model.Server, defaultConfig model.DefaultConfig, updateChan chan model.Update, wg *sync.WaitGroup) {
	defer wg.Done()

	updateChan <- model.Update{
		Host:       server.IP,
		Message:    "Trying to establish SSH connection...",
		IsComplete: false,
	}

	username := defaultConfig.Username
	if server.Username != "" {
		username = server.Username
	}

	authMethod, err := utils.PublicKeyFile(defaultConfig.IdentityFile)
	if err != nil {
		updateChan <- model.Update{
			Host:       server.IP,
			Message:    fmt.Sprintf("%s", err.Error()),
			IsComplete: false,
			IsError:    true,
		}
		return
	}
	sshConfig := &ssh.ClientConfig{
		User:            username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			authMethod,
		},
	}

	port := defaultConfig.Port
	if server.Port != "" {
		port = server.Port
	}

	sshAddress := net.JoinHostPort(server.IP, port)

	conn, err := ssh.Dial("tcp", sshAddress, sshConfig)

	if err != nil {
		updateChan <- model.Update{
			Host:       server.IP,
			Message:    fmt.Sprintf("Failed to create client. %s", err.Error()),
			IsComplete: false,
			IsError:    true,
		}
		return
	}

	defer conn.Close()

	session, err := conn.NewSession()

	if err != nil {
		updateChan <- model.Update{
			Host:       server.IP,
			Message:    fmt.Sprintf("Failed to create session. %s", err.Error()),
			IsComplete: false,
			IsError:    true,
		}
		return
	}
	defer session.Close()

	updateChan <- model.Update{
		Host:       server.IP,
		Message:    "SSH successfull!",
		IsComplete: false,
	}

	stdin, err := session.StdinPipe()

	if err != nil {
		updateChan <- model.Update{
			Host:       server.IP,
			Message:    fmt.Sprintf("Failed to create stdin pipe. %s", err.Error()),
			IsComplete: false,
			IsError:    true,
		}
		return
	}

	if err := session.Shell(); err != nil {
		updateChan <- model.Update{
			Host:       server.IP,
			Message:    fmt.Sprintf("Failed to start shell. %s", err.Error()),
			IsComplete: false,
			IsError:    true,
		}
		return
	}

	for _, command := range server.Commands {
		updateChan <- model.Update{
			Host:       server.IP,
			Message:    fmt.Sprintf("Running command: %s", command),
			IsComplete: false,
		}
		fmt.Fprintln(stdin, command)
	}
	fmt.Fprintln(stdin, "exit")

	if err := session.Wait(); err != nil {
		updateChan <- model.Update{
			Host:       server.IP,
			Message:    fmt.Sprintf("Commands didn't run smoothly, try again. %s", err.Error()),
			IsComplete: false,
			IsError:    true,
		}
		return
	}
	updateChan <- model.Update{
		Host:       server.IP,
		Message:    fmt.Sprintf("Commands for %s executed successfully!", server.IP),
		IsComplete: true,
	}

}
