package notify

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
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
	ContactName string `env:"CONTACTNAME"`
	AckLinkURL  string `env:"ACKLINKURL"`

	// generated
	ackLink string
}

func HostAlertFromEnv() (*HostAlertConfig, error) {
	var alert HostAlertConfig
	if err := env.Parse(&alert); err != nil {
		return nil, errors.Wrap(err, "parse host alert from env error")
	}
	alert.ackLink = getAckLink(alert.AckLinkURL, "", alert.HostName, alert.ContactName)
	return &alert, nil
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

	// generated
	ackLink      string
	icingaWebURL string
}

func ServiceAlertFromEnv() (*ServiceAlertConfig, error) {
	var alert ServiceAlertConfig
	if err := env.Parse(&alert); err != nil {
		return nil, errors.Wrap(err, "parse host alert from env error")
	}
	alert.icingaWebURL = getIcingaLink(alert.IcingaWebBaseURL, alert.HostName, alert.ServiceName)
	alert.ackLink = getAckLink(alert.AckLinkURL, "", alert.HostName, alert.ContactName)
	return &alert, nil
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
