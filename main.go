package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

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
	ReleaseBlog(dirname)
	fmt.Println("Done.")
}

func main() {
	server := hookserve.NewServer()
	server.Port = strcov.Atoi(os.Getenv("IONBLOG_PORT"))
	server.Secret = os.Getenv("IONBLOG_SECRET")
	server.GoListenAndServe()

	for {
		select {
		case event := <-server.Events:
			fmt.Println(event.Owner + " " + event.Repo + " " + event.Branch + " " + event.Commit)
			if event.Branch == "master" {
				PublishBlog()
			}
		default:
			time.Sleep(100)
		}
	}
}
