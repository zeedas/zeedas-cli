package offlineprint

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/zeedas/zeedas-cli/cmd/params"
	"github.com/zeedas/zeedas-cli/pkg/exitcode"
	"github.com/zeedas/zeedas-cli/pkg/heartbeat"
	"github.com/zeedas/zeedas-cli/pkg/offline"

	"github.com/spf13/viper"
)

// Run executes the print-offline-heartbeats command.
func Run(v *viper.Viper) (int, error) {
	queueFilepath, err := offline.QueueFilepath()
	if err != nil {
		return exitcode.ErrGeneric, fmt.Errorf(
			"failed to load offline queue filepath: %s",
			err,
		)
	}

	p := params.LoadOfflineParams(v)

	if p.QueueFile != "" {
		queueFilepath = p.QueueFile
	}

	hh, err := offline.ReadHeartbeats(queueFilepath, p.PrintMax)
	if err != nil {
		fmt.Println(err)
		return exitcode.ErrGeneric, fmt.Errorf("failed to read offline heartbeats: %w", err)
	}

	data, err := jsonWithoutEscaping(hh)
	if err != nil {
		fmt.Println(err)
		return exitcode.ErrGeneric, fmt.Errorf("failed to json marshal offline heartbeats: %w", err)
	}

	fmt.Print(string(data))

	return exitcode.Success, nil
}

// jsonWithoutEscaping returns a string representation of the given array of heartbeats.
// It does not escape the angle brackets "<", ">" and "&".
func jsonWithoutEscaping(hh []heartbeat.Heartbeat) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(hh)

	return buffer.Bytes(), err
}
