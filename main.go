package main

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"os"
)

func Auth() {
	cmd := &cli.Command{
		Name: "password-manager",
		Commands: []*cli.Command{
			{
				Name:  "auth",
				Usage: "Create an account, change your password and also login by using subcommands given in this command",
				Commands: []*cli.Command{
					{
						Name: "create",
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "username", Required: true, Usage: "Your username"},
							&cli.StringFlag{Name: "password", Required: true, Usage: "Your password"},
							&cli.StringFlag{Name: "email", Required: true, Usage: "Your email"},
						},
						Action: CreateUser,
					},
					{
						Name: "login",
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "username", Required: true, Usage: "Your username"},
							&cli.StringFlag{Name: "password", Required: true, Usage: "Your password"},
						},
						Action: LoginUser,
					},
				},
			},
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func main() {
	Auth()
}
