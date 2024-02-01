package apiclient

import (
	"context"
	"encoding/base64"
	"encoding/json"
	. "github.com/eatmoreapple/wxhelper/internal/models"
	"io"
)

type Client struct {
	transport *Transport
}

func (c *Client) GetUserInfo(ctx context.Context) (*Account, error) {
	resp, err := c.transport.GetUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	var r Result[*Account]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	if err = r.Err(); err != nil {
		return nil, err
	}
	return r.Data, nil
}

func (c *Client) CheckLogin(ctx context.Context) (bool, error) {
	resp, err := c.transport.CheckLogin(ctx)
	if err != nil {
		return false, err
	}
	defer func() { _ = resp.Body.Close() }()
	var r Result[bool]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return false, err
	}
	if err = r.Err(); err != nil {
		return false, err
	}
	return r.Data, nil
}

func (c *Client) GetContactList(ctx context.Context) (Members, error) {
	resp, err := c.transport.GetContactList(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	var r Result[Members]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	if err = r.Err(); err != nil {
		return nil, err
	}
	return r.Data, nil
}

func (c *Client) SendText(ctx context.Context, to, content string) error {
	resp, err := c.transport.SendText(ctx, to, content)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	var r Result[any]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}
	return r.Err()
}

func (c *Client) SendImage(ctx context.Context, to string, img io.Reader) error {
	// to base64
	data, err := io.ReadAll(img)
	if err != nil {
		return err
	}
	encoded := base64.StdEncoding.EncodeToString(data)

	resp, err := c.transport.SendImage(ctx, to, encoded)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	var r Result[any]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}
	return r.Err()
}

func (c *Client) SyncMessage(ctx context.Context) ([]*Message, error) {
	resp, err := c.transport.SyncMessage(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	var r Result[[]*Message]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	return r.Data, nil
}
