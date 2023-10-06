resource "aws_acm_certificate" "ServerCertificate" {
  domain_name       = "app.${var.domain_name}"
  validation_method = "DNS"

  tags = {
    Name = "server-acm-certificate"
  }
}

resource "aws_route53_record" "ValidationRecord" {
  zone_id = aws_route53_zone.HostedZone.zone_id
  name    = tolist(aws_acm_certificate.ServerCertificate.domain_validation_options)[0].resource_record_name
  type    = tolist(aws_acm_certificate.ServerCertificate.domain_validation_options)[0].resource_record_type
  records = [tolist(aws_acm_certificate.ServerCertificate.domain_validation_options)[0].resource_record_value]
  ttl     = 300
}

resource "aws_acm_certificate_validation" "CertificateValidation" {
  certificate_arn         = aws_acm_certificate.ServerCertificate.arn
  validation_record_fqdns = [aws_route53_record.ValidationRecord.fqdn]
}