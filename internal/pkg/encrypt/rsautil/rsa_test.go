package rsautil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testPublicKey = []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDT5+vNRR30ZWyAVPIaB6f9jI9p
wHY2ZY3QkHQ1dghaJiuxPXWZN+1lOIsfnTvsaQCdCaceo6d6sHkQ4jtl9r2GLCYT
ygZac3/MansQjxg6UhY5DRAzdLZqHfJcaN9Ih4PWqj9tqjqiwtxFFFEZUb1nH2iv
IRBNf5SFw8yVmkvvXQIDAQAB
-----END PUBLIC KEY-----`)

var testPkcs1Private = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDT5+vNRR30ZWyAVPIaB6f9jI9pwHY2ZY3QkHQ1dghaJiuxPXWZ
N+1lOIsfnTvsaQCdCaceo6d6sHkQ4jtl9r2GLCYTygZac3/MansQjxg6UhY5DRAz
dLZqHfJcaN9Ih4PWqj9tqjqiwtxFFFEZUb1nH2ivIRBNf5SFw8yVmkvvXQIDAQAB
AoGAN4vUoMMcXgL0FROvPqmBHJJqyWK82fd23BPxkk31VIQq8dPVbqtdXCodNdVG
bur7US7Fkt99OEjoA0f6H/k0pmufZ+6yMs8gGSDwxpeRLWvZZM1Wh5hIO6kRNGJs
sfcqMCx01y6TOxrFraaTNj5plhzS7na4m8dYSFvmORCBoxUCQQD6kS8k9aqI8xj9
Ck2cRMuWvJwItsqblq6EGG/ti3xw7uf6lAb2oZHfHq3lFMS1s1E2aRc2nyng9Y9q
cTfqfJLPAkEA2IAkMy7L4dz5Ra7/tqP4WcKaP82OW+uuaIcm26rCTLuoS3vW85Sh
RjNmGAVBd1aooNe5/Js29obEGNpND2zWEwJBAOFYSQn4VtKrrsGDzqDHzkFWhw3f
NwAO2Ay83YzJcbUvZzoYftq4HDSJpuLrdq3jAxroEJRzOHq03bJg+GTOfEkCQDtP
j5s9/LjZsqh2crN0ZDsi5uMHyzI/dL5KGEkhlK0007wqJw7/7tauig+WkQLCiNvX
fapIU1xiOyKb23SYWmUCQQDdXfk5bpDpkcR1UxPFnKqRQWJz+cFjoaZjUD699/HT
ZO3K5eNLWqcUvu8xuibDOG9LgPLFyhS8mVj3lX2nY8GK
-----END RSA PRIVATE KEY-----`)

var testPkcs8Private = []byte(`-----BEGIN PRIVATE KEY-----
MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBANPn681FHfRlbIBU
8hoHp/2Mj2nAdjZljdCQdDV2CFomK7E9dZk37WU4ix+dO+xpAJ0Jpx6jp3qweRDi
O2X2vYYsJhPKBlpzf8xqexCPGDpSFjkNEDN0tmod8lxo30iHg9aqP22qOqLC3EUU
URlRvWcfaK8hEE1/lIXDzJWaS+9dAgMBAAECgYA3i9SgwxxeAvQVE68+qYEckmrJ
YrzZ93bcE/GSTfVUhCrx09Vuq11cKh011UZu6vtRLsWS3304SOgDR/of+TSma59n
7rIyzyAZIPDGl5Eta9lkzVaHmEg7qRE0Ymyx9yowLHTXLpM7GsWtppM2PmmWHNLu
dribx1hIW+Y5EIGjFQJBAPqRLyT1qojzGP0KTZxEy5a8nAi2ypuWroQYb+2LfHDu
5/qUBvahkd8ereUUxLWzUTZpFzafKeD1j2pxN+p8ks8CQQDYgCQzLsvh3PlFrv+2
o/hZwpo/zY5b665ohybbqsJMu6hLe9bzlKFGM2YYBUF3Vqig17n8mzb2hsQY2k0P
bNYTAkEA4VhJCfhW0quuwYPOoMfOQVaHDd83AA7YDLzdjMlxtS9nOhh+2rgcNImm
4ut2reMDGugQlHM4erTdsmD4ZM58SQJAO0+Pmz38uNmyqHZys3RkOyLm4wfLMj90
vkoYSSGUrTTTvConDv/u1q6KD5aRAsKI29d9qkhTXGI7IpvbdJhaZQJBAN1d+Tlu
kOmRxHVTE8WcqpFBYnP5wWOhpmNQPr338dNk7crl40tapxS+7zG6JsM4b0uA8sXK
FLyZWPeVfadjwYo=
-----END PRIVATE KEY-----`)

func TestNewRsa(t *testing.T) {
	t.Run("case1", func(t *testing.T) {
		s := NewRsa(testPublicKey, testPkcs1Private)

		value, err := s.Encrypt([]byte("hello world"))
		assert.NoError(t, err)

		data, err := s.Decrypt(value)
		assert.NoError(t, err)
		assert.Equal(t, "hello world", string(data))
	})

	t.Run("case2", func(t *testing.T) {
		s := NewRsa(testPublicKey, testPkcs8Private)

		value, err := s.Encrypt([]byte("hello world"))
		assert.NoError(t, err)

		data, err := s.Decrypt(value)
		assert.NoError(t, err)
		assert.Equal(t, "hello world", string(data))
	})
}
