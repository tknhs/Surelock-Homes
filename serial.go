package main

import (
	"io"

	"github.com/tarm/goserial"
)

func SerialInit(serialConfig SerialPortConfig) (io.ReadWriteCloser, error) {
	c := &serial.Config{
		Name: serialConfig.Serial,
		Baud: 2400,
	}
	s, err := serial.OpenPort(c)
	if err != nil {
		return s, err
	}

	return s, nil
}

func SerialWrite(serialObject io.ReadWriteCloser, openClose string) error {
	_, err := serialObject.Write([]byte(openClose))
	if err != nil {
		return err
	}
	return nil
}
