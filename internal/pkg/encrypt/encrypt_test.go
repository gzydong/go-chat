package encrypt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase64Decode(t *testing.T) {

}

func TestBase64Encode(t *testing.T) {

}

func TestHashPassword(t *testing.T) {
	pwd := HashPassword("admin123sdgsdsgds")

	assert.Equal(t, true, len(pwd) == 60)
}

func TestMd5(t *testing.T) {
	assert.Equal(t, "c069d1fbc7bf8e994de8299110e68bc5", Md5("s6hqzp6j0kdfzh4n_cjq6b180000gn"))
}

func TestVerifyPassword(t *testing.T) {
	pwd := HashPassword("admin123")

	assert.Equal(t, true, VerifyPassword(pwd, "admin123"))
	assert.Equal(t, false, VerifyPassword(pwd, "admin1234453"))
}

func TestName(t *testing.T) {
	key := `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQCncDS9NJcpBMX0WYbyDeLETx29vhjgomegiGlZxrSzHTYlBra9
HQ6/bVOHRQNba1nlH+q6gM4Y2VwNdraDdqxhhDyiEs9bhUujUtMYmh2hcn5G8Fkx
JLiA8AE185iVOqojubmiatC5KKjWRFT+4z/iq5H26+VNmx3S9wnNW/5ydwIDAQAB
AoGAM49M3jqUla//mSf8cwstmk/Wk7g3Bu1bxcZb0qZqvIExTCOOIBwTj4UF5LCu
wPcEvpaefIHvdR1xyD+XIlJn8DvypO8hlJKE52bye4BkI3m8zeTTBf4QEXhAlMGY
KkVPbR+gABqz7PEWBDCzSlI0Xi8G0UweX1pi99huGepAyoECQQDTx0uSMPdm0pz9
wXAOnXXTyJ4t03D4kGe7pD3urPJcxxiQb7MWcmUuk7ycrJRc5C09EzIKVW15Jn+5
cKKC4jTBAkEAymawIE5gaefUuO2m731unLBph2KW8cqQUILPBLfu25T/rjcmywND
98/zHvGCCXN2kkr1Zx2wp7B88crQF/tdNwJBAM3GRSC8YXfQR2itTzN0PivVMBU4
8PkkXxbNBLxn4WrSrYSSdEHoT3ZNaKQXcGU99NL2VtYBocho5wwJbG6eW0ECQCM4
MzWr7cMAAFgdoorR/MlvOS3BzhpM8UfRO0zK5Nl41/Tsy+dPrigVG20rAUG7wco7
GPDUjcTgRR2d+Q/zQYkCQFi8hvWIJuWc3Q8K3uDqVaAt9SdgA99q1yFp0vM4Zcz3
0RnVO95cQPQP56ENk/oFXrstnNeRzEEzjUrua/McTZQ=
-----END RSA PRIVATE KEY-----`

	decrypt, _ := RSADecrypt("kNwo05Xqa+PMaY938jg9VMNnJkWwQaXmIF6q1IWChVlygaO9Me0IDIFEecR7UNTotqdip6gcNf4aLhTzgZNiZghCp93PF7sTqoBX55ZqJ7pOUNkswHUqMWAPxgzeoFuJYgzI+EM42gmBSA0UmoYfQrZmg5ObZjc7ahLu+UcFyVQ=", []byte(key))

	fmt.Println(decrypt)

}
