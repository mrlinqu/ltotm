package main

import "flag"

type cmdArgs struct {
	iface      string
	port       string
	tlsCert    string
	tlsKey     string
	storageDir string
}

func parseArgs() cmdArgs {
	ret := cmdArgs{}

	flag.StringVar(&ret.iface, "iface", "0.0.0.0", "listen interface")
	flag.StringVar(&ret.port, "port", "", "listen port")

	flag.StringVar(&ret.tlsCert, "cert", "", "tls cert")
	flag.StringVar(&ret.tlsKey, "key", "", "tls key")

	flag.StringVar(&ret.storageDir, "dir", "./storage", "path to storage dir")

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
