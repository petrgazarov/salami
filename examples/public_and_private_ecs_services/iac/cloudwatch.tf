resource "aws_cloudwatch_log_group" "ServerLogGroup" {
  name = "server-log-group"
}

resource "aws_cloudwatch_log_group" "PythonExecLogGroup" {
  name = "python-exec-log-group"
}