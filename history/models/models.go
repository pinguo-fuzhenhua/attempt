package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OperateHistory struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`        // 操作人姓名
	Email       string             `bson:"email"`       // 操作人邮箱
	OperateTime time.Time          `bson:"operateTime"` // 操作时间
	SystemName  string             `bson:"systemName"`  // 系统名称
	Scope       string             `bson:"scope"`       // 平台
	ModuleName  string             `bson:"moduleName"`  // 模块
	// OperateType int32     `bson:"operateType"` // 操作类型
	OperateType string   `bson:"operateType"` // 操作类型
	EntityType  string   `bson:"entityType"`  // 实体类型
	EntityName  []string `bson:"entityName"`  // 实体名称
	EntityID    []string `bson:"entityID"`    // 实体id
	//请求信息
	PageID         string `bson:"pageID"`         // 页面唯一标识
	RequestPath    string `bson:"requestPath"`    // 请求url
	RequestMethod  string `bson:"requestMethod"`  // 请求方法
	ClientIP       string `bson:"clientIP"`       // 客户端ip
	QueryParams    string `bson:"queryParams"`    // url参数
	RequestHeaders string `bson:"requestHeaders"` // 请求头信息
	RequestBody    string `bson:"requestBody"`    // 请求的body参数
	ResponseData   string `bson:"responseData"`   // 响应数据
	ResponseErrMsg string `bson:"responseErrMsg"` // 响应的错误信息
	// TraceID        string `bson:"traceId"`        // 请求跟踪id
	TraceID     string `bson:"traceId"`     // 请求跟踪id
	Environment string `bson:"environment"` // 环境
}

// 实体操作记录
type EntityOperate struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	OperateTime   time.Time          `bson:"operateTime"` // 操作时间
	EntityType    string             `bson:"entityType"`  // 实体类型 (展示操作记录使用)
	EntityName    []string           `bson:"entityName"`  // 实体名称（跳转展示使用）
	EntityID      []string           `bson:"entityID"`    // 实体id（跳转使用）
	TraceID       string             `bson:"traceID"`     // 请求跟踪id
	Scope         string             `bson:"scope"`       // 产品
	Environment   string             `bson:"environment"` // 环境
	CollID        string             `bson:"collID"`      // dbName|collectionName|_id
	OperateOption []*OperateOption   `bson:"operate"`     // 操作选项
	OperateType   string             `bson:"operateType"` // 操作类型
}

// 操作消息
type OperateOption struct {
	FieldName   string `bson:"fieldName"`   // 字段名
	BeforeValue string `bson:"beforeValue"` // 操作前的值
	AfterValue  string `bson:"afterValue"`  // 操作后的值
	ValueType   string `bson:"valueType"`   // 值的类型
}
