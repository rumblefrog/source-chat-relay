#include <sourcemod>
#include <morecolors>
#include <socket>
#include <smlib>
#include <bytebuffer>

#pragma semicolon 1

#define PLUGIN_AUTHOR "Fishy"
#define PLUGIN_VERSION "2.0.0"

#pragma newdecls required

char g_sHostname[64];
char g_sHost[64] = "127.0.0.1";
char g_sToken[64];
char g_sPrefix[8];

// Randomly selected port
int g_iPort = 57452;
int g_iFlag;

bool g_bFlag;

// EngineVersion eVer;

ConVar g_cHost;
ConVar g_cPort;
ConVar g_cPrefix;
ConVar g_cFlag;

Handle g_hSocket;
Handle g_hMessageSendForward;
Handle g_hMessageReceiveForward;

enum MessageType
{
	MessageInvalid = 0,
	MessageAuthenticate,
	MessageAuthenticateResponse,
	MessageChat,
	MessageEvent,
	MessageTypeCount,
}

enum AuthenticateResponse
{
	AuthenticateInvalid = 0,
	AuthenticateSuccess,
	AuthenticateDenied,
}

enum IdenticationType
{
	IdenticationInvalid = 0,
	IdenticationSteam,
	IdenticationDiscord,
}

/**
 * Base message structure
 * 
 * @note The type is declared on every derived message type
 * 
 * @field type - short(4) - The message type (enum MessageType)
 */
methodmap BaseMessage < ByteBuffer
{
	public BaseMessage()
	{
		return view_as<BaseMessage>(CreateByteBuffer());
	}

	property MessageType Type
	{
		public get()
		{
			return view_as<MessageType>(this.ReadInt());
		}
	}

	public void DataCursor()
	{
		// Skip the message type field
		this.Cursor = 4;
	}

	public int ReadDiscardString()
	{
		char cByte;

		for(int i = 0; i < MAX_BUFFER_LENGTH; i++) {
			cByte = this.ReadByte();
			
			if(cByte == '\0') {
				return i + 1;
			}
		}
		
		return MAX_BUFFER_LENGTH;
	}

	public void Dispatch()
	{
		char sDump[MAX_BUFFER_LENGTH];

		this.Dump(sDump, MAX_BUFFER_LENGTH);

		SocketSend(g_hSocket, sDump);

		this.Close();
	}
}

/**
 * Should only sent by clients
 * 
 * @field Hostname - string - The hostname
 * @field Token - string - The authentication token
 */
methodmap AuthenticateMessage < BaseMessage
{
	public int GetToken(char[] sToken, int iSize)
	{
		this.DataCursor();

		return this.ReadString(sToken, iSize);
	}

	public AuthenticateMessage(const char[] sHostname, const char[] sToken)
	{
		BaseMessage m = BaseMessage();

		m.WriteInt(view_as<int>(MessageAuthenticate));
		m.WriteString(sHostname);
		m.WriteString(sToken);

		return view_as<AuthenticateMessage>(m);
	}
}

/**
 * This message is only received from the server
 * 
 * @field Response - short(4) - The state of the authentication request (enum AuthenticateResponse)
 */
methodmap AuthenticateMessageResponse < BaseMessage
{
	property AuthenticateResponse Response
	{
		public get()
		{
			this.DataCursor();

			return view_as<AuthenticateResponse>(this.ReadInt());
		}
	}
}

/**
 * Bi-directional messaging structure
 * 
 * @field EntityName - string - The entity's name that it's sending from
 * @field IDType - short(4) - Type of ID (enum IdenticationType)
 * @field ID - string - The unique identication of the user (SteamID/Discord Snowflake/etc)
 * @field Username - string - The name of the user
 * @field Message - string - The message
 */
methodmap ChatMessage < BaseMessage
{
	public int GetEntityName(char[] sName, int iSize)
	{
		this.DataCursor();

		return this.ReadString(sName, iSize);
	}

	property IdenticationType IDType
	{
		public get()
		{
			this.DataCursor();

			// Skip EntityName
			this.ReadDiscardString();

			return view_as<IdenticationType>(this.ReadInt());
		}
	}

	public int GetUserID(char[] sID, int iSize)
	{
		this.DataCursor();

		// Skip EntityName
		this.ReadDiscardString();

		// Skip ID type
		this.Cursor += 4;

		return this.ReadString(sID, iSize);
	}

	public int GetUsername(char[] sUsername, int iSize)
	{
		this.DataCursor();

		// Skip EntityName
		this.ReadDiscardString();

		// Skip ID type
		this.Cursor += 4;

		// Skip UserID
		this.ReadDiscardString();

		return this.ReadString(sUsername, iSize);
	}

	public int GetMessage(char[] sMessage, int iSize)
	{
		this.DataCursor();

		// Skip EntityName
		this.ReadDiscardString();

		// Skip ID type
		this.Cursor += 4;

		// Skip UserID
		this.ReadDiscardString();

		// Skip Name
		this.ReadDiscardString();

		return this.ReadString(sMessage, iSize);
	}

	public ChatMessage(
		const char[] sEntityName,
		IdenticationType IDType,
		const char[] sUserID,
		const char[] sUsername,
		const char[] sMessage)
	{
		BaseMessage m = BaseMessage();

		m.WriteInt(view_as<int>(MessageChat));
		m.WriteString(sEntityName);
		m.WriteInt(view_as<int>(IDType));
		m.WriteString(sUserID);
		m.WriteString(sUsername);
		m.WriteString(sMessage);

		return view_as<ChatMessage>(m);
	}
}

