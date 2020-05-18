package http

import (
	"encoding/json"

	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/application"
	"github.com/PolarPanda611/trinitygo/httputil"
)

var _ FakeChartData = new(fakeChartDataImpl)

func init() {
	trinitygo.RegisterController("/v1/fake_data", fakeChartDataImpl{},
		application.NewRequestMapping(httputil.GET, "/analyze", "GetAnalyze"),
		application.NewRequestMapping(httputil.GET, "/tags", "GetTags"),
		application.NewRequestMapping(httputil.GET, "/activity", "GetActivity"),
		application.NewRequestMapping(httputil.GET, "/notices", "GetNotice"),
	)
}

type FakeChartData interface {
	GetAnalyze()
	GetTags()
	GetActivity()
	GetNotice()
}

type fakeChartDataImpl struct {
	Tctx application.Context `autowired:"true"`
}

func (c *fakeChartDataImpl) GetAnalyze() {
	res := `{"visitData":[{"x":"2020-04-25","y":7},{"x":"2020-04-26","y":5},{"x":"2020-04-27","y":4},{"x":"2020-04-28","y":2},{"x":"2020-04-29","y":4},{"x":"2020-04-30","y":7},{"x":"2020-05-01","y":5},{"x":"2020-05-02","y":6},{"x":"2020-05-03","y":5},{"x":"2020-05-04","y":9},{"x":"2020-05-05","y":6},{"x":"2020-05-06","y":3},{"x":"2020-05-07","y":1},{"x":"2020-05-08","y":5},{"x":"2020-05-09","y":3},{"x":"2020-05-10","y":6},{"x":"2020-05-11","y":5}],"visitData2":[{"x":"2020-04-25","y":1},{"x":"2020-04-26","y":6},{"x":"2020-04-27","y":4},{"x":"2020-04-28","y":8},{"x":"2020-04-29","y":3},{"x":"2020-04-30","y":7},{"x":"2020-05-01","y":2}],"salesData":[{"x":"1月","y":562},{"x":"2月","y":1196},{"x":"3月","y":866},{"x":"4月","y":708},{"x":"5月","y":513},{"x":"6月","y":1125},{"x":"7月","y":516},{"x":"8月","y":693},{"x":"9月","y":1024},{"x":"10月","y":353},{"x":"11月","y":796},{"x":"12月","y":988}],"searchData":[{"index":1,"keyword":"搜索关键词-0","count":890,"range":20,"status":1},{"index":2,"keyword":"搜索关键词-1","count":996,"range":1,"status":1},{"index":3,"keyword":"搜索关键词-2","count":554,"range":44,"status":1},{"index":4,"keyword":"搜索关键词-3","count":507,"range":95,"status":0},{"index":5,"keyword":"搜索关键词-4","count":551,"range":80,"status":1},{"index":6,"keyword":"搜索关键词-5","count":590,"range":38,"status":1},{"index":7,"keyword":"搜索关键词-6","count":519,"range":29,"status":1},{"index":8,"keyword":"搜索关键词-7","count":836,"range":79,"status":1},{"index":9,"keyword":"搜索关键词-8","count":753,"range":39,"status":0},{"index":10,"keyword":"搜索关键词-9","count":936,"range":91,"status":1},{"index":11,"keyword":"搜索关键词-10","count":653,"range":53,"status":0},{"index":12,"keyword":"搜索关键词-11","count":837,"range":65,"status":1},{"index":13,"keyword":"搜索关键词-12","count":402,"range":54,"status":1},{"index":14,"keyword":"搜索关键词-13","count":126,"range":25,"status":0},{"index":15,"keyword":"搜索关键词-14","count":18,"range":31,"status":1},{"index":16,"keyword":"搜索关键词-15","count":598,"range":61,"status":1},{"index":17,"keyword":"搜索关键词-16","count":313,"range":43,"status":1},{"index":18,"keyword":"搜索关键词-17","count":657,"range":0,"status":1},{"index":19,"keyword":"搜索关键词-18","count":382,"range":70,"status":1},{"index":20,"keyword":"搜索关键词-19","count":229,"range":99,"status":0},{"index":21,"keyword":"搜索关键词-20","count":908,"range":68,"status":1},{"index":22,"keyword":"搜索关键词-21","count":71,"range":76,"status":1},{"index":23,"keyword":"搜索关键词-22","count":646,"range":13,"status":0},{"index":24,"keyword":"搜索关键词-23","count":439,"range":45,"status":0},{"index":25,"keyword":"搜索关键词-24","count":210,"range":23,"status":0},{"index":26,"keyword":"搜索关键词-25","count":830,"range":23,"status":0},{"index":27,"keyword":"搜索关键词-26","count":770,"range":13,"status":1},{"index":28,"keyword":"搜索关键词-27","count":716,"range":99,"status":1},{"index":29,"keyword":"搜索关键词-28","count":949,"range":4,"status":0},{"index":30,"keyword":"搜索关键词-29","count":210,"range":17,"status":0},{"index":31,"keyword":"搜索关键词-30","count":885,"range":85,"status":0},{"index":32,"keyword":"搜索关键词-31","count":339,"range":87,"status":0},{"index":33,"keyword":"搜索关键词-32","count":674,"range":87,"status":1},{"index":34,"keyword":"搜索关键词-33","count":353,"range":53,"status":1},{"index":35,"keyword":"搜索关键词-34","count":927,"range":20,"status":1},{"index":36,"keyword":"搜索关键词-35","count":380,"range":14,"status":0},{"index":37,"keyword":"搜索关键词-36","count":478,"range":59,"status":0},{"index":38,"keyword":"搜索关键词-37","count":295,"range":6,"status":1},{"index":39,"keyword":"搜索关键词-38","count":699,"range":55,"status":1},{"index":40,"keyword":"搜索关键词-39","count":371,"range":9,"status":0},{"index":41,"keyword":"搜索关键词-40","count":291,"range":31,"status":0},{"index":42,"keyword":"搜索关键词-41","count":131,"range":5,"status":1},{"index":43,"keyword":"搜索关键词-42","count":822,"range":96,"status":1},{"index":44,"keyword":"搜索关键词-43","count":308,"range":96,"status":1},{"index":45,"keyword":"搜索关键词-44","count":55,"range":21,"status":1},{"index":46,"keyword":"搜索关键词-45","count":918,"range":81,"status":0},{"index":47,"keyword":"搜索关键词-46","count":157,"range":16,"status":0},{"index":48,"keyword":"搜索关键词-47","count":139,"range":36,"status":0},{"index":49,"keyword":"搜索关键词-48","count":890,"range":83,"status":1},{"index":50,"keyword":"搜索关键词-49","count":798,"range":65,"status":1}],"offlineData":[{"name":"Stores 0","cvr":0.7},{"name":"Stores 1","cvr":0.6},{"name":"Stores 2","cvr":0.2},{"name":"Stores 3","cvr":0.3},{"name":"Stores 4","cvr":0.4},{"name":"Stores 5","cvr":0.4},{"name":"Stores 6","cvr":0.5},{"name":"Stores 7","cvr":0.3},{"name":"Stores 8","cvr":0.2},{"name":"Stores 9","cvr":0.2}],"offlineChartData":[{"x":1587791395690,"y1":45,"y2":86},{"x":1587793195690,"y1":68,"y2":32},{"x":1587794995690,"y1":70,"y2":59},{"x":1587796795690,"y1":80,"y2":34},{"x":1587798595690,"y1":66,"y2":17},{"x":1587800395690,"y1":33,"y2":90},{"x":1587802195690,"y1":70,"y2":78},{"x":1587803995690,"y1":54,"y2":70},{"x":1587805795690,"y1":28,"y2":56},{"x":1587807595690,"y1":37,"y2":19},{"x":1587809395690,"y1":95,"y2":83},{"x":1587811195690,"y1":58,"y2":53},{"x":1587812995690,"y1":83,"y2":64},{"x":1587814795690,"y1":21,"y2":79},{"x":1587816595690,"y1":70,"y2":102},{"x":1587818395690,"y1":51,"y2":75},{"x":1587820195690,"y1":21,"y2":103},{"x":1587821995690,"y1":97,"y2":75},{"x":1587823795690,"y1":48,"y2":87},{"x":1587825595690,"y1":73,"y2":99}],"salesTypeData":[{"x":"家用电器","y":4544},{"x":"食用酒水","y":3321},{"x":"个护健康","y":3113},{"x":"服饰箱包","y":2341},{"x":"母婴产品","y":1231},{"x":"其他","y":1231}],"salesTypeDataOnline":[{"x":"家用电器","y":244},{"x":"食用酒水","y":321},{"x":"个护健康","y":311},{"x":"服饰箱包","y":41},{"x":"母婴产品","y":121},{"x":"其他","y":111}],"salesTypeDataOffline":[{"x":"家用电器","y":99},{"x":"食用酒水","y":188},{"x":"个护健康","y":344},{"x":"服饰箱包","y":255},{"x":"其他","y":65}],"radarData":[{"name":"个人","label":"引用","value":10},{"name":"个人","label":"口碑","value":8},{"name":"个人","label":"产量","value":4},{"name":"个人","label":"贡献","value":5},{"name":"个人","label":"热度","value":7},{"name":"团队","label":"引用","value":3},{"name":"团队","label":"口碑","value":9},{"name":"团队","label":"产量","value":6},{"name":"团队","label":"贡献","value":3},{"name":"团队","label":"热度","value":1},{"name":"部门","label":"引用","value":4},{"name":"部门","label":"口碑","value":1},{"name":"部门","label":"产量","value":6},{"name":"部门","label":"贡献","value":5},{"name":"部门","label":"热度","value":7}]}`
	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(res), &dat); err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	c.Tctx.HTTPResponseOk(dat, nil)
}

