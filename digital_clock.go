package main

import (
    "time"
    "github.com/JokerQyou/rpi"
    "github.com/JokerQyou/rpi/pcd8544"
)

const (
    SCLK = 4
    DIN = 3
    DC = 2
    CS = 1
    RST = 0
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

func get_time() string {
    t := time.Now()
    return t.Format("15:04:05")
}

func main() {
    keep_running := true
    for keep_running {
        pcd8544.LCDdrawrect(6 - 1, 6 - 1, pcd8544.LCDWIDTH - 6, pcd8544.LCDHEIGHT - 6, pcd8544.BLACK)
        pcd8544.LCDdrawstring(20, 12, get_time())
        pcd8544.LCDdisplay()
        // wait for 1 sec
        time.Sleep(time.Second)
    }
}