/**
 * Bi-directional event data
 * 
 * @field Event - string - The name of the event
 * @field Data - string - The data of the event
 */
methodmap EventMessage < BaseMessage
{
	public int GetEvent(char[] sEvent, int iSize)
	{
		this.DataCursor();

		return this.ReadString(sEvent, iSize);
	}

	public int GetData(char[] sData, int iSize)
	{
		this.DataCursor();

		// Skip event string
		this.ReadDiscardString();

		return this.ReadString(sData, iSize);
	}

	public EventMessage(const char[] sEvent, const char[] sData)
	{
		BaseMessage m = BaseMessage();

		m.WriteInt(view_as<int>(MessageEvent));
		m.WriteString(sEvent);
		m.WriteString(sData);

		return view_as<EventMessage>(m);
	}
}

public Plugin myinfo = 
{
	name = "Source Chat Relay",
	author = PLUGIN_AUTHOR,
	description = "Communicate between Discord & In-Game, monitor server without being in-game, control the flow of messages and user base engagement!",
	version = PLUGIN_VERSION,
	url = "https://keybase.io/RumbleFrog"
};

public APLRes AskPluginLoad2(Handle myself, bool late, char[] error, int err_max)
{
	RegPluginLibrary("Source-Chat-Relay");

	CreateNative("SCR_SendMessage", Native_SendMessage);

	return APLRes_Success;
}

public void OnPluginStart()
{
	CreateConVar("rf_scr_version", PLUGIN_VERSION, "Source Chat Relay Version", FCVAR_REPLICATED | FCVAR_SPONLY | FCVAR_DONTRECORD | FCVAR_NOTIFY);

	g_cHost = CreateConVar("rf_scr_host", "127.0.0.1", "Relay Server Host", FCVAR_PROTECTED);

	g_cPort = CreateConVar("rf_scr_port", "57452", "Relay Server Port", FCVAR_PROTECTED);
	
	g_cPrefix = CreateConVar("rf_scr_prefix", "", "Prefix required to send message to Discord. If empty, none is required.", FCVAR_NONE);
	
	g_cFlag = CreateConVar("rf_scr_flag", "", "If prefix is enabled, this admin flag is required to send message using the prefix", FCVAR_PROTECTED);

	AutoExecConfig(true, "Source-Server-Relay");
	
	// eVer = GetEngineVersion();

	g_hSocket = SocketCreate(SOCKET_TCP, OnSocketError);

	SocketSetOption(g_hSocket, SocketReuseAddr, 1);
	SocketSetOption(g_hSocket, SocketKeepAlive, 1);
	
	#if defined DEBUG
	SocketSetOption(g_hSocket, DebugMode, 1);
	#endif

	// ClientIndex, EntityName, ClientID, ClientName, Message
	g_hMessageSendForward = CreateGlobalForward(
		"SCR_OnMessageSend",
		ET_Event,
		Param_Cell,
		Param_String,
		Param_String,
		Param_String,
		Param_String);

	// EntityName, ClientID, ClientName, Message
	g_hMessageReceiveForward = CreateGlobalForward(
		"SCR_OnMessageReceive",
		ET_Event,
		Param_String,
		Param_String,
		Param_String,
		Param_String);

	LoadTranslations("sourcechatrelay.phrases");
}

public void OnConfigsExecuted()
{
	GetConVarString(FindConVar("hostname"), g_sHostname, sizeof g_sHostname);

	g_cHost.GetString(g_sHost, sizeof g_sHost);
	
	g_cPrefix.GetString(g_sPrefix, sizeof g_sPrefix);
	
	g_iPort = g_cPort.IntValue;
	
	char sFlag[8];
	
	g_cFlag.GetString(sFlag, sizeof sFlag);
	
	if (strlen(sFlag) != 0)
	{
		AdminFlag aFlag;
		
		g_bFlag = FindFlagByChar(sFlag[0], aFlag);
		
		g_iFlag = FlagToBit(aFlag);
	}
	
	File tFile;

	char sPath[PLATFORM_MAX_PATH], sIP[64];
	
	Server_GetIPString(sIP, sizeof sIP);
	
	BuildPath(Path_SM, sPath, sizeof sPath, "data/%s_%d.data", sIP, Server_GetPort());
	
	if (FileExists(sPath, false))
	{
		tFile = OpenFile(sPath, "r", false);
		
		tFile.ReadString(g_sToken, sizeof g_sToken, -1);
	} else
	{
		tFile = OpenFile(sPath, "w", false);
	
		GenerateRandomChars(g_sToken, sizeof g_sToken, 64);
	
		tFile.WriteString(g_sToken, true);
	}

	delete tFile;

	if (!SocketIsConnected(g_hSocket))
		ConnectRelay();
}

