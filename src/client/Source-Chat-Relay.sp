#pragma semicolon 1

#define DEBUG

#define PLUGIN_AUTHOR "Fishy"
#define PLUGIN_VERSION "0.0.1"

#include <sourcemod>
#include <socket>

#pragma newdecls required

#define HEADER_LEN 161

enum RelayFrame {
	Ping,
	Authenticate,
	AuthenticateReplyFrame,
	Message,
	FrameCount,
	Unknown
}

char sHostname[64];
//TODO: Change to 127.1
char sHost[64] = "192.168.86.106";

// Randomly selected port
int iPort = 57452;

ConVar cHost;
ConVar cPort;

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

	//TODO: Change to 127.1
	cHost = CreateConVar("scr_host", "192.168.86.106", "Relay Server Host", FCVAR_PROTECTED);

	cPort = CreateConVar("scr_port", "57452", "Relay Server Port", FCVAR_PROTECTED);

	AutoExecConfig(true, "Source-Server-Relay");

	hSocket = SocketCreate(SOCKET_TCP, OnSocketError);

	SocketSetOption(hSocket, SocketReuseAddr, 1);
	SocketSetOption(hSocket, SocketKeepAlive, 1);
	SocketSetOption(hSocket, DebugMode, 1);
}

public void OnConfigsExecuted()
{
	GetConVarString(FindConVar("hostname"), sHostname, sizeof sHostname);

	cHost.GetString(sHost, sizeof sHost);
	
	iPort = cPort.IntValue;
	
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
	if (SocketIsConnected(hSocket))
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
		
	PackMessage(client, sArgs);
}

void PackMessage(int client, const char[] message)
{
	// 64 - Hostname
	// 64 - ClientID (Snowflake / SteamID64)
	// 32 - ClientName
	// remaining
	
	int iMessageLen = strlen(message);
	
	int iFrameLen = HEADER_LEN + iMessageLen;
	
	char[] sFrame = new char[iFrameLen];
	
	char sName[32], sID64[64];
	
	Format(sFrame, iFrameLen, "%s%-64s", sFrame, sHostname);
	
	GetClientAuthId(client, AuthId_SteamID64, sID64, sizeof sID64);
	
	GetClientName(client, sName, sizeof sName);
	
	Format(sFrame, iFrameLen, "%s%-64s", sFrame, sID64);
	
	Format(sFrame, iFrameLen, "%s%-32s", sFrame, sName);
	
	Format(sFrame, iFrameLen, "%s%s", sFrame, message);
	
	PackFrame(Message, sFrame);
}

void PackFrame(RelayFrame opcode, const char[] payload)
{
	int iPayloadLen = strlen(payload);
	int iLen = iPayloadLen + 2;
	
	char[] sFrame = new char[iLen];
	
	switch (opcode)
	{
		case Ping:
		{
			sFrame[0] = '0';
		}
		case Authenticate:
		{
			sFrame[0] = '1';
		}
		case Message:
		{
			sFrame[0] = '4';
				
			Format(sFrame, iLen, "%s%s", sFrame, payload);
		}
	}
	
	SendFrame(sFrame);
	
	#if defined DEBUG
	#endif
}

void SendFrame(const char[] frame)
{
	SocketSend(hSocket, frame);
	
	PrintToConsoleAll(frame);
	
	// Testing if message can be de-constructed successfully
	ParseMessageFrame(frame);
}

void ParseMessageFrame(const char[] frame)
{
	if (frame[0] != '1')
		return;
	
	char hostname[64], id64[64], name[32];
	
	int iLen = strlen(frame);
	
	int iOffset = 2;
	
	for (int i = 0; i < 64; i++)
	{
		Format(hostname, sizeof hostname, "%s%c", hostname, frame[iOffset]);
		iOffset++;
	}
	
	CleanBuffer(hostname, sizeof hostname);
	
	for (int i = 0; i < 64; i++)
	{
		Format(id64, sizeof id64, "%s%c", id64, frame[iOffset]);
		iOffset++;
	}
	
	CleanBuffer(id64, sizeof id64);
	
	for (int i = 0; i < 32; i++)
	{
		Format(name, sizeof name, "%s%c", name, frame[iOffset]);
		iOffset++;
	}
	
	CleanBuffer(name, sizeof name);
	
	int iContentLen = iLen - iOffset;
	
	char[] sContent = new char[iContentLen];
	
	for (int i = 0; i < iContentLen; i++)
	{
		Format(sContent, iContentLen + 1, "%s%c", sContent, frame[iOffset]);
		iOffset++;
	}
	
	PrintToConsoleAll("===== PARSING =====");
	
	PrintToConsoleAll("hostname: %s", hostname);
	
	PrintToConsoleAll("id64: %s", id64);
	
	PrintToConsoleAll("name: %s", name);
	
	PrintToConsoleAll("sContentLen: %d", iContentLen);
	
	PrintToConsoleAll("sContent: %s", sContent);
	
	PrintToConsoleAll("===================");
}

void CleanBuffer(char[] buffer, int bufferlen)
{
	int iLen = strlen(buffer);

	int iOffset = iLen;
	
	for (int i = iLen; i > 0; i--)
	{
		if (buffer[i] != ' ')
		{
			iOffset = i;
			break;
		}
	}
	
	int iEnd = iOffset + 1;
	
	char[] sBuffer = new char[iEnd];
	
	for (int i = 0; i < iOffset; i++)
		Format(sBuffer, iEnd, "%s%c", sBuffer, buffer[i]);
	
	Format(buffer, bufferlen, "%s", sBuffer);
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
