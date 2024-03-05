import pulumi
import pulumi_teamcity as teamcity

my_random_resource = teamcity.Random("myRandomResource", length=24)
pulumi.export("output", {
    "value": my_random_resource.result,
})