func (c *fakeChartDataImpl) GetTags() {
	res := `{"visitData":[{"x":"2020-04-25","y":7},{"x":"2020-04-26","y":5},{"x":"2020-04-27","y":4},{"x":"2020-04-28","y":2},{"x":"2020-04-29","y":4},{"x":"2020-04-30","y":7},{"x":"2020-05-01","y":5},{"x":"2020-05-02","y":6},{"x":"2020-05-03","y":5},{"x":"2020-05-04","y":9},{"x":"2020-05-05","y":6},{"x":"2020-05-06","y":3},{"x":"2020-05-07","y":1},{"x":"2020-05-08","y":5},{"x":"2020-05-09","y":3},{"x":"2020-05-10","y":6},{"x":"2020-05-11","y":5}],"visitData2":[{"x":"2020-04-25","y":1},{"x":"2020-04-26","y":6},{"x":"2020-04-27","y":4},{"x":"2020-04-28","y":8},{"x":"2020-04-29","y":3},{"x":"2020-04-30","y":7},{"x":"2020-05-01","y":2}],"salesData":[{"x":"1月","y":562},{"x":"2月","y":1196},{"x":"3月","y":866},{"x":"4月","y":708},{"x":"5月","y":513},{"x":"6月","y":1125},{"x":"7月","y":516},{"x":"8月","y":693},{"x":"9月","y":1024},{"x":"10月","y":353},{"x":"11月","y":796},{"x":"12月","y":988}],"searchData":[{"index":1,"keyword":"搜索关键词-0","count":890,"range":20,"status":1},{"index":2,"keyword":"搜索关键词-1","count":996,"range":1,"status":1},{"index":3,"keyword":"搜索关键词-2","count":554,"range":44,"status":1},{"index":4,"keyword":"搜索关键词-3","count":507,"range":95,"status":0},{"index":5,"keyword":"搜索关键词-4","count":551,"range":80,"status":1},{"index":6,"keyword":"搜索关键词-5","count":590,"range":38,"status":1},{"index":7,"keyword":"搜索关键词-6","count":519,"range":29,"status":1},{"index":8,"keyword":"搜索关键词-7","count":836,"range":79,"status":1},{"index":9,"keyword":"搜索关键词-8","count":753,"range":39,"status":0},{"index":10,"keyword":"搜索关键词-9","count":936,"range":91,"status":1},{"index":11,"keyword":"搜索关键词-10","count":653,"range":53,"status":0},{"index":12,"keyword":"搜索关键词-11","count":837,"range":65,"status":1},{"index":13,"keyword":"搜索关键词-12","count":402,"range":54,"status":1},{"index":14,"keyword":"搜索关键词-13","count":126,"range":25,"status":0},{"index":15,"keyword":"搜索关键词-14","count":18,"range":31,"status":1},{"index":16,"keyword":"搜索关键词-15","count":598,"range":61,"status":1},{"index":17,"keyword":"搜索关键词-16","count":313,"range":43,"status":1},{"index":18,"keyword":"搜索关键词-17","count":657,"range":0,"status":1},{"index":19,"keyword":"搜索关键词-18","count":382,"range":70,"status":1},{"index":20,"keyword":"搜索关键词-19","count":229,"range":99,"status":0},{"index":21,"keyword":"搜索关键词-20","count":908,"range":68,"status":1},{"index":22,"keyword":"搜索关键词-21","count":71,"range":76,"status":1},{"index":23,"keyword":"搜索关键词-22","count":646,"range":13,"status":0},{"index":24,"keyword":"搜索关键词-23","count":439,"range":45,"status":0},{"index":25,"keyword":"搜索关键词-24","count":210,"range":23,"status":0},{"index":26,"keyword":"搜索关键词-25","count":830,"range":23,"status":0},{"index":27,"keyword":"搜索关键词-26","count":770,"range":13,"status":1},{"index":28,"keyword":"搜索关键词-27","count":716,"range":99,"status":1},{"index":29,"keyword":"搜索关键词-28","count":949,"range":4,"status":0},{"index":30,"keyword":"搜索关键词-29","count":210,"range":17,"status":0},{"index":31,"keyword":"搜索关键词-30","count":885,"range":85,"status":0},{"index":32,"keyword":"搜索关键词-31","count":339,"range":87,"status":0},{"index":33,"keyword":"搜索关键词-32","count":674,"range":87,"status":1},{"index":34,"keyword":"搜索关键词-33","count":353,"range":53,"status":1},{"index":35,"keyword":"搜索关键词-34","count":927,"range":20,"status":1},{"index":36,"keyword":"搜索关键词-35","count":380,"range":14,"status":0},{"index":37,"keyword":"搜索关键词-36","count":478,"range":59,"status":0},{"index":38,"keyword":"搜索关键词-37","count":295,"range":6,"status":1},{"index":39,"keyword":"搜索关键词-38","count":699,"range":55,"status":1},{"index":40,"keyword":"搜索关键词-39","count":371,"range":9,"status":0},{"index":41,"keyword":"搜索关键词-40","count":291,"range":31,"status":0},{"index":42,"keyword":"搜索关键词-41","count":131,"range":5,"status":1},{"index":43,"keyword":"搜索关键词-42","count":822,"range":96,"status":1},{"index":44,"keyword":"搜索关键词-43","count":308,"range":96,"status":1},{"index":45,"keyword":"搜索关键词-44","count":55,"range":21,"status":1},{"index":46,"keyword":"搜索关键词-45","count":918,"range":81,"status":0},{"index":47,"keyword":"搜索关键词-46","count":157,"range":16,"status":0},{"index":48,"keyword":"搜索关键词-47","count":139,"range":36,"status":0},{"index":49,"keyword":"搜索关键词-48","count":890,"range":83,"status":1},{"index":50,"keyword":"搜索关键词-49","count":798,"range":65,"status":1}],"offlineData":[{"name":"Stores 0","cvr":0.7},{"name":"Stores 1","cvr":0.6},{"name":"Stores 2","cvr":0.2},{"name":"Stores 3","cvr":0.3},{"name":"Stores 4","cvr":0.4},{"name":"Stores 5","cvr":0.4},{"name":"Stores 6","cvr":0.5},{"name":"Stores 7","cvr":0.3},{"name":"Stores 8","cvr":0.2},{"name":"Stores 9","cvr":0.2}],"offlineChartData":[{"x":1587791395690,"y1":45,"y2":86},{"x":1587793195690,"y1":68,"y2":32},{"x":1587794995690,"y1":70,"y2":59},{"x":1587796795690,"y1":80,"y2":34},{"x":1587798595690,"y1":66,"y2":17},{"x":1587800395690,"y1":33,"y2":90},{"x":1587802195690,"y1":70,"y2":78},{"x":1587803995690,"y1":54,"y2":70},{"x":1587805795690,"y1":28,"y2":56},{"x":1587807595690,"y1":37,"y2":19},{"x":1587809395690,"y1":95,"y2":83},{"x":1587811195690,"y1":58,"y2":53},{"x":1587812995690,"y1":83,"y2":64},{"x":1587814795690,"y1":21,"y2":79},{"x":1587816595690,"y1":70,"y2":102},{"x":1587818395690,"y1":51,"y2":75},{"x":1587820195690,"y1":21,"y2":103},{"x":1587821995690,"y1":97,"y2":75},{"x":1587823795690,"y1":48,"y2":87},{"x":1587825595690,"y1":73,"y2":99}],"salesTypeData":[{"x":"家用电器","y":4544},{"x":"食用酒水","y":3321},{"x":"个护健康","y":3113},{"x":"服饰箱包","y":2341},{"x":"母婴产品","y":1231},{"x":"其他","y":1231}],"salesTypeDataOnline":[{"x":"家用电器","y":244},{"x":"食用酒水","y":321},{"x":"个护健康","y":311},{"x":"服饰箱包","y":41},{"x":"母婴产品","y":121},{"x":"其他","y":111}],"salesTypeDataOffline":[{"x":"家用电器","y":99},{"x":"食用酒水","y":188},{"x":"个护健康","y":344},{"x":"服饰箱包","y":255},{"x":"其他","y":65}],"radarData":[{"name":"个人","label":"引用","value":10},{"name":"个人","label":"口碑","value":8},{"name":"个人","label":"产量","value":4},{"name":"个人","label":"贡献","value":5},{"name":"个人","label":"热度","value":7},{"name":"团队","label":"引用","value":3},{"name":"团队","label":"口碑","value":9},{"name":"团队","label":"产量","value":6},{"name":"团队","label":"贡献","value":3},{"name":"团队","label":"热度","value":1},{"name":"部门","label":"引用","value":4},{"name":"部门","label":"口碑","value":1},{"name":"部门","label":"产量","value":6},{"name":"部门","label":"贡献","value":5},{"name":"部门","label":"热度","value":7}]}`
	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(res), &dat); err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	c.Tctx.HTTPResponseOk(dat, nil)
}

