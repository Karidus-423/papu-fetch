package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func main() {
	layout()
}

func layout() {
	boxHeight, boxWidth := 30, 80
	var mySwagBorder = lipgloss.Border{
		Top:         "━",
		Bottom:      "━",
		Left:        "┃",
		Right:       "┃",
		TopLeft:     "┏",
		TopRight:    "┓",
		BottomLeft:  "┗",
		BottomRight: "┛",
	}

	var outerBox = lipgloss.NewStyle().
		BorderStyle(mySwagBorder).
		BorderForeground(lipgloss.Color("7")).
		Height(boxHeight).
		Width(boxWidth)
	var rightBox = lipgloss.NewStyle().
		BorderStyle(mySwagBorder).
		BorderForeground(lipgloss.Color("7")).
		Height(15).
		Width(20).Padding(0)
	var centerBox = lipgloss.NewStyle().
		BorderStyle(mySwagBorder).
		BorderForeground(lipgloss.Color("7")).
		Height(15).
		Width(40)
	var leftBox = lipgloss.NewStyle().
		BorderStyle(mySwagBorder).
		BorderForeground(lipgloss.Color("7")).
		Height(15).
		Width(20)
	var bottomLeft = lipgloss.NewStyle().
		BorderStyle(mySwagBorder).
		BorderForeground(lipgloss.Color("7")).
		Height(1).
		Width(1)
	var bottomRight = lipgloss.NewStyle().
		BorderStyle(mySwagBorder).
		BorderForeground(lipgloss.Color("7")).
		Height(1).
		Width(1)

	fmt.Println(outerBox.Render(
		lipgloss.JoinHorizontal(lipgloss.Center,
			rightBox.Render(),
			centerBox.Render(),
			leftBox.Render()),
		bottomLeft.Render(),
		bottomRight.Render(),
	))
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

func shell_name() string {
	var shellName string
	output := os.Getenv("SHELL")
	if strings.Contains(output, "/") {
		slashCnt := strings.Count(output, "/") + 1
		fmtOutput := strings.SplitAfterN(output, "/", slashCnt)

		shellName += fmtOutput[len(fmtOutput)-1] + "\n"
	}
	return shellName
}

func terminal_name() string {
	output := os.Getenv("TERM_PROGRAM")

	return output + "\n"
}

func host_name() string {
	username, hostcmd := os.Getenv("USER"), exec.Command("hostname")
	host, err := hostcmd.Output()
	if err != nil {
		print("Unable to find host name")
	}
	hostname := "@" + string(host)

	id := username + hostname

	return id
}
