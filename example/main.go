package main

import (
	"log"
	"os"

	"github.com/spoke-d/clui"
	"github.com/spoke-d/clui/autocomplete/fsys"
)

func main() {
	fsys := fsys.NewLocalFileSystem()

	cli := clui.New("example", "1.0.0", "EXAMPLE", clui.OptionFileSystem(fsys))
	cli.Add("version", versionCmdFn)
	cli.Add("config show", configShowCmdFn)
	cli.Add("config show something", configShowCmdFn)
	cli.Add("config show else", configShowCmdFn)
	cli.Add("config list", configShowCmdFn)

	code, err := cli.Run(os.Args[1:])
	if err != nil {
		log.Fatal(">>", err, code)
	}

	os.Exit(code.Code())
}
