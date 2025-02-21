package main

import (
	"bytes"
	"errors"
	"fmt"
	"syscall"

	clive "github.com/ASMfreaK/clive2"
	"github.com/gosthome/gosthome/components/api/frameshakers"
	"github.com/gosthome/gosthome/core/config"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

type Util struct {
	*clive.Command `cli:"usage:'Utilities for configuration'"`

	Subcommands struct {
		*Mac
		*Noise

		*HashPassword
	}
}

type Mac struct {
	*clive.Command `cli:"usage:'Generate a local MAC address to identify node'"`
}

func (*Mac) Action(ctx *cli.Context) error {
	m, err := config.GenerateMAC()
	if err != nil {
		return err
	}
	println(m.String())
	return nil
}

type Noise struct {
	*clive.Command `cli:"usage:'Generate a Noise PSK for encryption key'"`
}

func (*Noise) Action(ctx *cli.Context) error {
	n, err := frameshakers.GenerateEncryptionKey()
	if err != nil {
		return err
	}
	println(n.String())
	return nil
}

type HashPassword struct {
	*clive.Command `cli:"usage:'Hash your password for the config file'"`
}

func (*HashPassword) Action(ctx *cli.Context) error {
	fmt.Print("Password: ")
	bytepw, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return err
	}
	fmt.Print("\nConfirm password: ")
	confirm, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return err
	}
	fmt.Print("\n")
	if !bytes.Equal(bytepw, confirm) {
		return errors.New("passwords do not match")
	}
	hash, err := bcrypt.GenerateFromPassword(bytepw, 10)
	if err != nil {
		return err
	}
	fmt.Printf("Password for your config: %q\n", string(hash))
	return nil
}
