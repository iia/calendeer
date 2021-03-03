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
const ALT_FRI_ODD int = 1
const ALT_FRI_EVEN int = 2

func getIsAlternateFriday(ctxGin *gin.Context, oddEven int) {
    var h gin.H
    var h_stat int
    now := time.Now()
    _, _, day := now.Date()
    weekday := now.Weekday()

    if weekday != time.Friday {
        h = gin.H {
            "s": S_OK,
            "m": false,
        }
        h_stat = http.StatusOK
    } else {
        switch oddEven {
            case ALT_FRI_ODD: {
                if (weekday == time.Friday) && (day % 2 != 0) {
                    h = gin.H {
                        "s": S_OK,
                        "m": true,
                    }
                    h_stat = http.StatusOK
                } else {
                    h = gin.H {
                        "s": S_OK,
                        "m": false,
                    }
                    h_stat = http.StatusOK
                }

                break
            }

            case ALT_FRI_EVEN: {
                if (weekday == time.Friday) && (day % 2 == 0) {
                    h = gin.H {
                        "s": S_OK,
                        "m": true,
                    }
                    h_stat = http.StatusOK
                } else {
                    h = gin.H {
                        "s": S_OK,
                        "m": false,
                    }
                    h_stat = http.StatusOK
                }

                break
            }

            default: {
                h = gin.H {
                    "s": S_NOT_OK,
                    "m": "Unknown check option",
                }
                h_stat = http.StatusInternalServerError
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
     * download link of the iCalendar file.
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
        "tentis/get/is_alternate_friday/odd",
        func (ctxGin *gin.Context) {
            getIsAlternateFriday(ctxGin, ALT_FRI_ODD)
        },
    )

    ginSrv.GET(
        "tentis/get/is_alternate_friday/even",
        func (ctxGin *gin.Context) {
            getIsAlternateFriday(ctxGin, ALT_FRI_EVEN)
        },
    )

    ginSrv.Run(":5000")
}
