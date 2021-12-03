package notify

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

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
	HostDurationSec        string `env:"HOSTDURATIONSEC"`
	LongDateTime           string `env:"LONGDATETIME"`
	NotificationAuthorName string `env:"NOTIFICATIONAUTHORNAME"`
	NotificationComment    string `env:"NOTIFICATIONCOMMENT"`
	NotificationType       string `env:"NOTIFICATIONTYPE"`
	ContactName            string `env:"CONTACTNAME"`

	// custom
	AckLinkURL string `env:"ACKLINKURL"`

	// generated
	ackLink string
}

func HostAlertFromEnv() (*HostAlertConfig, error) {
	var alert HostAlertConfig
	if err := env.Parse(&alert); err != nil {
		return nil, errors.Wrap(err, "parse host alert from env error")
	}
	if alert.NotificationType == "PROBLEM" {
		alert.ackLink = getAckLink(alert.AckLinkURL, "", alert.HostName, alert.ContactName)
	}
	return &alert, nil
}

func (c *HostAlertConfig) Subject() string {
	return fmt.Sprintf("Host %s alert for %s(%s)!",
		c.HostState, c.HostName, c.HostAddress)
}

func (c *HostAlertConfig) Duration() string {
	du, _ := strconv.ParseFloat(c.HostDurationSec, 64)
	return (time.Duration(du) * time.Second).String()
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
	ServiceDurationSec     string `env:"SERVICEDURATIONSEC"`
	ContactName            string `env:"CONTACTNAME"`

	// custom
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
	if alert.NotificationType == "PROBLEM" {
		alert.ackLink = getAckLink(alert.AckLinkURL, "", alert.HostName, alert.ContactName)
	}
	return &alert, nil
}

func (c *ServiceAlertConfig) Subject() string {
	return fmt.Sprintf("%s - %s/%s is %s",
		c.NotificationType,
		c.HostDisplayName,
		c.ServiceDisplayName,
		c.ServiceState)
}

func (c *ServiceAlertConfig) Duration() string {
	du, _ := strconv.ParseFloat(c.ServiceDurationSec, 64)
	return (time.Duration(du) * time.Second).String()
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
