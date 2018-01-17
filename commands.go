package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"time"

	"github.com/urfave/cli"
)

var commands = []cli.Command{
	quickguide,
}

var quickguide = cli.Command{
	Name:    "quickguide",
	Aliases: []string{"q"},
	Action: func(ctx *cli.Context) error {

		speak("Hello! This is a quick guide to know how to use awsub.")
		speak("First, let's check if you have enough tools on your machine.")

		ng := 0

		if dkm, err := exec.LookPath("docker-machine"); err == nil {
			speak("✔\tThe command `docker-machine` is found in %v", dkm)
		} else {
			ng++
			speak("NG!\tdocker-machine: %v", err)
			speak("\tYou need to install `docker-machine`.")
			speak("\tIt's included in `Docker toolbox`, please go to https://docs.docker.com/toolbox/overview/ to install toolbox!")
		}

		speak("\nThen, let's check credentials so that you can access AWS or any cloud services.")

		if fpath, err := checkAWSCredentials(); err == nil {
			speak("✔\tThe AWS credential file is found at %v", fpath)
		} else {
			ng++
			speak("NG!\tAWS Credentials: %v", err)
		}

		if ng == 0 {
			speak("\nCongrats! It seems you are ready to use `awsub`.")
			speak("For the next step let's try following command.")
			speak("\n\tawsub --tasks ./test/tasks/word-count.csv --script ./test/script/word-count.sh\n")
		}

		return nil
	},
}

func checkAWSCredentials() (string, error) {
	myself, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("can't detect current user: %v", err)
	}
	fpath := filepath.Join(myself.HomeDir, ".aws", "credentials")
	stat, err := os.Stat(fpath)
	if err != nil {
		return "", fmt.Errorf("credential file can't be found at %v: %v", fpath, err)
	}
	if stat.IsDir() {
		return "", fmt.Errorf("credential file path seems to be a directory: %v", fpath)
	}
	return fpath, nil
}

func speak(format string, v ...interface{}) {
	dummysleep := time.Duration(rand.Intn(1000))
	time.Sleep(dummysleep * time.Millisecond)
	format += "\n"
	fmt.Printf(format, v...)
	rand.Seed(time.Now().UnixNano())
}