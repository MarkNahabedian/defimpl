package util

import "fmt"

// ImplName returns the name of the struct that defimpl will define
// for an interface of the specified name.
func ImplName(package_name string, interface_name string) string {
	if package_name == "" {
		return fmt.Sprintf("%sImpl", interface_name)
	}
	return fmt.Sprintf("%s.%sImpl", package_name, interface_name)
}

