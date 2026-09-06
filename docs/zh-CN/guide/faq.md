# 常见问题

## 通用问题

### 什么是 AirlineSim Autobuy？

AirlineSim Autobuy 是一个 Go 语言编写的自动化工具有可以持续监控 AirlineSim 游戏的二手飞机市场，根据用户自定义规则自动购买符合条件的飞机。它提供 Web 管理界面用于配置和监控。

### 这个工具是否违反游戏规则？

本工具仅用于个人学习和辅助游戏。使用者需自行承担使用自动化工具的风险。建议合理设置轮询间隔，避免对游戏服务器造成过大压力。

### 支持哪些平台？

程序使用 Go 语言编写，支持 Windows、Linux 和 macOS 三大平台。前端界面使用 Vue 3 构建，内嵌在 Go 二进制文件中，无需额外依赖。

### 是否需要付费？

AirlineSim Autobuy 是开源项目，基于 MIT 许可证发布，完全免费使用。

## 配置问题

### 如何配置多个服务器？

在 `configs/config.yaml` 中使用 `servers` 数组配置多个服务器：

```yaml
servers:
  - host: free1
    base_url: https://free1.airlinesim.aero
  - host: free2
    base_url: https://free2.airlinesim.aero
```

然后通过规则的 `server_id` 字段指定规则应用在哪个服务器上。

### 如何配置多个账户？

使用 `auths` 数组配置多个认证账户：

```yaml
auths:
  - username: "account1@example.com"
    password: "pass1"
    session_file: session1.json
  - username: "account2@example.com"
    password: "pass2"
    session_file: session2.json
```

每条规则可以通过 `auth_id` 指定使用哪个账户。

### 配置文件的密码是明文存储的吗？

是的，当前配置文件的密码以明文形式存储在 YAML 文件中。建议：

- 限制配置文件的访问权限（`chmod 600 config.yaml`）
- 不要在版本控制中提交包含真实密码的配置文件
- 后续版本可能会支持环境变量或加密存储

### 修改配置文件后需要重启吗？

不需要。程序内置了配置文件监控器，每 5 秒检查文件修改时间。检测到变化后会自动重新加载配置和规则。但修改服务器或认证信息后建议重启程序以确保完全生效。

## 技术问题

### 程序启动后没有输出？

检查以下可能的原因：

1. 配置文件路径是否正确（默认 `configs/config.yaml`）
2. 配置文件格式是否正确（YAML 缩进错误）
3. 网络连接是否正常（能否访问 `airlinesim.aero`）
4. 账户信息是否正确

### 登录失败怎么办？

登录失败的可能原因：

1. 用户名或密码错误——检查配置文件中的凭据
2. 网络问题——确保能够访问 `sar.simulogics.games`
3. 账户被锁定——尝试在浏览器中手动登录
4. API 变更——如果游戏更新了认证 API，程序可能需要更新

### Web 界面无法访问？

1. 确认程序正在运行
2. 检查 WebUI 配置是否正确：
   ```yaml
   webui:
     enabled: true
     host: 0.0.0.0
     port: 9090
   ```
3. 检查防火墙是否开放了对应端口
4. 如果使用云服务器，检查安全组规则
5. 尝试访问 `http://127.0.0.1:9090` 排除网络问题

### 如何查看运行日志？

- **控制台**：默认在控制台输出结构化日志
- **systemd**：`journalctl -u airlinesim-autobuy -f`
- **Docker**：`docker logs -f airlinesim-autobuy`
- **Web 界面**：日志页面提供实时事件查看

### 如何修改日志级别？

程序启动时默认使用 `Info` 级别。要查看更详细的日志，可以在 `main.go` 中修改日志级别：

```go
slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,  // 改为 Debug 级别
})))
```

## 故障排除

### 规则匹配但未购买

可能的原因：

1. 规则的 `auto_buy` 设置为 `false`——将 `auto_buy` 改为 `true`
2. 评分计算为 0——检查匹配条件是否过于严格
3. 余额不足——检查 `min_balance` 设置
4. 购买操作失败——查看日志中的具体错误信息

### 遇到 "URL validation failed" 错误

这是安全机制在阻止请求。程序只允许连接到以下域名的服务器：

- `airlinesim.aero` 及其子域名
- `sar.simulogics.games`
- `simulogics.games`

如果出现此错误，请检查 `base_url` 配置是否正确。

### 购买失败但没有错误信息

购买执行器会尝试：

1. 访问飞机详情页
2. 提交购买表单或出价
3. 检查响应状态码

如果购买失败，检查：

- 会话是否有效（可能已过期）
- 飞机是否已被其他人购买
- 账户余额是否充足
- 竞拍出价是否低于当前最高价

### 程序崩溃或 panic

如果遇到程序崩溃，请：

1. 查看完整的错误堆栈信息
2. 检查配置文件是否有语法错误
3. 确认网络连接稳定
4. 在 GitHub 提交 Issue 时附上完整的错误日志和配置（注意隐藏密码）

### 如何重置会话？

删除会话文件（默认为 `session.json`），程序会在下次运行时自动重新登录。

### 如何恢复默认配置？

删除 `configs/config.yaml` 文件，程序会在下次启动时自动生成包含默认配置的新文件。

### 如何报告问题或提出功能建议？

请在 GitHub 仓库提交 Issue：

- 问题报告：https://github.com/Maicarons/airlinesim-autobuy/issues
- 请附上：程序版本、操作系统、错误日志、配置文件（隐藏密码）

## 其他问题

### 程序会占用多少系统资源？

程序非常轻量：

- **CPU**：仅在轮询时短暂使用，大部分时间空闲
- **内存**：通常占用 10-30 MB
- **磁盘**：二进制文件约 12 MB，配置文件约 1 KB
- **网络**：每次轮询发送 1-2 个 HTTP 请求

### 可以同时监控多少个服务器？

理论上没有限制，但建议不要超过 5-10 个服务器，以免对游戏服务器造成过大压力。每个服务器会在独立的 goroutine 中运行。

### 程序会一直运行吗？

程序会持续运行直到收到中断信号（SIGINT/SIGTERM）或手动停止。可以通过 Web 界面的"启动/停止"按钮控制引擎的运行状态，但即使引擎停止，Web 界面本身仍然可用。