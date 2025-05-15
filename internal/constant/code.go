package constant

type ResultCode int

const (
	CodeSuccess        ResultCode = 20000
	CodeParamError     ResultCode = 30000
	CodeDBError        ResultCode = 40000
	CodeRecordNotFound ResultCode = 40001
	CodeRuntimeError   ResultCode = 50000
	CodeUnknownError   ResultCode = 60000
)

var ResultCodeMap = map[ResultCode]string{
	CodeSuccess:        "Success",
	CodeParamError:     "ParamError",
	CodeDBError:        "DBError",
	CodeRecordNotFound: "RecordNotFound",
	CodeRuntimeError:   "RuntimeError",
	CodeUnknownError:   "UnknownError",
}

func GetResultCodeName(code ResultCode) string {
	return ResultCodeMap[code]
}
