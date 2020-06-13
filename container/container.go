package container

// Client will hold connection will docker
type Client struct {
}

// PullImage will pull image
func (c *Client) PullImage(image string) error {
	return nil
}

// Create a new container and use net=host(v1), will return error about do not have enough port to deploy
// will pull image if it is not existed
func (c *Client) Create(image string) (string, error) {
	return "", nil
}

// Destroy will destroy container by name
func (c *Client) Destroy(name string) error {
	return nil
}
