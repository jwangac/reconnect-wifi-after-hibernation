// Copyright (c) 2020 Jiawei Wang

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package main

import (
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

func current_ssid() string {
	cmd := exec.Command("netsh", "wlan", "show", "interfaces")
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000}
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "    SSID") {
			pos := strings.Index(line, ": ")
			if pos > 0 {
				return strings.TrimSpace(line[pos+1:])
			}
		}
	}
	return ""
}

func disconnect() {
	cmd := exec.Command("netsh", "wlan", "disconnect")
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000}
	cmd.Output()
}

func connect(ssid string) {
	cmd := exec.Command("netsh", "wlan", "connect", ssid)
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000}
	cmd.Output()
}

func main() {
	// run in background
	if len(os.Args) == 1 {
		exec.Command("CMD", "/C", "START", "/B", os.Args[0], "forked").Run()
		return
	}
	syscall.NewLazyDLL("kernel32.dll").NewProc("FreeConsole").Call()

	ssid := ""
	for {
		tmp := current_ssid()
		if tmp != "" {
			ssid = tmp
		}

		start := time.Now()

		time.Sleep(60 * time.Second)

		if time.Since(start) > 120*time.Second {
			disconnect()

			connect(ssid)

			// retry after 5 seconds
			go func(ssid string) {
				time.Sleep(5 * time.Second)
				if current_ssid() == "" {
					connect(ssid)
				}
			}(ssid)
		}
	}
}
