package main

import (
	"fmt"
	"log"
	// "time"

	"github.com/libvirt/libvirt-go"
	"github.com/libvirt/libvirt-go-xml"
)

func main() {
	// Connect to libvirt
	conn, err := libvirt.NewConnect("qemu:///system")

	if err != nil {
		log.Fatalf("Failed to connect to libvirt: %v", err)
	}

	defer conn.Close()

	fmt.Println("Connected to libvirt")

	domainXML := &libvirtxml.Domain{
		Type: "kvm",
		Name: "vm-user-1",
		Memory: &libvirtxml.DomainMemory{
			Value: 1024 * 1024, // 1GB RAM
		},
		VCPU: &libvirtxml.DomainVCPU{
			Value: 2, // 2 vCPUs
		},
		OS: &libvirtxml.DomainOS{
			Type: &libvirtxml.DomainOSType{Arch: "x86_64", Machine: "pc-i440fx-2.11", Type: "hvm"},
		},
		Devices: &libvirtxml.DomainDeviceList{
			Disks: []libvirtxml.DomainDisk{
				{
					Device: "disk",
					Driver: &libvirtxml.DomainDiskDriver{Name: "qemu", Type: "raw"},
					Source: &libvirtxml.DomainDiskSource{File: &libvirtxml.DomainDiskSourceFile{File: "/var/lib/libvirt/images/ubuntu-18.04-server-cloudimg-amd64.img"}},
					Target: &libvirtxml.DomainDiskTarget{Dev: "vda", Bus: "virtio"},
				},
			},
			Interfaces: []libvirtxml.DomainInterface{
				{
					Model: &libvirtxml.DomainInterfaceModel{Type: "virtio"},
					Source: &libvirtxml.DomainInterfaceSource{
						Network: &libvirtxml.DomainInterfaceSourceNetwork{Network: "default", Bridge: "br0"},
					},
				},
			},
		},
	}

	// how to find domain by id
	// domains, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
	// if err != nil {
	// 	log.Fatalf("Failed to list domains: %v", err)
	// }
	// for _, domain := range domains {
	// 	name, err := domain.GetID()
	// 	if err != nil {
	// 		log.Fatalf("Failed to get domain id: %v", err)
	// 	}
	// 	fmt.Println(name)
	//
	xmlData, err := domainXML.Marshal()
	if err != nil {
		log.Fatalf("Failed to generate XML: %v", err)
	}

	fmt.Println(string(xmlData))

	domain, err := conn.DomainDefineXML(xmlData)
	if err != nil {
		log.Fatalf("Failed to define domain: %v", err)
	}
	defer domain.Free()

	// Start VM

	if err := domain.Create(); err != nil {
		log.Fatalf("Failed to start VM: %v", err)
	}

	// Get the VM's IP address (requires `virsh net-dhcp-leases default`)
	// time.Sleep(10 * time.Second) // Wait for DHCP to assign IP
	// leases, err := conn.NetworkLookupByName("default")

	// get ip address
	// network, err := conn.NetworkLookupByName("default")
	// if err != nil {
	// 	log.Fatalf("Failed to get network: %v", err)
	// }
	//
	// // Print VM info
	// fmt.Println("VM Created Successfully! Use SSH to access it:")
	// fmt.Println("ssh user@<vm-ip-address>")
}
