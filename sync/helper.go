package sync

import (
	"os"
	"io/ioutil"
	"log"
	"strings"
	"github.com/webdevops/go-shell"
	"fmt"
)

func NewShellCommand(command string, args ...string) *shell.Command {
	return shell.Cmd(ShellCommandInterfaceBuilder(command, args...)...)
}

func ShellCommandInterfaceBuilder(command string, args ...string) []interface{} {
	cmd := []string{
		command,
	}

	for _, val := range args {
		cmd = append(cmd, shell.Quote(val))
	}

	shellCmd := make([]interface{}, len(cmd))
	for i, v := range cmd {
		shellCmd[i] = v
	}

	return shellCmd
}

func DockerGetContainerId(container string) string {
	if strings.HasPrefix(container, "compose:") {
		containerName := strings.TrimPrefix(container, "compose:")

		cmd := shell.Cmd("docker-compose", "ps", "-q", containerName).Run()
		containerId := strings.TrimSpace(cmd.Stdout.String())

		return containerId
	}

	return container
}

func CreateTempfile() *os.File {
	tmpfile, err := ioutil.TempFile("", "gsync")
	if err != nil {
		log.Fatal(err)
	}

	return tmpfile
}

func CreateTempfileWithContent(content ...string) *os.File {
	tmpfile := CreateTempfile()

	if _, err := tmpfile.Write([]byte(strings.Join(content[:],"\n"))); err != nil {
		log.Fatal(err)
	}

	return tmpfile
}

func PathExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func FileExists(name string) bool {
	f, err := os.Stat(name);

	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	if f.IsDir() {
		return false
	}

	return true
}

func RsyncPath(name string) string {
	return strings.TrimRight(name, "/") + "/"
}

func ShellErrorHandler(recover interface{}) {
	process, ok := recover.(*shell.Process)
	if ok {
		p := process.ExitStatus
		p = 2
		if p != 0 {

			printMessage := func(header string, body string) {
				fmt.Println(header)
				fmt.Println(strings.Repeat("-", len(header)))
				fmt.Println("   " + strings.Replace(body, "\n", "\n   ", -1))
				fmt.Println()
			}

			fmt.Println("\n\n[!!!] Command execution failed")
			fmt.Println()

			printMessage("Command", process.Command.ToString())
			printMessage("Stdout", process.Stdout.String())
			printMessage("Stderr", process.Stderr.String())
			printMessage("Exit code", fmt.Sprintf("%d", p))
		}
	}
}
