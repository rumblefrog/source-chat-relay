/**
 * Player death/killfeed event test
 */

#include <sourcemod>
#include <Source-Chat-Relay>

public void OnPluginStart()
{
    HookEvent("player_death", Event_PlayerDeath, EventHookMode_Post);
}

public void Event_PlayerDeath(Event event, const char[] name, bool dontBroadcast)
{
    int iVictim = GetClientOfUserId(event.GetInt("userid"));
    int iAttacker = GetClientOfUserId(event.GetInt("attacker"));

    if (!Client_IsValid(iVictim) || !Client_IsValid(iAttacker))
        return;

    char sVictim[MAX_NAME_LENGTH], sAttacker[MAX_NAME_LENGTH];
	
    GetClientName(iVictim, sVictim, sizeof sVictim);
    GetClientName(iAttacker, sAttacker, sizeof sAttacker);
	
    SCR_SendEvent("Kill Feed", "%s killed %s", sAttacker, sVictim)
}

stock bool Client_IsValid(int client, bool checkConnected=true)
{
	if (client > 4096) {
		client = EntRefToEntIndex(client);
	}

	if (client < 1 || client > MaxClients) {
		return false;
	}

	if (checkConnected && !IsClientConnected(client)) {
		return false;
	}

	return true;
}