package main

func (c *Client) InputHandlerRegister() {
	c.inputHandlers["login"] = c.Login
	c.inputHandlers["add_friend"] = c.AddFriend
	c.inputHandlers["del_friend"] = c.DelFriend
	c.inputHandlers["chat_msg"] = c.SendChatMsg
}
