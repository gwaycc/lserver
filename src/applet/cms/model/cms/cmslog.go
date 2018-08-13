package cms

// 仅记录重要的操作日志
// 其他的由行为流水记录
func PutLog(username, kind, args, memo string) error {
	return cmsdb.PutLog(username, kind, args, memo)
}
