/**
 * Event as a service test
 */

#include <sourcemod>
#include <Source-Chat-Relay>

public Action SCR_OnEventReceive(char[] sEvent, char[] sData)
{
    if (!StrEqual(sEvent, "RequestPlayerCount"))
        return Plugin_Continue;

    char sCount[4];

    IntToString(GetClientCount(true), sCount, sizeof sCount);

    SCR_SendEvent("ResponsePlayerCount", sCount);

    return Plugin_Continue;
}