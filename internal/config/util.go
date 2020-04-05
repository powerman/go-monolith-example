package config

import "fmt"

func join(pfx, name string) string {
	if pfx == "" {
		return name
	}
	return fmt.Sprintf("%s.%s", pfx, name)
}
