package helpers

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func ParseEnvFile() error {
	env, err := os.Open(".env")

	if err != nil {
		return err
	}
	defer func() {
		err := env.Close()
		if err != nil {
			log.Fatalf("Can't close file: %s", err)
		}
	}()

	scanner := bufio.NewScanner(env)

	for scanner.Scan() {
		str := scanner.Text()
		keyNValue := strings.Split(str, "=")
		unquoted, err := strconv.Unquote(keyNValue[1])
		if err != nil {
			return err
		}

		err = os.Setenv(keyNValue[0], unquoted)
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("err: %s", err)
		return err
	}

	return nil
}
