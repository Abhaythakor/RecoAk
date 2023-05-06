package subdomain

import (
    "fmt"
    "log"
    "os"
    "os/exec"
    "strings"
    "time"
    "bufio"
    "bytes"
    "github.com/schollz/progressbar/v3"
    "io"

)

func amassEnum(path string) {
    subdomainFile := fmt.Sprintf("%s/Subdomain/subdomain.txt", path)
    outputDir := fmt.Sprintf("%s/Subdomain", path)


    // Time 
    t := time.Now()
    Hour := t.Format("15")
    Minute := t.Format("04")
    Second := t.Format("05")

    // Count the number of lines in the subdomain file
    subdomainFileHandle, err := os.Open(subdomainFile)
    if err != nil {
        log.Fatalf("Error opening subdomain file: %v", err)
    }
    defer subdomainFileHandle.Close()
    subdomainScanner := bufio.NewScanner(subdomainFileHandle)
    lineCount := 0
    for subdomainScanner.Scan() {
        lineCount++
    }

    // Create progress bar
    bar := progressbar.Default(int64(lineCount))

    // Run amass command
    cmd := exec.Command("amass", "enum", "-ip", "-df", subdomainFile, "-dir", outputDir, "-active", "-config", "/root/config.ini")
    var out bytes.Buffer
    cmd.Stdout = io.MultiWriter(&out, bar) // Use MultiWriter to write to both the buffer and the progress bar
    err = cmd.Run()
    if err != nil {
        log.Fatalf("Error running amass: %v", err)
    }
    fmt.Printf("Amass complete at (\033[34m%d:%d:%d\033[0m)\n", Hour, Minute, Second)

    // Process amass output
    lines := strings.Split(out.String(), "\n")
    for _, line := range lines {
        if strings.HasPrefix(line, "[") && strings.Contains(line, "IPv4") {
            fields := strings.Fields(line)
            if len(fields) >= 3 {
                ip := fields[2]
                fmt.Printf("%s\n", ip)
            }
        }
    }
}
