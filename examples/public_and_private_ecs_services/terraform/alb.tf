resource "aws_security_group" "AlbSecurityGroup" {
  name        = "alb-security-group"
  description = "Security group for ALB"
  vpc_id      = aws_vpc.MainVpc.id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_lb" "ServerAlb" {
  name               = "server-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.AlbSecurityGroup.id]
  subnets            = [aws_subnet.PublicSubnetA.id, aws_subnet.PublicSubnetB.id]

  enable_deletion_protection = true

  idle_timeout = 3600

  tags = {
    Name = "server-alb"
  }
}

resource "aws_lb_target_group" "ServerTargetGroup" {
  name     = "server-target-group"
  port     = 80
  protocol = "HTTP"
  vpc_id   = aws_vpc.MainVpc.id
  target_type = "ip"

  health_check {
    interval = 30
    path     = "/"
    protocol = "HTTP"
  }

  stickiness {
    type            = "lb_cookie"
    cookie_duration = 86400
    enabled         = true
  }
}

resource "aws_lb_listener" "ServerListener" {
  load_balancer_arn = aws_lb.ServerAlb.arn
  port     = 80
  protocol = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.ServerTargetGroup.arn
  }
}