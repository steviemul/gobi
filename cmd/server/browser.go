package main

import (
	"net"
	"os/exec"
	"runtime"

	"github.com/steviemul/gobi/internal/config"
)

// browserURL is the address to open for cfg. Wildcard and loopback listen
// addresses become "localhost"; a specific address is used as given.
func browserURL(cfg config.GobiConfig) string {
	host := cfg.Listen()

	switch host {
	case "", "0.0.0.0", "::", "127.0.0.1", "::1":
		host = "localhost"
	}

	return "http://" + net.JoinHostPort(host, cfg.Port()) + "/"
}

// browserCommand returns the command that opens url in the default browser on goos.
func browserCommand(goos string, url string) (name string, args []string) {
	switch goos {
	case "darwin":
		return "open", []string{url}
	case "windows":
		return "rundll32", []string{"url.dll,FileProtocolHandler", url}
	default:
		return "xdg-open", []string{url}
	}
}

// openBrowser starts the platform's opener for url and returns without waiting
// for the browser. The opener is reaped in the background so it doesn't linger
// as a zombie for the life of the server.
func openBrowser(url string) error {
	name, args := browserCommand(runtime.GOOS, url)

	cmd := exec.Command(name, args...)

	if err := cmd.Start(); err != nil {
		return err
	}

	go func() { _ = cmd.Wait() }()

	return nil
}
