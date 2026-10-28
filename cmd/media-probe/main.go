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
	// Keep the goroutine on one OS thread for COM initialization and cleanup.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	fmt.Println("Thread locked")

	// Initialize COM in multithreaded mode.
	// S_FALSE also indicates success and requires CoUninitialize.
	err := ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED)
	if err != nil {
		oleErr, ok := err.(*ole.OleError)

		if !ok || oleErr.Code() != sFalse {
			return fmt.Errorf("COM initialization failed: %w", err)
		}
	}

	defer ole.CoUninitialize()

	fmt.Println("COM initialized successfully")

	// Start an asynchronous request for the GSMTC session manager.
	operation, err := control.GlobalSystemMediaTransportControlsSessionManagerRequestAsync()
	if err != nil {
		return fmt.Errorf("request GSMTC manager: %w", err)
	}

	defer operation.Release()

	// Register a callback that sends the completion status to a Go channel.
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

	// Wait for the operation to finish before accessing its result.
	status := <-done
	fmt.Printf("GSMTC operation status: %v\n", status)

	if status != foundation.AsyncStatusCompleted {
		return fmt.Errorf("GSMTC operation failed with status: %v", status)
	}

	// Retrieve the session manager from the completed operation.
	result, err := operation.GetResults()
	if err != nil {
		return fmt.Errorf("get GSMTC manager: %w", err)
	}

	if result == nil {
		return fmt.Errorf("GSMTC manager is nil")
	}

	manager := (*control.GlobalSystemMediaTransportControlsSessionManager)(result)

	defer manager.Release()

	// Get the collection of available media sessions.
	sessions, err := manager.GetSessions()
	if err != nil {
		return fmt.Errorf("get GSMTC sessions: %w", err)
	}

	if sessions == nil {
		return fmt.Errorf("GSMTC sessions is nil")
	}

	defer sessions.Release()

	// Check how many media sessions are available.
	size, err := sessions.GetSize()
	if err != nil {
		return fmt.Errorf("get size of GSMTC sessions: %w", err)
	}

	fmt.Printf("Number of GSMTC sessions: %d\n", size)

	if size == 0 {
		fmt.Print("No media sessions found\n")
		return nil
	}

	for i := uint32(0); i < size; i++ {
		// Retrieve the session for inspection.
		sessionPtr, err := sessions.GetAt(i)
		if err != nil {
			return fmt.Errorf("get # %v session pointer: %w", i, err)
		}
		if sessionPtr == nil {
			return fmt.Errorf("Session %v is nil", i)
		}
		currentSession := (*control.GlobalSystemMediaTransportControlsSession)(sessionPtr)

		defer currentSession.Release()

		sourceAppUserModelID, err := currentSession.GetSourceAppUserModelId()
		if err != nil {
			return fmt.Errorf("get session id: %w", err)
		}

		fmt.Printf("Session %v: %q\n", i, sourceAppUserModelID)

	}

	return nil
}
