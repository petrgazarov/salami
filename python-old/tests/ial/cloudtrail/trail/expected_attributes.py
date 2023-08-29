resource_type = "cloudtrail.Trail"
logical_name = "CumuliCloudtrail"
text = (
    "Name: cumuli-cloudtrail\nGlobal service events: True\n"
    "Multi-region: True\nLogging enabled: True"
)
other_packages = []
uses = ["CumuliCloudtrailLogsBucket"]
exports = {}
