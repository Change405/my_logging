package main

import (
	"fmt"
	"os/exec"
)

func main() {
	cmd := exec.Command("ls", "-la")
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error when runnging test")
		fmt.Println(err.Error())
	}
	fmt.Println(string(stdout))

}