void ConnectRelay()
{	
	if (!SocketIsConnected(g_hSocket))
		SocketConnect(g_hSocket, OnSocketConnected, OnSocketReceive, OnSocketDisconnected, g_sHost, g_iPort);
	else
		PrintToServer("Socket is already connected?");
}

public Action Timer_Reconnect(Handle timer)
{
	ConnectRelay();
}

void StartReconnectTimer()
{
	if (SocketIsConnected(g_hSocket))
		SocketDisconnect(g_hSocket);
		
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
	AuthenticateMessage(g_sHostname, g_sToken).Dispatch();

	PrintToServer("Successfully Connected");
}

public int OnSocketReceive(Handle socket, const char[] receiveData, int dataSize, any arg)
{
	PrintToServer(receiveData);
	
	HandlePackets(receiveData, dataSize);
}

public void HandlePackets(const char[] sBuffer, int iSize)
{
	BaseMessage base = view_as<BaseMessage>(CreateByteBuffer(true, sBuffer, iSize));

	switch(base.Type)
	{
		case MessageChat:
		{
			ChatMessage m = view_as<ChatMessage>(base);

			Action aResult;

			char sEntity[64], sID[64], sName[MAX_NAME_LENGTH], sMessage[64];

			m.GetEntityName(sEntity, sizeof sEntity);
			m.GetUserID(sID, sizeof sID);
			m.GetUsername(sName, sizeof sName);
			m.GetMessage(sMessage, sizeof sMessage);

			Call_StartForward(g_hMessageReceiveForward);
			Call_PushString(sEntity);
			Call_PushString(sID);
			Call_PushString(sName);
			Call_PushString(sMessage);
			Call_Finish(aResult);

			if (aResult >= Plugin_Handled)
				return;

			PrintToChatAll("%T", "ChatMessage", sEntity, sName, sMessage);
		}
		case MessageEvent:
		{
			EventMessage m = view_as<EventMessage>(base);

			char sEvent[64], sData[64];

			m.GetEvent(sEvent, sizeof sEvent);
			m.GetData(sData, sizeof sData);

			PrintToChatAll("%T", "EventMessage", sEvent, sData);
		}
		case MessageAuthenticateResponse:
		{
			AuthenticateMessageResponse m = view_as<AuthenticateMessageResponse>(base);

			if (m.Response == AuthenticateDenied)
				SetFailState("Server denied our token. Stopping.");

			PrintToServer("Successfully authenticated");
		}
		default:
		{
			// They crazy
		}
	}

	base.Close();
}

public void OnClientSayCommand_Post(int client, const char[] command, const char[] sArgs)
{
	if (!Client_IsValid(client))
		return;
		
	if (!SocketIsConnected(g_hSocket))
		return;
		
	if (StrEqual(g_sPrefix, ""))
		DispatchMessage(client, sArgs);
	else
	{
		if (g_bFlag && !CheckCommandAccess(client, "arandomcommandthatsnotregistered", g_iFlag, true))
			return;

		if (StrContains(sArgs, g_sPrefix) != 0)
			return;
		
		char sBuffer[MAX_MESSAGE_LENGTH];
		
		for (int i = iLen; i < strlen(sArgs); i++)
			Format(sBuffer, sizeof sBuffer, "%s%c", sBuffer, sArgs[i]);
		
		DispatchMessage(client, sBuffer);
	}
}

void DispatchMessage(int iClient, const char[] sMessage)
{
	char sEntity[64], sID[64], sName[MAX_NAME_LENGTH];

	Action aResult;

	strcopy(sEntity, sizeof sEntity, g_sHostname);

	GetClientAuthId(iClient, AuthId_SteamID64, sID, sizeof sID);
	GetClientName(iClient, sName, sizeof sName);

	Call_StartForward(g_hMessageSendForward);
	Call_PushCell(iClient);
	Call_PushString(sEntity);
	Call_PushString(sID);
	Call_PushString(sName);
	Call_PushString(sMessage);
	Call_Finish(aResult);

	if (aResult >= Plugin_Handled)
		return;

	ChatMessage(g_sHostname, IdenticationSteam, sID, sName, sMessage).Dispatch();
}

stock void GenerateRandomChars(char[] buffer, int buffersize, int len)
{
	char charset[] = "adefghijstuv6789!@#$%^klmwxyz01bc2345nopqr&+=";
	
	for (int i = 0; i < len; i++)
		Format(buffer, buffersize, "%s%c", buffer, charset[GetRandomInt(0, sizeof charset)]);
}

public int Native_SendMessage(Handle plugin, int numParams)
{
	char sBuffer[512];
	int iWritten;

	int iClient = GetNativeCell(1);

	FormatNativeString(0, 2, 3, sizeof sBuffer, iWritten, sBuffer);

	DispatchMessage(iClient, sBuffer);
}
