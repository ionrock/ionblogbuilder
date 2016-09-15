package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/phayes/hookserve/hookserve"
	"github.com/urfave/cli"
)

func Do(parts ...string) *exec.Cmd {
	fmt.Printf("Running command: %s\n", parts)
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
		log.Print("Error checking out blog")
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
	log.Print("Checking out blog")
	dirname := CheckoutBlog()
	log.Print("Building the blog")
	BuildBlog(dirname)
	// ReleaseBlog(dirname)
	log.Print("Done.")
}

func ServeHook(port int, secret string) {
	server := hookserve.NewServer()
	server.Port = port
	server.Secret = secret
	log.Printf("Starting hook server on port %d", server.Port)
	server.ListenAndServe()

	for {
		select {
		case event := <-server.Events:
			log.Print(event.Owner + " " + event.Repo + " " + event.Branch + " " + event.Commit)
			if event.Branch == "master" {
				PublishBlog()
			}
		default:
			time.Sleep(5)
		}
	}
}

func ServeSite(port int) {
	r := http.NewServeMux()
	fs := http.FileServer(http.Dir("ionblog/blog/html"))
	r.Handle("/", fs)

	handler := handlers.LoggingHandler(os.Stdout, r)

	log.Printf("Webserver Listening on port %d", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
}

func IonblogServerAction(c *cli.Context) error {

	if c.String("secret") == "" {
		log.Print("A secret value is required!")
		return errors.New("A secret value is required!")
	}

	// PublishBlog()
	go ServeHook(c.Int("port"), c.String("secret"))
	go ServeSite(c.Int("webport"))

	Listen()
	return nil
}

func Listen() {
	SigQuit := make(chan os.Signal)
	signal.Notify(SigQuit, syscall.SIGINT, syscall.SIGTERM)
	SigStat := make(chan os.Signal)
	signal.Notify(SigStat, syscall.SIGUSR1)

forever:
	for {
		select {
		case s := <-SigQuit:
			log.Printf("Signal (%d) received, stopping", s)
			break forever
		}
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "ionblogbuilder"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "webport, P",
			Value: 80,
			Usage: "The webserver port",
		},
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
		cli.StringFlag{
			Name:   "output, o",
			Usage:  "Output directory",
			EnvVar: "IONBLOG_OUTPUT",
		},
	}

	app.Action = IonblogServerAction

	app.Run(os.Args)
}
