package utils

import (
	"bufio"
	"log/slog"
	"os"
	"strings"

	"github.com/GehirnInc/crypt"
	_ "github.com/GehirnInc/crypt/md5_crypt"
	"github.com/anjomro/kobra-unleashed/server/structs"
)

func CheckUser(user structs.LoginUser) bool {

	// Open /etc/shadow file and compare
	file, err := os.Open("/etc/shadow")
	if err != nil {
		slog.Error("Error opening file", err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, user.Username) {
			// Split line and compare
			slog.Info("Found user in shadow file:", "username", user.Username)
			shadowLine := strings.Split(line, ":")
			if CheckPassword(user.Password, shadowLine[1]) {
				return true
			}
		}
	}

	return false
}

func CheckPassword(password string, shadowHash string) bool {
	if !strings.HasPrefix(shadowHash, "$1$") {

		slog.Error("Password type not supported")
		return false
	}
	// $1$lwucTzHN$w7R4CU0Z2XpMLYKQlnuor.
	// Extract salt
	salt := strings.Split(shadowHash, "$")[2]
	// Generate hash
	hash, err := GenerateMD5Hash(password, salt)
	if err != nil {
		slog.Error("Error generating hash", err)
		return false
	}
	// Compare hash
	if hash == shadowHash {
		slog.Info("Password match")
		return true
	} else {
		slog.Info("Password mismatch")
		return false
	}
}

func GenerateMD5Hash(password string, salt string) (string, error) {

	crypt := crypt.MD5.New()

	hash, err := crypt.Generate([]byte(password), []byte("$1$"+salt))
	if err != nil {
		slog.Error("Error generating hash", err)
		return "", err
	}

	return string(hash), nil

}
