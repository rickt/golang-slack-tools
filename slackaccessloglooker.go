package main

import (
   "encoding/json"
   "fmt"
   "log"
   "net/http"
   "strings"
   "time"
)

type SlackAccessLog struct {
   Status     bool                  `json:"ok"`
   Logins     []SlackAccessLogEntry `json:"logins"`
   PagingData SlackAccessLogPaging  `json:"paging"`
}

type SlackAccessLogEntry struct {
   UserID    string `json:"user_id"`
   Username  string `json:"username"`
   DateFirst int64  `json:"date_first"`
   DateLast  int64  `json:"date_last"`
   Count     int    `json:"count"`
   IP        string `json:"ip"`
   UserAgent string `json:"user_agent"`
   ISP       string `json:"isp"`
   Country   string `json:"country"`
   Region    string `json:"region"`
}

type SlackAccessLogPaging struct {
   Count int `json:"count"`
   Total int `json:"total"`
   Page  int `json:"page"`
   Pages int `json:"pages"`
}

// helper func to do a case-insensitive search
func caseinsensitivecontains(a, b string) bool {
   return strings.Contains(strings.ToUpper(a), strings.ToUpper(b))
}

var (
   page     int    = 1
   pages    int    = 101
   token    string = "REDACTED"
   slackurl string = "https://slack.com/api/team.accessLogs"
)

func main() {
   var sal SlackAccessLog
   // 100 pages of JSON max from the team.accessLogs Slack API
   for page := 1; page < pages; page++ {
      // build the url
      url := fmt.Sprintf("%s?token=%s&page=%d", slackurl, token, page)
      // create the request
      req, err := http.NewRequest("GET", url, nil)
      if err != nil {
         log.Fatal("error: %s", err)
      }
      // create the http client
      client := &http.Client{}
      // get the response
      response, err := client.Do(req)
      if err != nil {
         log.Fatal("error: %s", err)
      }
      defer response.Body.Close()
      // decode the JSON response into our SlackAccessLog var
      if err := json.NewDecoder(response.Body).Decode(&sal); err != nil {
         log.Println(err)
      }
      // range through this page of the response and ignore Slack App/Android/iPhone useragents
      for _, dj := range sal.Logins {
         if !(caseinsensitivecontains(dj.UserAgent, "Slack_SSB") || caseinsensitivecontains(dj.UserAgent, "Android") || caseinsensitivecontains(dj.UserAgent, "iPhone") || caseinsensitivecontains(dj.UserAgent, "iPad") || caseinsensitivecontains(dj.UserAgent, "ApiApp")) {
            fmt.Printf("%s\t%s\t%s\n", dj.Username, time.Unix(dj.DateLast, 0), dj.UserAgent)
         }
      }
   }
}
