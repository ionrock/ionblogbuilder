package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/codegangsta/cli"
	"github.com/phayes/hookserve/hookserve"
)

func Do(parts ...string) *exec.Cmd {
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func CheckoutBlog() string {
	repo := "ionblog"
	endpoint := fmt.Sprintf("https://github.com/ionrock/%s", repo)
	action := "clone"

	if _, err := os.Stat(repo); !os.IsNotExist(err) {
		action = "pull"
	}
	cmd := Do("git", action, endpoint)

	if action == "pull" {
		cmd.Dir = repo
	}

	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	return repo
}

func BuildBlog(blogdir string) {
	cmd := Do("make", "build")
	cmd.Dir = blogdir
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func ReleaseBlog(blogdir string) {
	cmd := Do("make", "release")
	cmd.Dir = blogdir
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func PublishBlog() {
	dirname := CheckoutBlog()
	BuildBlog(dirname)
	fmt.Println("Done.")
}

func Serve(port int, secret string) {
	server := hookserve.NewServer()
	server.Port = port
	server.Secret = secret
	server.GoListenAndServe()

	for {
		select {
		case event := <-server.Events:
			fmt.Println(event.Owner + " " + event.Repo + " " + event.Branch + " " + event.Commit)
			if commit.Branch == "master" {
				PublishBlog()
			}
		default:
			time.Sleep(5)
		}
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "ionblogbuilder"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:   "port, p",
			Value:  5566,
			Usage:  "The port github will use to send webhooks",
			EnvVar: "PORT,IONBLOG_PORT",
		},
		cli.StringFlag{
			Name:   "secret, s",
			Usage:  "Github secret key",
			EnvVar: "IONBLOG_SECRET",
		},
	}

	app.Action = func(c *cli.Context) error {

		if c.String("secret") == "" {
			fmt.Println("A secret value is required!")
			return errors.New("A secret value is required!")
		}

		Serve(c.Int("port"), c.String("secret"))
		return nil
	}

	app.Run(os.Args)
}
