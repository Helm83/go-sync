package sync

import (
	"fmt"
	"errors"
	"os"
)

// General sync
func (filesystem *Filesystem) Deploy(server *Server) {
	switch server.Connection.GetType() {
	case "ssh":
		filesystem.deployRsync(server)
	case "docker":
		errors.New("Docker not supported")
	}
}

// Sync filesystem using rsync
func (filesystem *Filesystem) deployRsync(server *Server) {
	args := []string{"-rlptD", "--delete-after", "--progress", "--human-readable"}

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
	sourcePath := filesystem.localPath(server)
	targetPath := fmt.Sprintf("%s:%s", server.Connection.SshConnectionHostnameString(), filesystem.Path)

	// make sure source/target paths are using suffix slash
	args = append(args, RsyncPath(sourcePath), RsyncPath(targetPath))

	cmd := NewShellCommand("rsync", args...)
	cmd.Run()
}
