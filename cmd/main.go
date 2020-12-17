package main

import (
	"bufio"
	"fmt"
	"github.com/relex/aini"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var customers []string

func main() {
	fmt.Println("ok")
	//"../{{ platform or customer }}"
	//"../{{ platform }}/{{customer}}"
	projectDir := "/home/manu/Projects/ansible-kube"
	platform := "os"
	inventoryDir := fmt.Sprintf("%s/inventories/%s", projectDir, platform)

	files, err := ioutil.ReadDir(inventoryDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		customers = append(customers, f.Name())
	}

	fmt.Println(customers)


	for _, customer := range customers {
		inventoryData, err := parseInventoryWithCustomer(inventoryDir, customer)
		if err != nil {
			//fmt.Printf("Skipping the inventory %s cause: %s\n", customer, err.Error())
		//} else if filterByPlatform(inventoryData, "osf1") {
		//	fmt.Print("ansible-playbook -i ", inventoryDir, "/", customer, "/host.ini", " ", "kube-asset.yml", "\n")
		} else if filterByEnvironment(inventoryData, "hors_prod") {
			fmt.Print("ansible-playbook -i ", inventoryDir, "/", customer, "/host.ini", " ", "kube-asset.yml", "\n")
		}
	}

}

func parseInventory(inventoryPath string) (data *aini.InventoryData) {
	file, _ := os.Open(inventoryPath)
	reader := bufio.NewReader(file)
	data, _ = aini.Parse(reader)
	return
}

func parseInventoryWithCustomer(inventoryPath string, customer string) (data *aini.InventoryData, err error) {
	file, err := os.Open(inventoryPath + "/" + customer + "/hosts.ini")
	if err != nil {
		return
	}
	reader := bufio.NewReader(file)
	data, _ = aini.Parse(reader)
	return
}


func filterByPlatform(inventoryData *aini.InventoryData, plaformPrefix string) bool {

	for adminNode := range inventoryData.Groups["ansible_admin_node"].Hosts {
		if strings.Contains(adminNode,plaformPrefix) {
			return true
		}
	}
	return false
}

func filterByEnvironment(inventoryData *aini.InventoryData, environment string) bool {
	return strings.HasSuffix(inventoryData.Groups["all"].Vars["customer"], environment)
}