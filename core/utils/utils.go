package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
	"zrWorker/core/slog"
	"zrWorker/global"
)

// 返回一个32位md5加密后的字符串
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func GetTime() string {
	time := time.Now().Unix()
	s := strconv.FormatInt(time, 10)
	return s
}

// 获取系统当前日期
func GetDate() string {
	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	return time.Now().In(cstZone).Format("2006-01-02")
}

func RandInt(max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return r.Intn(max)
}

func ArrayToString(arr []string) string {
	var result string
	for _, i := range arr { //遍历数组中所有元素追加成string
		result += i + ","
	}
	return result
}

func MergeMap(x, y map[string]string) map[string]string {

	n := make(map[string]string)
	for i, v := range x {
		for j, w := range y {
			if i == j {
				n[i] = w

			} else {
				if _, ok := n[i]; !ok {
					n[i] = v
				}
				if _, ok := n[j]; !ok {
					n[j] = w
				}
			}
		}
	}

	return n
}

func GetLogPath(runTaskID string, name string) string {

	path := global.ServerSetting.LogPath
	if !PathExists(path) {
		err := os.MkdirAll(path, 0777)
		if err != nil {
			slog.Println(slog.WARN, "新建日志目录失败")
		}
	}
	//slog.Println(slog.INFO, path)
	path = global.ServerSetting.LogPath + "/" + name + runTaskID + ".json"
	return path
}

func GetScreenPath() string {

	path := global.ServerSetting.LogPath + "/img/"
	if !PathExists(path) {
		err := os.MkdirAll(path, 0777)
		if err != nil {
			slog.Println(slog.WARN, "新建日志目录失败")
		}
	}

	//slog.Println(slog.INFO, path)

	return path
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func PrinfI(k string, v any) {
	fmt.Printf(k+": %+v\n", v)
}

func GetDomainIp(domain string) string {
	addr, err := net.ResolveIPAddr("ip", domain)
	if err != nil {
		slog.Println(slog.DEBUG, err.Error())
		return ""
	}
	ip := addr.IP.String()

	return ip
}

// GetCurrentTimeText 获取当前时间format
func GetCurrentTimeText() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// BindArgsWithGin 绑定请求参数
func BindArgsWithGin(c *gin.Context, req interface{}) error {
	return c.ShouldBindWith(req, binding.Default(c.Request.Method, c.ContentType()))
}

// MakeMD5 MD5加密
func MakeMD5(data string) string {
	h := md5.New()
	h.Write([]byte(data)) // 需要加密的字符串为 123456
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr) // 输出加密结果
}

// Random 生成随机数
func Random(min, max int) int {
	if min == max {
		return max
	}
	max = max + 1
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return min + r.Intn(max-min)
}

// RandomStr 随机字符串
func RandomStr(l int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	seed := "1234567890QWERTYUIOPASDFGHJKLZXCVBNM"
	str := ""
	length := len(seed)
	for i := 0; i < l; i++ {
		point := r.Intn(length)
		str = str + seed[point:point+1]
	}
	return str
}

// BuildPassword 构建用户密码
func BuildPassword(password, salt string) string {
	return MakeMD5(password + salt)
}

// TernaryOperation 三元操作符
func TernaryOperation(exist bool, res, el interface{}) interface{} {
	if exist {
		return res
	}
	return el
}

// GetBeforeDate 获取n天前的时间
func GetDateFromNow(n int) time.Time {
	timer, _ := time.ParseInLocation("2006-01-02", time.Now().Format("2006-01-02"), time.Local)
	if n == 0 {
		return timer
	}
	return timer.AddDate(0, 0, n)
}

// StrArrExist 检测string数组中是否包含某个字符串
func StrArrExist(arr []string, check string) bool {
	for _, v := range arr {
		if v == check {
			return true
		}
	}
	return false
}

// RetryFunc 带重试的func
func RetryFunc(times int, f func() error) error {
	var (
		reTimes int
		err     error
	)
RETRY:
	if err = f(); err != nil {
		if reTimes == times {
			return err
		}
		time.Sleep(time.Duration(1) * time.Second)
		reTimes++
		goto RETRY
	}
	return nil
}

func GetNowTome() string {
	return time.Now().Format("20060102150405")
}

func GetRandomString(l int) string {
	str := "abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func WirteFileAppend(fileName string, newSubDomain []string) {
	fd, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	for _, st := range newSubDomain {
		buf := []byte(st + "\n")
		fd.Write(buf)
	}
	fd.Close()
}

func PathExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func SliceToString(slices []string) (result string) {
	b, err := json.Marshal(slices)
	if err != nil {
		return
	}
	result = string(b)
	return
}

func SliceSToString(slices [][]string) (result string) {
	b, err := json.Marshal(slices)
	if err != nil {
		return
	}
	result = string(b)
	return
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
