package wxclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eatmoreapple/env"
	. "github.com/eatmoreapple/wxhelper/internal/models"
	"io"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	transport *Transport
}

func (c *Client) CheckLogin(ctx context.Context) (bool, error) {
	resp, err := c.transport.CheckLogin(ctx)
	if err != nil {
		return false, err
	}
	defer func() { _ = resp.Body.Close() }()
	var r result[any]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return false, err
	}
	return r.Code == 1, nil
}

func (c *Client) GetUserInfo(ctx context.Context) (*Account, error) {
	resp, err := c.transport.GetUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	var r result[*Account]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	if r.Code != 1 {
		return nil, errors.New("get user info failed")
	}
	return r.Data, nil
}

func (c *Client) SendText(ctx context.Context, to string, content string) error {
	resp, err := c.transport.SendText(ctx, to, content)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	var r result[any]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}
	if r.Code == 0 {
		return errors.New("send text failed")
	}
	return nil
}

func (c *Client) GetContactList(ctx context.Context) (Members, error) {
	resp, err := c.transport.GetContactList(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	var r result[Members]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	return r.Data, nil
}

type HookSyncMsgOption struct {
	LocalURL *url.URL
	Timeout  time.Duration
}

func (c *Client) HTTPHookSyncMsg(ctx context.Context, o HookSyncMsgOption) error {
	opt := TransportHookSyncMsgOption{
		Url:        o.LocalURL.String(),
		EnableHttp: 1,
		Timeout:    strconv.Itoa(int(o.Timeout.Seconds()) * 100),
		Ip:         o.LocalURL.Hostname(),
		Port:       o.LocalURL.Port(),
	}
	resp, err := c.transport.HookSyncMsg(ctx, opt)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	var r result[any]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {

		return err
	}
	if r.Code != 0 {
		return errors.New("hook sync msg failed")
	}
	return nil
}

func (c *Client) UnhookSyncMsg(ctx context.Context) error {
	resp, err := c.transport.UnhookSyncMsg(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	var r result[any]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}
	if r.Code != 0 {
		return errors.New("unhook sync msg failed")
	}
	return nil
}

func (c *Client) SendImage(ctx context.Context, to string, img io.Reader) error {
	file, cb, err := readerToFile(img)
	if err != nil {
		return err
	}
	defer cb()
	// 转换成windows下c盘的路径
	filepath := fmt.Sprintf("C:\\%s", file.Name())
	resp, err := c.transport.SendImage(ctx, to, filepath)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	var r result[any]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}
	return nil
}

func (c *Client) SendFile(ctx context.Context, to string, img io.Reader) error {
	file, cb, err := readerToFile(img)
	if err != nil {
		return err
	}
	defer cb()
	// 转换成windows下c盘的路径
	filepath := fmt.Sprintf("C:\\%s", file.Name())
	resp, err := c.transport.SendFile(ctx, to, filepath)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	var r result[any]
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}
	if r.Code == 0 {
		return fmt.Errorf("send file failed with code %d", r.Code)
	}
	return nil
}

func New(transport *Transport) *Client {
	return &Client{transport: transport}
}

func Default() *Client {
	transport := NewTransport(env.Name("VIRTUAL_MACHINE_URL").String())
	return New(transport)
}
