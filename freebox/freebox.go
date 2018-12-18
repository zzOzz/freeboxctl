package freebox

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
)

// DeleteDownloadResponse TODO
type DeleteDownloadResponse struct {
	Success bool `json:"success"`
}

// AddDownloadResponse TODO
type AddDownloadResponse struct {
	Result struct {
		ID int `json:"id"`
	} `json:"result"`
	Success bool `json:"success"`
}

// APIVersion is returned by requesting `GET /api_version`
type APIVersion struct {
	FreeboxID  string `json:"uid,omitempty"`
	DeviceName string `json:"device_name,omitempty"`
	Version    string `json:"api_version,omitempty"`
	BaseURL    string `json:"api_base_url,omitempty"`
	DeviceType string `json:"device_type,omitempty"`
}

// App TODO
type App struct {
	Identifier string `json:"app_id,omitempty"`
	Name       string `json:"app_name,omitempty"`
	Version    string `json:"app_version,omitempty"`
	DeviceName string `json:"device_name,omitempty"`

	token        string `json:-`
	trackID      int
	status       string
	loggedIn     bool
	challenge    string
	password     string
	passwordSalt string
	sessionToken string
	permissions  []string
}

// apiResultDownloadsStats TODO
type apiResultDownloadsStats struct {
	NbRSSItemsUnread      int    `json:"nb_rss_items_unread"`
	NbTasks               int    `json:"nb_tasks"`
	NbTasksActive         int    `json:"nb_tasks_active"`
	NbTasksChecking       int    `json:"nb_tasks_checking"`
	NbTasksDone           int    `json:"nb_tasks_done"`
	NbTasksDownloading    int    `json:"nb_tasks_downloading"`
	NbTasksError          int    `json:"nb_tasks_error"`
	NbTasksExtracting     int    `json:"nb_tasks_extracting"`
	NbTasksQueued         int    `json:"nb_tasks_queue"`
	NbTasksRepairing      int    `json:"nb_tasks_repairing"`
	NbTasksSeeding        int    `json:"nb_tasks_seeding"`
	NbTasksStopped        int    `json:"nb_tasks_stopped"`
	NbTasksStopping       int    `json:"nb_tasks_stopping"`
	RXRate                int    `json:"rx_rate"`
	TXRate                int    `json:"tx_rate"`
	ThrottlingIsScheduled bool   `json:"throttling_is_scheduled"`
	ThrottlingMode        string `json:"throttling_mode"`
	ThrottlingRate        struct {
		RXRate int `json:"rx_rate"`
		TXRate int `json:"tx_rate"`
	} `json:"throttling_rate"`
}

// apiResultDownloadsFile TODO
type apiResultDownloadsFile struct {
	Path     string `json:"path"`
	ID       string `json:"id"`
	TaskID   string `json:"task_id"`
	Filepath string `json:"filepath"`
	Mimetype string `json:"mimetype"`
	Name     string `json:"name"`
	Rx       int    `json:"rx"`
	Status   string `json:"status"`
	Priority string `json:"priority"`
	Error    string `json:"error"`
	Size     int    `json:"size"`
}

// DownloadResult TODO
type DownloadResult struct {
	RxBytes         int    `json:"rx_bytes"`
	TxBytes         int    `json:"tx_bytes"`
	DownloadDir     string `json:"download_dir"`
	ArchivePassword string `json:"archive_password"`
	Eta             int    `json:"eta"`
	Status          string `json:"status"`
	IoPriority      string `json:"io_priority"`
	Type            string `json:"type"`
	PieceLength     int    `json:"piece_length"`
	QueuePos        int    `json:"queue_pos"`
	ID              int    `json:"id"`
	InfoHash        string `json:"info_hash"`
	CreatedTs       int    `json:"created_ts"`
	StopRatio       int    `json:"stop_ratio"`
	TxRate          int    `json:"tx_rate"`
	Name            string `json:"name"`
	TxPct           int    `json:"tx_pct"`
	RxPct           int    `json:"rx_pct"`
	RxRate          int    `json:"rx_rate"`
	Error           string `json:"error"`
	Size            int    `json:"size"`
}

