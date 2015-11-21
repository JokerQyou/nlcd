#include <wiringPi.h>
#include <stdio.h>
#include <stdlib.h>
#include "PCD8544.h"

int pin_sclk = 4;
int pin_din = 3;
int pin_dc = 2;
int pin_rst = 0;
int pin_ce = 1;

// This is already tuned
int lcd_contrast = 60;

void gpio_cleanup(void)
{
    pinMode(pin_sclk, INPUT);
    pinMode(pin_din, INPUT);
    pinMode(pin_dc, INPUT);
    pinMode(pin_rst, INPUT);
    pinMode(pin_ce, INPUT);
}

int main(int argc, char const *argv[])
{
    printf("nLCD tool\n");
    printf("Pin definations:\n");
    printf("CLK on %i\n", pin_sclk);
    printf("DIN on %i\n", pin_din);
    printf("DC on %i\n", pin_dc);
    printf("CE on %i\n", pin_ce);
    printf("RST on %i\n", pin_rst);

    if (wiringPiSetup() == -1)
    {
        printf("wiringPi error\n");
        exit(1);
    }

    LCDInit(pin_sclk, pin_din, pin_dc, pin_ce, pin_rst, lcd_contrast);
    LCDclear();

    int i;
    for (i = 0; i < 3; ++i)
    {
        LCDcommand(PCD8544_DISPLAYCONTROL | PCD8544_DISPLAYALLON);
        delay(1000);
        LCDcommand(PCD8544_DISPLAYCONTROL | PCD8544_DISPLAYNORMAL);
        delay(1000);
    }
    LCDclear();

    LCDclear();
    gpio_cleanup();
    return 0;
}
