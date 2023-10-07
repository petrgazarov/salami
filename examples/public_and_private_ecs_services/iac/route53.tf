resource "aws_route53_zone" "HostedZone" {
  name = var.domain_name
}

resource "aws_route53_record" "AppAlbRecord" {
  zone_id = aws_route53_zone.HostedZone.zone_id
  name    = "app.${var.domain_name}"
  type    = "A"

  alias {
    name                   = aws_lb.ServerAlb.dns_name
    zone_id                = aws_lb.ServerAlb.zone_id
    evaluate_target_health = true
  }
}

resource "aws_route53_record" "CertificateValidationRecord" {
  zone_id = aws_route53_zone.HostedZone.zone_id
  name    = tolist(aws_acm_certificate.ServerCertificate.domain_validation_options)[0].resource_record_name
  type    = tolist(aws_acm_certificate.ServerCertificate.domain_validation_options)[0].resource_record_type
  records = [tolist(aws_acm_certificate.ServerCertificate.domain_validation_options)[0].resource_record_value]
  ttl     = 300
}