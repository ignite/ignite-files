package hermes

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/ignite/ignite-files/hermes/data"
	"github.com/ignite/ignite-files/pkg/cmdrunner/exec"
	"github.com/ignite/ignite-files/pkg/localfs"
)

const (
	FlagHostChain        = "host-chain"
	FlagReferenceChain   = "reference-chain"
	FlagChainA           = "a-chain"
	FlagChainB           = "b-chain"
	FlagClientA          = "a-client"
	FlagClientB          = "b-client"
	FlagConnectionA      = "a-connection"
	FlagConnectionB      = "b-connection"
	FlagPortA            = "a-port"
	FlagPortB            = "b-port"
	FlagShowCounterparty = "show-counterparty"
	FlagChain            = "chain"
)

const (
	// CommandCreate is the Hermes create command.
	CommandCreate CommandName = "create"

	// CommandQuery is the Hermes query command.
	CommandQuery CommandName = "query"

	// CommandStart is the Hermes start command.
	CommandStart CommandName = "start"

	// CommandClient is the Hermes client command.
	CommandClient CreateCommand = "client"

	// CommandConnection is the Hermes connection command.
	CommandConnection CreateCommand = "connection"

	// CommandChannel is the Hermes channel command.
	CommandChannel CreateCommand = "channel"

	// CommandChannels  is the Hermes channels command.
	CommandChannels QueryCommand = "channels"
)

type (
	// CommandName represents a high level command under Hermes.
	CommandName string
	// QueryCommand represents the query command under Hermes.
	QueryCommand string
	// CreateCommand represents the create command under Hermes.
	CreateCommand string

	Hermes struct {
		path    string
		binary  []byte
		cleanup func()
	}

	// Option configures Generate configs.
	Option func(*configs)

	// Configs holds Generate configs.
	configs struct {
		flags map[string]interface{}
	}
)

// WithFlags assigns the command flags.
func WithFlags(flags map[string]interface{}) Option {
	return func(c *configs) {
		c.flags = flags
	}
}

// New returns the hermes binary executable.
func New() (*Hermes, error) {
	// untar the binary.
	gzr, err := gzip.NewReader(bytes.NewReader(data.Binary()))
	if err != nil {
		panic(err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	if _, err := tr.Next(); err != nil {
		return nil, err
	}

	binary, err := io.ReadAll(tr)
	if err != nil {
		return nil, err
	}

	path, cleanup, err := localfs.SaveBytesTemp(binary, "hermes", 0o755)
	if err != nil {
		return nil, err
	}

	return &Hermes{
		path:    path,
		binary:  binary,
		cleanup: cleanup,
	}, nil
}

func (h *Hermes) Cleanup() error {
	h.cleanup()
	h.binary = nil
	return os.RemoveAll(h.path)
}

func (h *Hermes) Create(ctx context.Context, cmd CreateCommand, options ...Option) error {
	return h.RunCmd(ctx, CommandCreate, string(cmd), options...)
}

func (h *Hermes) Query(ctx context.Context, cmd QueryCommand, options ...Option) error {
	return h.RunCmd(ctx, CommandQuery, string(cmd), options...)
}

func (h *Hermes) Start(ctx context.Context, options ...Option) error {
	return h.RunCmd(ctx, CommandStart, "", options...)
}

func (h *Hermes) RunCmd(ctx context.Context, command CommandName, subCommand string, options ...Option) error {
	c := configs{}
	for _, o := range options {
		o(&c)
	}

	cmd := []string{h.path, string(command)}
	if subCommand != "" {
		cmd = append(cmd, subCommand)
	}
	for flag, value := range c.flags {
		if _, ok := value.(bool); ok {
			cmd = append(cmd, fmt.Sprintf("--%s", flag))
		} else {
			cmd = append(cmd, fmt.Sprintf("--%s=%s", flag, value))
		}
	}

	// execute the command.
	return exec.Exec(ctx, cmd, exec.IncludeStdLogsToError())
}
