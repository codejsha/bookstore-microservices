package support

import "testing"

func TestNewOpensearchTransport_insecure_boundsDialHandshakeAndResponseHeaders(t *testing.T) {
	transport := newOpensearchTransport(true)

	if transport.TLSHandshakeTimeout != opensearchTLSHandshakeTimeout ||
		transport.ResponseHeaderTimeout != opensearchResponseHeaderTimeout ||
		transport.MaxIdleConnsPerHost != opensearchMaxIdleConnsPerHost {
		t.Fatalf("transport limits = handshake %v, response header %v, idle per host %d",
			transport.TLSHandshakeTimeout, transport.ResponseHeaderTimeout, transport.MaxIdleConnsPerHost)
	}
	if transport.TLSClientConfig == nil || !transport.TLSClientConfig.InsecureSkipVerify {
		t.Fatal("insecure transport must skip certificate verification")
	}
}

func TestNewOpensearchTransport_secure_keepsCertificateVerification(t *testing.T) {
	transport := newOpensearchTransport(false)

	if transport.TLSClientConfig != nil {
		t.Fatal("secure transport must use the default TLS verification")
	}
	if transport.ResponseHeaderTimeout != opensearchResponseHeaderTimeout {
		t.Fatalf("response header timeout = %v, want %v", transport.ResponseHeaderTimeout, opensearchResponseHeaderTimeout)
	}
}
