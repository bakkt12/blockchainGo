package main

import (
	"encoding/base64"
	"fmt"
)
/**
 地址就用 ：URLEncoding
字符用：StdEncoding
 */

func main() {
	//msg := "Hello, 世界"
	msg:="https://www.baidu.com/s?ie=utf-8&f=8&rsv_bp=1&tn=monline_4_dg&wd=go%20base64&oq=base58&rsv_pq=dca057310002e021&rsv_t=c67d4DrNLVnPQFDA0L%2FRm%2FdzLxiK5%2Fcg9i3l40h8SZ9JMGzL4wKbbU93IEzzLTG%2ByCfe&rqlang=cn&rsv_enter=0&rsv_dl=tb&inputT=576&rsv_sug3=79&rsv_sug1=43&rsv_sug7=101&bs=base58";
	encoded := base64.URLEncoding.EncodeToString([]byte(msg))
	//encoded := base64.StdEncoding.EncodeToString([]byte(msg))

	fmt.Println(encoded)
	decoded, err := base64.URLEncoding.DecodeString(encoded)
	//decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		fmt.Println("decode error:", err)
		return
	}
	fmt.Println(msg)
	fmt.Println(string(decoded))
}
