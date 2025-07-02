// pkg/logger/logger.go
package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

// Init 은 zap 로거를 초기화합니다.
// isProd 가 true 이면 JSON 포맷, false 이면 컬러 콘솔 포맷을 사용합니다.
func Init(isProd bool) {
	var core zapcore.Core

	// 로그 레벨 설정: Info 이상만 출력
	level := zapcore.InfoLevel

	// 로그 출력 포맷 설정
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "time",                                             // 로그 시간 필드 이름
		LevelKey:       "level",                                            // 로그 레벨 필드 이름 (INFO, ERROR 등)
		NameKey:        "logger",                                           // 로거 이름 필드 (사용 안 하면 생략 가능)
		CallerKey:      "caller",                                           // 로그 호출 위치 (파일명:라인번호)
		MessageKey:     "msg",                                              // 로그 메시지 필드 이름
		StacktraceKey:  "stacktrace",                                       // 에러 로그 시 스택트레이스 필드 이름
		LineEnding:     zapcore.DefaultLineEnding,                          // 줄바꿈 문자 설정 (보통 \n)
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,                   // 컬러 출력 (개발용)
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"), // 시간 포맷
		EncodeDuration: zapcore.SecondsDurationEncoder,                     // 수행 시간(seconds)
		EncodeCaller:   zapcore.ShortCallerEncoder,                         // 호출 위치 (short 형식)
	}

	if isProd {
		// 프로덕션 모드: JSON 로그
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg), // 1. JSON 포맷으로 인코딩 (key-value 구조)
			zapcore.Lock(os.Stdout),            // 2. 로그 출력 위치 (표준 출력, 락 처리됨)
			level,                              // 3. 출력 레벨 (InfoLevel 이상만 출력)
		)
	} else {
		// 개발 모드: 컬러 콘솔 로그
		core = zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderCfg), // 1. 콘솔용 인코더 (컬러 + 정해진 포맷)
			zapcore.Lock(os.Stdout),               // 2. 출력 위치: 표준 출력 (thread-safe)
			level,                                 // 3. 출력 로그 레벨 (InfoLevel 이상)
		)
	}

	// 로거 인스턴스 생성 (호출 위치 포함 + 에러에 대해 스택트레이스 포함)
	Log = zap.New(
		core,                                  // zapcore.Core (콘솔 or JSON 포맷 + 출력 대상 + 레벨 포함)
		zap.AddCaller(),                       // 호출 위치(file:line) 로그에 추가 (ex. app.go:123)
		zap.AddStacktrace(zapcore.ErrorLevel), // ERROR 이상일 때만 스택트레이스 포함
		// zap.Fields(),                          // 모든 로그에 공통 필드 추가
	)
}
