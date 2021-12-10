package flagvalue

import (
	"errors"
	"strconv"
)

type CustomFlag struct {
	address  string
	port     int
	Interval ArrayFlags
}

type ArrayFlags struct {
	start int
	end   int
}

//Yapıcı fonksiyon olduğu için pointer döndürür
func (c *CustomFlag) NewCustomFlag() *CustomFlag {
	c = &CustomFlag{}
	c.port = 0
	c.address = ""
	c.Interval.start = 1
	c.Interval.end = 65535
	return c
}

func (c *CustomFlag) GetAddress() string {
	return c.address
}

func (c *CustomFlag) SetAddress(value string) {
	if value == "" {
		panic("Please specify an address with '--address' flag")
	}
	c.address = value
}

func (c *CustomFlag) GetPort() int {
	return c.port
}

func (c *CustomFlag) SetPort(value int) {
	c.port = value
}

func (i *ArrayFlags) GetStart() int {
	return i.start
}

func (i *ArrayFlags) SetStart(value int) {
	i.start = value
}

func (i *ArrayFlags) GetEnd() int {
	return i.end
}

func (i *ArrayFlags) SetEnd(value int) {
	i.end = value
}

//Kuyruk flag'inin kontrolleri
func CheckInterval(stringInterval []string, portNum int) ([]int, error) {
	if portNum != 0 && len(stringInterval) != 0 {
		return nil, errors.New("you can't use '--port' flag when you are trying to scan a port interval")
	}
	if len(stringInterval) == 0 {
		slice := []int{1, 65535}
		return slice, nil
	}
	slice := []int{}
	for i := range stringInterval {
		text := stringInterval[i]
		number, err := strconv.Atoi(text)
		if err != nil {
			return nil, errors.New("only numeric values are valid as the port number")
		}
		slice = append(slice, number)
	}

	if len(slice) != 2 {
		return nil, errors.New("please type initial and final port number only")
	} else if slice[0] >= slice[1] {
		return nil, errors.New("initial port number cannot be greater or equal than final port number")
	} else if slice[0] < 1 {
		return nil, errors.New("initial port number cannot be less than 1")
	} else if slice[1] > 65535 {
		return nil, errors.New("final port number cannot be greater than 6553")
	}
	return slice, nil
}
