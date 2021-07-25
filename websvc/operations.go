package websvc

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cicovic-andrija/dante/atlas"
)

func (s *server) operationRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		s.log.info("signaling shutdown")
		s.signalShutdown()
	}
}

func (s *server) getCredits() string {
	// FIXME: handle error
	req, _ := http.NewRequest(http.MethodGet, "https://atlas.ripe.net:443/api/v2/credits/", nil)
	req.Header.Set("Authorization", "Key "+cfg.Atlas.Auth.Key)
	res, err := s.httpClient.Do(req)
	if err != nil {
		return "failed"
	}
	c := atlas.Credit{}
	err = json.NewDecoder(res.Body).Decode(&c)
	if err != nil {
		return "failed"
	}

	return fmt.Sprintf("remaining credits: %d", c.CurrentBalance)

	//ck := influxdb2.NewClient("http://localhost:8086", "B0LXp-bCs6-22sn3soDWsVapRJ5ofwwXNNhVtjSFTVTOmzJcPuTdIa3wv3eIERqttKlS3PWndtfLDERq5jvUMQ==")
	//writeApi := ck.WriteAPIBlocking("dante", "dante-bucket")
	//p := influxdb2.NewPoint("credit", nil, map[string]interface{}{"value": c.CurrentBalance}, time.Now())
	//writeApi.WritePoint(context.Background(), p)
	//ck.Close()
}
