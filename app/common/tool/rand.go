package tool

import (
	"math/rand"
)

// 生成随机字符串
func Randnum(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

// func Randnum() string {

// 	// 设置随机数生成器的种子，使用当前时间的纳秒级别时间戳
// 	rand.Seed(time.Now().UnixNano())

// 	// 生成一个六位数的随机整数
// 	min := 100000 // 最小值（六位数的最小值）
// 	max := 999999 // 最大值（六位数的最大值）
// 	randomInt := min + rand.Intn(max-min+1)
// 	//fmt.Printf("随机六位数: %06d\n", randomInt)
// 	return strconv.Itoa(randomInt)

// }
