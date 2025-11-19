package vo

import (
	"github.com/suhanyujie/go_utils/core/consts"
	_ "github.com/suhanyujie/go_utils/pkg/common_interface_vendor/pb_interface/api/platform/service/v1"
)

type BaseResp struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Status string `json:"status"`
}

func (_this *BaseResp) IsOk() bool {
	return _this.Code == consts.PlatformHttpApiOk
}

func (_this *BaseResp) IsFail() bool {
	return _this.Code != consts.PlatformHttpApiOk
}

type GameCallbackV2Req struct {
	BattleId             int64 `json:"battleId"`
	PlayerRankReqDtoList []any `json:"playerRankReqDtoList"`
}

// 因时常修改该结构，打算废弃，直接用 proto 中的结构
type PlayerRankReqDtoItem struct {
	Uid               string  `json:"uid"`
	NftId             int     `json:"nftId"` // 自己的 nftId
	RewardCoefficient float64 `json:"rewardCoefficient"`
	Rank              int     `json:"rank"`
	Stage             int     `json:"stage"`         // 当前关卡
	Point             int     `json:"point"`         // 当前进度，不同游戏对应的属性不同，可能是：波次，关卡，积分等
	ConsumeEnergy     int     `json:"consumeEnergy"` // 当前关卡消耗的体力值。只针对 pve
	WinStatus         int     `json:"winStatus"`     // 对于 uid，赢为1，输为 2
	GameTimeSec       int     `json:"gameTimeSec"`   // 游戏时长，秒数

	ChestMap map[string]int `json:"runes"`  // 奖励的符文，对应平台的 runes 字段
	Energy   int            `json:"energy"` // 奖励的体力值
	Token    float64        `json:"token"`  // 奖励的代币
}

type GameCallbackResp struct {
	BaseResp
	Data GameCallbackRespData `json:"data"`
}

type GameCallbackRespData struct {
	NftScore  int     `json:"nftScore"`
	EarnToken float32 `json:"earnToken"`

	//BattleId             int                           `json:"battleId"`
	//PlayerRankReqDtoList []GameCallbackRespDataDtoItem `json:"playerRankReqDtoList"`
}

type GameCallbackRespDataDtoItem struct {
	NftId             int            `json:"nftId"`
	Rank              int            `json:"rank"`
	RewardCoefficient float32        `json:"rewardCoefficient"` // 奖励系数
	Energy            int            `json:"energy"`            // 奖励的体力值
	Runes             map[string]int `json:"runes"`             // 奖励的符文，对应平台的 runes 字段
	Token             float32        `json:"token"`             // 奖励的 token
	Stage             int            `json:"stage"`             // 游戏进度 关卡
	Point             int            `json:"point"`             // 游戏进度 分数
	Uid               string         `json:"uid"`
}

type PlayerRankReqDtoItemReward struct {
	// 1：消耗金币。2:消耗体力值。3:消耗符文
	Type  int     `json:"type"`
	Value float64 `json:"value"`
}

type GameSessCallbackReq struct {
	BattleId          int64   `json:"battleId"`          // 对局 room
	Uid               string  `json:"uid"`               // 对局 uid
	RewardCoefficient float64 `json:"rewardCoefficient"` // 系数
	IsRecommend       bool    `json:"isRecommend"`       // 是否获取推荐的场次信息
}

type SessCallbackRespVo struct {
	//与平台结算的结果
	BaseResp
	Data *SessCallbackRespVoData `json:"data"`
}

// 推荐的场次信息
type SessCallbackRespVoData struct {
	LobbyId      int     `json:"lobbyId"`
	GameId       int     `json:"gameId"`   // 一款游戏的 id
	GameName     string  `json:"gameName"` // 一款游戏的名称
	Mode         string  `json:"mode"`
	ModeName     string  `json:"modeName"`
	ConsumeType  int     `json:"consumeType"`  // 推荐的场次的消耗类型，这里是门票
	ConsumeValue float64 `json:"consumeValue"` // 推荐的场次的门票数
	RewardType   int     `json:"rewardType"`   // 推荐的场次的奖励类型，这里是 token
	RewardValue  float64 `json:"rewardValue"`  // 推荐的场次的奖品数值
	BannerImg    string  `json:"bannerImg"`    // 推荐的场次的图片

	TaskInfos []any   `json:"taskInfos"` // 当前玩家的任务进度信息
	NftScore  int     `json:"nftScore"`  // 暂无，忽略
	EarnToken float32 `json:"earnToken"` // 暂无，忽略
}

type TaskInfoItem struct {
	// 任务文案
	TaskDetail string `json:"taskDetail" msgpack:"taskDetail"`
	// 已完成数量
	GoingNumber int32 `json:"goingNumber" msgpack:"goingNumber"`
	// 该任务的子任务总数
	Total int32 `json:"total" msgpack:"total"`
	// 任务是否完成
	IsCompleted bool `json:"isCompleted" msgpack:"isCompleted"`
}

type PvpKnockoutSettleReqVo struct {
	BattleId        int64  `json:"battleId"`
	Uid             string `json:"uid"`             // 玩家 uid，胜方的 uid
	CompetitorId    string `json:"competitorId"`    // 对手 uid，败方的 uid
	NftId           int    `json:"nftId"`           // 胜方对应的 nftId
	CompetitorNftId int    `json:"competitorNftId"` // 败方对应的 nftId
}

