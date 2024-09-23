package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

func main() {

	// Выполняем команду netstat
	cmd := exec.Command("netstat", "-u", "-n", "-a")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Ошибка при выполнении команды netstat:", err)
		return
	}

	// Преобразуем вывод команды в строку и разбиваем его на строки
	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")

	for _, line := range lines {
		if strings.Contains(line, "udp") {
			re := regexp.MustCompile(`(22[4-9]|23[0-9]\.\d+\.\d+\.\d+):(\d+)`)
			matches := re.FindStringSubmatch(line)
			if len(matches) == 3 {
				process(matches[1], matches[2])
			}
		}
	}

}
