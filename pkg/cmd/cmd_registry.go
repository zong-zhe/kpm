// Copyright 2023 The KCL Authors. All rights reserved.

package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/moby/term"
	"github.com/urfave/cli/v2"
	"kusionstack.io/kpm/pkg/reporter"
	"kusionstack.io/kpm/pkg/settings"
	"oras.land/oras-go/pkg/auth"
	dockerauth "oras.land/oras-go/pkg/auth/docker"
)

// NewRegistryCmd new a Command for `kpm registry`.
func NewRegCmd(settings *settings.Settings) *cli.Command {
	return &cli.Command{
		Hidden: false,
		Name:   "registry",
		Usage:  "run registry command.",
		Subcommands: []*cli.Command{
			{
				Name:  "login",
				Usage: "login to a registry",
				Flags: []cli.Flag{
					// The registry username.
					&cli.StringFlag{
						Name:    "username",
						Aliases: []string{"u"},
						Usage:   "registry username",
					},
					// The registry registry password or identity token.
					&cli.StringFlag{
						Name:    "password",
						Aliases: []string{"p"},
						Usage:   "registry password or identity token",
					},
					// Read password or identity token from stdin
					&cli.BoolFlag{
						Name:  "password-stdin",
						Value: false,
						Usage: "read password or identity token from stdin",
					},
				},
				Action: func(c *cli.Context) error {
					username, password, err := getUsernamePassword(c.String("username"), c.String("password"), c.Bool("password-stdin"))
					fmt.Printf("username: %v\n", username)
					fmt.Printf("password: %v\n", password)
					if err != nil {
						return err
					}

					if c.NArg() == 0 {
						reporter.Report("kpm: registry must be specified.")
						reporter.ExitWithReport("kpm: run 'kpm registry help' for more information.")
					}
					registry := c.Args().First()

					err = login(registry, username, password, settings)
					if err != nil {
						return err
					}

					return nil
				},
			},
		},
	}
}

func login(hostname, username, password string, setting *settings.Settings) error {

	authClient, err := dockerauth.NewClientWithDockerFallback(setting.CredentialsFile)

	if err != nil {
		return err
	}

	err = authClient.LoginWithOpts(
		[]auth.LoginOption{
			auth.WithLoginHostname(hostname),
			auth.WithLoginUsername(username),
			auth.WithLoginSecret(password),
		}...,
	)

	if err != nil {
		return err
	}

	reporter.Report("Login Succeeded")
	return nil
}

// Adapted from https://github.com/helm/helm
func getUsernamePassword(usernameOpt string, passwordOpt string, passwordFromStdinOpt bool) (string, string, error) {
	var err error
	username := usernameOpt
	password := passwordOpt

	if passwordFromStdinOpt {
		passwordFromStdin, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", "", err
		}
		password = strings.TrimSuffix(string(passwordFromStdin), "\n")
		password = strings.TrimSuffix(password, "\r")
	} else if password == "" {
		if username == "" {
			username, err = readLine("Username: ", false)
			if err != nil {
				return "", "", err
			}
			username = strings.TrimSpace(username)
		}
		if username == "" {
			password, err = readLine("Token: ", true)
			if err != nil {
				return "", "", err
			} else if password == "" {
				return "", "", errors.New("token required")
			}
		} else {
			password, err = readLine("Password: ", true)
			if err != nil {
				return "", "", err
			} else if password == "" {
				return "", "", errors.New("password required")
			}
		}
	} else {
		reporter.Report("kpm: Using --password via the CLI is insecure. Use --password-stdin.")
	}

	return username, password, nil
}

// Copied/adapted from https://github.com/helm/helm
func readLine(prompt string, silent bool) (string, error) {
	fmt.Print(prompt)
	if silent {
		fd := os.Stdin.Fd()
		state, err := term.SaveState(fd)
		if err != nil {
			return "", err
		}
		term.DisableEcho(fd, state)
		defer term.RestoreTerminal(fd, state)
	}

	reader := bufio.NewReader(os.Stdin)
	line, _, err := reader.ReadLine()
	if err != nil {
		return "", err
	}
	if silent {
		fmt.Println()
	}

	return string(line), nil
}
