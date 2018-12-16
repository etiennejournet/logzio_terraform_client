package logzio_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const updateAlertServiceUrl string = "%s/v1/alerts/%d"
const updateAlertServiceMethod string = http.MethodPut
const updateAlertMethodSuccess int = 200

func buildUpdateAlertRequest(alert CreateAlertType) map[string]interface{} {
	var createAlert = map[string]interface{}{}
	createAlert[fldAlertNotificationEndpoints] = alert.AlertNotificationEndpoints
	createAlert[fldDescription] = alert.Description
	if len(alert.Filter) > 0 {
		createAlert[fldFilter] = alert.Filter
	}
	createAlert[fldGroupByAggregationFields] = alert.GroupByAggregationFields
	createAlert[fldIsEnabled] = alert.IsEnabled
	createAlert[fldQueryString] = alert.QueryString
	createAlert[fldNotificationEmails] = alert.NotificationEmails
	createAlert[fldOperation] = alert.Operation
	createAlert[fldSearchTimeFrameMinutes] = alert.SearchTimeFrameMinutes
	createAlert[fldSeverityThresholdTiers] = alert.SeverityThresholdTiers
	createAlert[fldSuppressNotificationsMinutes] = alert.SuppressNotificationsMinutes
	createAlert[fldTitle] = alert.Title
	createAlert[fldValueAggregationField] = alert.ValueAggregationField
	createAlert[fldValueAggregationType] = alert.ValueAggregationType

	return createAlert
}

func buildUpdateApiRequest(apiToken string, alertId int64, jsonObject map[string]interface{}) (*http.Request, error) {
	jsonBytes, err := json.Marshal(jsonObject)
	if err != nil {
		return nil, err
	}

	logSomething("buildUpdateApiRequest", fmt.Sprintf("%s", jsonBytes))

	baseUrl := getLogzioBaseUrl()
	req, err := http.NewRequest(updateAlertServiceMethod, fmt.Sprintf(updateAlertServiceUrl, baseUrl, alertId), bytes.NewBuffer(jsonBytes))
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
	jsonBytes, _ := ioutil.ReadAll(resp.Body)
	logSomething("UpdateAlert::Response", fmt.Sprintf("%s", jsonBytes))

	if !checkValidStatus(resp, []int{updateAlertMethodSuccess}) {
		return nil, fmt.Errorf("API call %s failed with status code %d, data: %s", "UpdateAlert", resp.StatusCode, jsonBytes)
	}

	var target AlertType
	json.Unmarshal(jsonBytes, &target)

	return &target, nil
}
