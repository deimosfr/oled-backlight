# oled-backlight
Linux OLED backlight management for Linux intel cards

I've a Lenovo X1 Carbon Extreme Gen2 and recently the OLED support has been fixed (thanks to https://gitlab.freedesktop.org/drm/intel/issues/510)

Tested on Arch linux (Intel Corporation UHD Graphics 630):
* kernel 5.6.6-arch1-1
* xf86-video-intel 1:2.99.917+906+g846b53da-1

However common softs like xbacklight doesn't recognize the screen. So I've made this quick program in order to control backlight through the command line. See usage:


```
# print help
$ ./oled-backlight
Please add one of those argument: current|+|-|XY%

# return current percentage
$ ./oled-backlight current
50

# lighter +5%
$ sudo ./oled-backlight +
55

# darker -5%
$ sudo ./oled-backlight -
50

# set light percentage
$ sudo ./oled-backlight 60%
60
```

Note: it needs sudo permissions in order to modify the brightness
