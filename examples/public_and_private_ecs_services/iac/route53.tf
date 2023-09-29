resource "aws_route53_zone" "HostedZone" {
  name = var.domain_name
}

resource "aws_route53_record" "FramerARecords" {
  zone_id = aws_route53_zone.HostedZone.zone_id
  name    = var.domain_name
  type    = "A"
  records = ["52.223.52.2", "35.71.142.77"]
  ttl     = 1800
}

resource "aws_route53_record" "FramerCNAMERecord" {
  zone_id = aws_route53_zone.HostedZone.zone_id
  name    = "www.${var.domain_name}"
  type    = "CNAME"
  records = ["sites.framer.app."]
  ttl     = 1800
}

resource "aws_route53_record" "GoogleTXTRecord" {
  zone_id = aws_route53_zone.HostedZone.zone_id
  name    = var.domain_name
  type    = "TXT"
  ttl     = 1800
  records = ["google-site-verification=UzsIO5GYrHN8MrT_xe_PgmbydavWUcgC1cUfwJHpWB4"]
}