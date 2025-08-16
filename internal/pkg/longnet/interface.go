package longnet

import (
	"context"
	"time"
)

type IAuthorize func(ctx context.Context, token string) (int64, error)

type IClient interface {
	Connect(token string) error
	Write(data []byte) error
	Close()
	SetOnMessage(fn func(client IClient, data []byte))
}

type IConn interface {
	Network() string                                      // Network 网络协议类型
	Read() ([]byte, error)                                // Read 数据读取
	Write([]byte) error                                   // Write 数据写入
	Close() error                                         // Close 连接关闭
	SetCloseHandler(fn func(code int, text string) error) // SetCloseHandler 设置连接关闭回调事件
	SetReadDeadline(deadline time.Time) error             // SetReadDeadline 设置读超时
	SetWriteDeadline(deadline time.Time) error            // SetWriteDeadline 设置写超时
}

type ISession interface {
	ConnId() int64           // 连接ID
	UserId() int64           // 用户ID
	Read() ([]byte, error)   // 数据读取
	Write(data []byte) error // 写数据
	Close() error            // 关闭连接
	IsClosed() bool          // 是否存活
	Network() string         // 网络协议类型
	RefreshLastActiveAt()    // 刷新最后活跃时间
	LastActiveAt() int64     // 获取最后活跃时间
}

type ISessionManager interface {
	Options() *Options                         // 配置信息
	GenConnId() int64                          // 生成会话ID
	AllowAcceptConn() bool                     // 是否接受新连接
	NewSession(uid int64, conn IConn)          // 创建一个会话连接
	GetSession(connId int64) (ISession, error) // 获取连接
	GetSessionNum() int32                      // 获取连接总数
	GetSessionUserNum() int32                  // 获取连接总数
	GetConnIds(uid int64) []int64              // 获取用户在服务下的所有连接ID
	GetSessions(uid int64) []ISession          // 获取用户在服务下的所有连接
	Iterator() <-chan ISession                 // 迭代器(获取所有的连接)

	Assistant() IServerAssist
	Start(ctx context.Context) error // 启动服务
}

type IHeartbeat interface {
	Insert(connId int64, duration time.Duration)
	Cancel(connId int64)
}

type IHandler interface {
	OnOpen(smg ISessionManager, c ISession)                    // 打开连接事件
	OnMessage(smg ISessionManager, c ISession, message []byte) // 接收到数据事件
	OnClose(connId int64, uid int64)                           // 连接关闭事件
}

type IProcess interface {
	Start(ctx context.Context, s IServer) error // 启动服务
}

type IServer interface {
	ServerId() string // 服务ID

	SetCustomProcess(process IProcess) // 自定义处理任务
	SetAuthorize(cb IAuthorize)        // 鉴权
	SetHandler(h IHandler)             // 设置事件回调
	SetEncoder(h IEncoder)             // 设置数据编码解码器
	SetIdGenerator(gen IdGenerator)    // 设置连接ID生成器

	SessionManager() ISessionManager // 会话管理器
	Start(ctx context.Context) error // 启动服务
}

type IServerAssist interface {
	Handler() IHandler // 压缩器
	Encoder() IEncoder // 压缩器
	IdGenerator() IdGenerator
}

// ICompress 压缩器
type ICompress interface {
	Compress(data []byte) ([]byte, error)
	Decompress(data []byte) ([]byte, error)
}

type IEncrypter interface {
	Encrypt(data []byte) ([]byte, error) // 封包
	Decrypt(data []byte) ([]byte, error) // 封包
}

// IEncoder 封包解包器
type IEncoder interface {
	Pack(data *Packet) ([]byte, error)   // 封包
	UnPack(data []byte) (*Packet, error) // 解包
}
