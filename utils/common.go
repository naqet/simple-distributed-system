package utils

import (
	"net"
)

func GetServiceURL(port string) (string, error) {
    conn, err := net.Dial("udp", "0.0.0.0:8000")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	ip, _, err := net.SplitHostPort(conn.LocalAddr().String())

    if err != nil {
        return "", err
    }

    return "http://" + ip + ":" + port, nil
}
