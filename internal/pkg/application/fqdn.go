/*
#######
##       ___ _______ _  ___ ___
##      / _ `/ __/  ' \/ -_) _ \
##      \_,_/_/ /_/_/_/\__/_//_/
##
####### (c) 2020 Institut National de l'Audiovisuel ######################################## Archivage Numérique #######
*/

package application

import (
	"net"
	"os"
	"strings"

	"github.com/arnumina/failure"
)

func fqdn() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}

	addrs, err := net.LookupHost(hostname)
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		hosts, err := net.LookupAddr(addr)
		if err != nil || len(hosts) == 0 {
			continue
		}

		return strings.TrimSuffix(hosts[0], "."), nil
	}

	return "", nil
}

func getFQDN() (string, error) {
	fqdn, err := fqdn()
	if fqdn == "" || err != nil {
		return "",
			failure.New(err).
				Msg("impossible to retrieve the FQDN") /////////////////////////////////////////////////////////////////
	}

	return fqdn, nil
}

/*
######################################################################################################## @(°_°)@ #######
*/
