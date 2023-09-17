package utils

import "regexp"

func ExtractVersion(image string) string {
	// 使用正则表达式从镜像名中匹配版本号
	// 例如：registry.tuyuansu.com.cn:12316/nbqy/b_ctp:1.0.0-SNAPSHOT-20230509.4
	pattern := `:([\w.-]+)$`
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(image)
	if len(match) >= 2 {
		return match[1]
	}
	return ""
}
