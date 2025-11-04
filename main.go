package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
)

func main() {
	// Parse flags
	f := flag.NewFlagSet("Check certificates", flag.ExitOnError)
	a := f.String("accountname", "", "Account name for cert")
	fn := f.String("filename", "", "File name for cert")
	if err := f.Parse(os.Args[1:]); err != nil {
		fmt.Println("Couldn't parse flags")
		f.Usage()
	}

	if *a == "" && *fn == "" {
		fmt.Println("Must specify either --accountname or --filename")
		os.Exit(1)
	}
	if *a != "" && *fn != "" {
		fmt.Println("Must specify either --accountname or --filename")
		os.Exit(1)
	}

	// Get path of cert file
	curDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Couldn't get current working directory")
		os.Exit(2)
	}

	var filename string

	switch {
	case *a != "":
		// Account name is just $PWD/certs/acctName.cert
		filename = path.Join(curDir, "certs", fmt.Sprintf("%s.cert", *a))
	case *fn != "":
		filename = *fn
	}

	if _, err := os.Stat(filename); errors.Is(err, fs.ErrNotExist) {
		fmt.Printf("The cert file %s doesn't exist", filename)
		os.Exit(3)
	}

	fmt.Println("Filename: ", filename)

	// Run the command
	cmd := exec.Command("openssl",
		"x509",
		"-in", filename,
		"-noout",
		"-subject",
		"-dates",
		"-nameopt", "compat")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Couldn't run openssl command")
		os.Exit(3)
	}

	fmt.Printf(string(stdoutStderr))
}
