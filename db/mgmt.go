package db

import (
	"context"
	"strings"
)

const (
	notFoundString = "not found"
)

func EnsureOrganization(client *Client) (created bool, err error) {
	orgAPI := client.influxClient.OrganizationsAPI()
	lookupContext, cancelLookup := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancelLookup()
	_, err = orgAPI.FindOrganizationByName(lookupContext, client.Organization)
	if err != nil && strings.Contains(err.Error(), notFoundString) {
		createContext, cancelCreate := context.WithTimeout(context.Background(), DefaultTimeout)
		defer cancelCreate()
		_, err = orgAPI.CreateOrganizationWithName(createContext, client.Organization)
		return err == nil, err
	} else if err != nil {
		return false, err
	}
	return false, nil
}
