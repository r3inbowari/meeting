package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const configFileSizeLimit = 10 << 20

/**
 * 跨域资源Options
 */
func OptionsRet(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("content-type", "application/json")
		w.Header().Add("Access-Control-Allow-Headers", "Authorization")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT")
		return
	}
}

/**
 * Load File
 * @param path 文件路径
 * @param dist 存放目标
 */
func LoadConfig(path string, dist interface{}) error {
	configFile, err := os.Open(path)
	if err != nil {
		Fatal("Failed to open config file.", logrus.Fields{"path": path, "err": err})
		return err
	}

	fi, _ := configFile.Stat()
	if size := fi.Size(); size > (configFileSizeLimit) {
		Fatal("Config file size exceeds reasonable limited", logrus.Fields{"path": path, "size": size})
		return errors.New("limited")
	}

	if fi.Size() == 0 {
		Fatal("Config file is empty, skipping", logrus.Fields{"path": path, "size": 0})
		return errors.New("empty")
	}

	buffer := make([]byte, fi.Size())
	_, err = configFile.Read(buffer)
	buffer, err = StripComments(buffer)
	if err != nil {
		Fatal("Failed to strip comments from json", logrus.Fields{"err": err})
		return err
	}

	buffer = []byte(os.ExpandEnv(string(buffer)))

	err = json.Unmarshal(buffer, &dist)
	if err != nil {
		Fatal("Failed unmarshalling json", logrus.Fields{"err": err})
		return err
	}
	return nil
}

// StripComments
func StripComments(data []byte) ([]byte, error) {
	data = bytes.Replace(data, []byte("\r"), []byte(""), 0)
	lines := bytes.Split(data, []byte("\n"))
	filtered := make([][]byte, 0)

	for _, line := range lines {
		match, err := regexp.Match(`^\s*#`, line)
		if err != nil {
			return nil, err
		}
		if !match {
			filtered = append(filtered, line)
		}
	}
	return bytes.Join(filtered, []byte("\n")), nil
}

func GetIP(conn net.Conn) string {
	return strings.Split(conn.RemoteAddr().String(), ":")[0]
}

func PortInUse(port int) bool {
	checkStatement := fmt.Sprintf("lsof -i:%d ", port)
	output, _ := exec.Command("sh", "-c", checkStatement).CombinedOutput()
	if len(output) > 0 {
		return true
	}
	return false
}

/**
 * get map key
 */
func GetKeys(m map[int]int) []int {
	i := 0
	keys := make([]int, len(m))
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

/**
 * md5生成
 */
func CreateMD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

/**
 * json验证器初始化
 */
func InitValidator() {
	govalidator.SetFieldsRequiredByDefault(true)
}

/**
 * 序列化json
 */
func JsonBind(ptr interface{}, rq *http.Request) error {
	if rq.Body != nil {
		defer rq.Body.Close()
		err := json.NewDecoder(rq.Body).Decode(ptr)
		if err != nil && err != io.EOF {
			return err
		}
		return nil
	} else {
		return errors.New("empty request body")
	}
}

/**
 * auto split
 */
func GetAuthToken(r *http.Request) (string, error) {
	ss := strings.Split(r.Header.Get("Authorization"), " ")
	if len(ss) == 2 {
		return ss[1], nil
	}
	return "", errors.New("unauthorized")
}

/**
 * uuid create
 */
func CreateUUID() string {
	u1 := uuid.NewV4()
	return u1.String()
}
