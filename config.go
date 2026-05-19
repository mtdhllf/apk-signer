package main

import "time"

var (
	// AllowedPackages 指定要签名的包名,为空时无限制
	AllowedPackages []string

	KeystoreAlias = "gopher"
	StorePassword = "gopher"
	KeyPassword   = "gopher"
	// ValidUntil 有效时间
	// 无限制：ValidUntil = time.Time{} // time.Time 的零值
	// 有限制（例如：2026年5月19日14:30:00 UTC）：
	// ValidUntil = time.Date(2026, time.May, 19, 14, 30, 0, 0, time.UTC)
	ValidUntil time.Time = time.Time{}

	// Platform 平台，签名名体现
	Platform = "3588"
	// KeystoreFile 解压后的文件名称
	KeystoreFile = "3588.jks"
)
