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
						Usage:  "Create an account if you do not already have on by using this command. usage: create --useraname user_name --password your_password --email your@mail.com",
						Action: CreateUser,
					},
					{
						Name: "login",
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "username", Required: true, Usage: "Your username"},
							&cli.StringFlag{Name: "password", Required: true, Usage: "Your password"},
						},
						Usage:  "login into an already created account. usage: create --useraname user_name --password your_password",
						Action: LoginUser,
					},

					{
						Name:   "logout",
						Usage:  "Used to delete your keys stored in the cli. Use this to logout so that you can refresh tokens or create another account or do anything a logged out user could do",
						Action: Logout,
					},
				},
			},
			{
				Name:   "create",
				Usage:  "Create a key to be stored in the database. We store your keys fully encrypted in our database, possible to be decrypted soley by password in your posession unknown to any other sould",
				Action: CreateKey,
			},
			{
				Name:   "list",
				Usage:  "List all the stored keys so you can see which one exist in the database",
				Action: ListKeys,
			},
			{
				Name:   "get",
				Usage:  "get the key (decrypted). usage: get service_name ",
				Action: GetKey,
			},
			{
				Name:   "update",
				Usage:  "update the key stored in the database with the new one. usage: update service_name new_key ",
				Action: UpdateKey,
			},
			{
				Name:   "delete",
				Usage:  "delete the key stored in the database. usage: delete service_name",
				Action: DeleteKey,
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
