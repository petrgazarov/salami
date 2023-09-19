package target_file_manager

type TargetFileError struct {
	Message string
}

func (e *TargetFileError) Error() string {
	return "target file error: " + e.Message
}
