resource_type = "route53.Record"
logical_name = "CumuliValidationRecord"
text = "Ttl: 300"
other_packages = []
properties = [
    "The name, type and records are derived from the CumuliServerCertificate's domain validation option. "
    'E.g. domain_validation_options[0]["resource_record_name"]'
]
uses = ["CumuliHostedZone", "CumuliServerCertificate"]
exports = {}
