package core

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/kitex/pkg/utils"
)

// credFilePath : The filepath to store the credentials file
var credFilePath = filepath.Join(utils.GetConfDir(), "credentials")

// TODO storeCredentials : Store
func (m *GoLearn) StoreCredentials(username string, password string) error {
	if err := ioutil.WriteFile(credFilePath, []byte(username+"&"+password), os.ModePerm); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// TODO: deleteCredentials : Delete saved credentials
func (m *GoLearn) DeleteCredentials() {
	if err := os.Remove(credFilePath); err != nil {
		log.Println(err)
	}
}

// TODO: loadCredentials : Load saved credentials
func (m *GoLearn) LoadCredentials() (string, string, error) {
	content, err := ioutil.ReadFile(credFilePath)
	if err != nil {
		log.Println(err)
		return "", "", err
	}
	stringContent := string(content)
	username := strings.Split(stringContent, "&")[0]
	password := strings.Split(stringContent, "&")[1]
	return username, password, nil
}

// TODO: restoreSession : Check if the user's credentials have been stored before
func (m *GoLearn) RestoreSession() error {
	if _, err := os.Stat(credFilePath); err != nil {
		fmt.Println("No past session, redirecting to login...")
		return err
	}
	fmt.Println("Found past session, attempting to restore...")

	// Try to read stored credentials
	username, password, err := m.LoadCredentials()
	if err != nil {
		fmt.Println("Trouble reading credentials")
	}

	//  Attempt to login with stored credentials
	fmt.Println("Attempting to log in")
	if err := App.Client.Auth.Login(username, password); err != nil {
		fmt.Printf("Error trying to log	in: %s\n", err.Error())
		return err
	}

	return nil
}
