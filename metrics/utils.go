package metrics

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type certificate struct {
	index    int
	notAfter *time.Time
}

func processCertificates(content string) ([]certificate, error) {
	res := []certificate{}
	data := []byte(content)
	for idx := 0; len(data) != 0; idx++ {
		block, rest := pem.Decode(data)
		dataStr := strings.TrimSpace(string(rest))
		data = []byte(dataStr)
		if block == nil || block.Bytes == nil {
			err := fmt.Errorf("invalid pem decode")
			log.WithError(err).Errorf("error while reading certificate '%s': ", content)
			res = append(res, certificate{
				index: idx,
			})
			return res, err
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			log.WithError(err).WithField("certificate", content).Error("error while reading certificate")
			res = append(res, certificate{
				index: idx,
			})
			return res, err
		}
		res = append(res, certificate{
			index:    idx,
			notAfter: &cert.NotAfter,
		})
	}
	return res, nil
}
