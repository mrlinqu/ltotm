package main

import (
	"flag"
	"os"
)

type cmdArgs struct {
	iface      string
	port       string
	tlsCert    string
	tlsKey     string
	storageDir string
}

func parseArgs() cmdArgs {
	ret := cmdArgs{}

	flag.StringVar(&ret.iface, "iface", envOrDefault("LTOTM_IFACE", "0.0.0.0"), "listen interface")
	flag.StringVar(&ret.port, "port", envOrDefault("LTOTM_PORT", ""), "listen port")

	flag.StringVar(&ret.tlsCert, "cert", envOrDefault("LTOTM_TLS_CERT", ""), "tls cert")
	flag.StringVar(&ret.tlsKey, "key", envOrDefault("LTOTM_TLS_KEY", ""), "tls key")

	flag.StringVar(&ret.storageDir, "dir", envOrDefault("LTOTM_STORAGE", "./data"), "path to storage dir")

	flag.Parse()

	if ret.port == "" {
		if ret.tlsCert == "" || ret.tlsKey == "" {
			ret.port = "80"
		} else {
			ret.port = "443"
		}
	}

	return ret
}

func envOrDefault(name, defaultValue string) string {
	if value, ok := os.LookupEnv(name); ok {
		return value
	}
	return defaultValue
}
