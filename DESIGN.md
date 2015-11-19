Current plan is to build a binary file which could work like this:

* Read configuration from a text file
* Determine GPIO pins the LCD hardware mounted on
* Run the binary will disply a `hello world` text on the LCD, for 2 seconds
* Pass text content to the binary to display it
* You can specify how long the text should stay on the screen

General usage is like this:

```bash
nlcd -t 5 "Hi there"  # Cause the text "Hi there" to be displayed on the LCD for 5 seconds
nlcd --time 5 "Hi there"  # Same as above
nlcd "Hi" -c /home/john/lcd.config  # Use specified LCD in given config file
nlcd "Hi" --config /home/john/lcd.config  # Same as above
nlcd -h  # Display help text
nlcd --help  # Same as above
nlcd -v  # Display version text
nlcd --version  # Same as above
```

You should specify the physical pin number of the following wires:

* RST
* CS / CE
* DC
* DIN
* CLK / SCLK

You can put your config file to `/etc/nlcd/default`, so that you can omit the config file in the commandline arguments.
