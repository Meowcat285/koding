package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/ec2"
	"github.com/mitchellh/goamz/elb"
)

var (
	flagHost     = flag.String("env", "prod", "env name")
	flagFiltered = flag.Bool("filtered", false, "filter by ec2 tags")
)

var ELBS = map[string]string{
	// app
	"prod":       "awseb-e-x-AWSEBLoa-2AG3XORA8JXC",
	"latest":     "awseb-e-3-AWSEBLoa-1S2VPBAQXDRW9",
	"sandbox":    "awseb-e-p-AWSEBLoa-1POHSLP6A7STY",
	"monitoring": "awseb-e-j-AWSEBLoa-AQBHUYZM5ZX6",
	// proxy
	"proxy-eu-west-1":      "awseb-e-s-AWSEBLoa-1LOTB5BKTJJBW",
	"proxy-us-east-1":      "awseb-e-a-AWSEBLoa-RTLJ62SKJY5G",
	"proxy-us-west-2":      "awseb-e-7-AWSEBLoa-1V808KG9PDQH5",
	"proxy-ap-southeast-1": "awseb-e-u-AWSEBLoa-15H1DQTBBUMG",
	"proxy-dev-us-e-1":     "awseb-e-p-AWSEBLoa-16QIN1OI6WYNK",
}

var Tags = map[string]string{
	// app
	"prod":    "prod.koding.com",
	"latest":  "latest.koding.com",
	"sandbox": "sandbox.koding.com",
}

var ELB2Region = map[string]aws.Region{
	// proxies
	"awseb-e-s-AWSEBLoa-1LOTB5BKTJJBW": aws.EUWest,
	"awseb-e-a-AWSEBLoa-RTLJ62SKJY5G":  aws.USEast,
	"awseb-e-7-AWSEBLoa-1V808KG9PDQH5": aws.USWest2,
	"awseb-e-u-AWSEBLoa-15H1DQTBBUMG":  aws.APSoutheast,

	// app ELBs
	"awseb-e-j-AWSEBLoa-AQBHUYZM5ZX6":  aws.USEast,
	"awseb-e-x-AWSEBLoa-2AG3XORA8JXC":  aws.USEast,
	"awseb-e-3-AWSEBLoa-1S2VPBAQXDRW9": aws.USEast,
	"awseb-e-p-AWSEBLoa-1POHSLP6A7STY": aws.USEast,
	"awseb-e-p-AWSEBLoa-16QIN1OI6WYNK": aws.USEast,
}

var auth = aws.Auth{
	AccessKey: "AKIAI7CKP5SNHCBUEDXQ",
	SecretKey: "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",
}

func getEC2(elbName string) *ec2.EC2 {
	reg, ok := ELB2Region[elbName]
	if !ok {
		panic(elbName)
	}

	return ec2.New(auth, reg)
}

func getELB(elbName string) *elb.ELB {
	reg, ok := ELB2Region[elbName]
	if !ok {
		panic(elbName)
	}

	return elb.New(auth, reg)
}

// func join() chan error {
// 	c := make(chan error, 2)

// 	copy := func(i int) {
// 		time.Sleep(time.Millisecond * time.Duration(100*i))
// 		select {
// 		case c <- fmt.Errorf("%d", i):
// 		default:
// 			fmt.Println("i-->", i)
// 		}
// 	}

// 	go copy(1)
// 	go copy(2)

// 	return c
// }
//

type KodingHstore gorm.Hstore

