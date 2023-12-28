package encoding

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
)

func TestEncode(t *testing.T) {
	t.Skip()

	var pkg = bytes.NewBuffer(nil)
	for i := 0; i < 100; i++ {
		data, err := NewEncode([]byte(fmt.Sprintf(`{"date":"21日星期三%d","sunrise":"06:19","high":"高温 11.0℃","low":"低温 1.0℃","sunset":"18:26","aqi":85,"fx":"南风","fl":"<3级","type":"多云","notice":"阴晴之间，谨防紫外线侵扰"}`, i)))
		if err != nil {
			panic(err)
		}
		pkg.Write(data)
	}

	bu := bufio.NewReader(pkg)

	for i := 0; i < 100; i++ {
		data, err := NewDecode(bu)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(string(data))
	}
}
