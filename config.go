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

func uuid() string {
	uuidgen, err := exec.Command("uuidgen").Output()
	if err != nil {
		panic(err)
	}
	return string(uuidgen[:len(uuidgen)-1])
}

func randomNum(i int64) string {
	// 0 - 65535
	rand.Seed(time.Now().Unix() * i)
	return strconv.Itoa(rand.Intn(math.MaxUint16))
}

func scan() string {
	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()
	input := sc.Text()
	return input
}

// get config file path
func getConfigPath() (string, error) {
	// check the home directory
	homeDir := ""
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
func setConfig() error {
	var config SLHConfig

	var serialportConfig SerialPortConfig
	fmt.Printf("Input the serial port: ")
	serialportConfig.Serial = scan()

	var bluetoothConfig BluetoothConfig
	bluetoothConfig.Uuid = uuid()
	bluetoothConfig.Major = randomNum(1)
	bluetoothConfig.Minor = randomNum(2)

	var twitterConfig TwitterConfig
	fmt.Printf("Input the twitter account: ")
	twitterConfig.Account = scan()
	fmt.Printf("Input the twitter consumer-key: ")
	twitterConfig.ConsumerKey = scan()
	fmt.Printf("Input the twitter consumer-secret: ")
	twitterConfig.ConsumerSecret = scan()
	fmt.Printf("Input the twitter access-token: ")
	twitterConfig.AccessToken = scan()
	fmt.Printf("Input the twitter access-token-secret: ")
	twitterConfig.AccessTokenSecret = scan()

	config.SerialPort = serialportConfig
	config.Bluetooth = bluetoothConfig
	config.Twitter = twitterConfig

	var buffer bytes.Buffer
	encoder := toml.NewEncoder(&buffer)
	if err := encoder.Encode(config); err != nil {
		return err
	}

	// create the config file
	if file, err := getConfigPath(); err != nil {
		return err
	} else {
		ioutil.WriteFile(file, []byte(buffer.String()), 0755)
		fmt.Printf("Created: \"" + file + "\"\n\n")
		fmt.Println(buffer.String())
	}

	return nil
}

// get config data
func getConfig() (SLHConfig, error) {
	var config SLHConfig

	if file, err := getConfigPath(); err != nil {
		return config, err
	} else {
		if _, err := toml.DecodeFile(file, &config); err != nil {
			return config, err
		} else {
			return config, nil
		}
	}
}
