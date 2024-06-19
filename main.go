package main

import (
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	print(gpu_name() + "\n")
	print(cpu_name() + "\n")
}

func cpu_name() string {
	var cpu string
	var Cpu string

	runcmd := exec.Command("lscpu")
	output, err := runcmd.Output()
	if err != nil {
		println("lscpu is not present.")
	}

	for _, line := range strings.Split(string(output), "\n") {
		if strings.Contains(line, "Model name") {
			cpu += strings.TrimPrefix(line, "Model name:")
			whtSpaces := strings.Count(cpu, " ")
			Cpu += strings.Replace(cpu, " ", "", whtSpaces-4)
		}
	}

	return Cpu
}

func gpu_name() string {
	var bruh string
	var gpu string

	cmd := exec.Command("lspci", "-v")
	shell, err := cmd.Output()
	if err != nil {
		println("lspci is not present.")
	}

	for _, line := range strings.Split(
		strings.TrimSuffix(string(shell), "\n"), "\n") {
		if strings.Contains(line, "VGA") {
			bruh += line[strings.Index(line, ": ")+2 : strings.Index(line, " (")]
		}
	}

	var rgx = regexp.MustCompile(`\[(.*?)\]`)

	gpuRegex := rgx.FindStringSubmatch(bruh)

	for _, char := range gpuRegex[1] {
		gpu += string(char)
	}

	return gpu
}
