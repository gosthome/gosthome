package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	clive "github.com/ASMfreaK/clive2"
	_ "github.com/gosthome/gosthome/components"
	"github.com/gosthome/gosthome/core"
	"github.com/gosthome/gosthome/core/config"
	"github.com/urfave/cli/v2"
)

type app struct {
	*clive.Command `cli:"usage:'Control your (embedded) systems by simple yet powerful configuration files and remotely through encrypted API'"`
	Verbose        bool
	Subcommands    struct {
		*Run
		*Util
	}
}

func (a *app) Before(ctx *cli.Context) error {
	if a.Verbose {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})))
	}
	return nil
}

func (*app) Version() string {
	c := core.Commit()
	if c != "" {
		c = "-" + c
	}
	return core.Version() + c
}

type Run struct {
	*clive.Command `cli:"usage:'Run a configuration'"`

	Config string `cli:"usage:'config file to read',required"`
}

func (r *Run) Action(ctx *cli.Context) error {

	f, err := os.Open(r.Config)
	if err != nil {
		return fmt.Errorf("error reading file %w", err)
	}
	cfg, err := config.LoadConfig(f)
	if err != nil {
		return fmt.Errorf("error loading configuration from %s: %w", r.Config, err)
	}
	n, err := core.NewNode(ctx.Context, cfg)
	if err != nil {
		return fmt.Errorf("error initalizing node: %w", err)
	}
	defer func() {
		err := n.Close()
		if err != nil {
			slog.Error("Error stopping node", "err", err)
		}
	}()
	n.Start()
	<-ctx.Context.Done()
	return nil
}

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	err := clive.Build(&app{}).RunContext(ctx, os.Args)
	if err != nil {
		println(err.Error())
	}
}
