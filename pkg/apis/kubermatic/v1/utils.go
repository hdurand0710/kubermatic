/*
Copyright 2020 The Kubermatic Kubernetes Platform contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	"fmt"
	"math"
	"math/big"
	"net"
	"strings"

	netutils "k8s.io/utils/net"
)

// AllExposeStrategies is a set containing all the ExposeStrategy.
var AllExposeStrategies = NewExposeStrategiesSet(ExposeStrategyNodePort, ExposeStrategyLoadBalancer, ExposeStrategyTunneling)

// ExposeStrategyFromString returns the expose strategy which String
// representation corresponds to the input string, and a bool saying whether a
// match was found or not.
func ExposeStrategyFromString(s string) (ExposeStrategy, bool) {
	es := ExposeStrategy(s)
	return es, AllExposeStrategies.Has(es)
}

// String returns the string representation of the ExposeStrategy.
func (e ExposeStrategy) String() string {
	return string(e)
}

// ExposeStrategiesSet is a set of ExposeStrategies.
type ExposeStrategiesSet map[ExposeStrategy]struct{}

// NewByte creates a ExposeStrategiesSet from a list of values.
func NewExposeStrategiesSet(items ...ExposeStrategy) ExposeStrategiesSet {
	es := ExposeStrategiesSet{}
	for _, item := range items {
		es[item] = struct{}{}
	}
	return es
}

// Has returns true if and only if item is contained in the set.
func (e ExposeStrategiesSet) Has(item ExposeStrategy) bool {
	_, contained := e[item]
	return contained
}

// Has returns true if and only if item is contained in the set.
func (e ExposeStrategiesSet) String() string {
	es := make([]string, 0, len(e))
	for k := range e {
		es = append(es, string(k))
	}
	// can be easily optimized in terms of allocations by using a bytes buffer
	// with some more verbosity, but this is not supposed to be called in a
	// perf critical path.
	return fmt.Sprintf("[%s]", strings.Join(es, ", "))
}

func (e ExposeStrategiesSet) Items() []string {
	var items []string
	for s := range e {
		items = append(items, string(s))
	}
	return items
}

// IsIPv4Only returns true if the cluster networking is IPv4-only.
func (c *Cluster) IsIPv4Only() bool {
	return len(c.Spec.ClusterNetwork.Pods.CIDRBlocks) == 1 && netutils.IsIPv4CIDRString(c.Spec.ClusterNetwork.Pods.CIDRBlocks[0])
}

// IsIPv6Only returns true if the cluster networking is IPv6-only.
func (c *Cluster) IsIPv6Only() bool {
	return len(c.Spec.ClusterNetwork.Pods.CIDRBlocks) == 1 && netutils.IsIPv6CIDRString(c.Spec.ClusterNetwork.Pods.CIDRBlocks[0])
}

// IsDualStack returns true if the cluster networking is dual-stack (IPv4 + IPv6).
func (c *Cluster) IsDualStack() bool {
	res, err := netutils.IsDualStackCIDRStrings(c.Spec.ClusterNetwork.Pods.CIDRBlocks)
	if err != nil {
		return false
	}
	return res
}

// Validate validates the network ranges. Returns nil if valid, error otherwise.
func (r *NetworkRanges) Validate() error {
	if r == nil {
		return nil
	}
	for _, cidr := range r.CIDRBlocks {
		if _, _, err := net.ParseCIDR(cidr); err != nil {
			return fmt.Errorf("unable to parse CIDR %q: %w", cidr, err)
		}
	}
	return nil
}

// GetIPv4CIDR returns the first found IPv4 CIDR in the network ranges, or an empty string if no IPv4 CIDR is found.
func (r *NetworkRanges) GetIPv4CIDR() string {
	for _, cidr := range r.CIDRBlocks {
		if netutils.IsIPv4CIDRString(cidr) {
			return cidr
		}
	}
	return ""
}

// GetIPv4CIDRs returns all IPv4 CIDRs in the network ranges, or an empty string if no IPv4 CIDR is found.
func (r *NetworkRanges) GetIPv4CIDRs() (res []string) {
	for _, cidr := range r.CIDRBlocks {
		if netutils.IsIPv4CIDRString(cidr) {
			res = append(res, cidr)
		}
	}
	return
}

// HasIPv4CIDR returns true if the network ranges contain any IPv4 CIDR, false otherwise.
func (r *NetworkRanges) HasIPv4CIDR() bool {
	return r.GetIPv4CIDR() != ""
}

// GetIPv6CIDR returns the first found IPv6 CIDR in the network ranges, or an empty string if no IPv6 CIDR is found.
func (r *NetworkRanges) GetIPv6CIDR() string {
	for _, cidr := range r.CIDRBlocks {
		if netutils.IsIPv6CIDRString(cidr) {
			return cidr
		}
	}
	return ""
}

// GetIPv6CIDRs returns all IPv6 CIDRs in the network ranges, or an empty string if no IPv6 CIDR is found.
func (r *NetworkRanges) GetIPv6CIDRs() (res []string) {
	for _, cidr := range r.CIDRBlocks {
		if netutils.IsIPv6CIDRString(cidr) {
			res = append(res, cidr)
		}
	}
	return
}

// HasIPv6CIDR returns true if the network ranges contain any IPv6 CIDR, false otherwise.
func (r *NetworkRanges) HasIPv6CIDR() bool {
	return r.GetIPv6CIDR() != ""
}

// GetIpNets returns all network ranges in the IPNet form.
func (r *NetworkRanges) GetIpNets() ([]*net.IPNet, error) {
	ipnets := make([]*net.IPNet, 0)
	for _, cidr := range r.CIDRBlocks {
		_, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, err
		}
		ipnets = append(ipnets, ipnet)
	}
	return ipnets, nil
}

// SplitCIDRByNumberOfAddresses splits each CIDR string into subnets with the closest number of addresses per subnet.
func (r *NetworkRanges) SplitByNumberOfAddresses(number uint32) (*NetworkRanges, error) {
	ipnets, err := r.GetIpNets()
	if err != nil {
		return nil, err
	}

	res := NetworkRanges{CIDRBlocks: make([]string, 0)}
	for _, ipnet := range ipnets {
		splittedIpNets, err := splitIPNetByNumberOfAddresses(ipnet, number)
		if err != nil {
			return nil, err
		}
		for _, splittedIpNet := range splittedIpNets {
			res.CIDRBlocks = append(res.CIDRBlocks, splittedIpNet.String())
		}
	}

	return &res, nil
}

// splitIPNetByNumberOfAddresses splits an IPNet into subnets with the closest number of addresses per subnet.
func splitIPNetByNumberOfAddresses(ipnet *net.IPNet, number uint32) ([]*net.IPNet, error) {
	ipsNumber := addressCountIpnet(ipnet)

	// truncate result to nearest uint64
	optimalSplit := int(ipsNumber / uint64(number))
	return splitIPNetIntoN(ipnet, optimalSplit)
}

// adressCountIpnet returns the numbers of addresses for this IPNet.
func addressCountIpnet(network *net.IPNet) uint64 {
	prefixLen, bits := network.Mask.Size()
	return 1 << (uint64(bits) - uint64(prefixLen))
}

// splitIPNetIntoN splits an IPNet into n subranges.
func splitIPNetIntoN(iprange *net.IPNet, n int) ([]*net.IPNet, error) {
	var err error

	// invalid value
	if n <= 1 || addressCountIpnet(iprange) < uint64(n) {
		return nil, fmt.Errorf("Cannot split %s into %d subranges", iprange, n)
	}
	// power of two
	if isPowerOfTwo(n) || isPowerOfTwoPlusOne(n) {
		return splitIPNet(iprange, n)
	}

	var closestMinorPowerOfTwo int
	// find the closest power of two in a stupid way
	for i := n; i > 0; i-- {
		if isPowerOfTwo(i) {
			closestMinorPowerOfTwo = i
			break
		}
	}

	subnets, err := splitIPNet(iprange, closestMinorPowerOfTwo)
	if err != nil {
		return nil, err
	}
	for len(subnets) < n {
		var newSubnets []*net.IPNet
		level := 1
		for i := len(subnets) - 1; i >= 0; i-- {
			divided, err := divideIPNet(subnets[i])
			if err != nil {
				return nil, err
			}
			newSubnets = append(newSubnets, divided...)
			if len(subnets)-level+len(newSubnets) == n {
				reverseIPNet(newSubnets)
				subnets = subnets[:len(subnets)-level]
				subnets = append(subnets, newSubnets...)
				return subnets, nil
			}
			level++
		}
		reverseIPNet(newSubnets)
		subnets = newSubnets
	}
	return subnets, nil
}

// reverseIPNet reverses a slice of IPNet.
func reverseIPNet(ipnets []*net.IPNet) {
	for i, j := 0, len(ipnets)-1; i < j; i, j = i+1, j-1 {
		ipnets[i], ipnets[j] = ipnets[j], ipnets[i]
	}
}

// isPowerOfTwo returns if a number is a power of 2.
func isPowerOfTwo(x int) bool {
	return x != 0 && (x&(x-1)) == 0
}

// isPowerOfTwoPlusOne returns if a number is a power of 2 plus 1.
func isPowerOfTwoPlusOne(x int) bool {
	return isPowerOfTwo(x - 1)
}

// nextPowerOfTwo returns the next power of 2.
func nextPowerOfTwo(v uint32) uint32 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}

// closestPowerOfTwo returns the closest power of 2.
func closestPowerOfTwo(v uint32) uint32 {
	next := nextPowerOfTwo(v)
	if prev := next / 2; (v - prev) < (next - v) {
		next = prev
	}
	return next
}

// splitIPNet into approximate n counts.
func splitIPNet(ipnet *net.IPNet, n int) ([]*net.IPNet, error) {
	var err error
	subnets := make([]*net.IPNet, 0, n)

	maskBits, _ := ipnet.Mask.Size()
	closestPow2 := int(closestPowerOfTwo(uint32(n)))
	pow2 := int(math.Log2(float64(closestPow2)))

	wantedMaskBits := maskBits + pow2

	currentSubnet, err := currentSubnet(ipnet, wantedMaskBits)
	if err != nil {
		return nil, err
	}
	subnets = append(subnets, currentSubnet)
	nxtSubnet := currentSubnet
	for i := 0; i < closestPow2-1; i++ {
		nxtSubnet, err = nextSubnet(nxtSubnet, wantedMaskBits)
		if err != nil {
			return nil, err
		}
		subnets = append(subnets, nxtSubnet)
	}

	if len(subnets) < n {
		lastSubnet := subnets[len(subnets)-1]
		subnets = subnets[:len(subnets)-1]
		ipnets, err := divideIPNet(lastSubnet)
		if err != nil {
			return nil, err
		}
		subnets = append(subnets, ipnets...)
	}
	return subnets, nil
}

// divideIPNet divides an IPNet into two IPNet structures.
func divideIPNet(ipnet *net.IPNet) ([]*net.IPNet, error) {
	subnets := make([]*net.IPNet, 0, 2)

	maskBits, _ := ipnet.Mask.Size()
	wantedMaskBits := maskBits + 1

	currentSubnet, err := currentSubnet(ipnet, wantedMaskBits)
	if err != nil {
		return nil, err
	}
	subnets = append(subnets, currentSubnet)
	nextSubnet, err := nextSubnet(currentSubnet, wantedMaskBits)
	if err != nil {
		return nil, err
	}
	subnets = append(subnets, nextSubnet)

	return subnets, nil
}

func currentSubnet(network *net.IPNet, prefixLen int) (*net.IPNet, error) {
	currentFirst, _, err := AddressRange(network)
	if err != nil {
		return nil, err
	}
	mask := net.CIDRMask(prefixLen, 8*len(currentFirst))
	return &net.IPNet{IP: currentFirst.Mask(mask), Mask: mask}, nil
}

// nextSubnet returns the next subnet for an ipnet.
func nextSubnet(network *net.IPNet, prefixLen int) (*net.IPNet, error) {
	_, currentLast, err := AddressRange(network)
	if err != nil {
		return nil, err
	}
	mask := net.CIDRMask(prefixLen, 8*len(currentLast))
	currentSubnet := &net.IPNet{IP: currentLast.Mask(mask), Mask: mask}
	_, last, err := AddressRange(currentSubnet)
	if err != nil {
		return nil, err
	}
	last = inc(last)
	next := &net.IPNet{IP: last.Mask(mask), Mask: mask}
	if last.Equal(net.IPv4zero) || last.Equal(net.IPv6zero) {
		return next, nil
	}
	return next, nil
}

// AddressRange returns the first and last addresses in the given CIDR range.
func AddressRange(network *net.IPNet) (firstIP, lastIP net.IP, err error) {
	firstIP = network.IP

	prefixLen, bits := network.Mask.Size()
	if prefixLen == bits {
		lastIP := make([]byte, len(firstIP))
		copy(lastIP, firstIP)
		return firstIP, lastIP, nil
	}

	firstIPInt, bits, err := ipToInteger(firstIP)
	if err != nil {
		return nil, nil, err
	}
	hostLen := uint(bits) - uint(prefixLen)
	lastIPInt := big.NewInt(1)
	lastIPInt.Lsh(lastIPInt, hostLen)
	lastIPInt.Sub(lastIPInt, big.NewInt(1))
	lastIPInt.Or(lastIPInt, firstIPInt)
	lastIP = integerToIP(lastIPInt, bits)
	return
}

// inc increments an IP address to the next IP in the subnet.
func inc(ip net.IP) net.IP {
	incIP := make([]byte, len(ip))
	copy(incIP, ip)
	for j := len(incIP) - 1; j >= 0; j-- {
		incIP[j]++
		if incIP[j] > 0 {
			break
		}
	}
	return incIP
}

// ipToInteger converts an IP address to its integer representation.
// It supports both IPv4 as well as IPv6 addresses.
func ipToInteger(ip net.IP) (*big.Int, int, error) {
	val := &big.Int{}
	val.SetBytes([]byte(ip))

	switch len(ip) {
	case net.IPv4len:
		return val, 32, nil
	case net.IPv6len:
		return val, 128, nil
	default:
		return nil, 0, fmt.Errorf("unsupported address length %d", len(ip))
	}
}

// integerToIP converts an Integer IP address to net.IP format.
func integerToIP(ipInt *big.Int, bits int) net.IP {
	ipBytes := ipInt.Bytes()
	ret := make([]byte, bits/8)
	for i := 1; i <= len(ipBytes); i++ {
		ret[len(ret)-i] = ipBytes[len(ipBytes)-i]
	}
	return net.IP(ret)
}
