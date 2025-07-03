package model

type StatusCode int

const (
	// Success Codes (1xxx)
	Success               StatusCode = 1000
	SuccessPasteCreated   StatusCode = 1001
	SuccessPasteRetrieved StatusCode = 1002
	SuccessPasteUpdated   StatusCode = 1003

	// Client Error Codes (4xxx)
	ErrorBadRequest       StatusCode = 4000
	ErrorValidationFailed StatusCode = 4001
	ErrorNoUpdatableField StatusCode = 4002
	ErrorPasteNotFound    StatusCode = 4404
	ErrorRequestFailed    StatusCode = 4999

	// Server Error Codes (5xxx)
	ErrorInternalServer    StatusCode = 5000
	ErrorDatabase          StatusCode = 5001
	ErrorPasteCreateFailed StatusCode = 5500
)

var messageMap = map[StatusCode]string{
	// Success
	Success:               "요청이 성공적으로 처리되었습니다.",
	SuccessPasteCreated:   "Paste가 성공적으로 생성되었습니다.",
	SuccessPasteRetrieved: "Paste 조회 성공.",
	SuccessPasteUpdated:   "Paste가 성공적으로 수정되었습니다.",

	// Client Errors
	ErrorBadRequest:       "잘못된 요청입니다.",
	ErrorValidationFailed: "입력값 유효성 검사에 실패했습니다.",
	ErrorNoUpdatableField: "업데이트할 필드가 없습니다.",
	ErrorPasteNotFound:    "해당 Paste를 찾을 수 없습니다.",
	ErrorRequestFailed:    "요청에 실패했습니다.",

	// Server Errors
	ErrorInternalServer:    "서버 내부 오류가 발생했습니다.",
	ErrorDatabase:          "데이터베이스 처리 중 오류가 발생했습니다.",
	ErrorPasteCreateFailed: "Paste 생성 중 오류가 발생했습니다.",
}

// 메시지 가져오기
func GetMessage(code StatusCode) string {
	if msg, ok := messageMap[code]; ok {
		return msg
	}
	return "정의되지 않은 상태 코드입니다."
}
