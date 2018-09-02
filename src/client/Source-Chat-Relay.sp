#pragma semicolon 1

#define DEBUG

#define PLUGIN_AUTHOR "Fishy"
#define PLUGIN_VERSION "0.0.1"

#include <sourcemod>
#include <socket>

#pragma newdecls required

enum FrameType
{
	START, // 1
	OPCODE, // 1
	USERID, // 8
	PAYLOAD, // 256
	CLOSE, // 1
	FRAMECOUNT,
}

char sHostname[64];
char sHost[64] = "127.0.0.1";
char sToken[128];

// Randomly selected port
int iPort = 57452;
int iChannel = 1;
int iBindings[128];
int iTotalBindings;

char sBindings[64];
char pBindings[128][16];

ConVar cHost;
ConVar cPort;
ConVar cToken;
ConVar cChannel;
ConVar cBindings;

Handle hSocket;

public Plugin myinfo = 
{
	name = "Source Chat Relay",
	author = PLUGIN_AUTHOR,
	description = "Prototype Stage",
	version = PLUGIN_VERSION,
	url = "https://keybase.io/RumbleFrog"
};

public void OnPluginStart()
{
	CreateConVar("sm_scr_version", PLUGIN_VERSION, "Source Chat Relay Version", FCVAR_REPLICATED | FCVAR_SPONLY | FCVAR_DONTRECORD | FCVAR_NOTIFY);

	cHost = CreateConVar("scr_host", "127.0.0.1", "Relay Server Host", FCVAR_PROTECTED);

	cPort = CreateConVar("scr_port", "57452", "Relay Server Port", FCVAR_PROTECTED);

	cToken = CreateConVar("scr_token", "", "Relay Server Token", FCVAR_PROTECTED);

	cChannel = CreateConVar("scr_channel", "1", "Channel to send messages on", FCVAR_NONE);

	cBindings = CreateConVar("scr_bindings", "", "Channel(s) to listen for messages on; delimited by comma", FCVAR_NONE);

	AutoExecConfig(true, "Source-Server-Relay");

	hSocket = SocketCreate(SOCKET_TCP, OnSocketError);

	SocketSetOption(hSocket, SocketReuseAddr, 1);
	SocketSetOption(hSocket, SocketKeepAlive, 1);
}

public void OnConfigsExecuted()
{
	GetConVarString(FindConVar("hostname"), sHostname, sizeof sHostname);

	cHost.GetString(sHost, sizeof sHost);
	
	iPort = cPort.IntValue;
	
	cToken.GetString(sToken, sizeof sToken);
	
	iChannel = cChannel.IntValue;
	
	cBindings.GetString(sBindings, sizeof sBindings);

	iTotalBindings = ExplodeString(sBindings, ",", pBindings, sizeof pBindings, sizeof pBindings[]);

	for (int i = 0; i < iTotalBindings; i++)
		iBindings[i] = StringToInt(pBindings[i]);
	
	if (!SocketIsConnected(hSocket))
		ConnectRelay();
}

void ConnectRelay()
{
	if (!SocketIsConnected(hSocket))
		SocketConnect(hSocket, OnSocketConnected, OnSocketReceive, OnSocketDisconnected, sHost, iPort);
	else
		PrintToServer("Socket is already connected?");
}

public Action Timer_Reconnect(Handle timer)
{
	ConnectRelay();
}

void StartReconnectTimer()
{
	SocketDisconnect(hSocket);
	CreateTimer(10.0, Timer_Reconnect);
}

public int OnSocketDisconnected(Handle socket, any arg)
{	
	StartReconnectTimer();
	
	PrintToServer("Socket disconnected");
}

public int OnSocketError(Handle socket, int errorType, int errorNum, any ary)
{
	StartReconnectTimer();
	
	LogError("Socket error %i (errno %i)", errorType, errorNum);
}

public int OnSocketConnected(Handle socket, any arg)
{
	PrintToServer("Successfully Connected");
}

public int OnSocketReceive(Handle socket, const char[] receiveData, int dataSize, any arg)
{
	
}

bool IsListening(int channel)
{
	for (int i = 0; i < iTotalBindings; i++)
		if (iBindings[i] == channel)
			return true;
			
	return false;
}

void PreProcessFrame()
{
	// START - 1 byte
	// OPCODE - 1 byte
	// USERID - 8 bytes
	// PAYLOAD - 256 - bytes
	// TERMINATE - 1 byte
}

void PackFrame(int iFrame, const char payload)
{

}

stock bool Client_IsValid(int iClient, bool bAlive = false)
{
	if (iClient >= 1 &&
	iClient <= MaxClients &&
	IsClientConnected(iClient) &&
	IsClientInGame(iClient) &&
	!IsFakeClient(iClient) &&
	(bAlive == false || IsPlayerAlive(iClient)))
	{
		return true;
	}

	return false;
}
