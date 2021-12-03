package notify

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	ICINGA_LINK = "%s/monitoring/service/show?host=%s&service=%s"
)

type HostAlertConfig struct {
	HostAddress            string `env:"HOSTADDRESS"`
	HostDisplayName        string `env:"HOSTDISPLAYNAME"`
	HostName               string `env:"HOSTNAME"`
	HostOutput             string `env:"HOSTOUTPUT"`
	HostState              string `env:"HOSTSTATE"`
	LongDateTime           string `env:"LONGDATETIME"`
	NotificationAuthorName string `env:"NOTIFICATIONAUTHORNAME"`
	NotificationComment    string `env:"NOTIFICATIONCOMMENT"`
	NotificationType       string `env:"NOTIFICATIONTYPE"`

	// custom
	ContactName      string `env:"CONTACTNAME"`
	IcingaWebBaseURL string `env:"ICINGAWEBBASEURL"`
	AckLinkURL       string `env:"ACKLINKURL"`
}

type ServiceAlertConfig struct {
	HostAddress            string `env:"HOSTADDRESS"`
	HostDisplayName        string `env:"HOSTDISPLAYNAME"`
	HostName               string `env:"HOSTNAME"`
	LongDateTime           string `env:"LONGDATETIME"`
	NotificationAuthorName string `env:"NOTIFICATIONAUTHORNAME"`
	NotificationComment    string `env:"NOTIFICATIONCOMMENT"`
	NotificationType       string `env:"NOTIFICATIONTYPE"`
	ServiceDisplayName     string `env:"SERVICEDISPLAYNAME"`
	ServiceName            string `env:"SERVICENAME"`
	ServiceOutput          string `env:"SERVICEOUTPUT"`
	ServiceState           string `env:"SERVICESTATE"`

	// custom
	ContactName      string `env:"CONTACTNAME"`
	IcingaWebBaseURL string `env:"ICINGAWEBBASEURL"`
	ServiceWiki      string `env:"SERVICEWIKI"`
	AckLinkURL       string `env:"ACKLINKURL"`
}

type AckInfo struct {
	Service string `json:"service"`
	Host    string `json:"host"`
	User    string `json:"user"`
}

func getAckLink(apiURL, service, host, user string) string {
	if apiURL == "" {
		return ""
	}

	resp, err := http.PostForm(apiURL, url.Values{
		"service":  {service},
		"host":     {host},
		"user":     {user},
		"ack_type": {"1"},
	})
	if err != nil {
		return fmt.Sprintf("get ack link error: %s", err)
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("read ack link api response error: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("ack api %d: %s", resp.StatusCode, string(result))
	}
	return string(result)
}

func getIcingaLink(baseURL, host, service string) string {
	if baseURL == "" {
		return ""
	}

	return fmt.Sprintf(ICINGA_LINK, baseURL, host, service)
}
