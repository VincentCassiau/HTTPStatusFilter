package main

import (
    "bufio"
    "flag"
    "fmt"
    "net/http"
    "os"
    "strconv"
    "strings"
    "sync"
    "time"
)

const version = "1.0.0"

var (
    file        string
    status      string
    delay       int
    concurrency int
    showVersion bool
    showHelp    bool
)

func init() {
    flag.StringVar(&file, "f", "", "File containing URLs (one per line)")
    flag.StringVar(&status, "s", "", "Expected status codes (comma-separated) or ranges (e.g. 200,201,300-399)")
    flag.IntVar(&delay, "d", 0, "Delay between requests in milliseconds (optional)")
    flag.IntVar(&concurrency, "c", 5, "Number of concurrent requests (optional)")
    flag.BoolVar(&showVersion, "v", false, "Show version")
    flag.BoolVar(&showHelp, "h", false, "Show help")
}

func main() {
    flag.Parse()

    if showVersion {
        fmt.Println("HTTPStatusFilter version:", version)
        return
    }

    if showHelp {
        printUsage()
        return
    }

    if file == "" || status == "" {
        printUsage()
        return
    }

    expectedStatuses := parseStatusCodes(status)
    urls, err := readLines(file)
    if err != nil {
        fmt.Println("Error reading file:", err)
        return
    }

    urlChan := make(chan string, len(urls))
    resultChan := make(chan string, len(urls))
    var wg sync.WaitGroup

    for i := 0; i < concurrency; i++ {
        wg.Add(1)
        go worker(urlChan, resultChan, expectedStatuses, delay, &wg)
    }

    for _, url := range urls {
        urlChan <- url
    }
    close(urlChan)

    go func() {
        wg.Wait()
        close(resultChan)
    }()

    for result := range resultChan {
        fmt.Println(result)
    }
}

func parseStatusCodes(statusStr string) map[int]bool {
    statuses := make(map[int]bool)
    parts := strings.Split(statusStr, ",")
    for _, part := range parts {
        if strings.Contains(part, "-") {
            rangeParts := strings.Split(part, "-")
            start, err1 := strconv.Atoi(rangeParts[0])
            end, err2 := strconv.Atoi(rangeParts[1])
            if err1 == nil && err2 == nil {
                for i := start; i <= end; i++ {
                    statuses[i] = true
                }
            }
        } else {
            status, err := strconv.Atoi(part)
            if err == nil {
                statuses[status] = true
            }
        }
    }
    return statuses
}

func readLines(file string) ([]string, error) {
    f, err := os.Open(file)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    var lines []string
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}

func worker(urls <-chan string, results chan<- string, expectedStatuses map[int]bool, delay int, wg *sync.WaitGroup) {
    defer wg.Done()
    for url := range urls {
        time.Sleep(time.Duration(delay) * time.Millisecond)
        resp, err := http.Get(url)
        if err != nil {
            fmt.Println("Error:", err)
            continue
        }
        defer resp.Body.Close()

        if expectedStatuses[resp.StatusCode] {
            results <- url
        }
    }
}

func printUsage() {
    fmt.Println("Usage: HTTPStatusFilter -f <url_file> -s <status_codes> [-d <milliseconds>] [-c <number>]")
    fmt.Println("Options:")
    fmt.Println("  -f : file containing URLs (required)")
    fmt.Println("  -s : expected status codes (required)")
    fmt.Println("  -d : delay between requests in milliseconds (optional)")
    fmt.Println("  -c : number of concurrent requests (optional)")
    fmt.Println("  -v : show version")
    fmt.Println("  -h : show help")
}
