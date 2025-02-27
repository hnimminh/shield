package nftcli

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"
)

func NftOverShell(cmdstring string) (string, error) {
	if !strings.HasPrefix(cmdstring, "nft ") || len(cmdstring) < 5 {
		return "", errors.New("ForbiddenCommand")
	}
	args := strings.Split(cmdstring[4:], " ")
	excmd := exec.Command("nft", args...)

	var out bytes.Buffer
	excmd.Stdout = &out

	err := excmd.Run()
	if err != nil {
		return out.String(), err
	}
	return out.String(), nil
}

func StreamNftOverShell(stream string) error {
	r, w, err := os.Pipe()
	if err != nil {
		return err
	}
	defer r.Close()

	echo := exec.Command("echo", stream)
	echo.Stdout = w
	err = echo.Start()
	if err != nil {
		return err
	}
	defer echo.Wait()
	w.Close()

	nft := exec.Command("nft", "-f", "-")
	nft.Stdin = r
	nft.Stdout = os.Stdout
	return nft.Run()
}
