package flowdock

import (
	"encoding/json"
	"time"
	"fmt"
	cmd "github.com/milanaleksic/flowdock_stats/cmd_colors"
	"net/http"
	"strconv"
)

type User struct {
	Nick string `json: "nick"`
}

type Message struct {
	User      string          `json:"user"`
	Edited    int             `json:"edited"`
	Content   json.RawMessage `string:"content,omitempty"`
	Id        int             `json:"id"`
	CreatedAt CustomTime      `json:"created_at"`
}

type CustomTime struct {
	time.Time
}

var nilTime = (time.Time{}).UnixNano()

const jsonFormatDate = "2006-01-02T15:04:05.000Z"

func (ct *CustomTime) IsSet() bool {
	return ct.UnixNano() != nilTime
}

func (d *CustomTime) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b) - 1] == '"' {
		b = b[1 : len(b) - 1]
	}
	d.Time, err = time.Parse(jsonFormatDate, string(b))
	return err
}

type FlowdockApi struct {
	ApiToken string
}

// Which company's flow should be scanned for messages? Also, until when do we want to get the messages.
// In case you only want the last block of the messages, until should be set to -1.
func (api *FlowdockApi) GetMessages(company string, flow string, until int) (messages []Message, err error) {
	flowLocation := fmt.Sprintf("https://api.flowdock.com/flows/%s/%s/messages", company, flow)
	request, err := http.NewRequest("GET", flowLocation, nil)
	if err != nil {
		cmd.Warn(fmt.Sprintf("Error encountered while fetching: %v", err))
		return
	}

	request.SetBasicAuth(api.ApiToken, "")
	values := request.URL.Query()
	if until != -1 {
		values.Add("until_id", strconv.Itoa(until))
	}
	values.Add("limit", "100")
	request.URL.RawQuery = values.Encode()

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		cmd.Warn(fmt.Sprintf("Error encountered while fetching: %v", err))
		return
	}
	if resp.StatusCode != 200 {
		cmd.Warn(fmt.Sprintf("Server responded with a non 200: %v", resp.StatusCode))
		return
	}
	defer resp.Body.Close()
	messages = make([]Message, 100)
	err = json.NewDecoder(resp.Body).Decode(&messages)
	if err != nil {
		cmd.Warn(fmt.Sprintf("Error encountered while parsing: %v", err))
		return
	}
	return
}

// Get a "Nick name" of a user in Flowdock based on his/hers numeric ID
func (api *FlowdockApi) GetUser(userId string) (userNick string, err error) {
	request, err := http.NewRequest("GET", "https://api.flowdock.com/users/" + userId, nil)
	if err != nil {
		cmd.Warn(fmt.Sprintf("Error encountered while fetching: %v", err))
		return
	}
	request.SetBasicAuth(api.ApiToken, "")
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		cmd.Warn(fmt.Sprintf("Error encountered while fetching: %v", err))
		return
	}
	defer resp.Body.Close()
	result := User{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		cmd.Warn(fmt.Sprintf("Error encountered while parsing: %v", err))
		return
	}
	userNick = result.Nick
	return
}

func printOutContents(msg Message) {
	var content string
	err := json.Unmarshal(msg.Content, &content)
	if err != nil {
		cmd.Warn(fmt.Sprint("Could not parse contents: ", err))
	}
	cmd.Info(fmt.Sprintf(msg.User, " said: ", content, " id=", msg.Id, " at ", msg.CreatedAt))
}