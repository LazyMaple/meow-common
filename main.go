package meow_common

import (
	"embed"
	"fmt"
	"io"
	"os"
	"os/exec"

	jsoniter "github.com/json-iterator/go"
	"github.com/justincormack/go-memfd"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger *zap.Logger
	json   = jsoniter.ConfigCompatibleWithStandardLibrary

	ossAPI    = "https://s3.bitiful.net"
	ossAK     = ""
	ossSK     = ""
	ossRegion = "cn-east-1"
)

func init() {
	Logger = loadLog()
}

func loadLog() *zap.Logger {
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(os.Stdout),
		zap.DebugLevel,
	)
	Logger = zap.New(core)
	defer Logger.Sync()
	return Logger
}

func SetOSSAPI(s string) {
	ossAPI = s
}

func SetOSSAK(s string) {
	ossAK = s
}

func SetOSSSK(s string) {
	ossSK = s
}

func SetOSSRegion(s string) {
	ossRegion = s
}

func GetOSSClient() (*minio.Client, error) {
	return minio.New(ossAPI, &minio.Options{
		Creds:  credentials.NewStaticV4(ossAK, ossSK, ""),
		Secure: true,
		Region: ossRegion,
	})
}

func WarpLogError(err error) zap.Field {
	return zap.Error(err)
}

func LogError(msg string, err error) {
	Logger.Error(msg, WarpLogError(err))
}

func GetCommand(fs embed.FS, filename string, args ...string) (*exec.Cmd, func() error, error) {
	f, err := fs.Open(filename)
	if err != nil {
		return &exec.Cmd{}, nil, err
	}
	defer f.Close()

	mfd, err := memfd.Create()
	if err != nil {
		return &exec.Cmd{}, nil, err
	}
	defer mfd.Unmap()

	_, err = io.Copy(mfd, f)
	if err != nil {
		mfd.Close()
		return &exec.Cmd{}, nil, err
	}

	cmd := exec.Command(fmt.Sprintf("/proc/self/fd/%d", mfd.Fd()), args...)
	return cmd, mfd.Close, nil
}

func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