func main() {
	// go func() {
	// 	c := join()
	// 	fmt.Println(<-c)
	// 	c = nil
	// 	close(c)
	// }()

	// time.Sleep(time.Second * 2)
	// fmt.Println("Hello, playground")
	// endOfDay := now.EndOfDay().UTC()
	// difference := time.Now().UTC().Sub(endOfDay)
	// fmt.Println("difference-->", difference)

	// endOfDay = now.EndOfDay().UTC()
	// difference = endOfDay.Sub(time.Now().UTC())
	// fmt.Println("difference-->", difference)
	// return

	flag.Parse()

	currentELBInstances := make([]string, 0)

	elbNames, ok := ELBS[*flagHost]
	if !ok {
		log.Fatalf("%s not found", *flagHost)
	}

	if *flagFiltered {
		tagName := Tags[*flagHost]
		var err error
		currentELBInstances, err = GetchInstancesByTag(elbNames, tagName)
		if err != nil {
			log.Fatal(err.Error())
		}
	} else {
		for _, elbName := range strings.Split(elbNames, ",") {
			currentInstances, err := fetchProdELBAttachedInstances(elbName)
			if err != nil {
				log.Fatal(err.Error())
			}
			currentELBInstances = append(currentELBInstances, currentInstances...)
		}
	}
	fmt.Println("currentELBInstances-->", currentELBInstances)
	fmt.Println("strings.Fields(createI2csshString(currentELBInstances))...-->", strings.Fields(createI2csshString(*flagHost, currentELBInstances)))
	_, err := exec.Command("i2cssh", strings.Fields(createI2csshString(*flagHost, currentELBInstances))...).Output()
	if err != nil {
		log.Fatal(err.Error())
	}

}

var paramToKey = map[string]string{
	"prod":       "/Users/siesta/Documents/koding/credential/private_keys/koding-eb-deployment-us-east-1-2015-06.pem",
	"latest":     "/Users/siesta/Documents/koding/credential/private_keys/koding-eb-deployment-us-east-1-2015-06.pem",
	"sandbox":    "/Users/siesta/Documents/koding/credential/private_keys/koding-eb-deployment-us-east-1-2015-06.pem",
	"monitoring": "/Users/siesta/Documents/koding/credential/private_keys/koding-eb-deployment-us-east-1-2015-06.pem",

	"proxy-eu-west-1":      "/Users/siesta/Documents/koding/credential/private_keys/koding-eb-deployment-eu-west-1-2015-06.pem",
	"proxy-us-east-1":      "/Users/siesta/Documents/koding/credential/private_keys/koding-eb-deployment-us-east-1-2015-06.pem",
	"proxy-us-west-2":      "/Users/siesta/Documents/koding/credential/private_keys/koding-eb-deployment-us-west-2-2015-06.pem",
	"proxy-ap-southeast-1": "/Users/siesta/Documents/koding/credential/private_keys/koding-eb-deployment-ap-southeast-1-2015-06.pem",
	"proxy-dev-us-e-1":     "/Users/siesta/Documents/koding/credential/private_keys/koding-eb-deployment-us-east-1-2015-06.pem",
}

func createI2csshString(param string, instances []string) string {
	return fmt.Sprintf(
		"--forward-agent --login ec2-user --rows 3 --broadcast -Xi=%s  --machines %s",
		paramToKey[param],
		strings.Join(instances, ","),
	)
}

// fetch currently attached instances to the prod ELB
func fetchProdELBAttachedInstances(elbName string) ([]string, error) {
	dlr := &elb.DescribeLoadBalancer{
		Names: []string{elbName},
	}

	res, err := getELB(elbName).DescribeLoadBalancers(dlr)
	if err != nil {
		return nil, err
	}

	if len(res.LoadBalancers) < 1 {
		return nil, fmt.Errorf("%s ELB not found!", elbName)
	}

	attachedInstances := make([]string, 0)
	for _, inst := range res.LoadBalancers[0].Instances {
		attachedInstances = append(attachedInstances, inst.InstanceId)
	}

	resp, err := getEC2(elbName).Instances(attachedInstances, ec2.NewFilter())
	if err != nil {
		return nil, err
	}

	instanceIps := make([]string, 0)
	for _, reservations := range resp.Reservations {
		for _, instance := range reservations.Instances {
			if instance.PublicIpAddress != "" {
				instanceIps = append(instanceIps, instance.PublicIpAddress)
			}
		}
	}

	return instanceIps, nil
}

func GetchInstancesByTag(elbName string, environmentTag string) ([]string, error) {

	filter := ec2.NewFilter()
	filter.Add("tag-key", "environment")
	filter.Add("tag-value", environmentTag)

	resp, err := getEC2(elbName).Instances(nil, filter)
	if err != nil {
		return nil, err
	}

	instances := make([]string, 0)
	for _, reservations := range resp.Reservations {
		for _, instance := range reservations.Instances {
			if instance.PublicIpAddress != "" {
				instances = append(instances, instance.PublicIpAddress)
			}
		}
	}

	return instances, nil
}
