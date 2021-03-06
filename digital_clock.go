package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/JokerQyou/rpi"
	"github.com/JokerQyou/rpi/pcd8544"
	"github.com/kidoman/embd"
	"github.com/kidoman/embd/sensor/bmp085"

	_ "github.com/kidoman/embd/host/rpi"
)

const (
	SCLK     = 0
	DIN      = 1
	DC       = 2
	CS       = 3
	RST      = 4
	CONTRAST = 60
)

func init() {
	rpi.WiringPiSetup()
	pcd8544.LCDInit(SCLK, DIN, DC, CS, RST, CONTRAST)
	pcd8544.LCDclear()
	pcd8544.LCDdisplay()

}

func gpio_cleanup() {
	pcd8544.LCDclear()
	pcd8544.LCDdisplay()

	rpi.PinMode(SCLK, rpi.INPUT)
	rpi.PinMode(DIN, rpi.INPUT)
	rpi.PinMode(DC, rpi.INPUT)
	rpi.PinMode(CS, rpi.INPUT)
	rpi.PinMode(RST, rpi.INPUT)
}

func get_time() (string, string) {
	t := time.Now()
	return t.Format("15:04:05"), t.Format("01-02 Mon")
}

func main() {
	sig := make(chan os.Signal, 1)
	if err := embd.InitI2C(); err != nil {
		panic(err)
	}
	defer embd.CloseI2C()

	bus := embd.NewI2CBus(1)
	baro := bmp085.New(bus)
	defer baro.Close()

	// Draw the outline
	pcd8544.LCDdrawrect(6-1, 6-1, pcd8544.LCDWIDTH-6, pcd8544.LCDHEIGHT-6, pcd8544.BLACK)
	// date string changes slowly
	var (
		date_str string
		temp_str string
	)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-sig
		fmt.Println(s)
		gpio_cleanup()
		os.Exit(0)
	}()

	keep_running := true
	for keep_running {
		// Get temperature
		temp, err := baro.Temperature()
		if err != nil {
			gpio_cleanup()
			panic(err)
		}
		_temp_str := fmt.Sprint(strconv.FormatFloat(temp, 'f', 2, 64), " C")

		time_str, _date_str := get_time()
		pcd8544.LCDdrawstring(20, 12, time_str)
		if date_str != _date_str {
			date_str = _date_str
			pcd8544.LCDdrawstring(18, 24, date_str)
		}
		if temp_str != _temp_str {
			temp_str = _temp_str
			pcd8544.LCDdrawstring(20, 36, temp_str)
		}
		pcd8544.LCDdisplay()
		// wait for 1 sec
		time.Sleep(time.Second)
	}
}
