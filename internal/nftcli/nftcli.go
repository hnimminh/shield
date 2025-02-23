package nftcli

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

func NftOverShell(cmdstring string) (string, error) {
	if !strings.HasPrefix(cmdstring, "nft ") || len(cmdstring) < 5 {
		return "", errors.New("ForbiddenCommand")
	}

	excmd := exec.Command("nft", cmdstring[5:])

	var out bytes.Buffer
	excmd.Stdout = &out

	err := excmd.Run()
	if err != nil {
		return out.String(), err
	}
	return out.String(), nil
}
