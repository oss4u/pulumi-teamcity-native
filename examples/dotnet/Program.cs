using System.Collections.Generic;
using System.Linq;
using Pulumi;
using Teamcity = Oss4u.Teamcity;

return await Deployment.RunAsync(() => 
{
    var myRandomResource = new Teamcity.Random("myRandomResource", new()
    {
        Length = 24,
    });

    return new Dictionary<string, object?>
    {
        ["output"] = 
        {
            { "value", myRandomResource.Result },
        },
    };
});

