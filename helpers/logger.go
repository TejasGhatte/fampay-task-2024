package helpers

import (
	"fmt"

	"github.com/TejasGhatte/fampay-task-2024/initializers"
)

func LogDatabaseError(customString string, err error, path string) {
	if err == nil {
		err = fmt.Errorf("no error description provided")
	}
	initializers.Logger.Warnw(customString, "Message", err.Error(), "Path", path, "Error", err)
}

func LogServerError(customString string, err error, path string) {
	if err == nil {
		err = fmt.Errorf("no error description provided")
	}

	initializers.Logger.Errorw(customString, "Message", err.Error(), "Path", path, "Error", err)
}