//            "path": "L0Rpc3F1ZSBkdXIvRW5yZWdpc3RyZW1lbnRz",
//            "filecount": 0,
//            "link": false,
//            "modification": 1362005535,
//            "foldercount": 0,
//            "name": "Enregistrements",
//            "index": 1,
//            "mimetype": "inode/directory",
//            "hidden": false,
//            "type": "dir",
//            "size": 4096
type FileResult struct {
	Path         	string	`json:"path"`
	FileCount    	int		`json:"filecount"`
	Link	     	bool 	`json:"link"`
	Modification	int 	`json:"modification"`
	FolderCount     int    	`json:"foldercount"`
	Name			string 	`json:"name"`
	Index			int 	`json:"index"`
	MimeType	    string 	`json:"mimetype"`
	Hidden          bool 	`json:"hidden"`
	Type			string 	`json:"type"`
	Size			int    	`json:"size"`
}
// FilesResult TODO
type apiResponseFile struct {
	Success	bool        	`json:"success"`
	Result	[]FileResult 	`json:"result"`
}

// apiResponseDownload TODO
type apiResponseDownload struct {
	Success bool             `json:"success"`
	Result  []DownloadResult `json:"result"`
}

// apiResponseDownloadsFile TODO
type apiResponseDownloadsFile struct {
	Success bool                     `json:"success"`
	Result  []apiResultDownloadsFile `json:"result"`
}

// apiResponseDownloadsStats TODO
type apiResponseDownloadsStats struct {
	Success bool                    `json:"success"`
	Result  apiResultDownloadsStats `json:"result"`
}

// apiRequestLoginSession TODO
type apiRequestLoginSession struct {
	AppID      string `json:"app_id,omitempty"`
	AppVersion string `json:"app_version,omitempty"`
	Password   string `json:"password,omitempty"`
}

// apiResponseLoginAuthorize TODO
type apiResponseLoginAuthorize struct {
	Success bool `json:"success"`
	Result  struct {
		AppToken string `json:"app_token,omitempty"`
		TrackID  int    `json:"track_id,omitempty"`
	} `json:"result"`
}

// apiResponseLoginSession TODO
type apiResponseLoginSession struct {
	Success bool `json:"success"`
	Result  struct {
		SessionToken string          `json:"session_token,omitempty"`
		Challenge    string          `json:",omitempty"`
		PasswordSalt string          `json:",omitempty"`
		Permissions  map[string]bool `json:",omitempty"`
	} `json:"result"`
}

// apiResponseLoginAuthorizeTrack TODO
type apiResponseLoginAuthorizeTrack struct {
	Success bool `json:"success"`
	Result  struct {
		Status       string `json:"status,omitempty"`
		Challenge    string `json:"challenge,omitempty"`
		PasswordSalt string `json:"password_salt,omitempty"`
	} `json:"result"`
}

// apiResponseLogin TODO
type apiResponseLogin struct {
	Success bool `json:"success"`
	Result  struct {
		LoggedIn     bool   `json:"logged_in,omitempty"`
		Challenge    string `json:"challenge,omitempty"`
		PasswordSalt string `json:"password_salt,omitempty"`
	} `json:"result"`
}

// apiResponseTasks TODO
type apiResponseTasks struct {
	Success bool `json:"success"`
	Result []struct {
		Id     	int   `json:"id,omitempty"`
		From   	string `json:"from,omitempty"`
		To 		string `json:"to,omitempty"`
		Type	string `json:"type,omitempty"`
		State	string `json:"state,omitempty"`
	} `json:"result"`
}


// Client is the Freebox API client
type Client struct {
	URL string

	App        App
	apiVersion APIVersion
	client     *http.Client
}


// apiExtractPayload TODO
type apiExtractPayload struct {
	Src      string `json:"src"`
	Dst string `json:"dst"`
	Password   string `json:"password,omitempty"`
	DeleteArchive   bool `json:"delete_archive,omitempty"`
	OverWrite   bool `json:"overwrite,omitempty"`
}
var instance *Client
var once sync.Once

// GetInstance getting App
func GetInstance() *Client {
	once.Do(func() {
		instance = New()
		err := instance.Connect()
		if err != nil {
			logrus.Fatalf("fbx.Connect(): %v", err)
		}

		err = instance.Authorize()
		if err != nil {
			logrus.Fatalf("fbx.Authorize(): %v", err)
		}

		err = instance.Login()
		if err != nil {
			logrus.Fatalf("fbx.Login(): %v", err)
		}
	})
	return instance
}

// New returns a `Client` object with standard configuration
func New() *Client {
	client := Client{
		URL:    "http://mafreebox.free.fr/",
		client: &http.Client{},
		App: App{
			Identifier: "go-freebox",
			Name:       "Go Freebox",
			Version:    "0.1.0",
			DeviceName: "Golang",
		},
	}
	if os.Getenv("GOFBX_TOKEN") != "" {
		client.App.token = os.Getenv("GOFBX_TOKEN")
	}
	if os.Getenv("GOFBX_URL") != "" {
		client.URL = os.Getenv("GOFBX_URL")
	}
	return &client
}

