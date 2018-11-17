package logzio_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const updateServiceUrl string = "%s/v1/alerts/%d"
const updateServiceMethod string = "DELETE"

func buildUpdateAlertRequest(alert CreateAlertType) map[string]interface{} {
	var createAlert = map[string]interface{}{}
	createAlert["title"] = alert.Title
	createAlert["description"] = alert.Description
	if len(alert.Filter) > 0 {
		createAlert["filter"] = alert.Filter
	}
	createAlert["query_string"] = alert.QueryString
	createAlert["operation"] = alert.Operation
	createAlert["severityThresholdTiers"] = alert.SeverityThresholdTiers
	createAlert["searchTimeFrameMinutes"] = alert.SearchTimeFrameMinutes
	createAlert["notificationEmails"] = alert.NotificationEmails
	createAlert["isEnabled"] = alert.IsEnabled
	createAlert["suppressNotificationMinutes"] = alert.SuppressNotificationMinutes
	createAlert["valueAggregationType"] = alert.ValueAggregationType
	createAlert["valueAggregationField"] = alert.ValueAggregationField
	createAlert["groupByAggregationFields"] = alert.GroupByAggregationFields
	createAlert["alertNotificationEndpoints"] = alert.AlertNotificationEndpoints
	return createAlert
}

func buildUpdateApiRequest(apiToken string, alertId int64, jsonObject map[string]interface{}) (*http.Request, error) {

	jsonBytes, err := json.Marshal(jsonObject)
	if err != nil {
		return nil, err
	}

	jsonStr, _ := prettyprint(jsonBytes)
	log.Printf("%s::%s", "buildUpdateApiRequest", jsonStr)

	baseUrl := getLogzioBaseUrl()
	req, err := http.NewRequest(updateServiceMethod, fmt.Sprintf(updateServiceUrl, baseUrl, alertId), bytes.NewBuffer(jsonBytes))
	addHttpHeaders(apiToken, req)

	return req, err
}

func (c *Client) UpdateAlert(alertId int64, alert CreateAlertType) (*AlertType, error) {

	err := validateCreateAlertRequest(alert)
	if err != nil {
		return nil, err
	}

	createAlert := buildUpdateAlertRequest(alert)
	req, _ := buildUpdateApiRequest(c.apiToken, alertId, createAlert)

	var client http.Client
	resp, _ := client.Do(req)

	data, _ := ioutil.ReadAll(resp.Body)
	s, _ := prettyprint(data)

	log.Printf("%s::%s", "UpdateAlert::Response", s)

	if !checkValidStatus(resp, []int{200}) {
		return nil, fmt.Errorf("API call %s failed with status code %d, data: %s", "UpdateAlert", resp.StatusCode, s)
	}

	var target AlertType
	json.Unmarshal(data, &target)

	return &target, nil
}