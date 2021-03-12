package fs

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func GetOrCreateDir(path string) (os.FileInfo, error) {
	if _, err := os.Stat(path); err != nil {
		if pe, ok := err.(*os.PathError); !ok || pe.Err.Error() != "no such file or directory" {
			return nil, err
		}

		if err := os.MkdirAll(path, os.ModeDir | os.ModePerm); err != nil {
			fmt.Printf("Can't create '%s' dir: '%s'\n", path, err)
			return nil, err
		}
	}

	return os.Stat(path)
}

func GetOrCreateFile(path string) (os.FileInfo, error) {
	if _, err := os.Stat(path); err != nil {
		if pe, ok := err.(*os.PathError); !ok || pe.Err.Error() != "no such file or directory" {
			return nil, err
		}

		if _, err := os.Create(path); err != nil {
			fmt.Printf("Can't create '%s' file: '%s'\n", path, err)
			return nil, err
		}
	}

	return os.Stat(path)
}

func Exist(path string) bool {
	// @todo other errors
	_, err := os.Stat(path)

	return err == nil
}

func Remove(path string) {
	_ = os.RemoveAll(path)
}

func ContinueIn(path string, cmd string) error {
	return execute(path, cmd, os.Stdout, os.Stdin)
}

func ContinueInArgs(path string, cmd ...string) error {
	return execute(path, strings.Join(cmd, " "), os.Stdout, os.Stdin)
}

func ExecIn(path string, cmd string) (string, error) {
	buf := &bytes.Buffer{}

	err := execute(path, cmd, buf, os.Stdin)

	return buf.String(), err
}

func ExecInArgs(path string, cmd ...string) (string, error) {
	return ExecIn(path, strings.Join(cmd, " "))
}

func execute(path string, cmd string, out io.Writer, in io.Reader) error {
	fmt.Printf("Exec `%s` in `%s`\n", cmd, path)

	cmdSlice := strings.Split(cmd, " ")
	shell := exec.Command(cmdSlice[0], cmdSlice[1:]...)
	shell.Dir = path
	shell.Stdout = out
	shell.Stderr = out
	shell.Stdin = in

	err := shell.Run()

	return err
}


