#if defined _Source_Chat_Relay_included
 #endinput
#endif
#define _Source_Chat_Relay_included

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
