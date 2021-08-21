package main

// https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go

import (
	"reflect"
	"regexp"
	"testing"
	"time"
)

var (
	resultString      = ""
	resultBool        = false
	resultInt         = 0
	resultSlice4Int   = [4]int{}
	resultSlice2Int   = [2]int{}
	resultSlice3Float = [3]float64{}
)

func TestGetHostname(t *testing.T) {
	hostname := getHostname()

	if reflect.TypeOf(hostname) != reflect.TypeOf("") {
		t.Error("getHostname() was incorrect, didn't get a string.")
	}

	if len(hostname) == 0 {
		t.Error("getHostname() was incorrect, got zero-length string.")
	}
}

func BenchmarkGetHostname(b *testing.B) {
	var res string

	for n := 0; n < b.N; n++ {
		res = getHostname()
	}

	resultString = res
}

func TestGetIPAddress(t *testing.T) {
	IPAddress := getIPAddress()
	IPAddresLen := len(IPAddress)
	minLen, maxLen := 7, 15
	if reflect.TypeOf(IPAddress) != reflect.TypeOf("") {
		t.Error("getIPAddress() was incorrect, didn't get a string.")
	}

	if IPAddresLen < minLen {
		t.Errorf("getIPAddress() was incorrect, IP address was too short, got: %d, want minimum of %d.", IPAddresLen, minLen)
	}
	if IPAddresLen > maxLen {
		t.Errorf("getIPAddress() was incorrect, IP address was too long, got: %d, want mazimum of %d.", IPAddresLen, maxLen)
	}

	// https://www.regular-expressions.info/ip.html
	// https://regex101.com/
	regex := regexp.MustCompile(`\b(?:(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.){3}(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\b`)
	if regex.Match([]byte(IPAddress)) == false {
		t.Error("getIPAddress() was incorrect, valid IP address was not found.")
	}
}

func BenchmarkGetIPAddress(b *testing.B) {
	var res string

	for n := 0; n < b.N; n++ {
		res = getIPAddress()
	}

	resultString = res
}

func TestGetUsername(t *testing.T) {
	username := getUsername()

	if reflect.TypeOf(username) != reflect.TypeOf("") {
		t.Error("getUsername() was incorrect, didn't get a string.")
	}

	if len(username) == 0 {
		t.Error("getUsername() was incorrect, got zero-length string.")
	}
}

func BenchmarkGetUsername(b *testing.B) {
	var res string

	for n := 0; n < b.N; n++ {
		res = getUsername()
	}

	resultString = res
}

func TestGetHostUptime(t *testing.T) {
	uptime := getHostUptime()

	if reflect.TypeOf(uptime) != reflect.TypeOf("") {
		t.Error("getHostUptime() was incorrect, didn't get a string.")
	}

	if len(uptime) == 0 {
		t.Error("getHostUptime() was incorrect, got zero-length string.")
	}

	// \d+h\d+m\d+s matches e.g. 348h24m0s
	regex := regexp.MustCompile(`\d+h\d+m\d+s`)
	if regex.Match([]byte(uptime)) == false {
		t.Error("getHostUptime() was incorrect, valid uptime was not found.")
	}
}

func BenchmarkGetHostUptime(b *testing.B) {
	var res string

	for n := 0; n < b.N; n++ {
		res = getHostUptime()
	}

	resultString = res
}

func TestGetOSVersion(t *testing.T) {
	version := getOSVersion()

	if reflect.TypeOf(version) != reflect.TypeOf("") {
		t.Error("getOSVersion() was incorrect, didn't get a string.")
	}

	if len(version) == 0 {
		t.Error("getOSVersion() was incorrect, got zero-length string.")
	}
}

func BenchmarkGetOSVersion(b *testing.B) {
	var res string

	for n := 0; n < b.N; n++ {
		res = getOSVersion()
	}

	resultString = res
}

func TestGetRebootRequired(t *testing.T) {
	// reboot := getRebootRequired()
	// if len(reboot) == 0 {
	// 	t.Error("getRebootRequired() was incorrect, got zero-length string.")
	// }
}

func BenchmarkGetRebootRequired(b *testing.B) {
	var res bool

	for n := 0; n < b.N; n++ {
		res = getRebootRequired()
	}

	resultBool = res
}

func TestGetMemoryDetails(t *testing.T) {
	memory := getMemoryDetails()

	if reflect.TypeOf(memory) != reflect.TypeOf([4]int{}) {
		t.Error("getMemoryDetails() was incorrect, didn't get a slice of four ints.")
	}

	if len(memory) != 4 {
		t.Error("getMemoryDetails() was incorrect, didn't get a slice of four elements.")
	}

	if memory[0] == 0 {
		t.Error("getMemoryDetails() was incorrect, first slice was zero.")
	}
	if memory[1] == 0 {
		t.Error("getMemoryDetails() was incorrect, second slice was zero.")
	}
}

func BenchmarkGetMemoryDetails(b *testing.B) {
	var res [4]int

	for n := 0; n < b.N; n++ {
		res = getMemoryDetails()
	}

	resultSlice4Int = res
}

