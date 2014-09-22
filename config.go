package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
)

type SLHConfig struct {
	SerialPort SerialPortConfig `toml:"serialport"`
	Bluetooth  BluetoothConfig  `toml:"bluetooth"`
	Twitter    TwitterConfig    `toml:"twitter"`
}
type SerialPortConfig struct {
	Serial string `toml:"serial"`
}
type BluetoothConfig struct {
	Uuid  string `toml:"uuid"`
	Major string `toml:"major"`
	Minor string `toml:"minor"`
}
type TwitterConfig struct {
	Account           string `toml:"account"`
	ConsumerKey       string `toml:"consumer_key"`
	ConsumerSecret    string `toml:"consumer_secret"`
	AccessToken       string `toml:"access_token"`
	AccessTokenSecret string `toml:"access_token_secret"`
}

func UuidGenerate() (string, error) {
	uuidgen, err := exec.Command("uuidgen").Output()
	if err != nil {
		return "", err
	}
	return string(uuidgen[:len(uuidgen)-1]), nil
}

func RandomNum(i int64) string {
	// 0 - 65535
	rand.Seed(time.Now().Unix() * i)
	return strconv.Itoa(rand.Intn(math.MaxUint16))
}

func ScanInput() string {
	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()
	input := sc.Text()
	return input
}

// get config file path
func GetConfigPath() (string, error) {
	// check the home directory
	var homeDir string
	switch runtime.GOOS {
	case "darwin", "linux":
		homeDir = os.Getenv("HOME")
	default:
		return "", errors.New("don't support this platform")
	}

	// create the config directory
	// $HOME/.config/surelock-homes
	configDir := filepath.Join(homeDir, ".config", "surelock-homes")
	if _, err := os.Stat(configDir); err != nil {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return "", err
		}
	}

	// create the config filepath
	file := filepath.Join(configDir, "config.tml")
	return file, nil
}

// set config data
func SetConfig() error {
	var config SLHConfig

	var serialportConfig SerialPortConfig
	fmt.Printf("Input the serial port: ")
	serialportConfig.Serial = ScanInput()

	var bluetoothConfig BluetoothConfig
	uuid, err := UuidGenerate()
	if err != nil {
		// failed to generate a uuid
		return err
	}
	bluetoothConfig.Uuid = uuid
	bluetoothConfig.Major = RandomNum(1)
	bluetoothConfig.Minor = RandomNum(2)

	var twitterConfig TwitterConfig
	fmt.Printf("Input the twitter account: ")
	twitterConfig.Account = ScanInput()
	fmt.Printf("Input the twitter consumer-key: ")
	twitterConfig.ConsumerKey = ScanInput()
	fmt.Printf("Input the twitter consumer-secret: ")
	twitterConfig.ConsumerSecret = ScanInput()
	fmt.Printf("Input the twitter access-token: ")
	twitterConfig.AccessToken = ScanInput()
	fmt.Printf("Input the twitter access-token-secret: ")
	twitterConfig.AccessTokenSecret = ScanInput()

	config.SerialPort = serialportConfig
	config.Bluetooth = bluetoothConfig
	config.Twitter = twitterConfig

	var buffer bytes.Buffer
	encoder := toml.NewEncoder(&buffer)
	if err := encoder.Encode(config); err != nil {
		return err
	}

	// create the config file
	if file, err := GetConfigPath(); err != nil {
		return err
	} else {
		ioutil.WriteFile(file, []byte(buffer.String()), 0755)
		fmt.Printf("Created: \"" + file + "\"\n\n")
		fmt.Println(buffer.String())
	}

	return nil
}

// get config data
func GetConfig() (SLHConfig, error) {
	var config SLHConfig

	if file, err := GetConfigPath(); err != nil {
		return config, err
	} else {
		if _, err := toml.DecodeFile(file, &config); err != nil {
			return config, err
		} else {
			return config, nil
		}
	}
}