// APIVersion returns an `ApiVersion` structure field with the configuration fetched during `Connect()`
func (c *Client) APIVersion() *APIVersion {
	return &c.apiVersion
}

// APICode TODO
func (a *APIVersion) APICode() string {
	return "v" + strings.Split(a.Version, ".")[0]
}

// DleteResource performs low-level DELETE request on the Freebox API
func (c *Client) DeleteResource(resource string, authenticated bool) ([]byte, error) {
	var url string
	if authenticated {
		url = fmt.Sprintf("%s%s%s/%s", strings.TrimRight(c.URL, "/"), c.apiVersion.BaseURL, c.apiVersion.APICode(), resource)
	} else {
		url = fmt.Sprintf("%s%s", c.URL, resource)
	}
	logrus.Debugf(">>> DELETE  %q", url)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if authenticated {
		req.Header.Set("X-Fbx-App-Auth", c.App.sessionToken)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	logrus.Debugf("<<< %s", body)

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("Status code: %d", resp.StatusCode)
	}

	return body, err
}

// GetResource performs low-level GET request on the Freebox API
func (c *Client) GetResource(resource string, authenticated bool) ([]byte, error) {
	var url string
	if authenticated {
		url = fmt.Sprintf("%s%s%s/%s", strings.TrimRight(c.URL, "/"), c.apiVersion.BaseURL, c.apiVersion.APICode(), resource)
	} else {
		url = fmt.Sprintf("%s%s", c.URL, resource)
	}
	logrus.Debugf(">>> GET  %q", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if authenticated {
		req.Header.Set("X-Fbx-App-Auth", c.App.sessionToken)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	logrus.Debugf("<<< %s", body)

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("Status code: %d", resp.StatusCode)
	}

	return body, err
}

func (a *APIVersion) authBaseURL() string {
	return strings.TrimLeft(a.BaseURL, "/") + a.APICode() + "/"
}

// PostResource post data and returns body
func (c *Client) PostResource(resource string, data interface{}, authenticated bool) ([]byte, error) {
	var url string
	if authenticated {
		url = fmt.Sprintf("%s%s%s", c.URL, c.apiVersion.authBaseURL(), resource)
	} else {
		url = fmt.Sprintf("%s%s", c.URL, resource)
	}

	payload := new(bytes.Buffer)
	encoder := json.NewEncoder(payload)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}

	payloadString := strings.TrimSpace(fmt.Sprintf("%s", payload))
	logrus.Debugf(">>> POST %s payload=%s", url, payloadString)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if authenticated {
		req.Header.Set("X-Fbx-App-Auth", c.App.sessionToken)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	logrus.Debugf("<<< %s", body)

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("Status code: %d", resp.StatusCode)
	}

	return body, err
}

// Connect tries to contact the Freebox API, and fetches API versions
func (c *Client) Connect() error {
	body, err := c.GetResource("api_version", false)
	if err != nil {
		return err
	}
	logrus.Debugf("API response: %s", string(body))

	err = json.Unmarshal(body, &c.apiVersion)
	if err != nil {
		return err
	}

	logrus.Debugf("API version: FreeboxID=%q DeviceName=%q Version=%q DeviceType=%q", c.apiVersion.FreeboxID, c.apiVersion.DeviceName, c.apiVersion.Version, c.apiVersion.DeviceType)

	return nil
}

// Authorize retrieve token
func (c *Client) Authorize() error {
	if c.App.token != "" {
		logrus.Debugf("App already registered, skiping `Authorize()`")
		return nil
	}
	body, err := c.PostResource(c.apiVersion.authBaseURL()+"login/authorize", c.App, false)
	if err != nil {
		return err
	}

	var response apiResponseLoginAuthorize
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	c.App.token = response.Result.AppToken
	c.App.trackID = response.Result.TrackID
	logrus.Infof("Authorize: Token=%q TrackID=%d", c.App.token, c.App.trackID)

	for {
		body, err := c.GetResource(c.apiVersion.authBaseURL()+fmt.Sprintf("login/authorize/%d", c.App.trackID), false)
		if err != nil {
			return err
		}

		var response apiResponseLoginAuthorizeTrack
		err = json.Unmarshal(body, &response)
		if err != nil {
			return err
		}

		if response.Result.Status == "denied" {
			return fmt.Errorf("App denied")
		}

		if response.Result.Status == "granted" {
			break
		}
		time.Sleep(5 * time.Second)
	}

	return nil
}

