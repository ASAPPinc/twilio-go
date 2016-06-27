package twilio

import (
	//"fmt"

	"fmt"
	"net/url"
	"strconv"
)

const pathPart = "Messages"

type MessageService struct {
	client *Client
}

type Message struct {
	Body  string
	From  string
	To    string
	Price string
	Sid   string
}

type MessageDetails struct {
	Sid          string `json:"sid"`
	DateCreated  string `json:"date_created"`
	DateUpdated  string `json:"date_updated"`
	DateSent     string `json:"date_sent"`
	AccountSid   string `json:"account_sid"`
	To           string `json:"to"`
	From         string `json:"from"`
	Body         string `json:"body"`
	Status       string `json:"status"`
	NumSegments  string `json:"num_segments"`
	NumMedia     string `json:"num_media"`
	Direction    string `json:"direction"`
	ApiVersion   string `json:"api_version"`
	Price        string `json:"price"`
	PriceUnit    string `json:"price_unit"`
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Uri          string `json:"uri"`
}

type MessagePage struct {
	Messages        []MessageDetails `json:"messages"`
	Uri             string           `json:"uri"`
	PageSize        int              `json:"page_size"`
	Start           int              `json:"start"`
	NextPageUri     string           `json:"next_page_uri"`
	NumPages        int              `json:"num_pages"`
	Total           int              `json:"total"`
	LastPageUri     string           `json:"last_page_uri"`
	Page            int              `json:"page"`
	FirstPageUri    string           `json:"first_page_uri"`
	End             int              `json:"end"`
	PreviousPageUri string           `json:"previous_page_uri"`
}

type MessageIterator struct {
	pos      int
	messages []MessageDetails
	params   url.Values
	client   *Client
}

func (m *MessageService) Create(data url.Values) (Message, error) {
	msg := new(Message)
	resp, err := m.client.MakeRequest("POST", pathPart, data, msg)
	if err != nil {
		return *msg, err
	}
	if resp.StatusCode < 200 && resp.StatusCode > 299 {
		return *msg, fmt.Errorf("request not successful, status=%s, statusCode=%d", resp.Status, resp.StatusCode)
	}
	return *msg, nil
}

func (m *MessageService) SendMessage(from string, to string, body string, mediaUrls []url.URL) (Message, error) {
	v := url.Values{}
	v.Set("Body", body)
	v.Set("From", from)
	v.Set("To", to)
	if mediaUrls != nil {
		for _, mediaUrl := range mediaUrls {
			v.Add("MediaUrl", mediaUrl.String())
		}
	}
	return m.Create(v)
}

func (m *MessageService) ListMessages(messagesPerPage int) (iterator *MessageIterator) {
	params := url.Values{"PageSize": []string{strconv.Itoa(messagesPerPage)}}
	return &MessageIterator{0, nil, params, m.client}
}

func (m *MessageIterator) fetch() (err error) {
	var page MessagePage
	_, err = m.client.ListResource(pathPart, m.params, &page)
	if err != nil {
		return
	}
	nextPageUrl, err := url.Parse(page.NextPageUri)
	if err != nil {
		return
	}
	m.params = nextPageUrl.Query()
	m.messages = page.Messages
	m.pos = 0
	return
}

// Returns nil when the list is complete.
func (m *MessageIterator) Next() (*MessageDetails, error) {
	if m.pos >= len(m.messages) {
		err := m.fetch()
		if err != nil {
			return nil, err
		}
	}
	if len(m.messages) == 0 {
		return nil, nil
	}
	message := &m.messages[m.pos]
	m.pos += 1
	return message, nil
}
