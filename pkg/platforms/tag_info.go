package platforms

type TagInfo struct {
	TagName string // 标签名称，例如 "v1.0.0"
	SHA     string // 标签对应的提交对象的 SHA 值
}
