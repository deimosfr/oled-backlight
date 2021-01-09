package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "strconv"
    "strings"
)

func getCurrentBacklight(brightnessPath string) string {
    data, err := ioutil.ReadFile(brightnessPath)
    if err != nil {
        fmt.Println("File reading error", err)
        return "0"
    }
    return strings.TrimSuffix(string(data), "\n")
}

func getCurrentBacklightPercentage(backlightNumber int) string {
    max := 2000.0
    if backlightNumber == 0 {
        return "0%"
    }
    backlightNumberFloat := float64(backlightNumber)
    percentage := int((backlightNumberFloat / max) * 100)
    return strconv.Itoa(percentage)
}

func setCurrentBacklightPercentage(backlightNumber string) int {
    percentage, _ := strconv.Atoi(strings.Split(backlightNumber, "%")[0])
    if percentage >= 100 {
        percentage = 100
    } else if percentage <= 0 {
        percentage = 0
    }
    return percentage * 20
}

func help() {
    fmt.Print("Please add one of those argument: current|+|-|XY%\n--pretty: used with 'current' argument prints percentage and a lamp")
}

func main() {
    var currentBrightness, setBrightness int
    brightnessPath := "/sys/class/backlight/intel_backlight/brightness"

    if len(os.Args) < 2 {
        help()
        os.Exit(1)
    }

    currentBrightness, _ = strconv.Atoi(getCurrentBacklight(brightnessPath))
    if os.Args[1] == "current" {
        toPrint := getCurrentBacklightPercentage(currentBrightness)
        if len(os.Args) == 3 && os.Args[2] == "--pretty" {
            toPrint = "\uF0EB " + toPrint + "%"
        } else if len(os.Args) == 3 && os.Args[2] == "--pretty2" {
            toPrint = "\uF822 " + toPrint + "%"
        }
        fmt.Println(toPrint)
        os.Exit(0)
    } else if os.Args[1] == "+" {
        if currentBrightness >= 2000 {
            setBrightness = 2000
        } else if currentBrightness < 100 {
            setBrightness = currentBrightness + 20
        } else {
            setBrightness = currentBrightness + 100
        }
    } else if os.Args[1] == "-" {
        if currentBrightness <= 0 {
            setBrightness = 0
        } else if currentBrightness <= 100 {
            setBrightness = currentBrightness - 20
        } else {
            setBrightness = currentBrightness - 100
        }
    } else if strings.Contains(os.Args[1], "%") {
        setBrightness = setCurrentBacklightPercentage(os.Args[1])
    } else {
       help()
       os.Exit(1)
    }

    fmt.Println(getCurrentBacklightPercentage(setBrightness))
    f, err := os.Create(brightnessPath)
    if err != nil {
        fmt.Println(err)
        return
    }

    _, err = f.WriteString(strconv.Itoa(setBrightness))
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
}
