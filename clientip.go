package clientip

import (
	"net"
	"net/http"
	"strings"
)

const headerXForwardedFor = "x-forwarded-for"

var headerIPs = []string{
	"x-client-ip", "cf-connecting-ip", "fastly-client-ip", "true-client-ip",
	"x-real-ip", "x-cluster-client-ip", "x-forwarded", "forwarded-for",
}

// FromRequest returns the client IP address from the HTTP request
func FromRequest(r *http.Request) net.IP {
	// Load-balancers (AWS ELB) or proxies.
	if ip := fromXForwardedFor(r.Header.Get("x-forwarded-for")); ip != nil {
		return ip
	}

	for _, header := range headerIPs {
		if ip := net.ParseIP(r.Header.Get(header)); ip != nil {
			return ip
		}
	}

	remoteAddr := r.RemoteAddr
	if raddr, ok := splitHostPort(remoteAddr); ok {
		remoteAddr = raddr
	}

	return net.ParseIP(remoteAddr)
}

func fromXForwardedFor(xfwdfor string) net.IP {
	// x-forwarded-for may return multiple IP addresses in the format:
	// "client IP, proxy 1 IP, proxy 2 IP"
	// Therefore, the right-most IP address is the IP address of the most recent proxy
	// and the left-most IP address is the IP address of the originating client.
	// source: http://docs.aws.amazon.com/elasticloadbalancing/latest/classic/x-forwarded-headers.html
	// Azure Web App's also adds a port for some reason, so we'll only use the first part (the IP)
	for _, ip := range strings.Split(xfwdfor, ",") {
		ip = strings.TrimSpace(ip)
		if raddr, ok := splitHostPort(ip); ok {
			ip = raddr
		}

		// Sometimes IP addresses in this header can be 'unknown' (http://stackoverflow.com/a/11285650).
		// Therefore taking the left-most IP address that is not unknown
		// A Squid configuration directive can also set the value to "unknown" (http://www.squid-cache.org/Doc/config/forwarded_for/)
		if parsedIP := net.ParseIP(ip); parsedIP != nil {
			return parsedIP
		}
	}

	return nil
}

func splitHostPort(addr string) (string, bool) {
	raddr, _, err := net.SplitHostPort(addr)
	return raddr, raddr != "" && err == nil
}
