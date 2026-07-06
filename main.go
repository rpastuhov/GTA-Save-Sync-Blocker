package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/moutend/go-hook/pkg/keyboard"
	"github.com/moutend/go-hook/pkg/types"
)

const (
	ruleName     = "GTA_SYNC"
	blockIP      = "192.81.241.171"
	instructions = `
========================================
              HOW TO USE
========================================

1. Start a heist or mission.

2. Press Ctrl+F9 at any point between loading into the job and the moment
   you receive the payout from Madrazo's people (or the mission finale).

3. After loading into the session, switch to Story Mode.

4. Wait for the cutscene to finish, then press Ctrl+F12.

5. Join a session via invite.

6. Deposit the money into your bank account.

7. Force a save (the easiest way is to change your outfit in the interaction
   menu - you can pick the same outfit and press Enter).

IMPORTANT:
- Cayo Perico limit: 4 heists per hour from the moment you receive your
  first payout (exceeding this limit is not recommended).
- Your cut is calculated as: potential take * 0.88 * your percentage.
- Contract limit: 2 per 30 minutes from the moment you receive your first
  payout (a third one won't pay out, and may result in a ban).

HOTKEYS:
- Ctrl+F9  - Enable save-blocking mode
- Ctrl+F12 - Disable save-blocking mode
`
)

var (
	firewallInitiallyActive bool
	blockingActive          bool
)

func main() {
	if !isAdmin() {
		fmt.Println("ERROR: Administrator privileges required!")
		fmt.Println("Please run the program as administrator.")
		fmt.Println("\nPress Enter to exit...")
		fmt.Scanln()
		os.Exit(1)
	}

	fmt.Print(instructions)

	firewallInitiallyActive = isFirewallEnabled()
	if !firewallInitiallyActive {
		fmt.Println("Firewall is disabled. Enabling...")
		enableFirewall()
	}

	setupExitHandler()

	keyboardChan := make(chan types.KeyboardEvent, 100)

	if err := keyboard.Install(nil, keyboardChan); err != nil {
		fmt.Printf("Error installing keyboard hook: %v\n", err)
		fmt.Scanln()
		cleanup()
		os.Exit(1)
	}
	defer keyboard.Uninstall()

	fmt.Println("Program started. Waiting for key presses...")

	ctrlPressed := false

	for event := range keyboardChan {
		if event.Message != types.WM_KEYDOWN && event.Message != types.WM_KEYUP {
			continue
		}

		isDown := event.Message == types.WM_KEYDOWN

		switch event.VKCode {
		case types.VK_CONTROL, types.VK_LCONTROL, types.VK_RCONTROL:
			ctrlPressed = isDown
		case types.VK_F9:
			if isDown && ctrlPressed && !blockingActive {
				activateBlocking()
			}
		case types.VK_F12:
			if isDown && ctrlPressed && blockingActive {
				deactivateBlocking()
			}
		}
	}
}

func isAdmin() bool {
	cmd := exec.Command("net", "session")
	err := cmd.Run()
	return err == nil
}

func isFirewallEnabled() bool {
	cmd := exec.Command("netsh", "advfirewall", "show", "allprofiles", "state")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "ON") || strings.Contains(string(output), "ВКЛЮЧИТЬ")
}

func enableFirewall() {
	profiles := []string{"domainprofile", "privateprofile", "publicprofile"}
	for _, profile := range profiles {
		cmd := exec.Command("netsh", "advfirewall", "set", profile, "state", "on")
		cmd.Run()
	}
}

func disableFirewall() {
	profiles := []string{"domainprofile", "privateprofile", "publicprofile"}
	for _, profile := range profiles {
		cmd := exec.Command("netsh", "advfirewall", "set", profile, "state", "off")
		cmd.Run()
	}
}

func activateBlocking() {
	cmd := exec.Command("netsh", "advfirewall", "firewall", "add", "rule",
		fmt.Sprintf("name=%s", ruleName),
		"dir=out",
		"action=block",
		fmt.Sprintf("remoteip=%s", blockIP))

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error adding rule: %v\n", err)
		return
	}

	blockingActive = true
	showTooltip("[ACTIVE] Server sync disabled")
}

func deactivateBlocking() {
	cmd := exec.Command("netsh", "advfirewall", "firewall", "delete", "rule",
		fmt.Sprintf("name=%s", ruleName))

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error deleting rule: %v\n", err)
		return
	}

	blockingActive = false
	showTooltip("[DISABLED] Server sync restored")
}

func cleanup() {
	fmt.Println("\nShutting down...")

	if blockingActive {
		cmd := exec.Command("netsh", "advfirewall", "firewall", "delete", "rule",
			fmt.Sprintf("name=%s", ruleName))
		cmd.Run()
	}

	if !firewallInitiallyActive {
		fmt.Println("Restoring original firewall state...")
		disableFirewall()
	}
}

func setupExitHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(0)
	}()
}

func showTooltip(message string) {
	go func() {
		fmt.Printf("[%s] %s\n", time.Now().Format("15:04:05"), message)
		time.Sleep(3 * time.Second)
	}()
}
