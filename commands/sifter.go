package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type ConsulEvent struct {
	Id            string `json:"ID"`
	Name          string `json:"Name"`
	Payload       string `json:"Payload,omitempty"`
	NodeFilter    string `json:"NodeFilter,omitempty"`
	ServiceFilter string `json:"ServiceFilter"`
	TagFilter     string `json:"TagFilter"`
	Version       int    `json:"Version"`
	LTime         int    `json:"LTime"`
}

func runCommand(command string) bool {
	parts := strings.Fields(command)
	cli := parts[0]
	args := parts[1:len(parts)]
	cmd := exec.Command(cli, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("exec='error' message='%v'", err)
		return false
	} else {
		return true
	}
}

func getHostname() string {
	hostname, _ := os.Hostname()
	return hostname
}

func createKey(event string) string {
	hostname := getHostname()
	return fmt.Sprintf("sifter/%s/%s", event, hostname)
}

func readStdin() string {
	bytes, _ := ioutil.ReadAll(os.Stdin)
	stdin := string(bytes)
	if stdin == "" || stdin == "[]\n" || stdin == "\n" {
		return ""
	} else {
		// TODO: Yes this is a gross hack and only works if
		// there is a single event in the payload.
		stdin = strings.TrimPrefix(stdin, "[")
		stdin = strings.TrimSuffix(stdin, "]\n")
		return stdin
	}
}

func decodeStdin(data string) (string, int64) {
	var events ConsulEvent
	err := json.Unmarshal([]byte(data), &events)
	if err != nil {
		Log(fmt.Sprintf("error: %s", data), "info")
		os.Exit(1)
	}
	name := string(events.Name)
	lTime := int64(events.LTime)
	Log(fmt.Sprintf("decoded event='%s' ltime='%d'", name, lTime), "info")
	return name, lTime
}
