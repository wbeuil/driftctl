package backend

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const BackendKeyCloud = "tfcloud"
const TerraformCloudAPI = "https://app.terraform.io/api/v2"

type TFCloudAttributes struct {
	HostedStateDownloadUrl string `json:"hosted-state-download-url"`
}

type TFCloudData struct {
	Attributes TFCloudAttributes `json:"attributes"`
}

type TFCloudBody struct {
	Data TFCloudData `json:"data"`
}

func NewCloudReader(workspaceId string, opts *Options) (*HTTPBackend, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/workspaces/%s/current-state-version", TerraformCloudAPI, workspaceId), nil)

	req.Header.Add("Content-Type", "application/vnd.api+json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", opts.TerraformCloudToken))

	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if res.StatusCode < 200 || res.StatusCode >= 400 {
		return nil, errors.Errorf("error requesting Cloud backend state: status code: %d", res.StatusCode)
	}

	if err != nil {
		return nil, err
	}

	bodyBytes, _ := ioutil.ReadAll(res.Body)

	body := TFCloudBody{}
	err = json.Unmarshal(bodyBytes, &body)

	if err != nil {
		return nil, err
	}

	rawURL := body.Data.Attributes.HostedStateDownloadUrl
	logrus.WithFields(logrus.Fields{"hosted-state-download-url": rawURL}).Trace("Cloud backend response")

	opt := Options{}
	return NewHTTPReader(rawURL, &opt)
}
