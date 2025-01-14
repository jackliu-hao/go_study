package sms

import "context"

type Service interface {
	// Send biz : 业务
	Send(ctx context.Context, biz string, args []string, phoneNumbers ...string) error
}

// 1、提高可用性 ：重试机制、客户端限流、failover(轮询、实时监测)
//  1.1 实时检测：
//    1.1.1 基于超时的实时检测（连续超时）
//    1.1.2 基于响应时间的实时检测 （比如说,平均响应时间上升20%）
//    1.1.3 基于长尾请求的实时监测（响应时间超过 1s 的请求占比超过了10%）
//    1.1.4 错误率

//2. 提高安全性
// 2.1 完整的资源申请和审批流程
// 2.2 鉴权:
//  2.2.1 静态token
//  2.2.2 动态token
//3. 提高可观测性： 日志、metriecs,tracing ，丰富的排查手段

// 4. 提高性能