type BaseSucResp struct {
	BaseResp
	Data interface{} `json:"data"`
}

type GetUserInfoResp struct {
	BaseResp
	Data GetUserInfoRespData `json:"data"`
}

type GetUserInfoRespData struct {
	UserId string `json:"userId"`
	NftId  int    `json:"nftId"`
}

type CreateOrderReq struct {
	GameId       int             `json:"gameId"`
	Uid          string          `json:"uid"`
	Goods        map[int32]int32 `json:"goods"`
	CurrencyType int             `json:"currencyType"`
	Price        float64         `json:"price"`
}

type CreateOrderResp struct {
	BaseResp
	Data CreateOrderRespData `json:"data"`
}

type CreateOrderRespData struct {
	// 订单号
	OrderNum string `json:"orderNum"`
}

type QueryOrderStatusReq struct {
	GameId   int    `json:"gameId"`
	Uid      string `json:"uid"`
	OrderNum string `json:"orderNum"`
}

type QueryOrderStatusResp struct {
	BaseResp
	Data QueryOrderStatusRespData `json:"data"`
}

type QueryOrderStatusRespData struct {
	// 订单状态
	Status int `json:"status"`
}

type QueryTaskStateReq struct {
	GameId    int    `json:"gameId"`
	Uid       string `json:"uid"`
	TaskIdArr []int  `json:"taskIdArr"`
}

type QueryTaskStateResp struct {
	BaseResp
	Data QueryTaskStateRespData `json:"data"`
}

type QueryTaskStateRespData struct {
	StateInfo map[int]int `json:"stateInfo"`
}

type GetPropListReq struct {
	GameId int    `json:"gameId"`
	Uid    string `json:"uid"`
}

type GetPropListResp struct {
	BaseResp
	Data GetPropListRespData `json:"data"`
}

type GetPropListRespData struct {
	PropList []PropItem `json:"propList"`
}

type PropItem struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Num  int    `json:"num"`
}

type ConsumeEnergyReq struct {
	GameId  string `json:"gameName"`
	LevelId int64  `json:"levelId"`
	Energy  int    `json:"needEnergy"`
}

type ConsumeEnergyResp struct {
	BaseResp
	Data any `json:"data"`
}

type GamePveNextReq struct {
	BattleId int64 `json:"battleId"`
}

type GamePveNextResp struct {
	BaseResp
	Data any `json:"data"`
}

type FarmEventPushReq struct {
	// 类型（1：收获 2：种植 3：天灾人祸 4:开垦土地 5:购买礼包 6:与神秘商人交易）
	Type  int    `json:"type"`
	Value string `json:"value"`
}

type CreateMatchReq struct {
	GameId         int                    `json:"gameId"`
	BattleId       int64                  `json:"battleId"` // 对局 id
	Mode           string                 `json:"mode"`
	TeamNum        int                    `json:"teamNum"`
	PlayersPerTeam int                    `json:"playersPerTeam"`
	Difficulty     int                    `json:"difficulty"` // 跑酷的难度参数，1~5
	Stage          int                    `json:"stage"`
	Players        []CreateMatchReqPlayer `json:"players"`
	Timeout        int                    `json:"timeout"`
	// GameType 0:pve 1:1v1 2:1vn 3:淘汰赛
	GameType int `json:"gameType"`
}

type CreateMatchReqPlayer struct {
	Token      string                          `json:"token"` //Uid-> 对应的token
	IsRobot    bool                            `json:"isRobot"`
	RobotLevel int                             `json:"robotLevel"`
	Uid        string                          `json:"uid"`       // 如果是机器人，这里就是机器人的UID
	Legends    []string                        `json:"legends"`   // nft id 列表
	TaskInfos  []*CreateMatchReqPlayerTaskItem `json:"taskInfos"` //玩家的任务进度信息
}

type CreateMatchReqPlayerTaskItem struct {
	// 任务文案
	TaskDetail string `json:"taskDetail" msgpack:"taskDetail"`
	// 已完成数量
	GoingNumber int32 `json:"goingNumber" msgpack:"goingNumber"`
	// 该任务的子任务总数
	Total int32 `json:"total" msgpack:"total"`
	// 任务是否完成
	IsCompleted bool `json:"isCompleted" msgpack:"isCompleted"`
}

// ref: https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN#f62e72d5
type Send2FsGroupReq struct {
	MsgType string              `json:"msg_type"`
	Card    Send2FsGroupReqCard `json:"card"`
}

type Send2FsGroupReqCard struct {
	Header   Send2FsGroupReqCardHeader `json:"header"`
	Elements []Send2FsGroupReqElement  `json:"elements"`
}

type Send2FsGroupReqCardHeader struct {
	Title Send2FsGroupReqCardHeaderTitle `json:"title"`
}

type Send2FsGroupReqCardHeaderTitle struct {
	Tag     string `json:"tag"` // plain_text
	Content string `json:"content"`
}

type Send2FsGroupReqElement struct {
	Tag  string                     `json:"tag"` // div
	Text Send2FsGroupReqElementText `json:"text"`
}

type Send2FsGroupReqElementText struct {
	Tag     string `json:"tag"`     // lark_md
	Content string `json:"content"` // 放置 md 文档内容
}
