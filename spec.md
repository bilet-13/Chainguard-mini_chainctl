# mychainctl - Chainguard CLI PoC Specification

## 1. 專案概述與功能範圍 (Capabilities)

本專案為概念驗證 (PoC) 型的指令列工具，旨在模擬 Chainguard `chainctl` 的核心體驗，展示與 OCI Registry API 的無縫整合。開發目標為在 24 小時內快速交付具備核心展示功能的可用版本。

**支援指令清單：**

- `mychainctl version`：印出當前 CLI 版本與編譯資訊。
- `mychainctl images list <repository>`：查詢特定儲存庫（如 `chainguard/python`）下所有可用的 Image Tags。支援終端機表格輸出或 JSON 格式。
- `mychainctl images inspect <image:tag>`：獲取特定映像檔的詳細架構資訊（如 Digest、OS/Arch、Media Type）。

## 2. 工具與技術堆疊 (Tech Stack)

| 工具 / 套件                       | 用途說明                                                         |
| :-------------------------------- | :--------------------------------------------------------------- |
| **Golang 1.21+**                  | 核心開發語言。                                                   |
| **`spf13/cobra`**                 | 構建現代化 CLI 的標準框架，負責指令解析、Flag 路由與 Help 說明。 |
| **`google/go-containerregistry`** | (簡稱 ggcr) 處理 OCI 標準 Registry 的 API 請求與底層協議處理。   |
| **`text/tabwriter`**              | Go 內建套件，用於終端機的自動對齊表格輸出。                      |
| **`encoding/json`**               | 處理 REST API 回傳資料與 `--output=json` 的序列化需求。          |

## 3. 專案架構與檔案職責 (File Structure)

```text
mychainctl/
├── main.go               # 程式進入點，初始化 cobra
├── Makefile              # 定義 `make build` 與 `make run` 加速開發與展示
├── cmd/                  # 放置 Cobra 的各個指令邏輯
│   ├── root.go           # 根指令設定 (定義 --output 等全域 Flag)
│   ├── version.go        # version 指令
│   └── images.go         # images 相關指令 (包含 list, inspect)
└── pkg/                  # 核心業務邏輯與 API Client
    └── registry/
        └── client.go     # 封裝使用 ggcr 打 API 的邏輯
```

- **`main.go`**: 僅負責實例化 Root Command 並執行 `Execute()`。
- **`cmd/root.go`**: 關閉預設的錯誤提示 (`cmd.SilenceUsage = true`, `cmd.SilenceErrors = true`)，註冊全域 Flags。
- **`cmd/images.go`**: 處理命令列引數，呼叫 `pkg/registry/client.go`，並根據 `--output` Flag 決定使用 `tabwriter` 還是 `json` 輸出。
- **`pkg/registry/client.go`**: 純粹的業務邏輯層。不處理終端機輸出，僅負責與 API 溝通並回傳資料結構或錯誤。

## 4. 狀態與方法定義 (States & Methods)

```go
// 狀態定義：封裝與 Registry 互動所需的基礎設定
type Client struct {
    RegistryURL string
}

// 初始化方法
func NewClient(registryURL string) *Client

// 核心行為定義
func (c *Client) ListTags(ctx context.Context, repository string) ([]string, error)
func (c *Client) InspectImage(ctx context.Context, imageRef string) (*ImageMetadata, error)

// 資料載體 (DTO)
type ImageMetadata struct {
    Digest    string `json:"digest"`
    MediaType string `json:"mediaType"`
    Platform  string `json:"platform"`
}
```

## 5. 錯誤處理與防禦機制 (Error Handling & Timeouts)

為符合快速開發原則並確保基礎穩定性，採取以下極簡防禦策略：

1. **簡單暴力的 Timeout 控制**：在 `cmd` 層呼叫 `Client` 方法前，強制使用 `context.WithTimeout(context.Background(), 10*time.Second)`。若 API 請求超過 10 秒未回應，直接取消請求避免死鎖。
2. **Error 往上拋 (Bubble Up)**：`pkg/registry` 遇到任何錯誤（網路異常、解析失敗）不進行重試，直接 return 給 `cmd` 層。
3. **終端機錯誤輸出**：`cmd` 層透過 Cobra 的 `RunE` 接收錯誤，由 main 函數統一攔截並輸出乾淨的錯誤訊息。

## 6. 目標 API (Target OCI Registry APIs)

本專案直接對接 Chainguard 的公共映像檔庫 (`cgr.dev`)，透過標準 OCI Distribution API 展示資料獲取能力。底層網路請求與 JSON 解析將由 `ggcr` 套件處理。

- **List Tags API (獲取標籤列表)**
  - **Endpoint:** `GET https://cgr.dev/v2/<repository>/tags/list`
  - **用途:** 供 `images list` 指令使用，回傳特定映像檔（如 `chainguard/python`）的所有可用版本標籤。
- **Manifest API (獲取映像檔清單與架構資訊)**
  - **Endpoint:** `GET https://cgr.dev/v2/<repository>/manifests/<tag_or_digest>`
  - **Header 需求:** 需帶入 `Accept` header 指定 OCI Manifest 或 Docker v2 的 media type (由 `ggcr` 自動處理)。
  - **用途:** 供 `images inspect` 指令使用，解析 Image 的架構、OS、Media Type 與 Digest 雜湊值。
