package dockers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

// Docker is the docker struct
type Docker struct {
	Container string `json:"Id"`
}

// Payload is a struct
type Payload struct {
	Image string   `json:"Image"`
	Cmd   []string `json:"Cmd"`
}

// RunContainer runs a container
func (d Docker) RunContainer(c echo.Context, image string, cmd string) error {

	containerID, err := d.CreateContainer(image, cmd)
	if err != nil {
		fmt.Println("Error creating the container:", err)
	}

	err = d.StartContainer(containerID)
	if err != nil {
		fmt.Println("Error starting the container:", err)
	}

	err = d.WaitContainer(containerID)
	if err != nil {
		fmt.Println("Error waiting the container:", err)
	}

	output := d.ReadOutput(containerID)

	return c.String(http.StatusOK, "Container Output: "+output)
}

// CreateContainer creates a container and returns its ID
func (d Docker) CreateContainer(image string, cmd string) (string, error) {

	dockerHost := os.Getenv("DOCKER_HOST")
	payload := Payload{
		Image: image,
		Cmd:   []string{"/bin/sh", "-c", cmd},
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error in JSON Marshal.")
		return "", err
	}
	req, err := http.NewRequest("POST", "http://"+dockerHost+"/v1.24/containers/create", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error in POST to create a container:", err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading the body response of POST to create the container:", err)
		return "", err
	}

	err = json.Unmarshal(body, &d)
	if err != nil {
		fmt.Println("Error reading container ID:", err)
		return "", err
	}

	return d.Container, err
}

// StartContainer starts a container and returns its error
func (d Docker) StartContainer(containerID string) error {
	dockerHost := os.Getenv("DOCKER_HOST")
	URL := "http://" + dockerHost + "/v1.24/containers/" + containerID + "/start"
	resp, err := http.Post(URL, "", nil)
	if err != nil {
		fmt.Println("Error in POST to start the container:", err)
	}
	defer resp.Body.Close()
	return err
}

// WaitContainer returns when container finishes executing cmd
func (d Docker) WaitContainer(containerID string) error {
	dockerHost := os.Getenv("DOCKER_HOST")
	URL := "http://" + dockerHost + "/v1.24/containers/" + containerID + "/wait"
	resp, err := http.Post(URL, "", nil)
	if err != nil {
		fmt.Println("Error in GET /wait:", err)
	}
	defer resp.Body.Close()
	return err
}

// ReadOutput returns the command ouput of a given containerID
func (d Docker) ReadOutput(containerID string) string {
	dockerHost := os.Getenv("DOCKER_HOST")
	URL := "http://" + dockerHost + "/v1.24/containers/" + containerID + "/logs?stdout=1"
	resp, err := http.Get(URL)
	if err != nil {
		fmt.Println("Error in GET to get the command output of the container:", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading the body response of GET to read the command output:", err)
	}
	return string(body)
}

// PullImage pulls an image, like docker pull
func (d Docker) PullImage(image string) error {
	dockerHost := os.Getenv("DOCKER_HOST")
	URL := "http://" + dockerHost + "/v1.24/images/create?fromImage=" + image
	resp, err := http.Post(URL, "", nil)
	if err != nil {
		fmt.Println("Error in POST to start the container:", err)
	}
	defer resp.Body.Close()
	return err
}

// ListImages returns the docker images, like docker image ls
func (d Docker) ListImages() string {
	dockerHost := os.Getenv("DOCKER_HOST")
	URL := "http://" + dockerHost + "/v1.24/images/json"
	resp, err := http.Get(URL)
	if err != nil {
		fmt.Println("Error in GET to get the images list:", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading the body response of GET to read the command output:", err)
	}
	return string(body)
}