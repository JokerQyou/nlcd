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
char timeString[9];

void cleanup(void)
{
    pinMode(pin_sclk, INPUT);
    pinMode(pin_din, INPUT);
    pinMode(pin_dc, INPUT);
    pinMode(pin_rst, INPUT);
    pinMode(pin_ce, INPUT);
}

char * get_time(void)
{
    time_t current_time;
    struct tm * time_info;
    time(&current_time);
    time_info = localtime(&current_time);
    strftime(timeString, sizeof(timeString), "%H:%M:%S", time_info);
    // printf("%s\n", &timeString);
    return &timeString;
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
    LCDcommand(PCD8544_DISPLAYCONTROL | PCD8544_DISPLAYALLON);
    LCDdisplay();
    LCDclear();
    LCDcommand(PCD8544_DISPLAYCONTROL | PCD8544_DISPLAYNORMAL);

    while (1)
    {
        LCDdrawrect(6 - 1, 6 - 1, LCDWIDTH - 6, LCDHEIGHT - 6, BLACK);
        LCDdrawstring(12, 12, get_time());
        LCDdisplay();
        delay(1000);
        // LCDclear();
    }

    LCDclear();
    cleanup();
    return 0;
}
