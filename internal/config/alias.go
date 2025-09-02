package config

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func SetupAlias() {
	fmt.Printf("%s - Setting up alias\n", AppName)
	executable, err := os.Executable()
	if err != nil {
		fmt.Printf("failed to get executable: %v\n", err)
		return
	}
	profile := path.Join(config.UserHomeDir, ".profile")

	content, err := os.ReadFile(profile)
	if err != nil {
		fmt.Printf("failed to read file: %v", err)
		return
	}
	if strings.Contains(string(content), "alias govuln=") {
		fmt.Printf("alias already exists in file: %s\n", profile)
		return
	}
	if err := os.WriteFile(profile+".bak", content, 0644); err != nil {
		fmt.Printf("failed to backup file: %v\n", err)
		return
	}
	content = fmt.Appendf(content, `
# govulncheck-hook
alias govuln=%s
`, executable)

	if err := os.WriteFile(profile, content, 0644); err != nil {
		fmt.Printf("failed to write updated file: %v\n", err)
		return
	}
	fmt.Printf(`successfully updated file: %s .
A backup copy has been created at: %s.bak\n`, profile, profile)
	fmt.Println("Please restart your shell or run the following command to use the alias:")
	fmt.Printf("\nsource %s\n", profile)
	fmt.Println("\nAfter restarting your shell, you can run the tool using the 'govuln' command.")
}
