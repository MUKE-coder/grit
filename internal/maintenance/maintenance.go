package maintenance

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const maintenanceFile = ".maintenance"

// Enable puts the application in maintenance mode by creating a .maintenance file.
func Enable(projectRoot string) error {
	path := filepath.Join(projectRoot, maintenanceFile)
	content := fmt.Sprintf(`{"time":"%s"}`, time.Now().Format(time.RFC3339))
	return os.WriteFile(path, []byte(content), 0644)
}

// Disable brings the application back online by removing the .maintenance file.
func Disable(projectRoot string) error {
	path := filepath.Join(projectRoot, maintenanceFile)
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("application is not in maintenance mode")
		}
		return err
	}
	return nil
}

// IsEnabled checks if maintenance mode is active.
func IsEnabled(projectRoot string) bool {
	path := filepath.Join(projectRoot, maintenanceFile)
	_, err := os.Stat(path)
	return err == nil
}
