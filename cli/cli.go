package cli

import (
	"fmt"

	"github.com/docopt/docopt-go"
	"github.com/smartystreets/raptr/messages"
)

func ReadMessage() (interface{}, string) {
	usage := `
Usage:
    raptr upload
        --name=N
        --path=P
        --storage=S
        --category=C
        [--distribution=D]
        [--config=CONF]
    raptr link
        --name=N
        --version=V
        --storage=S
        --category=C
        --source-distribution=SD
        --target-distribution=TD
        [--config=CONF]
    raptr unlink
        --name=N
        --version=V
        --storage=S
        --category=C
        --distribution=D
        [--config=CONF]
    raptr clean
        --storage=S
        [--config=CONF]
    raptr -h

Examples:
    raptr upload --name=liveaddress-logging --path=/package/files --storage=s3 --category=operations --distribution=staging
    raptr link --name=liveaddress-logging --version=1.0.7 --storage=s3 --category=liveaddress --source-distribution=staging --target-distribution=production
    raptr unlink --name=liveaddress-logging --version=1.0.7 --storage=s3 --category=liveaddress --distribution=production
    raptr clean --storage=named-storage


Options:
    --name=N                  Name of the package
    --version=V               Version of the package
    --path=P                  Path to package files
    --storage=S               Storage mechanism to use (eg: s3, filesystem, etc)
    --category=C              Category the package falls under
    --distribution=D          Distribution to link to (optional for upload, if exists, links package to distribution)
    --source-distribution=SD  Source distribution
    --target-distribution=TD  Target distribution
    --config=CONF             Path to a config file (optional)
    -h  Shows this screen
`

	//func Parse(doc string, argv []string, help bool, version string, optionsFirst bool, exit ...bool) (map[string]interface{}, error)
	arguments, err := docopt.Parse(usage, nil, true, "raptr (Remote Apt Repository) version 1.2.3", false, true)
	if err != nil {
		fmt.Println("error:", err)
	}

	configFile := ""
	if c, _ := arguments["--config"]; c != nil && len(c.(string)) > 0 {
		configFile = c.(string)
	}
	distribution := ""
	if d, _ := arguments["--distribution"]; d != nil {
		distribution = arguments["--distribution"].(string)
	}

	if arguments["upload"] == true {
		return messages.UploadCommand{
			PackageName:  arguments["--name"].(string), // .([]string)[0],
			PackagePath:  arguments["--path"].(string),
			StorageName:  arguments["--storage"].(string),
			Category:     arguments["--category"].(string),
			Distribution: distribution, // optional
		}, configFile
	} else if arguments["link"] == true {
		return messages.LinkCommand{
			PackageName:        arguments["--name"].(string),
			PackageVersion:     arguments["--version"].(string),
			StorageName:        arguments["--storage"].(string),
			Category:           arguments["--category"].(string),
			SourceDistribution: arguments["--source-distribution"].(string),
			TargetDistribution: arguments["--target-distribution"].(string),
		}, configFile
	} else if arguments["unlink"] == true {
		return messages.UnlinkCommand{
			PackageName:    arguments["--name"].(string),
			PackageVersion: arguments["--version"].(string),
			StorageName:    arguments["--storage"].(string),
			Category:       arguments["--category"].(string),
			Distribution:   distribution,
		}, configFile
	} else if arguments["clean"] == true {
		return messages.CleanCommand{
			StorageName: arguments["--storage"].(string),
		}, configFile
	}

	return nil, ""
}
