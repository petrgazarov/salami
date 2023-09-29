resource "aws_security_group" "ALBSecurityGroup" {
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
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_lb" "ServerAlb" {
  name               = "server-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.ALBSecurityGroup.id]
  subnets            = [aws_subnet.PublicSubnetA.id, aws_subnet.PublicSubnetB.id]
  enable_deletion_protection = true
  idle_timeout = 3600
}

resource "aws_lb_target_group" "ServerTargetGroup" {
  name     = "server-target-group"
  port     = 80
  protocol = "HTTP"
  vpc_id   = aws_ec2_vpc.MainVpc.id
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
  load_balancer_arn = aws_lb.LoadBalancer.ServerAlb.arn
  port     = 443
  protocol = "HTTPS"
  ssl_policy = "ELBSecurityPolicy-TLS13-1-2-2021-06"
  certificate_arn = aws_acm_certificate.ServerCertificate.arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.ServerTargetGroup.arn
  }
}

resource "aws_route53_record" "AppALBRecord" {
  zone_id = aws_route53_zone.HostedZone.zone_id
  name    = "app.${var.domain_name}"
  type    = "A"

  alias {
    name                   = aws_lb.LoadBalancer.ServerAlb.dns_name
    zone_id                = aws_lb.LoadBalancer.ServerAlb.zone_id
    evaluate_target_health = true
  }
}