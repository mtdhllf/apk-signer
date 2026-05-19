# APK Sign Tool (Go)

这是一个使用 Go 语言开发的 APK 签名工具，它将 `aapt.exe`, `apksigner.jar`, `zipalign.exe` 和你的 `.jks` 签名文件嵌入到单个二进制文件中，方便分发和使用。

## 功能特性

*   将所有必要的工具和签名文件打包成一个独立的 Windows 可执行文件。
*   支持拖拽 APK 文件进行签名。
*   可配置允许签名的应用包名列表。
*   可配置签名密钥信息。
*   可设置程序的有效期。
*   签名后自动进行文件对齐，并清理临时文件。

## 配置要点 (config.go)

所有核心配置都集中在 `config.go` 文件中。

### 1. `AllowedPackages` (允许的包名列表)

```go
AllowedPackages = []string{
    "com.example.app",
    "com.test.app",
}
```
*   **作用：** 定义允许进行签名的 APK 包名列表。
*   **配置：** 在切片中添加你允许的包名字符串。
*   **特殊情况：** 如果 `AllowedPackages` 为空切片（`[]string{}` 或不初始化），则表示不限制包名，所有 APK 都可以进行签名。

### 2. `KeystoreAlias`, `StorePassword`, `KeyPassword` (密钥信息)

```go
KeystoreAlias = "gopher"
StorePassword = "gopher"
KeyPassword   = "gopher"
```
*   **作用：** 配置你的 `.jks` 签名文件的别名、密钥库密码和密钥密码。
*   **配置：** 根据你的 `.jks` 文件实际信息进行修改。

### 3. `KeystoreFile` (签名文件名称)

```go
KeystoreFile = "3588.jks"
```
*   **作用：** 定义程序内部在临时目录中使用的 `.jks` 签名文件名称。
*   **重要提示：**
    *   在编译 Go 程序之前，你实际使用的 `.jks` 签名文件（例如 `3588.jks` 或 `other.jks`）**必须重命名为 `key.jks`**，并放置在项目根目录下。这是因为 `go:embed` 只能嵌入固定名称的文件。
    *   `KeystoreFile` 变量仅用于指定程序运行时，从嵌入的 `keystore.jks` 中提取出来后，在临时目录中使用的文件名。

### 4. `ValidUntil` (程序有效期)

```go
// 无限制：
ValidUntil time.Time = time.Time{} // time.Time 的零值，表示无限制

// 有时间限制（例如：2026年5月19日14:30:00 UTC+8）：
ValidUntil time.Time = time.Date(2026, time.May, 19, 14, 30, 0, 0, time.FixedZone("UTC+8", 8*60*60))
```
*   **作用：** 设置程序的过期时间。
*   **配置：**
    *   设置为 `time.Time{}` 表示程序永不过期。
    *   使用 `time.Date()` 函数指定具体的过期日期和时间。最后一个参数 `time.FixedZone("UTC+8", 8*60*60)` 用于指定时区（这里是 UTC+8）。

### 5. `Platform` (输出文件名后缀)

```go
Platform = "3588"
```
*   **作用：** 定义签名后输出 APK 文件名中的平台标识。
*   **配置：** 修改为任何你想要的字符串。例如，如果 `Platform` 是 `"my_platform"`，输出文件将是 `原文件名_my_platform_sign.apk`。

## 使用方法

1.  **配置：** 根据上述说明修改 `config.go` 文件。
2.  **准备签名文件：** 将你的 `.jks` 签名文件重命名为 `keystore.jks`，并放置在项目根目录下。
3.  **编译：** 在项目根目录运行 `go build -o apk-signer.exe .` 命令。
4.  **运行：** 将需要签名的 APK 文件拖拽到生成的 `apk-signer.exe` 可执行文件上。
5.  **输出：** 签名后的 APK 文件将生成在原 APK 文件所在的目录，命名格式为 `原文件名_Platform_sign.apk`。

## 注意事项

*   确保你的系统已安装 Java 运行环境 (JRE)，因为 `apksigner.jar` 依赖 Java。
*   `aapt.exe`, `apksigner.jar`, `zipalign.exe` 和 `key.jks` 会被嵌入到最终的 `apk-signer.exe` 中。