// sg-update project main.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
	"golang.org/x/crypto/ssh"
)

type Versions struct {
	Package_name    string
	Package_version string
}

type JsonObject struct {
	Defaults []Versions
}

func GetDefaultsFromFile(file string) JsonObject {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		panic("Missing defaults.json file " + err.Error())
	}
	var default_json JsonObject
	err = json.Unmarshal(content, &default_json)
	if err != nil {
		panic("Unmarshalling error from Json " + err.Error())
	}
	return default_json
}

func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func SshGetVersion(dnsname string) {

	cmd := "nmap"
	args := []string{"--script", "ssl-cert,ssl-enum-ciphers-2", " -p 443 ", dnsname}
	cmdOut, err := exec.Command(cmd, args...).Output()
	if err != nil {
		fmt.Println(dnsname, " : nmap error..skipping")
	} else {
		fmt.Println(dnsname, " : ", string(cmdOut))
	}
}

func main() {
	fmt.Println("Start time: ", time.Now())
	//	defaults := GetDefaultsFromFile("./defaults.json")
	svc := elb.New(session.New(), aws.NewConfig().WithRegion("us-east-1"))

	params := &elb.DescribeLoadBalancersInput{
	//		LoadBalancerNames: []*string{
	//			aws.String("AccessPointName"),
	//		},
	}

	instances, err := svc.DescribeLoadBalancers(params)
	if err != nil {
		panic(err)
	}
	// fmt.Println(instances.LoadBalancerDescriptions)
	// fmt.Println("> Number of instances: ", len(instances.LoadBalancerDescriptions))
	for idx := 0; idx < len(instances.LoadBalancerDescriptions); idx++ {
		// for _, inst := range instances.LoadBalancerDescriptions[idx] {
		//			hostapp := "UNKNOWN"
		//			for _, tag := range inst.Tags {
		//				if *tag.Key == "Name" {
		//					hostapp = *tag.Value
		//				}
		//			}
		// fmt.Println(*instances.LoadBalancerDescriptions[idx].DNSName)
		go SshGetVersion(*instances.LoadBalancerDescriptions[idx].DNSName)

		//}
	}
	// fmt.Println("> first element: ", instances.Reservations[0])
	time.Sleep(120 * 1e9)
	fmt.Println("End time: ", time.Now())
}
