package main

func (c *Client) MessageHandlerRegister() {
	c.messageHandlers[111] = c.OnLoginRsp
	c.messageHandlers[222] = c.OnAddFriendRsp
	c.messageHandlers[333] = c.OnDelFriendRsp
	c.messageHandlers[444] = c.OnSendChatMsgRsp

}
