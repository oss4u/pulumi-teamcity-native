import * as pulumi from "@pulumi/pulumi";
import * as teamcity from "@oss4u/teamcity";

const myRandomResource = new teamcity.Random("myRandomResource", {length: 24});
export const output = {
    value: myRandomResource.result,
};
