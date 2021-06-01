package main

import (
    "fmt"
    "io/ioutil"
    "math"
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
    maxKernelContent, err := ioutil.ReadFile(maxBrightnessPath)
    if err != nil {
        return 0, fmt.Errorf("can't read max brightness path %s", maxBrightnessPath)
    }

    max, err := strconv.ParseFloat(strings.TrimSuffix(string(maxKernelContent), "\n"), 64)
    if err != nil {
        return 0, fmt.Errorf("can't read %s content: %s", maxBrightnessPath, err)
    }

    return max, nil
}

func calculatePercentage(backlightNumber float64) (int, error) {
    max, err := getMaxBacklight()
    if err != nil {
        return 0, err
    }

    if backlightNumber == 0 {
        return 0, nil
    } else if backlightNumber == max {
        return 100, nil
    }
    
    percentage := int((backlightNumber / max) * 100)

    return percentage, nil
}

func setCurrentBacklight(backlightNumber string, brightnessUnit float64, maxBrightness float64) float64 {
    var calculatedBrightness float64
    percentage, _ := strconv.ParseFloat(strings.Split(backlightNumber, "%")[0], 64)

    if percentage > 99 || percentage + brightnessUnit >= 100 {
        calculatedBrightness = maxBrightness
    } else if percentage <= 0 || percentage < brightnessUnit {
        calculatedBrightness = 0
    } else {
        calculatedBrightness = percentage * brightnessUnit
    }

    return calculatedBrightness
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
    currentBrightnessPercentage, err := calculatePercentage(currentBrightness)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    if os.Args[1] == "current" {
        percentageToPrint, err := calculatePercentage(currentBrightness)
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        toPrint := fmt.Sprintf("%d", percentageToPrint)
        if len(os.Args) == 3 && os.Args[2] == "--pretty" {
            toPrint = fmt.Sprintf("\uF0EB %d%%", percentageToPrint)
        } else if len(os.Args) == 3 && os.Args[2] == "--pretty2" {
            toPrint = fmt.Sprintf("\uF822 %d%%", percentageToPrint)
        }

        fmt.Println(toPrint)
        os.Exit(0)
    } else if os.Args[1] == "+" {
        if currentBrightness >= maxBrightness64 - brightnessUnit {
            setBrightness = maxBrightness64
        } else if currentBrightnessPercentage < 5 {
            setBrightness = currentBrightness + brightnessUnit
        } else {
            setBrightness = currentBrightness + (brightnessUnit * 5)
            if setBrightness > maxBrightness64 {
                setBrightness = maxBrightness64
            }
        }
    } else if os.Args[1] == "-" {
        if currentBrightnessPercentage < 1 {
            setBrightness = 0
        } else if currentBrightnessPercentage <= 5 {
            setBrightness = currentBrightness - brightnessUnit
        } else {
            setBrightness = currentBrightness - (brightnessUnit * 5)
            brightnessPercentage, _ := calculatePercentage(setBrightness)
            if brightnessPercentage < 6 {
                setBrightness = math.Ceil(5 * brightnessUnit)
            }
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
        err := f.Close()
        if err != nil {
            return
        }
        return
    }
    err = f.Close()
    if err != nil {
        fmt.Println(err)
        return
    }

    current, err := calculatePercentage(setBrightness)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    fmt.Printf("%d", current)
}