func (c *fakeChartDataImpl) GetActivity() {
	res := `
	[
        {
            "group": {
                "link": "http://github.com/",
                "name": "徐家汇站"
            },
            "id": "trend-1",
            "project": {
                "link": "http://github.com/",
                "name": "出发"
            },
            "template": "在 @{group}  @{project}",
            "updatedAt": "2020-04-25T05:09:55.617Z",
            "user": {
                "avatar": "https://gw.alipayobjects.com/zos/antfincdn/XAosXuNZyF/BiazfanxmamNRoxxVxka.png",
                "name": "M1111"
            }
        },
        {
            "group": {
                "link": "http://github.com/",
                "name": "徐家汇站"
            },
            "id": "trend-2",
            "project": {
                "link": "http://github.com/",
                "name": "进站"
            },
            "template": "在 @{group}  @{project}",
            "updatedAt": "2020-04-25T05:09:55.617Z",
            "user": {
                "avatar": "https://gw.alipayobjects.com/zos/rmsportal/cnrhVkzwxjPwAaCfPbdc.png",
                "name": "M2222"
            }
        },
        {
            "group": {
                "link": "http://github.com/",
                "name": "徐家汇站"
            },
            "id": "trend-3",
            "project": {
                "link": "http://github.com/",
                "name": "出发"
            },
            "template": "在 @{group}  @{project}",
            "updatedAt": "2020-04-25T05:09:55.617Z",
            "user": {
                "avatar": "https://gw.alipayobjects.com/zos/rmsportal/gaOngJwsRYRaVAuXXcmB.png",
                "name": "M2222"
            }
        }
	]
	`
	var dat []map[string]interface{}
	if err := json.Unmarshal([]byte(res), &dat); err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	c.Tctx.HTTPResponseOk(dat, nil)
}
func (c *fakeChartDataImpl) GetNotice() {
	res := `
	[
        {
			"index":1,
            "href": "",
            "metro_code": "M1111",
            "next_station":"漕宝路",
            "out": "2020-04-25T05:09:55.617Z",
            "out_comment": "上海南站拥堵",
            "estimate_in":"2020-04-25T05:14:55.617Z"
        },
        {
			"index":2,
            "description": "希望是一个好东西，也许是最好的，好东西是不会消亡的",
            "href": "",
            "member": "全组都是吴彦祖",
            "memberLink": "",
"next_station":"上海体育馆",
            "metro_code": "M2222",
            "out_comment": "列车满载",
            "out": "2020-04-25T05:10:55.617Z",
 "estimate_in":"2020-04-25T05:15:55.617Z"
        },
        {
			"index":3,
            "description": "城镇中有那么多的酒馆，她却偏偏走进了我的酒馆",
            "href": "",
            "member": "中二少女团",
            "memberLink": "",
            "metro_code": "M3333",
"next_station":"徐家汇",
            "out_comment": "需要增加列车",
            "out": "2020-04-25T05:11:55.617Z",
            "estimate_in":"2020-04-25T05:16:55.617Z"
        },
        {
			"index":4,
            "description": "那时候我只会想自己想要什么，从不想自己拥有什么",
            "href": "",
            "member": "程序员日常",
            "memberLink": "",
            "metro_code": "M4444",
"next_station":"衡山知路",
            "out_comment": "需要增加列车",
            "out": "2020-04-25T05:12:55.617Z",
 "estimate_in":"2020-04-25T05:17:55.617Z"
        },
        {
			"index":5,
            "description": "凛冬将至",
            "href": "",
            "member": "高逼格设计天团",
            "memberLink": "",
            "metro_code": "M5555",
"next_station":"常熟路",
"out_comment": "需要增加列车",

            "out": "2020-04-25T05:13:55.617Z",
 "estimate_in":"2020-04-25T05:18:55.617Z"
        },
        {
			"index":6,
            "description": "生命就像一盒巧克力，结果往往出人意料",
            "href": "",
            "member": "骗你来学计算机",
            "memberLink": "",
"next_station":"陕西南路",
            "metro_code": "M6666",
"out_comment": "需要增加列车",
            "out":  "2020-04-25T05:14:55.617Z",
			"estimate_in":"2020-04-25T05:19:55.617Z"
        }
    ]
	`
	var dat []map[string]interface{}
	if err := json.Unmarshal([]byte(res), &dat); err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	c.Tctx.HTTPResponseOk(dat, nil)
}
