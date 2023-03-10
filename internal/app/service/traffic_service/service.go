package traffic_service

import (
	"log"
	"time"
	"where-do-i-sit/internal/app"
	serror "where-do-i-sit/internal/app/error"
	"where-do-i-sit/internal/app/sk_api"
	"where-do-i-sit/pkg/cache"
)

type TrafficService struct {
	cache cache.Cache
}

func New(cache cache.Cache) TrafficService {
	return TrafficService{
		cache,
	}
}

func (t TrafficService) GetStationList() ([]app.Station, error) {
	return sk_api.GetStationList()
}

func (t TrafficService) GetStationByName(s string, line string) (station app.Station, err error) {
	var stations []app.Station
	res, exists := t.cache.Get("stationList")
	if !exists {
		stations, err = t.GetStationList()
		t.cache.Set("stationList", stations, 24*time.Hour)
		if err != nil {
			return
		}
	} else {
		var ok bool
		stations, ok = res.([]app.Station)
		if !ok {
			log.Println("assertion failed on stationList")
			err = serror.ErrInternal
			return
		}
	}

	for _, st := range stations {
		if st.Name == s && st.Line == line {
			station = st
			return
		}
	}

	err = serror.ErrNoSuchStation
	return
}

func (t TrafficService) GetStatisticCongestion(stationCode, prevStationCode string, time time.Time) (ret []app.Congestion, err error) {
	ret, err = sk_api.GetStatisticCongestion(stationCode, prevStationCode, time)
	return
}

func (t TrafficService) GetRealtimeCongestion(stationCode, prevStationCode string) (app.Congestion, error) {
	//TODO implement me
	panic("implement me")
}
