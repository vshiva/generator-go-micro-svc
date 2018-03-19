<%=licenseText%>
package util

import "gopkg.in/urfave/cli.v1"

// ConcatFlags concatenates all flags provided to a single array. Currently it
// does not inspect if an flag already exists, or sort by name etc.
func ConcatFlags(flags ...[]cli.Flag) []cli.Flag {
	result := []cli.Flag{}

	for _, f := range flags {
		result = append(result, f...)
	}

	return result
}
