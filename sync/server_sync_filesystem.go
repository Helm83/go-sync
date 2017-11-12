package sync

import (
	"os"
	"fmt"
	"errors"
	"strings"
	"github.com/webdevops/go-shell"
	"github.com/webdevops/go-shell/commandbuilder"
)

// General sync
func (filesystem *Filesystem) Sync() {
	switch filesystem.Connection.GetInstance().GetType() {
	case "local":
		fallthrough
	case "ssh":
		filesystem.syncRsync()
	case "docker":
		errors.New("Docker not supported")
	}
}

// Sync filesystem using rsync
func (filesystem *Filesystem) syncRsync() {
	connection := filesystem.Connection.GetInstance()

	args := []string{"-rlptD", "--delete-after", "--progress", "--human-readable"}

	// add custom options
	if filesystem.Options.Rsync != nil {
		args = append(args, filesystem.Options.Rsync.Array()...)
	}

	if filesystem.Connection.GetInstance().IsSsh() {
		args = append(args, "-e", shell.Quote("ssh " + strings.Join(commandbuilder.ConnectionSshArguments, " ")))
	}

	// include filter
	if len(filesystem.Filter.Include) > 0 {
		includeTempFile := CreateTempfileWithContent(filesystem.Filter.Include...)
		args = append(args, fmt.Sprintf("--files-from=%s", includeTempFile.Name()))

		// remove file after run
		defer os.Remove(includeTempFile.Name())
	}

	// exclude filter
	if len(filesystem.Filter.Exclude) > 0 {
		excludeTempFile := CreateTempfileWithContent(filesystem.Filter.Exclude...)
		args = append(args, fmt.Sprintf("--exclude-from=%s", excludeTempFile.Name()))

		// remove file after run
		defer os.Remove(excludeTempFile.Name())
	}

	// build source and target paths
	sourcePath := ""
	switch connection.GetType() {
	case "ssh":
		sourcePath = fmt.Sprintf("%s:%s", connection.SshConnectionHostnameString(), filesystem.Path)
	case "local":
		sourcePath = filesystem.Path
	}

	targetPath := filesystem.localPath()

	// make sure source/target paths are using suffix slash
	args = append(args, RsyncPath(sourcePath), RsyncPath(targetPath))

	cmd := shell.NewCmd("rsync", args...)
	cmd.Run()
}
