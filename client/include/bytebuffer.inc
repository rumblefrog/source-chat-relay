#if defined _bytebuffer_included
 #endinput
#endif
#define _bytebuffer_included

#define BYTEBUFFER_VERSION		2

#if !defined MAX_BUFFERS
 #define MAX_BUFFERS 32
#endif

#if !defined MAX_BUFFER_LENGTH
 #define MAX_BUFFER_LENGTH 1024
#endif

enum ByteBuffer {
	INVALID_BUFFER = -1
};

static char g_Buffer[MAX_BUFFERS][MAX_BUFFER_LENGTH];
static bool g_InUse[MAX_BUFFERS];
static bool g_IsLE[MAX_BUFFERS];
static int g_Cursor[MAX_BUFFERS];
static int g_Length[MAX_BUFFERS];

methodmap ByteBuffer {
	property int Cursor {
		public get() {
			CheckBuffer(this);
			return g_Cursor[this];
		}
		
		public set(int value) {
			CheckBuffer(this);
			g_Cursor[this] = value;
		}
	}
	
	property bool LittleEndian {
		public get() {
			CheckBuffer(this);
			return g_IsLE[this];
		}
		
		public set(bool value) {
			CheckBuffer(this);
			g_IsLE[this] = value;
		}
	}
	
	property bool BigEndian {
		public get() {
			CheckBuffer(this);
			return !g_IsLE[this];
		}
		
		public set(bool value) {
			CheckBuffer(this);
			g_IsLE[this] = !value;
		}
	}
	
	public void WriteByte(int data) {
		CheckBuffer(this);
		g_Buffer[this][this.Cursor] = data & 0x000000FF;
		this.Cursor++;
	}
	
	public void WriteShort(int data) {
		CheckBuffer(this);
		if(this.LittleEndian) {
			this.WriteByte(data & 0x000000FF);
			this.WriteByte((data & 0x0000FF00) >> 8);
		} else {
			this.WriteByte((data & 0x0000FF00) >> 8);
			this.WriteByte(data & 0x000000FF);
		}
	}
	
	public void WriteInt(int data) {
		CheckBuffer(this);
		if(this.LittleEndian) {
			this.WriteByte(data & 0x000000FF);
			this.WriteByte((data & 0x0000FF00) >> 8);
			this.WriteByte((data & 0x00FF0000) >> 16);
			this.WriteByte((data & 0xFF000000) >> 24);
		} else {
			this.WriteByte((data & 0xFF000000) >> 24);
			this.WriteByte((data & 0x00FF0000) >> 16);
			this.WriteByte((data & 0x0000FF00) >> 8);
			this.WriteByte(data & 0x000000FF);
		}
	}
	
	public void WriteString(const char[] input) {
		CheckBuffer(this);
		int length = strlen(input);
		for(int i = 0; i < length; i++) {
			this.WriteByte(view_as<int>(input[i]));
		}
		
		this.WriteByte(0);
	}
	
	public int ReadByte() {
		CheckBuffer(this);
		int value = g_Buffer[this][this.Cursor] & 0xFF;
		this.Cursor++;
		return value;
	}
	
	public int ReadShort() {
		CheckBuffer(this);
		int a, b;
		
		if(this.LittleEndian) {
			a = this.ReadByte();
			b = (this.ReadByte() << 8) & 0x0000FF00;
		} else {
			a = (this.ReadByte() << 8) & 0x0000FF00;
			b = this.ReadByte();
		}
		
		return (a | b);
	}
	
	public int ReadInt() {
		CheckBuffer(this);
		int a, b, c, d;
		
		if(this.LittleEndian) {
			a = this.ReadByte();
			b = (this.ReadByte() << 8) & 0x0000FF00;
			c = (this.ReadByte() << 16) & 0x00FF0000;
			d = (this.ReadByte() << 24) & 0xFF000000;
		} else {
			a = (this.ReadByte() << 24) & 0xFF000000;
			b = (this.ReadByte() << 16) & 0x00FF0000;
			c = (this.ReadByte() << 8) & 0x0000FF00;
			d = this.ReadByte();
		}
		
		return (a | b | c | d);
	}
	
	public int ReadString(char[] output, int maxlen) {
		CheckBuffer(this);
		for(int i = 0; i < maxlen; i++) {
			output[i] = this.ReadByte();
			
			if(output[i] == '\0') {
				return i + 1;
			}
		}
		
		return maxlen;
	}
	
	public int Dump(char[] output, int maxlen) {
		CheckBuffer(this);
		// Dump from start to cursor
		if(maxlen > this.Cursor) {
			maxlen = this.Cursor;
		}
		
		int i;
		for(i = 0; i < maxlen; i++) {
			output[i] = g_Buffer[this][i];
		}
		
		return i;
	}
	
	public int Append(ByteBuffer buffer) {
		CheckBuffer(this);
		int length = buffer.Cursor;
		char[] dumped = new char[length];
		buffer.Dump(dumped, length);
		
		for(int i = 0; i < length; i++) {
			this.WriteByte(dumped[i]);
		}
		
		return length;
	}
	
	public ByteBuffer Extract(int length) {
		CheckBuffer(this);
		ByteBuffer output = CreateByteBuffer(this.LittleEndian);
		
		for(int i = 0; i < length; i++) {
			output.WriteByte(this.ReadByte());
		}
		
		output.Cursor = 0;
		return output;
	}
	
	public void Reset() {
		CheckBuffer(this);
		this.Cursor = 0;
	}
	
	public void Close() {
		CheckBuffer(this);
		g_InUse[this] = false;
	}
}

ByteBuffer CreateByteBuffer(bool littleEndian = true, const char data[] = "", int dataSize = 0) {
	int id = -1;
	for(int i = 0; i < MAX_BUFFERS; i++) {
		if(!g_InUse[i]) {
			id = i;
			break;
		}
	}
	
	if(id == -1) {
		ThrowError("Cannot allocate a new buffer: all %d buffers in use", MAX_BUFFERS);
	}
	
	g_InUse[id] = true;
	g_IsLE[id] = littleEndian;
	g_Cursor[id] = 0;
	g_Length[id] = 0;
	
	if(dataSize > 0) {
		// Write starting data
		if(dataSize >= MAX_BUFFER_LENGTH) {
			dataSize = MAX_BUFFER_LENGTH;
		}
		
		for(int i = 0; i < dataSize; i++) {
			g_Buffer[id][i] = data[i];
		}
	}
	
	return view_as<ByteBuffer>(id);
}

static void CheckBuffer(ByteBuffer index) {
	if(!g_InUse[index]) {
		ThrowError("Cannot read/write from closed buffer %d", index);
	}
}
