package main

import (
    "fmt"
    "net/http"
    "github.com/gin-gonic/gin"
)

func getTrashCalendar(ctxGin *gin.Context) {
    client := &http.Client{}
    req, err := http.NewRequest(
        "GET",
        "https://api.abfall.io/",
        nil,
    )

    if err != nil {
        fmt.Println(err)

        ctxGin.Status(http.StatusServiceUnavailable)

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

    queryURL.Add("t", "ics")
    queryURL.Add("s", "57a5732bbba87512418093fdde1497df")
    queryURL.Add("kh", "DaA02103019b46345f1998698563DaAd")

    req.URL.RawQuery = queryURL.Encode()

    res, err := client.Do(req)

    if err != nil {
        fmt.Println(err)

        ctxGin.Status(http.StatusServiceUnavailable)

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

    ginSrv.Run(":5000")
}
