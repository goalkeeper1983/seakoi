package tools

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"

	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 总配置
type tomlConfig struct {
	Title     string
	LogConfig logConfig `toml:"Log"`
}

type logConfig struct {
	Debug       bool   `toml:"debug"`       //是否开启调试
	MaxSize     int    `toml:"maxSize"`     //日志文件最大多少兆
	MaxDays     int    `toml:"maxDays"`     //日志文件保留天数
	MaxBackups  int    `toml:"maxBackups"`  //保留文件数
	FileName    string `toml:"fileName"`    //日志名字
	Compress    bool   `toml:"compress"`    //日志生成压缩包,大幅降低磁盘空间,必要时使用
	RotateByDay bool   `toml:"rotateByDay"` //每天轮转一次,如果开启,maxBackups的值需要>=maxDays
}

var (
	Log        *zap.Logger
	TomlConfig *tomlConfig
)

func defaultConfig() *logConfig {
	return &logConfig{
		Debug:       true,              //是否开启调试
		MaxSize:     1024,              //日志文件最大多少兆
		MaxDays:     29,                //日志文件保留天数
		MaxBackups:  30,                //保留文件数
		FileName:    "log/default.log", //日志名字
		Compress:    true,              //日志生成压缩包,大幅降低磁盘空间,必要时使用
		RotateByDay: false,             //每天轮转一次,如果开启,maxBackups的值需要>=maxDays
	}
}

func init() {
	TomlConfig = &tomlConfig{
		LogConfig: *defaultConfig(),
	}

	//加入日期轮转
	l := &lumberjack.Logger{
		Filename:   TomlConfig.LogConfig.FileName,
		MaxSize:    TomlConfig.LogConfig.MaxSize,    // megabytes 兆字节
		MaxBackups: TomlConfig.LogConfig.MaxBackups, //保留文件数
		MaxAge:     TomlConfig.LogConfig.MaxDays,    // days
		LocalTime:  true,
		Compress:   TomlConfig.LogConfig.Compress, //日志生成压缩包,大幅降低磁盘空间,必要时使用
	}

	//开启24小时轮转一次
	if TomlConfig.LogConfig.RotateByDay {
		err := l.Rotate()
		if err != nil {
			log.Panicln(err.Error())
		}
		go func() {
			rotateStartTime := time.Now().Format("2006-01-02") + " 23:59:59"
			rotateLeftTime, _ := time.ParseInLocation("2006-01-02 15:04:05", rotateStartTime, time.Local)
			remainSecond := time.Duration(rotateLeftTime.Unix() - time.Now().Unix() + 1)
			for {
				<-time.After(time.Second * remainSecond)
				err = l.Rotate()
				if err != nil {
					log.Panicln(err.Error())
				}
				remainSecond = 24 * 60 * 60
			}
		}()
	}

	//写入磁盘
	w := zapcore.AddSync(l)

	//初始化encoder配置
	encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoderCfg.EncodeCaller = shortCallerEncoder //日志调用方法
	encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder

	var core zapcore.Core
	//调试模式
	if TomlConfig.LogConfig.Debug {
		consoleErrors := zapcore.Lock(os.Stderr)
		core = zapcore.NewTee(
			zapcore.NewCore(zapcore.NewConsoleEncoder(encoderCfg), consoleErrors, zap.DebugLevel),
			zapcore.NewCore(zapcore.NewConsoleEncoder(encoderCfg), w, zap.DebugLevel),
		)
	} else {
		core = zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderCfg), // zapcore.NewJSONEncoder(encoderCfg) //json格式
			w,
			zap.DebugLevel,
		)
	}
	Log = zap.New(core, zap.AddCaller())
	Log.Info(RunFuncName(), zap.Any("cpuNum", runtime.NumCPU()))
}

func shortCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	// TODO: consider using a byte-oriented API to save an allocation.
	enc.AppendString("[" + caller.TrimmedPath() + "]")
}

type GormLogger struct {
}

func (This *GormLogger) Print(v ...interface{}) {
	switch v[0] {
	case "sql":
		sql := v[3].(string)
		switch reflect.TypeOf(v[4]).Kind() {
		case reflect.Slice:
			args := reflect.ValueOf(v[4])
			for i := 0; i != args.Len(); i++ {
				index := strings.Index(sql, "?")
				sql = sql[0:index] + fmt.Sprintf("'%v'", args.Index(i)) + sql[index+1:]
			}
		}

		Log.Info("sql",
			//zap.String("module", "gorm"),
			//zap.Any("type", "sql"),
			zap.Any("sql", sql),
			zap.Any("values", v[4]),
			zap.Any("query", v[3]),
			zap.Any("duration", v[2]),
			zap.Any("rows_returned", v[5]),
			zap.Any("src", v[1]))

	case "log":
		Log.Info("log", zap.Any("gorm", v[2]))
	}
}

type GormLoggerV2 struct {
}

// LogMode log mode
func (l *GormLoggerV2) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	return &newlogger
}

// Info print info
func (l *GormLoggerV2) Info(ctx context.Context, msg string, data ...interface{}) {
	Log.Info(utils.FileWithLineNum(), zap.Any("msg", msg), zap.Any("data", data))
}

// Warn print warn messages
func (l *GormLoggerV2) Warn(ctx context.Context, msg string, data ...interface{}) {
	Log.Warn(utils.FileWithLineNum(), zap.Any("msg", msg), zap.Any("data", data))

}

// Error print error messages
func (l *GormLoggerV2) Error(ctx context.Context, msg string, data ...interface{}) {
	Log.Error(utils.FileWithLineNum(), zap.Any("msg", msg), zap.Any("data", data))
}

// Trace print sql message
func (l *GormLoggerV2) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	switch {
	case err != nil:
		Log.Error(utils.FileWithLineNum(), zap.Any("err", err.Error()), zap.Any("elapsed", fmt.Sprintf("%vms", float64(elapsed.Nanoseconds())/1e6)), zap.Any("rows", rows), zap.Any("sql", sql))
	default:
		Log.Info(utils.FileWithLineNum(), zap.Any("elapsed", fmt.Sprintf("%vms", float64(elapsed.Nanoseconds())/1e6)), zap.Any("rows", rows), zap.Any("sql", sql))
	}
}
