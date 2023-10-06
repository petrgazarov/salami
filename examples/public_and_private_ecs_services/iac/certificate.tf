resource "aws_acm_certificate" "ServerCertificate" {
  domain_name       = "app.${var.domain_name}"
  validation_method = "DNS"

  tags = {
    Name = "server-acm-certificate"
  }
}

resource "aws_route53_record" "ValidationRecord" {
  for_each = aws_acm_certificate.ServerCertificate.domain_validation_options
  zone_id  = aws_route53_zone.HostedZone.zone_id
  name     = each.value.resource_record_name
  type     = each.value.resource_record_type
  records  = [each.value.resource_record_value]
  ttl      = 300
}

resource "aws_acm_certificate_validation" "CertificateValidation" {
  certificate_arn         = aws_acm_certificate.ServerCertificate.arn
  validation_record_fqdns = [for r in aws_route53_record.ValidationRecord : r.fqdn]
}