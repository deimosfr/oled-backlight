package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "strconv"
    "strings"
)

const brightnessPath = "/sys/class/backlight/intel_backlight/brightness"
const maxBrightnessPath = "/sys/class/backlight/intel_backlight/max_brightness"

func getCurrentBacklight(brightnessPath string) string {
    data, err := ioutil.ReadFile(brightnessPath)
    if err != nil {
        fmt.Println("File reading error", err)
        return "0"
    }
    return strings.TrimSuffix(string(data), "\n")
}

func getMaxBacklight() (float64, error) {
    max_kernel_content, err := ioutil.ReadFile(maxBrightnessPath)
    if err != nil {
        return 0, fmt.Errorf("Can't read max brightness path %s", maxBrightnessPath)
    }

    max, err := strconv.ParseFloat(strings.TrimSuffix(string(max_kernel_content), "\n"), 64)
    if err != nil {
        return 0, fmt.Errorf("Can't read %s content: %s", maxBrightnessPath, err)
    }

    return max, nil
}

func getCurrentBacklightPercentage(backlightNumber float64) (string, error) {
    max, err := getMaxBacklight()
    if err != nil {
        return "", err
    }

    if backlightNumber == 0 {
        return "0%", nil
    }

    backlightNumberFloat := backlightNumber
    percentage := int((backlightNumberFloat / max) * 100)

    return strconv.Itoa(percentage), nil
}

func setCurrentBacklight(backlightNumber string, brightnessUnit float64, maxBrightness float64) float64 {
    var calculated_brightness float64
    percentage, _ := strconv.ParseFloat(strings.Split(backlightNumber, "%")[0], 64)

    if percentage > 99 || percentage + brightnessUnit >= 100 {
        calculated_brightness = maxBrightness
    } else if percentage <= 0 || percentage < brightnessUnit {
        calculated_brightness = 0
    } else {
        calculated_brightness = percentage * brightnessUnit
    }

    return calculated_brightness
}

func help() {
    fmt.Print("Please add one of those argument: current|+|-|XY%\n--pretty: used with 'current' argument prints percentage and a lamp")
}

func main() {
    var setBrightness float64

    if len(os.Args) < 2 {
        help()
        os.Exit(1)
    }

    currentBrightness, _ := strconv.ParseFloat(getCurrentBacklight(brightnessPath), 64)
    maxBrightness64, err := getMaxBacklight()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    brightnessUnit := maxBrightness64 / 100
    currentBrightnessPercentageString, err := getCurrentBacklightPercentage(currentBrightness)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    currentBrightnessPercentage, _ := strconv.Atoi(currentBrightnessPercentageString)

    if os.Args[1] == "current" {
        toPrint, err := getCurrentBacklightPercentage(currentBrightness)
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        if len(os.Args) == 3 && os.Args[2] == "--pretty" {
            toPrint = "\uF0EB " + toPrint + "%"
        } else if len(os.Args) == 3 && os.Args[2] == "--pretty2" {
            toPrint = "\uF822 " + toPrint + "%"
        }

        fmt.Println(toPrint)
        os.Exit(0)
    } else if os.Args[1] == "+" {
        if currentBrightness >= maxBrightness64 - brightnessUnit {
            setBrightness = maxBrightness64
        } else if currentBrightnessPercentage <= 5 {
            setBrightness = currentBrightness + brightnessUnit
        } else {
            setBrightness = currentBrightness + (brightnessUnit * 5)
        }
    } else if os.Args[1] == "-" {
        if currentBrightnessPercentage < 1 {
            setBrightness = 0
        } else if currentBrightnessPercentage <= 5 {
            setBrightness = currentBrightness - brightnessUnit
        } else {
            setBrightness = currentBrightness - + (brightnessUnit * 5)
        }
    } else if strings.Contains(os.Args[1], "%") {
        setBrightness = setCurrentBacklight(os.Args[1], brightnessUnit, maxBrightness64)
    } else {
       help()
       os.Exit(1)
    }

    f, err := os.Create(brightnessPath)
    if err != nil {
        fmt.Println(err)
        return
    }

    setBrightnessString := fmt.Sprintf("%d", int(setBrightness))
    _, err = f.WriteString(setBrightnessString)
    if err != nil {
        fmt.Println(err)
        f.Close()
        return
    }
    err = f.Close()
    if err != nil {
        fmt.Println(err)
        return
    }

    current, err := getCurrentBacklightPercentage(setBrightness)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    fmt.Println(current)
}
