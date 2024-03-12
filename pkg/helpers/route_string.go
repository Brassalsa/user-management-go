package helpers

import "fmt"

func RouteStrClosure(prefix string) func(method, postfix string) string {
	return func(method, route string) string{
		return fmt.Sprintf("%s %s%s", method, prefix, route)
	}
}