//go:build tray

package cmd

import "os"

func init() {
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "sync")
	}
}
