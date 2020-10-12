package client

import (
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/xerrors"
)

// Hope to implement using github.com/goodwithtech/dockle/pkg

type DockleClient struct{}

func (c *DockleClient) Do(ctx context.Context, image string) ([]byte, error) {
	tmpfile, err := ioutil.TempFile("", "*.json")
	if err != nil {
		return nil, xerrors.Errorf("failed to create tmpfile: %w", err)
	}
	filename := tmpfile.Name()

	defer tmpfile.Close()
	defer os.Remove(filename)

	if _, err := exec.CommandContext(ctx, "dockle", "-o", filename, "-f", "json", image).CombinedOutput(); err != nil {
		return nil, xerrors.Errorf("failed to execute dockle: %w", err)
	}
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, xerrors.Errorf("failed to read tmpfile: %w", err)
	}
	return body, nil
}

type DockleResponse struct {
	Target  string
	Summary DockleSummary  `json:"summary"`
	Details []DockleDetail `json:"details"`
}

func (dr *DockleResponse) ExtractImage() string {
	return strings.Split(dr.Target, " ")[0]
}

type DockleSummary struct {
	Fatal int `json:"fatal"`
	Warn  int `json:"warn"`
	Info  int `json:"info"`
	Skip  int `json:"skip"`
	Pass  int `json:"pass"`
}

type DockleDetail struct {
	Code   string   `json:"code"`
	Title  string   `json:"title"`
	Level  string   `json:"level"`
	Alerts []string `json:"alerts"`
}
