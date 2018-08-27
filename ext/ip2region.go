package ext

import "github.com/lionsoul2014/ip2region/binding/golang/ip2region"

// Region region
type Region struct {
	IP2Region *ip2region.Ip2Region
}

// NewRegion new ip2region
func NewRegion(path string) (*Region, error) {
	ip2Region, err := ip2region.New(path)
	if err != nil {
		return nil, err
	}
	defer ip2Region.Close()

	return &Region{
		IP2Region: ip2Region,
	}, nil
}

// Query query ip
func (r *Region) Query(ipList []string, mode string) (map[string]ip2region.IpInfo, error) {
	var err error
	var ipinfo = make(map[string]ip2region.IpInfo)
	for _, ip := range ipList {
		switch mode {
		case "memory":
			ipinfo[ip], err = r.IP2Region.MemorySearch(ip)
			if err != nil {
				return nil, err
			}
		case "binary":
			ipinfo[ip], err = r.IP2Region.BinarySearch(ip)
			if err != nil {
				return nil, err
			}
		case "btree":
			ipinfo[ip], err = r.IP2Region.BtreeSearch(ip)
			if err != nil {
				return nil, err
			}
		}
	}

	return ipinfo, nil
}
