/**
 * Forward interception test
 */

#include <sourcemod>
#include <Source-Chat-Relay>

public Action SCR_OnMessageSend(int iClient, char[] sClientName, char[] sMessage)
{
    Format(sClientName, MAX_NAME_LENGTH, "%s", "Bob");
    Format(sMessage, MAX_COMMAND_LENGTH, "%s", "This is bob");

    return Plugin_Changed;
}