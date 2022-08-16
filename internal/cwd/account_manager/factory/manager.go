package factory

import (
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/cansulting/elabox-system-tools/foundation/path"
	"github.com/cansulting/elabox-system-tools/foundation/perm"
	"github.com/cansulting/elabox-system-tools/internal/cwd/account_manager/data"
)

const PASSWD_FILE = "/etc/passwd"

// use to check if account exist
func IsAccountExist(username string) (bool, error) {
	contents, err := os.ReadFile(PASSWD_FILE)
	if err != nil {
		return false, err
	}
	return strings.Contains(string(contents), username+":"), nil
}

func LoadAccount(username string) *data.Account {
	return nil
}

func CreateAccount(pass string, acc data.Account) error {
	// save account profile
	if err := SaveAccount(acc); err != nil {
		return err
	}
	home := path.PATH_USERS + "/" + acc.Username
	if err := os.MkdirAll(home, perm.PRIVATE); err != nil {
		return err
	}
	// generates pass
	cmd := exec.Command("openssl", "passwd", "-1", pass)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(output) + "." + err.Error())
	}
	s_passwd := string(output)
	s_passwd = strings.TrimSuffix(s_passwd, "\n")
	// creates user
	cmd = exec.Command("useradd", "-p", s_passwd, "-m", "-d", home, acc.Username)
	output, err = cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(output) + "." + err.Error())
	}
	return nil
}

func SaveAccount(account data.Account) error {
	return nil
}

func DeleteAccount(username string) error {
	cmd := exec.Command("userdel", "-f", username)
	_, err := cmd.CombinedOutput()
	return err
}