// Login freebox
func (c *Client) Login() error {
	if c.App.token == "" {
		return fmt.Errorf("You need to get a token with `Authorize()` first")
	}

	body, err := c.GetResource(c.apiVersion.authBaseURL()+"login", false)
	if err != nil {
		return err
	}

	var response apiResponseLogin
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	c.App.loggedIn = response.Result.LoggedIn
	c.App.challenge = response.Result.Challenge
	c.App.passwordSalt = response.Result.PasswordSalt
	hash := hmac.New(sha1.New, []byte(c.App.token))
	hash.Write([]byte(c.App.challenge))
	c.App.password = fmt.Sprintf("%x", hash.Sum(nil))

	if !c.App.loggedIn {
		data := apiRequestLoginSession{
			AppID:      c.App.Identifier,
			AppVersion: c.App.Version,
			Password:   c.App.password,
		}
		body, err := c.PostResource(c.apiVersion.authBaseURL()+"login/session/", data, false)
		if err != nil {
			return err
		}

		var response apiResponseLoginSession
		err = json.Unmarshal(body, &response)
		if err != nil {
			return err
		}

		c.App.sessionToken = response.Result.SessionToken
	}
	return nil
}

// DownloadsStats retrieve download stats
func (c *Client) DownloadsStats() (*apiResultDownloadsStats, error) {
	body, err := c.GetResource("downloads/stats", true)
	if err != nil {
		return nil, err
	}

	var response apiResponseDownloadsStats
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response.Result, nil
}

// DownloadsStatByID retrieve status for a task
func (c *Client) DownloadsStatByID(taskID int) (*[]apiResultDownloadsFile, error) {
	body, err := c.GetResource("downloads/"+strconv.Itoa(taskID)+"/files", true)
	if err != nil {
		return nil, err
	}

	var response apiResponseDownloadsFile
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response.Result, nil
}

// Downloads retrives downloads
func (c *Client) Downloads() (*[]DownloadResult, error) {
	body, err := c.GetResource("downloads/", true)
	if err != nil {
		return nil, err
	}

	var response apiResponseDownload
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response.Result, nil
}

// Files retrives files
func (c *Client) Files(path string) (*[]FileResult, error) {
	body, err := c.GetResource("fs/ls/" + path, true)
	if err != nil {
		return nil, err
	}

	var response apiResponseFile
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response.Result, nil
}

// Files retrives files
func (c *Client) DownloadFile(path string) ([]byte, error) {
	body, err := c.GetResource("dl/" + path, true)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// AddDownload add download url
func (c *Client) AddDownload(downloadURL string) (interface{}, error) {

	apiURL := c.URL
	resource := fmt.Sprintf("%s%s", c.apiVersion.authBaseURL(), "downloads/add")
	data := url.Values{}
	data.Set("download_url", downloadURL)

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String() // 'https://api.com/user/'

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Set("X-Fbx-App-Auth", c.App.sessionToken)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	// fmt.Println(resp.Status)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	logrus.Debugf("<<< %s", body)

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("Status code: %d", resp.StatusCode)
	}
	// return body, err

	var response interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, err

}

// ExtractFile
func (c *Client) ExtractFile(file string, dest string) (interface{}, error) {
	var data apiExtractPayload
	data.Src = file
	data.Dst = dest
	body, err := c.PostResource("fs/extract/", data,true)
	if err != nil {
		return nil, err
	}

	var response interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, err
}

// FSTasks retrieve download stats
func (c *Client) FSTasks() (*apiResponseTasks, error) {
	body, err := c.GetResource("fs/tasks", true)
	if err != nil {
		return nil, err
	}
	var response apiResponseTasks
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// RemoveFSTasks remove fs task
func (c *Client) RemoveTask(taskID int) (interface{}, error) {

	body, err := c.DeleteResource("fs/tasks/"+strconv.Itoa(taskID), true)
	if err != nil {
		return nil, err
	}

	var response interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// RemoveAllFSTasks remove fs task
func (c *Client) RemoveAllFSTask() (error) {

	tasks, err := c.FSTasks()
	if err != nil {
		return err
	}
	for _, task := range tasks.Result {
		c.RemoveTask(task.Id)
	}

	return nil
}

// RemoveDownload add download url
func (c *Client) RemoveDownload(downloadID int) (interface{}, error) {

	body, err := c.DeleteResource("downloads/"+strconv.Itoa(downloadID), true)
	if err != nil {
		return nil, err
	}

	var response DeleteDownloadResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// RemoveDownload add download url
func (c *Client) CleanDownloads() (error) {
	downloads, err := c.Downloads()
	if err != nil {
		return err
	}
	for _, down := range *downloads {
		c.RemoveDownload(down.ID)
	}

	return nil
}