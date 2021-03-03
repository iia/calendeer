package main

import (
    "os"
    "fmt"
    "time"
    "net/http"
    "strings"
    "path/filepath"
    "github.com/gin-gonic/gin"
)

// Constants.

const S_OK uint8 = 1
const S_NOT_OK uint8 = 2
const P_FILE_PTH_ALT_FRIDAY string = ".persist/ALT_FRIDAY"

func postAlternateFriday(ctxGin *gin.Context, oddEven string) {
    var err error

    workDir, err := os.Getwd()

    if err != nil {
        fmt.Println(err)

        ctxGin.JSON(
            http.StatusInternalServerError,
            gin.H {
                "s": S_NOT_OK,
                "m": "Can't get CWD",
            },
        )

        return
    }

    fp, err := os.OpenFile(
        filepath.Join(workDir, P_FILE_PTH_ALT_FRIDAY),
        os.O_TRUNC | os.O_CREATE | os.O_WRONLY,
        0664,
    )

    if err != nil {
        fmt.Println(err)

        ctxGin.JSON(
            http.StatusInternalServerError,
            gin.H {
                "s": S_NOT_OK,
                "m": "File open failed",
            },
        )

        return
    }

    _, err = fp.WriteString(oddEven)

    if err != nil {
        fmt.Println(err)

        ctxGin.JSON(
            http.StatusInternalServerError,
            gin.H {
                "s": S_NOT_OK,
                "m": "File write failed",
            },
        )

        return
    }

    ctxGin.JSON(
        http.StatusOK,
        gin.H {
            "s": S_OK,
            "m": true,
        },
    )
}

func getIsAlternateFriday(ctxGin *gin.Context) {
    var h_stat int
    var h gin.H

    now := time.Now()
    _, _, day := now.Date()
    weekday := now.Weekday()

    workDir, err := os.Getwd()

    if err != nil {
        fmt.Println(err)

        ctxGin.JSON(
            http.StatusInternalServerError,
            gin.H {
                "s": S_NOT_OK,
                "m": "Can't get CWD",
            },
        )

        return
    }

    data, err := os.ReadFile(
        filepath.Join(workDir, P_FILE_PTH_ALT_FRIDAY),
    )

    if err != nil {
        fmt.Println(err)

        ctxGin.JSON(
            http.StatusInternalServerError,
            gin.H {
                "s": S_NOT_OK,
                "m": "File read failed",
            },
        )

        return
    }

    switch strings.Trim(string(data), "\n") {
        case "ODD": {
            if weekday == time.Friday && day % 2 != 0 {
                h_stat = http.StatusOK
                h = gin.H {
                    "s": S_OK,
                    "m": true,
                }
            } else {
                h_stat = http.StatusOK
                h = gin.H {
                    "s": S_OK,
                    "m": false,
                }
            }

            break
        }

        case "EVEN": {
            if weekday == time.Friday && day % 2 == 0 {
                h_stat = http.StatusOK
                h = gin.H {
                    "s": S_OK,
                    "m": true,
                }
            } else {
                h_stat = http.StatusOK
                h = gin.H {
                    "s": S_OK,
                    "m": false,
                }
            }

            break
        }

        default: {
            h_stat = http.StatusInternalServerError
            h = gin.H {
                "s": S_NOT_OK,
                "m": "Unknown check option",
            }
        }
    }

    ctxGin.JSON(h_stat, h)
}

func getTrashCalendar(ctxGin *gin.Context) {
    client := &http.Client{}
    req, err := http.NewRequest(
        "GET",
        "https://api.abfall.io/",
        nil,
    )

    if err != nil {
        fmt.Println(err)

        ctxGin.JSON(
            http.StatusInternalServerError,
            gin.H {
                "s": S_NOT_OK,
                "m": "Failed to form HTTP request",
            },
        )

        return
    }

    req.Header.Set(
        "User-Agent",
        "Mozilla/5.0 (X11; Linux x86_64) " +
        "AppleWebKit/537.36 (KHTML, like Gecko) " +
        "Chrome/88.0.4324.182 " +
        "Safari/537.36",
    )

    queryURL := req.URL.Query()

    /*
     * Abfall.io API stuff found in the offered
     * download link of iCalendar file.
     */
    queryURL.Add("t", "ics")
    queryURL.Add("s", "57a5732bbba87512418093fdde1497df")
    queryURL.Add("kh", "DaA02103019b46345f1998698563DaAd")

    req.URL.RawQuery = queryURL.Encode()

    res, err := client.Do(req)

    if err != nil {
        ctxGin.JSON(
            http.StatusInternalServerError,
            gin.H {
                "s": S_NOT_OK,
                "m": "Failed to do HTTP request",
            },
        )

        return
    }

    reader := res.Body

    defer reader.Close()

    extraHeaders := map[string]string {
        "Content-Disposition": `attachment; filename="trash_calendar.ics"`,
    }

    ctxGin.DataFromReader(
        http.StatusOK,
        res.ContentLength,
        res.Header.Get("Content-Type"),
        reader,
        extraHeaders,
    )
}

func main() {
    ginSrv := gin.Default()

    ginSrv.GET(
        "tentis/get/trash_calendar",
        func (ctxGin *gin.Context) {
            getTrashCalendar(ctxGin)
        },
    )

    ginSrv.GET(
        "tentis/get/is_alternate_friday",
        func (ctxGin *gin.Context) {
            getIsAlternateFriday(ctxGin)
        },
    )

    ginSrv.POST(
        "tentis/post/alternate_friday/odd",
        func (ctxGin *gin.Context) {
            postAlternateFriday(ctxGin, "ODD")
        },
    )

    ginSrv.POST(
        "tentis/post/alternate_friday/even",
        func (ctxGin *gin.Context) {
            postAlternateFriday(ctxGin, "EVEN")
        },
    )

    ginSrv.Run(":5000")
}
