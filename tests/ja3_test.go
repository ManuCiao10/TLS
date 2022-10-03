package tests

import (
	"testing"

	utls "github.com/bogdanfinn/utls"
	"github.com/stretchr/testify/assert"

	tls_client "github.com/bogdanfinn/tls-client"
)

func TestJA3(t *testing.T) {
	t.Log("testing ja3 chrome")
	ja3_chrome_105(t)
	t.Log("testing ja3 firefox")
	ja3_firefox_105(t)
	t.Log("testing ja3 opera")
	ja3_opera_91(t)
}

func ja3_chrome_105(t *testing.T) {
	input := browserFingerprints[chrome][utls.HelloChrome_105.Str()][ja3String]

	specFunc, err := tls_client.GetSpecFactorFromJa3String(input)

	if err != nil {
		t.Fatal(err)
	}

	spec, err := specFunc()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(spec.CipherSuites), 15, "Client should have 15 CipherSuites")
	assert.Equal(t, len(spec.Extensions), 16, "Client should have 16 extensions")
}

func ja3_firefox_105(t *testing.T) {
	input := browserFingerprints[firefox][utls.HelloFirefox_105.Str()][ja3String]

	specFunc, err := tls_client.GetSpecFactorFromJa3String(input)

	if err != nil {
		t.Fatal(err)
	}

	spec, err := specFunc()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(spec.CipherSuites), 17, "Client should have 17 CipherSuites")
	assert.Equal(t, len(spec.Extensions), 15, "Client should have 15 extensions")
}

func ja3_opera_91(t *testing.T) {
	input := browserFingerprints[opera][utls.HelloOpera_91.Str()][ja3String]

	specFunc, err := tls_client.GetSpecFactorFromJa3String(input)

	if err != nil {
		t.Fatal(err)
	}

	spec, err := specFunc()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(spec.CipherSuites), 15, "Client should have 15 CipherSuites")
	assert.Equal(t, len(spec.Extensions), 16, "Client should have 16 extensions")
}
