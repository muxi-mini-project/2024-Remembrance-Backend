package study

// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// )

// const (
// 	appId       = "YOUR_APP_ID"
// 	appSecret   = "YOUR_APP_SECRET"
// 	redirectURI = "YOUR_REDIRECT_URI"
// )

// // 微信授权码请求 URL
// func wechatAuthURL() string {
// 	return fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_userinfo&state=STATE#wechat_redirect", appId, redirectURI)
// }

// // 使用授权码获取 Access Token
// func getAccessToken(code string) (string, error) {
// 	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code", appId, appSecret, code)
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer resp.Body.Close()
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return "", err
// 	}

// 	var result map[string]interface{}
// 	json.Unmarshal(body, &result)

// 	return result["access_token"].(string), nil
// }

// // 主函数
// func main() {
// 	http.HandleFunc("/wx_login", func(w http.ResponseWriter, r *http.Request) {
// 		// 获取授权码
// 		code := r.URL.Query().Get("code")
// 		if code != "" {
// 			// 使用授权码获取 Access Token
// 			token, err := getAccessToken(code)
// 			if err != nil {
// 				fmt.Fprintf(w, "Error: %s", err)
// 				return
// 			}
// 			fmt.Fprintf(w, "Access Token: %s", token)
// 		} else {
// 			// 重定向到微信授权页面
// 			http.Redirect(w, r, wechatAuthURL(), http.StatusFound)
// 		}
// 	})

// 	fmt.Println("Server is running...")
// 	http.ListenAndServe(":8080", nil)
// }
