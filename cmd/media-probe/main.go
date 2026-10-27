package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/go-ole/go-ole"
	"github.com/saltosystems/winrt-go"
	"github.com/saltosystems/winrt-go/windows/foundation"
	"github.com/saltosystems/winrt-go/windows/media/control"
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

	operation, err := control.GlobalSystemMediaTransportControlsSessionManagerRequestAsync()
	if err != nil {
		return fmt.Errorf("request GSMTC manager: %w", err)
	}

	defer operation.Release()

	done := make(chan foundation.AsyncStatus, 1)

	handlerIID := ole.NewGUID(winrt.ParameterizedInstanceGUID(
		foundation.GUIDAsyncOperationCompletedHandler,
		control.SignatureGlobalSystemMediaTransportControlsSessionManager,
	))

	handler := foundation.NewAsyncOperationCompletedHandler(
		handlerIID,
		func(
			_ *foundation.AsyncOperationCompletedHandler,
			_ *foundation.IAsyncOperation,
			status foundation.AsyncStatus,
		) {
			done <- status
		},
	)

	defer handler.Release()

	if err := operation.SetCompleted(handler); err != nil {
		return fmt.Errorf("register completion handler: %w", err)
	}

	status := <-done
	fmt.Printf("GSMTC operation status: %v\n", status)

	if status != foundation.AsyncStatusCompleted {
		return fmt.Errorf("GSMTC operation failed with status: %v", status)
	}

	result, err := operation.GetResults()
	if err != nil {
		return fmt.Errorf("get GSMTC manager: %w", err)
	}

	if result == nil {
		return fmt.Errorf("GSMTC manager is nil")
	}

	manager := (*control.GlobalSystemMediaTransportControlsSessionManager)(result)

	defer manager.Release()

	sessions, err := manager.GetSessions()
	if err != nil {
		return fmt.Errorf("get GSMTC sessions: %w", err)
	}

	if sessions == nil {
		return fmt.Errorf("GSMTC sessions is nil")
	}

	defer sessions.Release()

	size, err := sessions.GetSize()
	if err != nil {
		return fmt.Errorf("get size of GSMTC sessions: %w", err)
	}

	fmt.Printf("Number of GSMTC sessions: %d\n", size)

	return nil
}
