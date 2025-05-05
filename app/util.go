package app

import (
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/abmpio/abmp/pkg/log"
	"github.com/abmpio/app/host"
	"github.com/abmpio/libx/net/ipaddr"
)

func setHealthCheckEnv() {
	healthcheck := host.GetHostEnvironment().GetEnvString(host.ENV_Healthcheck)
	if len(healthcheck) > 0 {
		// set by other logic
		return
	}

	http := host.GetHostEnvironment().GetEnvString(host.ENV_HTTP)
	if len(http) <= 0 {
		return
	}
	if !strings.HasPrefix("http://", http) {
		// if value is 127.0.0.1:8080, then append http scheme
		http = fmt.Sprintf("http://%s", http)
	}
	url, err := url.Parse(http)
	if err != nil {
		log.Logger.Warn("无效的http参数配置")
		return
	}
	advertiseHost := host.GetHostEnvironment().GetEnvString(host.ENV_AdvertiseHost)
	if len(advertiseHost) <= 0 {
		// used http
		advertiseHost = url.Hostname()
		// ip := net.ParseIP(advertiseHost)
		needParse := false
		if advertiseHost == "localhost" || ipaddr.IsAny(advertiseHost) || advertiseHost == "127.0.0.1" {
			needParse = true
		}
		ip := net.ParseIP(advertiseHost)
		if ip == nil {
			needParse = true
		} else if ip.IsLoopback() || ipaddr.IsAny(ip) {
			needParse = true
		}
		if needParse {
			// is localhost?
			addrList, _ := ipaddr.GetPrivateIPv4()
			if len(addrList) > 0 {
				advertiseHost = addrList[0].IP.String()
			}
		}
	}
	var healthcheckUrl string
	if len(url.Port()) > 0 {
		healthcheckUrl = fmt.Sprintf("%s://%s:%s/%s", url.Scheme, advertiseHost, url.Port(), "api/health/check")
	} else {
		healthcheckUrl = fmt.Sprintf("%s://%s/%s", url.Scheme, advertiseHost, "api/health/check")
	}
	host.GetHostEnvironment().SetEnv(host.ENV_Healthcheck, healthcheckUrl)
	host.GetHostEnvironment().SetEnv(host.ENV_AdvertiseHost, advertiseHost)
}
