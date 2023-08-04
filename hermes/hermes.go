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
		flags map[string]string
	}
)

// WithFlags assigns the command flags.
func WithFlags(flags map[string]string) Option {
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

// hermes create client --host-chain ibc-1 --reference-chain ibc-0
// hermes create client --host-chain ibc-0 --reference-chain ibc-1
// hermes create connection --a-chain ibc-0 --a-client 07-tendermint-0 --b-client 07-tendermint-0
// hermes create channel --a-chain ibc-0 --a-connection connection-0 --a-port transfer --b-port transfer
// hermes query channels --show-counterparty --chain ibc-0
// hermes start
// hermes tx ft-transfer --timeout-seconds 1000 --dst-chain ibc-1 --src-chain ibc-0 --src-port transfer --src-channel channel-0 --amount 100000
// hermes tx ft-transfer --timeout-seconds 10000 --denom ibc/C1840BD16FCFA8F421DAA0DAAB08B9C323FC7685D0D7951DC37B3F9ECB08A199 --dst-chain ibc-0 --src-chain ibc-1 --src-port transfer --src-channel channel-0 --amount 100000

func (h *Hermes) Create(ctx context.Context, cmd CreateCommand, options ...Option) error {
	return h.RunCmd(ctx, CommandCreate, string(cmd), options...)
}

func (h *Hermes) Query(ctx context.Context, cmd QueryCommand, options ...Option) error {
	return h.RunCmd(ctx, CommandQuery, string(cmd), options...)
}

func (h *Hermes) RunCmd(ctx context.Context, command CommandName, subCommand string, options ...Option) error {
	c := configs{}
	for _, o := range options {
		o(&c)
	}

	cmd := []string{h.path, string(command), subCommand}
	for flag, value := range c.flags {
		cmd = append(cmd, fmt.Sprintf("--%s=%s", flag, value))
	}

	// execute the command.
	return exec.Exec(ctx, cmd, exec.IncludeStdLogsToError())
}
