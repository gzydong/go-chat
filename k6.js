import ws from 'k6/ws';
import { check, sleep } from 'k6';
import { Counter, Rate, Trend } from 'k6/metrics';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

// 自定义指标
const wsConnections = new Counter('ws_connections');                // 连接计数
const wsConnectionErrors = new Counter('ws_connection_errors');     // 连接错误计数
const wsMessagesSent = new Counter('ws_messages_sent');             // 发送消息计数
const wsMessagesReceived = new Counter('ws_messages_received');     // 接收消息计数
const wsMessageErrors = new Counter('ws_message_errors');           // 消息错误计数
const wsConnectionDuration = new Trend('ws_connection_duration');   // 连接持续时间
const wsMessageSendTime = new Trend('ws_message_send_time');        // 消息发送时间
const wsMessageReceiveTime = new Trend('ws_message_receive_time');  // 消息接收时间
const wsConnectionSuccessRate = new Rate('ws_connection_success');  // 连接成功率

// 配置参数
export let options = {
    // 阶梯式加压测试
    stages: [
        { duration: '30s', target: 50 },    // 30秒内增加到50个并发用户
        { duration: '1m', target: 100 },    // 1分钟内增加到100个并发用户
        { duration: '3m', target: 100 },    // 保持100个并发用户3分钟
        { duration: '1m', target: 200 },    // 1分钟内增加到200个并发用户
        { duration: '3m', target: 200 },    // 保持200个并发用户3分钟
        { duration: '1m', target: 0 },      // 1分钟内减少到0个并发用户
    ],
    // 限制每秒请求数，防止瞬间打爆服务器
    rps: 1000,
    // 启用阈值检查
    thresholds: {
        'ws_connection_success': ['rate>0.95'], // 连接成功率应大于95%
        'ws_connection_duration': ['p(95)<1000'], // 95%的连接时间应小于1秒
        'ws_message_send_time': ['p(95)<100'],  // 95%的消息发送时间应小于100ms
        'ws_message_receive_time': ['p(95)<200'], // 95%的消息接收时间应小于200ms
    },
};

// 生成随机用户ID (5-6位数字)
function getRandomUserId() {
    return Math.floor(Math.random() * 900000) + 100000;
}

// 生成随机消息大小 (1KB - 10KB)
function getRandomMessage(size) {
    return randomString(size);
}

// 主测试函数
export default function() {
    const userId = getRandomUserId();
    const url = `ws://localhost:9501/wss?uid=${userId}`;

    // 记录开始时间
    const startTime = new Date().getTime();

    // 连接WebSocket
    const res = ws.connect(url, null, function(socket) {
        wsConnections.add(1);
        wsConnectionSuccessRate.add(true);

        // 连接成功事件
        socket.on('open', function() {
            console.log(`连接成功: 用户 ${userId}`);

            // 上行数据测试 - 发送多条不同大小的消息
            const messageSizes = [1024, 5120, 10240]; // 1KB, 5KB, 10KB
            const messageCount = 100; // 每种大小发送20条

            // 存储消息ID和发送时间，用于计算RTT
            const sentMessages = new Map();

            // 发送消息
            for (let size of messageSizes) {
                for (let i = 0; i < messageCount; i++) {
                    const messageId = `${userId}-${size}-${i}`;
                    const payload = {
                        cmd: "message",
                        id: messageId,
                        timestamp: new Date().getTime(),
                        body: {
                            type: "text",
                            size: size,
                            content: getRandomMessage(size)
                        }
                    };

                    const sendStart = new Date().getTime();
                    sentMessages.set(messageId, sendStart);

                    try {
                        socket.send(JSON.stringify(payload));
                        wsMessagesSent.add(1);
                        wsMessageSendTime.add(new Date().getTime() - sendStart);
                    } catch (e) {
                        console.error(`发送消息失败: ${e}`);
                        wsMessageErrors.add(1);
                    }

                    // 短暂休眠，避免消息拥塞
                    sleep(0.05); // 50ms
                }
            }

            // 下行数据测试 - 请求服务器推送大量数据
            const downloadRequest = {
                cmd: "download_test",
                id: `download-${userId}`,
                timestamp: new Date().getTime(),
                body: {
                    size: 50240, // 请求约50KB的数据
                    chunks: 5    // 分5块发送
                }
            };

            try {
                socket.send(JSON.stringify(downloadRequest));
            } catch (e) {
                console.error(`发送下行测试请求失败: ${e}`);
            }

            // 设置一个定时器，在适当的时间后关闭连接
            setTimeout(function() {
                try {
                    socket.close();
                } catch (e) {
                    console.error(`关闭连接失败: ${e}`);
                }
            }, 15000); // 15秒后关闭连接
        });

        // 接收消息事件
        socket.on('message', function(data) {
            wsMessagesReceived.add(1);

            try {
                const message = JSON.parse(data);
                const receiveTime = new Date().getTime();

                // 检查是否是我们发送的消息的响应
                if (message.id && message.cmd === "ack") {
                    const originalId = message.id;
                    if (sentMessages.has(originalId)) {
                        const sendTime = sentMessages.get(originalId);
                        const rtt = receiveTime - sendTime;
                        wsMessageReceiveTime.add(rtt);

                        // 记录RTT日志（仅在开发时使用，生产环境可以注释掉）
                        // console.log(`消息 ${originalId} RTT: ${rtt}ms`);

                        // 从Map中删除已处理的消息
                        sentMessages.delete(originalId);
                    }
                }

                // 检查下行数据测试响应
                if (message.cmd === "download_chunk") {
                    // 记录下行数据块接收
                    console.log(`接收到下行数据块: ${message.chunk_id}/${message.total_chunks}, 大小: ${message.body.content.length}`);
                }

            } catch (e) {
                console.error(`解析消息失败: ${e}, 原始数据: ${data}`);
                wsMessageErrors.add(1);
            }
        });

        // 错误事件
        socket.on('error', function(e) {
            console.error(`WebSocket错误: ${e}`);
            wsConnectionErrors.add(1);
        });

        // 关闭事件
        socket.on('close', function() {
            const duration = new Date().getTime() - startTime;
            wsConnectionDuration.add(duration);
            console.log(`连接关闭: 用户 ${userId}, 持续时间: ${duration}ms`);
        });
    });

    // 检查连接是否成功
    check(res, { 'status is 101': (r) => r && r.status === 101 });

    // 如果连接失败
    if (!res || res.status !== 101) {
        wsConnectionSuccessRate.add(false);
        wsConnectionErrors.add(1);
        console.error(`连接失败: 用户 ${userId}, 状态码: ${res ? res.status : 'unknown'}`);
    }

    // 等待连接完成
    sleep(20); // 等待20秒，确保有足够时间完成WebSocket交互
}

// 测试完成后的清理函数
export function teardown() {
    console.log('测试完成，清理资源...');
}
