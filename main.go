package main

import (
	"fmt"
	"log"

	"github.com/libvirt/libvirt-go"
	"github.com/libvirt/libvirt-go-xml"
)

func uintPtr(value uint) *uint {
	return &value
}

func main() {

	conn, err := libvirt.NewConnect("qemu:///system")

	if err != nil {
		log.Fatalf("Failed to connect to libvirt: %v", err)
	}

	defer conn.Close()

	fmt.Println("Connected to libvirt")

	domainXML := &libvirtxml.Domain{
		Type: "kvm",
		Name: "unbuntu20",
		UUID: "b3cca077-4e9f-4a9d-af69-bb17437ed6cd", // Using the provided UUID
		Memory: &libvirtxml.DomainMemory{
			Value: 2048,
			Unit:  "MiB",
		},
		CurrentMemory: &libvirtxml.DomainCurrentMemory{
			Value: 2048,
			Unit:  "MiB",
		},
		VCPU: &libvirtxml.DomainVCPU{
			Placement: "static",
			Value:     2, // 2 CPUs
		},
		OS: &libvirtxml.DomainOS{
			Type: &libvirtxml.DomainOSType{
				Arch:    "x86_64",
				Machine: "pc-q35-6.2",
				Type:    "hvm",
			},
			BootDevices: []libvirtxml.DomainBootDevice{
				{Dev: "cdrom"}, // Boot from ISO first
				{Dev: "hd"},    // Then hard disk
			},
		},
		Features: &libvirtxml.DomainFeatureList{
			ACPI: &libvirtxml.DomainFeature{},
		},
		CPU: &libvirtxml.DomainCPU{
			Mode:       "host-passthrough",
			Check:      "none",
			Migratable: "on",
		},
		Clock: &libvirtxml.DomainClock{
			Offset: "utc",
		},
		OnPoweroff: "destroy",
		OnReboot:   "destroy",
		OnCrash:    "destroy",
		Devices: &libvirtxml.DomainDeviceList{
			Disks: []libvirtxml.DomainDisk{
				// {
				// 	Device: "disk",
				// 	Driver: &libvirtxml.DomainDiskDriver{
				// 		Name:    "qemu",
				// 		Type:    "qcow2",
				// 		Discard: "unmap",
				// 	},
				// 	Source: &libvirtxml.DomainDiskSource{
				// 		File: &libvirtxml.DomainDiskSourceFile{
				// 			File: "/var/lib/libvirt/images/debian-12-genericcloud-amd64.qcow2",
				// 		},
				// 		Index: 1,
				// 	},
				// 	Target: &libvirtxml.DomainDiskTarget{
				// 		Dev: "vda",
				// 		Bus: "virtio",
				// 	},
				// },
				{
					Device: "disk",
					Driver: &libvirtxml.DomainDiskDriver{
						Name: "qemu",
						Type: "qcow2",
					},
					Source: &libvirtxml.DomainDiskSource{
						File: &libvirtxml.DomainDiskSourceFile{
							File: "/var/lib/libvirt/images/alexng.img",
						},
						Index: 1,
					},
					Target: &libvirtxml.DomainDiskTarget{
						Dev: "vda",
						Bus: "virtio",
					},
					Alias: &libvirtxml.DomainAlias{
						Name: "virtio-disk0",
					},
					Address: &libvirtxml.DomainAddress{
						PCI: &libvirtxml.DomainAddressPCI{
							Domain:   uintPtr(0),
							Bus:      uintPtr(4),
							Slot:     uintPtr(0),
							Function: uintPtr(0),
						},
					},

					ReadOnly: &libvirtxml.DomainDiskReadOnly{},
				},
				// {
				// 	Device: "cdrom",
				// 	Driver: &libvirtxml.DomainDiskDriver{
				// 		Name: "qemu",
				// 		Type: "raw",
				// 	},
				// 	Source: &libvirtxml.DomainDiskSource{
				// 		File: &libvirtxml.DomainDiskSourceFile{
				// 			File: "/var/lib/libvirt/images/cloud-init.iso",
				// 		},
				// 		Index: 2,
				// 	},
				// 	Target: &libvirtxml.DomainDiskTarget{
				// 		Dev: "sda",
				// 		Bus: "sata",
				// 	},
				// 	ReadOnly: &libvirtxml.DomainDiskReadOnly{},
				// },
			},
			Interfaces: []libvirtxml.DomainInterface{
				{
					MAC: &libvirtxml.DomainInterfaceMAC{
						Address: "52:54:00:b7:a5:c2",
					},

					Source: &libvirtxml.DomainInterfaceSource{
						Network: &libvirtxml.DomainInterfaceSourceNetwork{
							Network: "default",
							Bridge:  "virbr0",
						},
					},
					Model: &libvirtxml.DomainInterfaceModel{
						Type: "virtio",
					},
				},
			},
			Graphics: []libvirtxml.DomainGraphic{
				{
					VNC: &libvirtxml.DomainGraphicVNC{
						Port:   -1,
						Listen: "0.0.0.0",
					},

					Spice: &libvirtxml.DomainGraphicSpice{
						Port:   -1,
						Listen: "0.0.0.0",
					},
				},
			},
		},
	}

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

	if err := domain.Create(); err != nil {
		log.Fatalf("Failed to start VM: %v", err)
	}

	domain.GetInterfaceParameters("vnet1", 0)
}
