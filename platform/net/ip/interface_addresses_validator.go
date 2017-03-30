package ip

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type InterfaceAddressesValidator interface {
	Validate(desiredInterfaceAddresses []InterfaceAddress) error
}

type interfaceAddressesValidator struct {
	interfaceAddrsProvider InterfaceAddressesProvider
}

func NewInterfaceAddressesValidator(interfaceAddrsProvider InterfaceAddressesProvider) InterfaceAddressesValidator {
	return &interfaceAddressesValidator{
		interfaceAddrsProvider: interfaceAddrsProvider,
	}
}

func (i *interfaceAddressesValidator) Validate(desiredInterfaceAddresses []InterfaceAddress) error {
	//default via 159.8.154.33 dev eth1
	//10.0.0.0/8 via 10.113.205.1 dev eth0
	//10.112.166.128/26 dev eth0  proto kernel  scope link  src 10.112.166.152
	//10.113.205.0/24 dev eth0  proto kernel  scope link  src 10.113.205.134
	//159.8.154.32/27 dev eth1  proto kernel  scope link  src 159.8.154.62
	//161.26.0.0/16 via 10.113.205.1 dev eth0

	//{ip.simpleInterfaceAddress{interfaceName:"lo", ip:"127.0.0.1"},
	// ip.simpleInterfaceAddress{interfaceName:"eth0", ip:"10.113.205.134"},
	// ip.simpleInterfaceAddress{interfaceName:"eth0", ip:"10.112.166.152"},
	// ip.simpleInterfaceAddress{interfaceName:"eth1", ip:"159.8.154.62"}}
	systemInterfaceAddresses, err := i.interfaceAddrsProvider.Get()
	if err != nil {
		return bosherr.WrapError(err, "Getting network interface addresses")
	}

	for _, desiredInterfaceAddress := range desiredInterfaceAddresses {
		ifaceName := desiredInterfaceAddress.GetInterfaceName()
		iface, found := i.findInterfaceByName(ifaceName, systemInterfaceAddresses)
		if !found {
			return bosherr.WrapErrorf(err, "Validating network interface '%s' IP addresses, no interface configured with that name", ifaceName)
		}
		desiredIP, _ := desiredInterfaceAddress.GetIP()
		actualIP, _ := iface.GetIP()
		if desiredIP != actualIP {
			return bosherr.WrapErrorf(err, "Validating network interface '%s' IP addresses, expected: '%s', actual: '%s'", ifaceName, desiredIP, actualIP)
		}
	}

	return nil
}

func (i *interfaceAddressesValidator) findInterfaceByName(ifaceName string, ifaces []InterfaceAddress) (InterfaceAddress, bool) {
	for _, iface := range ifaces {
		if iface.GetInterfaceName() == ifaceName {
			return iface, true
		}
	}

	return nil, false
}
