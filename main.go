package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"net"
	"os"
	"strings"
	"time"
)

// Alive stands for boolean whether the proxy responded in Resolve function.
// ResponseTime is time it took for the instance to respond.
// Address is (for now) a string representing the IP address of the proxy to be connected.
type Email struct {
	Alive        bool
	Address      string
	Records      []*net.MX
}

// Resolve makes TCP connection to given host with given timeout.
// The function broadcasts to given channel about the response of the connection.
// Even if the connection is refused or it times out, the channel will receive
// new proxy instance with filled fields.
func Resolve(status chan Email, throttle chan bool, host string) {

	var email Email
	email.Address = host

	records, err := net.LookupMX(strings.Split(host, "@")[1])
	email.Records = records
	if err != nil || len(records) == 0 {
		if err.Error() == "dial udp 8.8.4.4:53: too many open files" {
			throttle <- true
		}
		email.Alive = false
		status <- email
		return
	}

	email.Alive = true
	status <- email
	return
}

// ReadFile returns string array of IP addresses in a file. It uses scanner.Scan()
// to parse IP's from the given file.
func ReadFile(filename string) ([]string, error) {

	var proxies []string
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return proxies, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		proxies = append(proxies, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return proxies, err
	}

	return proxies, nil
}

// WriteFile writes given buffer to given filename. It creates the file
// it is not yet created. The created file uses 0600 permissions by default.
func WriteLine(f *os.File, line string) error {
	if _, err := f.WriteString(line); err != nil {
		return err
	}
	return nil
}

func main() {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	proxies, err := ReadFile(dir+"/input.txt")
	if err != nil {
		fmt.Println(err)
	}

	f, err := os.OpenFile(dir+"/output.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	proxychan := make(chan Email, len(proxies))
	throttle := make(chan bool)
	lines := 0
	amountOf := 0
	alive := 0
	te := time.Now()

	go func() {
		for _, proxy := range proxies {
			time.Sleep(50*time.Millisecond)
			lines++
			amountOf++
			go Resolve(proxychan, throttle, proxy)
		}
		for _ = range throttle {
			fmt.Println("Throttling triggered.")
			time.Sleep(3*time.Second)
		}
	}()

	for proxy := range proxychan {
		if proxy.Alive {
			alive++
			err := WriteLine(f, proxy.Address + "\n")
			if err != nil {
				fmt.Println(err)
			}
		}
		fmt.Printf("\rProgress: %v/%v/%v. Time elapsed: %v", amountOf, alive, len(proxies), time.Now().Sub(te))
	}
}
