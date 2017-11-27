package managers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/bigtable"
)

const (
	beaconsTableName      = "beacons"
	systemFamilyName      = "system"
	applicationFamilyName = "application"
	locationFamilyName    = "location"
	requestFamilyName     = "request"
	headerFamilyName      = "header"
	cookieFamilyName      = "cookie"
	dataFamilyName        = "data"

	maxResponseSize = 1000000
)

var projectID string
var instanceID string

func init() {
	projectID = os.Getenv("BIGTABLE_PROJECT_ID")
	instanceID = os.Getenv("BIGTABLE_INSTANCE_ID")
}

type BeaconManager interface {
	Init(context.Context) error
	Save(context.Context, string, *http.Request, []byte) error
}

type beaconManager struct{}

func (beaconManager) Init(ctx context.Context) error {
	adminClient, err := bigtable.NewAdminClient(ctx, projectID, instanceID)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	return createTables(ctx, adminClient, beaconsTableName, systemFamilyName, applicationFamilyName, locationFamilyName, requestFamilyName, headerFamilyName, cookieFamilyName, dataFamilyName)
}

func (beaconManager) Save(ctx context.Context, fileType string, r *http.Request, body []byte) error {
	// hostname -- bestbuy.com
	hostName := strings.TrimSpace(strings.ToLower(r.FormValue("_h")))

	// must have an event time
	eventTimeStr := r.FormValue("_t")
	eventTime, err := strconv.ParseInt(eventTimeStr, 10, 64)

	// check the skew... if > skew then this time is outside our range
	now := time.Now().UnixNano()
	skew := now - (eventTime * 1000000)
	now = now - skew

	// create a client
	client, err := bigtable.NewClient(ctx, projectID, instanceID)
	if err != nil {
		return err
	}
	defer client.Close()

	// row key for insert
	rowKey := fmt.Sprintf("%s#%d", hostName, now)

	// type of event captured
	event := r.FormValue("_e")
	if event == "" {
		event = "ping"
	}

	// the rest are free form fields. ignore anything starting with an "_"
	args := make(map[string]string)
	for key := range r.Form {
		switch key {
		case "_t", "_e", "_h", "_n":
			continue
		default:
			args[key] = r.FormValue(key)
		}
	}

	tbl := client.Open(beaconsTableName)

	mut := bigtable.NewMutation()
	ts := bigtable.Now()
	// system event details
	mut.Set(systemFamilyName, "time", ts, []byte(eventTimeStr))
	mut.Set(systemFamilyName, "event", ts, []byte(event))
	mut.Set(systemFamilyName, "type", ts, []byte(fileType))

	// application
	for key, value := range args {
		mut.Set(applicationFamilyName, key, ts, []byte(value))
	}

	if r.Header.Get("X-Appengine-Citylatlong") != "" {
		mut.Set(locationFamilyName, "city", ts, []byte(r.Header.Get("X-Appengine-City")))
		mut.Set(locationFamilyName, "region", ts, []byte(r.Header.Get("X-Appengine-Region")))
		mut.Set(locationFamilyName, "country", ts, []byte(r.Header.Get("X-Appengine-Country")))
		mut.Set(locationFamilyName, "latlong", ts, []byte(r.Header.Get("X-Appengine-Citylatlong")))
	}

	// request
	mut.Set(requestFamilyName, "url", ts, []byte(r.URL.String()))
	mut.Set(requestFamilyName, "host", ts, []byte(r.Host))
	mut.Set(requestFamilyName, "proto", ts, []byte(r.Proto))
	mut.Set(requestFamilyName, "method", ts, []byte(r.Method))
	mut.Set(requestFamilyName, "address", ts, []byte(r.RemoteAddr))

	// headers
	for key := range r.Header {
		key = strings.ToLower(key)

		if strings.HasPrefix(key, "x-") {
			continue
		}

		if strings.EqualFold(key, "cookie") {
			continue
		}
		mut.Set(headerFamilyName, key, ts, []byte(r.Header.Get(key)))
	}

	// cookie
	for _, cookie := range r.Cookies() {
		mut.Set(cookieFamilyName, cookie.Name, ts, []byte(cookie.Value))
	}

	//
	if r.Method == http.MethodPost && strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		var obj interface{}
		err = json.NewDecoder(bytes.NewReader(body)).Decode(&obj)
		if err != nil {
			return err
		}
		mut.Set(dataFamilyName, "json", ts, body)
	}

	// write row
	return tbl.Apply(ctx, rowKey, mut)
}

func NewBeaconManager() BeaconManager {
	return beaconManager{}
}

func createTables(ctx context.Context, adminClient *bigtable.AdminClient, tableName string, familyName ...string) error {
	tables, err := adminClient.Tables(ctx)
	if err != nil {
		return err
	}

	if !sliceContains(tables, tableName) {
		if err := adminClient.CreateTable(ctx, tableName); err != nil {
			return err
		}
	}

	tblInfo, err := adminClient.TableInfo(ctx, tableName)
	if err != nil {
		return err
	}
	for _, name := range familyName {
		if !sliceContains(tblInfo.Families, name) {
			if err := adminClient.CreateColumnFamily(ctx, tableName, name); err != nil {
				return err
			}
		}
	}

	return nil
}

func sliceContains(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

func getHostname(host string) string {
	i := strings.Index(host, ":")
	if i != -1 {
		return host[:i]
	}
	return host
}
