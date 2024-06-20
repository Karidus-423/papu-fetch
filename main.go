package main

import (
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	print(gpu_name() + "\n")
	print(cpu_name() + "\n")
	print(storage_name() + "\n")
	print(ram_name() + "\n")
	print(os_name() + "\n")
	print(desktopEnv_name())
	print(theme_name())
	print(get_shell())
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
	var filter string
	var gpu string

	runcmd := exec.Command("lspci", "-v")
	output, err := runcmd.Output()
	if err != nil {
		println("lspci is not present.")
	}

	splitting := strings.TrimSuffix(string(output), "\n")

	for _, line := range strings.Split(splitting, "\n") {
		if strings.Contains(line, "VGA") {
			filter += line[strings.Index(line, ": ")+2 : strings.Index(line, " (")]
		}
	}

	var rgx = regexp.MustCompile(`\[(.*?)\]`)

	gpuRegex := rgx.FindStringSubmatch(filter)

	for _, char := range gpuRegex[1] {
		gpu += string(char)
	}

	return gpu
}

func storage_name() string {
	var storage string

	runcmd := exec.Command("lsblk", "-o", "MODEL", "-A")
	output, err := runcmd.Output()
	if err != nil {
		println("lsblk is not present.")
	}

	for drvnum, line := range strings.Split(string(output), "\n") {
		if drvnum == 1 {
			storage += strings.TrimSuffix(line, "/\n")
		}
	}

	return storage
}

func ram_name() string {
	var ram string

	runcmd := exec.Command("free", "--giga")
	output, err := runcmd.Output()
	if err != nil {
		println("free is not present.")
	}

	var rgx = regexp.MustCompile(`Mem:\s+(\d+)`)

	outputArray := strings.Split(string(output), "\n")
	for _, line := range outputArray {
		if strings.Contains(line, "Mem:") {
			ramRgx := rgx.FindStringSubmatch(line)
			ram += ramRgx[1]
		}
	}
	ram += " GB"
	return ram
}

func os_name() string {
	var distro string
	file, err := os.ReadFile("/etc/os-release")
	if err != nil {
		println("Unable to find file /etc/os-release")
	}

	var rgx = regexp.MustCompile(`\"(.*?)\"`)

	file_lines := strings.Split(string(file), "\n")
	for _, line := range file_lines {
		if strings.Contains(line, "PRETTY_NAME") {
			pretty_name := rgx.FindStringSubmatch(line)
			distro += pretty_name[1]
		}
	}

	return distro
}

func desktopEnv_name() string {
	output, err := exec.Command("printenv", "XDG_CURRENT_DESKTOP").Output()
	if err != nil {
		print("$XDG_CURRENT_DESKTOP variable is not set")
	}

	return string(output)
}

func theme_name() string {
	output, err := exec.Command("printenv", "GTK_THEME").Output()
	if err != nil {
		print("GTK_THEME variable is not set.")
	}

	gtkTheme := string(output)

	return gtkTheme
}

func get_shell() string {
	output := os.Getenv("SHELL")
	if strings.Contains(output, "/") {
		slashCnt := strings.Count(output, "/")
		print(slashCnt)
	}
	return output
}
