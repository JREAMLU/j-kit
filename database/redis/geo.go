package redis

import (
	"errors"

	"github.com/JREAMLU/j-kit/constant"
)

// Geo redis geo
type Geo struct {
	Structure
}

const (
	// GEOADD geoadd
	GEOADD = "GEOADD"
	// GEOPOS geopos
	GEOPOS = "GEOPOS"
	// GEORADIUS georadius
	GEORADIUS = "GEORADIUS"
)

// NewGeo new geo
func NewGeo(instanceName, keyPrefixFmt string) Geo {
	return Geo{
		Structure: NewStructure(instanceName, keyPrefixFmt),
	}
}

// Add add
func (g *Geo) Add(keySuffix string, longitude, latitude, member interface{}) (int64, error) {
	return g.Int64(MASTER, GEOADD, g.InitKey(keySuffix), longitude, latitude, member)
}

// Adds add multi
func (g *Geo) Adds(keySuffix string, longitudes, latitudes, members []interface{}) (int64, error) {
	if len(longitudes) != len(latitudes) || len(longitudes) != len(members) {
		return constant.ZeroInt64, errors.New("param error: slice len must be equal")
	}

	params := make([]interface{}, len(longitudes)*3+1)
	params[0] = g.InitKey(keySuffix)
	n := 1
	for i := range members {
		params[n] = longitudes[i]
		n++

		params[n] = latitudes[i]
		n++

		params[n] = members[i]
		n++
	}

	return g.Int64(MASTER, GEOADD, params...)
}

// Pos longitude and latitude (x,y)
func (g *Geo) Pos(keySuffix string, member interface{}) (longitude float64, latitude float64, err error) {
	results, err := g.Float64Slice(SLAVE, GEOPOS, g.InitKey(keySuffix), member)
	if err != nil {
		return longitude, latitude, err
	}

	if len(results) <= 0 {
		return longitude, latitude, errors.New("Float64Slice reply is nil")
	}

	item := results[0]
	if len(item) > 0 {
		longitude = item[0]
	}

	if len(item) > 1 {
		latitude = item[1]
	}

	return longitude, latitude, err
}

// MultiPos MultiPos
func (g *Geo) MultiPos(keySuffix string, members ...interface{}) ([][]float64, error) {
	params := make([]interface{}, len(members)+1)
	params[0] = g.InitKey(keySuffix)

	n := 1
	for i := range members {
		params[n] = members[i]
		n++
	}

	return g.Float64Slice(SLAVE, GEOPOS, params...)
}

// Radius GEORADIUS
func (g *Geo) Radius(keySuffix string, longitude, latitude, radius interface{}, distanceUnit, orderBy string, count interface{}) ([][]string, error) {
	params := []interface{}{
		g.InitKey(keySuffix),
		longitude,
		latitude,
		radius,
		distanceUnit,
		"WITHDIST",
		orderBy,
		"COUNT",
		count,
	}

	var results [][]string

	replies, err := g.MultiBulk(MASTER, GEORADIUS, params...)
	if err != nil {
		return nil, err
	}

	for _, reply := range replies {
		switch reply := reply.(type) {
		case []interface{}:
			result := make([]string, len(reply))
			for i := range reply {
				if reply[i] == nil {
					continue
				}

				if p, ok := reply[i].(string); ok {
					result[i] = p
				}
			}

			results = append(results, result)
		default:
			return nil, errors.New("wrong type, not []interface{}")
		}
	}

	return results, nil
}
