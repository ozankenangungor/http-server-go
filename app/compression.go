package main

import "strings"

// supportedEncodings lists the compression schemes this server can produce,
// in order of preference.
var supportedEncodings = []string{"gzip"}

// negotiateEncoding returns the first scheme in supportedEncodings that also
// appears in the client's (possibly comma-separated) Accept-Encoding
// header, or an empty string if none match.
func negotiateEncoding(acceptEncoding string) string {
	requested := strings.Split(acceptEncoding, ",")
	for _, supported := range supportedEncodings {
		for _, scheme := range requested {
			if strings.TrimSpace(scheme) == supported {
				return supported
			}
		}
	}
	return ""
}
