package svn

import (
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
)

type ChangedFile struct {
	Path string `xml:"path,attr"`
}

type Changelist struct {
	Name  string        `xml:"name,attr"`
	Files []ChangedFile `xml:"entry"`
}

type SVNStatus struct {
	Changelists []Changelist `xml:"changelist"`
}

func GetStatus() (*SVNStatus, error) {
	statusCommand := exec.Command("svn", "status", "--xml")
	statusCommand.Stderr = os.Stderr
	statusCommand.Stdin = os.Stdin
	outPipe, err := statusCommand.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("error getting status out pipe: %w", err)
	}
	defer outPipe.Close()

	if err := statusCommand.Start(); err != nil {
		return nil, fmt.Errorf("error starting status command: %w", err)
	}
	xmlDecoder := xml.NewDecoder(outPipe)
	var status SVNStatus
	if err := xmlDecoder.Decode(&status); err != nil {
		return nil, fmt.Errorf("error decoding xml: %w", err)
	}

	return &status, nil
}
