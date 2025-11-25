package until

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// ------------------ 更新信号 ------------------

func UpdateSignal() error {
	script := "entrypoint.sh"

	cmd := exec.Command("ps", "-eo", "pid,args")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return errors.New("更新服务查找失败")
	}

	lines := strings.Split(out.String(), "\n")
	found := false

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		pidStr, args := fields[0], strings.Join(fields[1:], " ")
		if strings.Contains(args, script) {
			pid, err := strconv.Atoi(pidStr)
			if err != nil {
				continue
			}

			proc, err := os.FindProcess(pid)
			if err != nil {
				continue
			}

			if err := proc.Signal(syscall.SIGUSR1); err != nil {
				return errors.New("更新信号发送失败")
			}

			log.Println("更新信号发送成功,更新中请稍候...")
			found = true
		}
	}

	if !found {
		return errors.New("未找到更新监测脚本进程")
	}
	return nil
}

// ------------------ GitHub Release 版本 ------------------

type githubRelease struct {
	TagName     string    `json:"tag_name"`
	Prerelease  bool      `json:"prerelease"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func isNewer(newVer, oldVer string) bool {
	newVer = strings.TrimPrefix(newVer, "v")
	oldVer = strings.TrimPrefix(oldVer, "v")

	newParts := strings.Split(newVer, ".")
	oldParts := strings.Split(oldVer, ".")

	for len(newParts) < 4 {
		newParts = append(newParts, "0")
	}
	for len(oldParts) < 4 {
		oldParts = append(oldParts, "0")
	}

	for i := 0; i < 4; i++ {
		var n, o int
		fmt.Sscanf(newParts[i], "%d", &n)
		fmt.Sscanf(oldParts[i], "%d", &o)
		if n > o {
			return true
		}
		if n < o {
			return false
		}
	}
	return false
}

func fetchLatestStableRelease(owner, repo string) (*githubRelease, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", owner, repo)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("无法获取版本信息，HTTP 状态码：%d", resp.StatusCode)
	}

	var releases []githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, err
	}

	var latest *githubRelease
	for _, r := range releases {
		if r.Prerelease {
			continue
		}
		if latest == nil || r.PublishedAt.After(latest.PublishedAt) {
			latest = &r
		}
	}

	if latest == nil {
		return nil, errors.New("未找到正式版本发布")
	}
	return latest, nil
}

// 检查远端版本是否比本地新
func CheckNewVer(untilVersion string) (bool, string, error) {
	latest, err := fetchLatestStableRelease("wz1st", "go-iptv")
	if err != nil {
		return false, "", err
	}
	return isNewer(latest.TagName, untilVersion), latest.TagName, nil
}

// ------------------ 下载与校验 ------------------

func downloadFile(url, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载 %s 返回状态 %s", url, resp.Status)
	}

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func fileSHA256(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func verifySHA256(file string, sums map[string]string) bool {
	hash, err := fileSHA256(file)
	if err != nil {
		return false
	}
	expected, ok := sums[file]
	if !ok {
		return false
	}
	return strings.EqualFold(hash, expected)
}

// 下载并校验最新 Release 文件，arch 用于选择 iptv 和 license
func DownloadAndVerify(arch string) (bool, string, error) {
	release, err := fetchLatestStableRelease("wz1st", "go-iptv")
	if err != nil {
		return false, "", err
	}

	log.Println("最新版本:", release.TagName)

	iptvFile := fmt.Sprintf("iptv_%s", arch)
	licenseFile := fmt.Sprintf("license_%s", arch)
	files := []string{iptvFile, licenseFile, "updata.sh", "SHA256SUMS.txt"}
	fileURLs := make(map[string]string)

	for _, asset := range release.Assets {
		for _, fname := range files {
			if asset.Name == fname {
				fileURLs[fname] = asset.BrowserDownloadURL
			}
		}
	}

	for _, f := range files {
		if _, ok := fileURLs[f]; !ok {
			return false, "", fmt.Errorf("%s 不存在于 release", f)
		}
	}

	// 下载文件
	for _, f := range files {
		fmt.Println("下载:", f)
		if err := downloadFile(fileURLs[f], f); err != nil {
			return false, "", fmt.Errorf("下载 %s 失败: %v", f, err)
		}
	}

	// 读取 SHA256SUMS.txt
	sumsFile := "SHA256SUMS.txt"
	sums := make(map[string]string)
	f, err := os.Open(sumsFile)
	if err != nil {
		return false, "", err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) != 2 {
			continue
		}
		sums[parts[1]] = parts[0]
	}

	// 校验文件 SHA256
	for _, f := range []string{iptvFile, licenseFile, "updata.sh"} {
		if ok := verifySHA256(f, sums); !ok {
			return false, "", fmt.Errorf("%s 校验失败", f)
		}
	}

	fmt.Println("所有文件下载并校验通过")
	return true, release.TagName, nil
}
