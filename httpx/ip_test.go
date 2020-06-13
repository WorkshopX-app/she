package httpx

import "testing"

func TestIsPrivate(t *testing.T) {
	cases := map[string]bool{
		"127.0.0.1":     true,
		"fd00::":        true,
		"10.0.100.10":   true,
		"172.31.46.100": true,
	}

	for addr, expect := range cases {
		yes, err := IsPrivateIP(addr)
		if err != nil {
			t.FailNow()
		}
		if yes != expect {
			t.Fail()
		}
	}
}
