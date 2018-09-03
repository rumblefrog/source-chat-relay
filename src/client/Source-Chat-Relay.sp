#pragma semicolon 1

#define DEBUG

#define PLUGIN_AUTHOR "Fishy"
#define PLUGIN_VERSION "0.0.1"

#define MAX_PAYLOAD_LEN 999

#include <sourcemod>
#include <socket>

#pragma newdecls required

enum PayloadType
{
	Ping,
	Message,
	Terminate
}

enum RelayFrame
{
	PayloadType:OPCODE,
	PAYLOADLEN,
	TERMINATE,
	FRAMECOUNT
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

public void OnClientSayCommand_Post(int client, const char[] command, const char[] sArgs)
{
	if (!Client_IsValid(client))
		return;
		
	if (!SocketIsConnected(hSocket))
		return;
		
	ProcessFrame(Message, "Test");
}

bool IsListening(int channel)
{
	for (int i = 0; i < iTotalBindings; i++)
		if (iBindings[i] == channel)
			return true;
			
	return false;
}

void ProcessFrame(PayloadType payloadT, const char[] payload)
{
	int vFrame[RelayFrame];

	vFrame[OPCODE] = payloadT;

	vFrame[PAYLOADLEN] = strlen(payload);

	vFrame[TERMINATE] = 1;

	PackFrame(vFrame, payload);
}

bool PackFrame(int vFrame[RelayFrame], const char[] payload)
{
	// OPCODE - 1 byte
	// PAYLOADLEN - 3 bytes
	// PAYLOAD <- INJECTED <- X bytes
	// TERMINATE - 2 bytes

	int iLen = vFrame[PAYLOADLEN] + 6;

	char[] sFrame = new char[iLen];

	switch (vFrame[OPCODE])
	{
		case Ping:
		{
			sFrame[0] = 1;
		}
		case Message:
		{
			sFrame[0] = 2;
		}
		default:
		{
			LogError("Invalid OPCODE %d", vFrame[OPCODE]);
			return false;
		}
	}

	if (vFrame[PAYLOADLEN] > MAX_PAYLOAD_LEN)
	{
		LogError("Payload length exceeds 3 bytes");
		return false;
	}

	char sLen[3];

	IntToString(vFrame[PAYLOADLEN], sLen, 3);

	strcopy(sFrame[5], 3, sLen);

	strcopy(sFrame[8], iLen - 5, payload);

	Format(sFrame[iLen - 2], 2, "\0");

	return true;
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
