resource "aws_acm_certificate" "ServerCertificate" {
  domain_name       = "app.${var.domain_name}"
  validation_method = "DNS"

  tags = {
    Name = "server-acm-certificate"
  }
}

resource "aws_acm_certificate_validation" "CertificateValidation" {
  certificate_arn         = aws_acm_certificate.ServerCertificate.arn
  validation_record_fqdns = [aws_route53_record.CertificateValidationRecord.fqdn]
}