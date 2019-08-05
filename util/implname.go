package util

import "fmt"

// ImplName returns the name of the struct that defimpl will define
// for an interface of the specified name.
func ImplName(package_name string, interface_name string) string {
	return fmt.Sprintf("%s.%s", package_name, interface_name)
}

