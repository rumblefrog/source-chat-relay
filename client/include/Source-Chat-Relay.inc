#if defined _Source_Chat_Relay_included
 #endinput
#endif
#define _Source_Chat_Relay_included


/**
 * Called upon sending a message via chat triggers
 * 
 * @param iClient 		Index of client that is sending the message
 * @param sEntityName	Hostname of the server
 * @param sClientID 	SteamID64 of the client sending the message
 * @param sClientName	Name of the client sending the message
 * @param sMessage		Message content
 */
forward Action SCR_OnMessageSend(int iClient, char[] sEntityName, char[] sClientID, char[] sClientName, char[] sMessage);

/**
 * Called upon receiving a message
 * 
 * @param sEntityName	Entity name/hostname of the sending client
 * @param sClientName	Name of the client that sent the message
 * @param sMessage		Message content
 */
forward Action SCR_OnMessageReceive(char[] sEntityName, char[] sClientName, char[] sMessage);

/**
 * Sends a message to the router
 * 
 * @param iClient   Client ID to use as name display (If left at 0, it will display "CONSOLE")
 * @param fmt       Format string
 * @param ...       Format arguments
 */
native void SCR_SendMessage(int iClient = 0, const char[] fmt, any ...);

#if !defined REQUIRE_PLUGIN
public __pl_Source_Chat_Relay_SetNTVOptional()
{
	MarkNativeAsOptional("SCR_SendMessage");
}
#endif

public SharedPlugin __pl_Source_Chat_Relay =
{
	name = "Source-Chat-Relay",
	file = "Source-Chat-Relay.smx",
	#if defined REQUIRE_PLUGIN
	required = 1,
	#else
	required = 0,
	#endif
};
