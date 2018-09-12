#pragma semicolon 1

#define DEBUG

#define PLUGIN_AUTHOR "Fishy"
#define PLUGIN_VERSION "0.0.1"

#include <sourcemod>
#include <socket>

#pragma newdecls required

enum RelayFrame {
	Ping = 2,
	Message = 6
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
	
	// if (!SocketIsConnected(hSocket))
	// 	ConnectRelay();
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
		
	// if (!SocketIsConnected(hSocket))
	// 	return;
		
	PackMessage(client, sArgs);
}

bool IsListening(int channel)
{
	for (int i = 0; i < iTotalBindings; i++)
		if (iBindings[i] == channel)
			return true;
			
	return false;
}

void PackBuffer(char[] buffer, int bufferlen, int fieldlen, const char[] content)
{	
	int iPadLen = fieldlen - strlen(content);
	
	// strcopy(buffer, bufferlen, content);
	
	Format(buffer, bufferlen, "%s%s", buffer, content);
	
	for (int i = 0; i < iPadLen; i++)
		Format(buffer, bufferlen, "%s ", buffer);
}

// sm plugins reload Source-Chat-Relay
void PackMessage(int client, const char[] message)
{
	// 15
	// 64
	// 32
	// remaining
	
	const int HeaderLen = 112;
	
	int iMessageLen = strlen(message);
	
	int iFrameLen = HeaderLen + iMessageLen;
	
	char[] sFrame = new char[iFrameLen];
	
	char sIP[15];
	
	PackServerIP(sIP, sizeof sIP);
	
	PackBuffer(sFrame, iFrameLen, 15, sIP);

	PackBuffer(sFrame, iFrameLen, 64, sHostname);
	
	char sName[32];
	
	GetClientName(client, sName, sizeof sName);
	
	PackBuffer(sFrame, iFrameLen, 32, sName);
	
	PackBuffer(sFrame, iFrameLen, iMessageLen, message);
	
	PackFrame(Message, sFrame);
}

void PackFrame(RelayFrame opcode, const char[] payload)
{
	int iPayloadLen = strlen(payload);
	int iLen = iPayloadLen + view_as<int>(opcode);
	
	char[] sFrame = new char[iLen];
	
	switch (opcode)
	{
		case Ping:
		{
			sFrame[0] = '0';
		}
		case Message:
		{
			sFrame[0] = '1';
			
			if (!PackLength(iPayloadLen, sFrame, iLen))
				return;
				
			Format(sFrame, iLen, "%s%s", sFrame, payload);
		}
	}
	
	SendFrame(sFrame);
	
	#if defined DEBUG
	#endif
}

bool PackLength(int len, char[] buffer, int buffersize)
{
	if (len > 9999)
	{
		LogError("Payload length exceeds maximum length");
		return false;
	}
	
	char sLen[4];
	
	IntToString(len, sLen, sizeof sLen);
	
	int iPadLen = 4 - strlen(sLen);
	
	for (int i = 0; i < iPadLen; i++)
		Format(buffer, buffersize, "%s%s", buffer, "0");
	
	Format(buffer, buffersize, "%s%s", buffer, sLen);
	
	return true;
}

void SendFrame(const char[] frame)
{
	PrintToConsoleAll(frame);
	
	// Testing if message can be de-constructed successfully
	ParseMessageFrame(frame);
}

void PackServerIP(char[] IP, int size)
{
	int pieces[4];
	int longip = GetConVarInt(FindConVar("hostip"));
	
	pieces[0] = (longip >> 24) & 0x000000FF;
	pieces[1] = (longip >> 16) & 0x000000FF;
	pieces[2] = (longip >> 8) & 0x000000FF;
	pieces[3] = longip & 0x000000FF
	
	Format(IP, size, "%d.%d.%d.%d", pieces[0], pieces[1], pieces[2], pieces[3]);
}

void ParseMessageFrame(const char[] frame)
{
	if (frame[0] != '1')
		return;
	
	char ip[15], hostname[64], name[32], len[4];
	
	Format(len, sizeof len, "%c%c%c%c", frame[1], frame[2], frame[3], frame[4]);
	
	int iLen = strlen(frame);
	
	int iOffset = 5;
	
	for (int i = 0; i < 15; i++)
	{
		Format(ip, sizeof ip, "%s%c", ip, frame[iOffset]);
		iOffset++;
	}
	
	CleanBuffer(ip, sizeof ip);
	
	for (int i = 0; i < 64; i++)
	{
		Format(hostname, sizeof hostname, "%s%c", hostname, frame[iOffset]);
		iOffset++;
	}
	
	CleanBuffer(hostname, sizeof hostname);
	
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
	
	PrintToConsoleAll("ip: %s", ip);
	
	PrintToConsoleAll("hostname: %s", hostname);
	
	PrintToConsoleAll("name: %s", name);
	
	PrintToConsoleAll("sContentLen: %d", iContentLen);
	
	PrintToConsoleAll("sContent: %s", sContent);
	
	PrintToConsoleAll("===================");
}

// sm plugins reload Source-Chat-Relay
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
