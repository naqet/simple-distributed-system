package utils

import "os"

func GetPort(defaultPort string) string {
    port := os.Getenv("PORT")

    if port != "" {
        return port
    }

    return defaultPort
}
