package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/go-ole/go-ole"
)

const (
	sFalse = 1
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	fmt.Println("Thread locked")

	err := ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED)
	if err != nil {
		oleErr, ok := err.(*ole.OleError)

		if !ok || oleErr.Code() != sFalse {
			return fmt.Errorf("COM initialization failed: %w", err)
		}
	}

	defer ole.CoUninitialize()

	fmt.Println("COM initialized successfully")
	return nil
}