func TestGetDiskDetails(t *testing.T) {
	disk := getDiskDetails()

	if reflect.TypeOf(disk) != reflect.TypeOf([2]int{}) {
		t.Error("getDiskDetails() was incorrect, didn't get a slice of two ints.")
	}

	if len(disk) != 2 {
		t.Error("getDiskDetails() was incorrect, didn't get a slice of two elements.")
	}

	if disk[0] == 0 {
		t.Error("getDiskDetails() was incorrect, first slice was zero.")
	}
	if disk[1] == 0 {
		t.Error("getDiskDetails() was incorrect, second slice was zero.")
	}
}

func BenchmarkGetDiskDetails(b *testing.B) {
	var res [2]int

	for n := 0; n < b.N; n++ {
		res = getDiskDetails()
	}

	resultSlice2Int = res
}

func TestGetLoadAvgs(t *testing.T) {
	load := getLoadAvgs()

	if reflect.TypeOf(load) != reflect.TypeOf([3]float64{}) {
		t.Error("getLoadAvgs() was incorrect, didn't get a slice of three floats.")
	}

	if len(load) != 3 {
		t.Error("getLoadAvgs() was incorrect, didn't get a slice of three elements.")
	}
}

func BenchmarkGetLoadAvgs(b *testing.B) {
	var res [3]float64

	for n := 0; n < b.N; n++ {
		res = getLoadAvgs()
	}

	resultSlice3Float = res
}

func TestGetLastSuccessfulSend(t *testing.T) {
	// now := time.Now()
	var now time.Time
	lastSend := getLastSuccessfulSend(now)

	if reflect.TypeOf(lastSend) != reflect.TypeOf("") {
		t.Error("getLastSuccessfulSend() was incorrect, didn't get a string.")
	}

	regex := regexp.MustCompile(`Last successful data send`)
	if regex.Match([]byte(lastSend)) == false {
		t.Error("getLastSuccessfulSend() was incorrect, common string not found.")
	}

	regex = regexp.MustCompile(`Never`)
	if regex.Match([]byte(lastSend)) == false {
		t.Error("getLastSuccessfulSend() was incorrect, 'never' not found.")
	}
}

func BenchmarkGetLastSuccessfulSend(b *testing.B) {
	var res string

	for n := 0; n < b.N; n++ {
		res = getLastSuccessfulSend(time.Now())
	}

	resultString = res
}

func TestGetAppUptime(t *testing.T) {
	var now time.Time
	uptime := getAppUptime(now)

	if reflect.TypeOf(uptime) != reflect.TypeOf("") {
		t.Error("getAppUptime() was incorrect, didn't get a string.")
	}
}

func BenchmarkGetAppUptime(b *testing.B) {
	var res string

	for n := 0; n < b.N; n++ {
		res = getAppUptime(time.Now())
	}

	resultString = res
}

func TestGetCPUDetails(t *testing.T) {
	count, speed := getCPUDetails()

	if reflect.TypeOf(count) != reflect.TypeOf(0) {
		t.Error("getCPUDetails() was incorrect, didn't get an int.")
	}
	if reflect.TypeOf(speed) != reflect.TypeOf("") {
		t.Error("getCPUDetails() was incorrect, didn't get a string.")
	}
}

func BenchmarkGetCPUDetails(b *testing.B) {
	var (
		count int
		speed string
	)

	for n := 0; n < b.N; n++ {
		count, speed = getCPUDetails()
	}

	resultInt = count
	resultString = speed
}

func TestGetRandomHostname(t *testing.T) {
	hostname := getRandomHostname()

	if reflect.TypeOf(hostname) != reflect.TypeOf("") {
		t.Error("getRandomHostname() was incorrect, didn't get a string.")
	}

	if len(hostname) == 0 {
		t.Error("getRandomHostname() was incorrect, got zero-length string.")
	}
}

func BenchmarkGetRandomHostname(b *testing.B) {
	var res string

	for n := 0; n < b.N; n++ {
		res = getRandomHostname()
	}

	resultString = res
}

func TestGetRandomIPAddress(t *testing.T) {
	IPAddress := getRandomIPAddress()

	if reflect.TypeOf(IPAddress) != reflect.TypeOf("") {
		t.Error("getRandomIPAddress() was incorrect, didn't get a string.")
	}

	if len(IPAddress) == 0 {
		t.Error("getRandomIPAddress() was incorrect, got zero-length string.")
	}
}

func BenchmarkGetRandomIPAddress(b *testing.B) {
	var res string

	for n := 0; n < b.N; n++ {
		res = getRandomIPAddress()
	}

	resultString = res
}

func TestGetRandomUsername(t *testing.T) {
	username := getRandomUsername()

	if reflect.TypeOf(username) != reflect.TypeOf("") {
		t.Error("getRandomUsername() was incorrect, didn't get a string.")
	}

	if len(username) == 0 {
		t.Error("getRandomUsername() was incorrect, got zero-length string.")
	}
}

func BenchmarkGetRandomUsername(b *testing.B) {
	var res string

	for n := 0; n < b.N; n++ {
		res = getRandomUsername()
	}

	resultString = res
}
