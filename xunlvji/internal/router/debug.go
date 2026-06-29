package router

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// debugPage 开发预览 HTML 页面
const debugPage = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>寻旅记 API 预览</title>
<style>
  * { margin: 0; padding: 0; box-sizing: border-box; }
  body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; background: #f5f7fa; color: #2c3e50; }
  .header { background: linear-gradient(135deg, #4EC3F7 0%, #2196F3 100%); color: #fff; padding: 24px; text-align: center; }
  .header h1 { font-size: 24px; margin-bottom: 6px; }
  .header p { font-size: 13px; opacity: 0.9; }
  .container { max-width: 1200px; margin: 20px auto; padding: 0 16px; }
  .grid { display: grid; grid-template-columns: 320px 1fr; gap: 16px; }
  .sidebar { background: #fff; border-radius: 12px; padding: 16px; box-shadow: 0 2px 8px rgba(0,0,0,0.06); height: fit-content; position: sticky; top: 16px; }
  .sidebar h2 { font-size: 15px; color: #2c3e50; margin-bottom: 12px; padding-bottom: 8px; border-bottom: 2px solid #4EC3F7; }
  .api-group { margin-bottom: 16px; }
  .api-group-title { font-size: 12px; color: #95a5a6; font-weight: 600; margin-bottom: 6px; text-transform: uppercase; letter-spacing: 0.5px; }
  .api-item { display: block; width: 100%; text-align: left; padding: 8px 12px; margin-bottom: 4px; border: none; background: #f8f9fa; border-radius: 6px; cursor: pointer; font-size: 13px; color: #2c3e50; transition: all 0.2s; }
  .api-item:hover { background: #e3f2fd; color: #4EC3F7; }
  .api-item.active { background: #4EC3F7; color: #fff; }
  .method { display: inline-block; padding: 1px 5px; border-radius: 3px; font-size: 10px; font-weight: 700; margin-right: 6px; }
  .method.GET { background: #31C27C; color: #fff; }
  .method.POST { background: #FFB020; color: #fff; }
  .main { background: #fff; border-radius: 12px; padding: 20px; box-shadow: 0 2px 8px rgba(0,0,0,0.06); min-height: 600px; }
  .endpoint-info { margin-bottom: 16px; }
  .endpoint-info .path { font-family: 'Courier New', monospace; font-size: 14px; color: #4EC3F7; word-break: break-all; }
  .endpoint-info .desc { font-size: 13px; color: #7f8c8d; margin-top: 4px; }
  .params { background: #f8f9fa; border-radius: 8px; padding: 12px; margin-bottom: 16px; }
  .params-title { font-size: 12px; font-weight: 600; color: #2c3e50; margin-bottom: 8px; }
  .param-row { display: flex; gap: 8px; margin-bottom: 6px; align-items: center; font-size: 12px; }
  .param-row label { min-width: 100px; color: #7f8c8d; }
  .param-row input, .param-row select { flex: 1; padding: 4px 8px; border: 1px solid #ddd; border-radius: 4px; font-size: 12px; }
  .btn-send { background: #4EC3F7; color: #fff; border: none; padding: 10px 24px; border-radius: 8px; font-size: 14px; font-weight: 600; cursor: pointer; transition: background 0.2s; }
  .btn-send:hover { background: #2196F3; }
  .btn-send:disabled { background: #bdc3c7; cursor: not-allowed; }
  .response-section { margin-top: 20px; }
  .response-meta { display: flex; gap: 12px; margin-bottom: 8px; font-size: 12px; }
  .status-badge { padding: 2px 8px; border-radius: 4px; font-weight: 600; }
  .status-2xx { background: #d4edda; color: #155724; }
  .status-4xx { background: #fff3cd; color: #856404; }
  .status-5xx { background: #f8d7da; color: #721c24; }
  .json-viewer { background: #2c3e50; color: #ecf0f1; padding: 16px; border-radius: 8px; font-family: 'Courier New', monospace; font-size: 12px; line-height: 1.6; overflow-x: auto; max-height: 500px; overflow-y: auto; white-space: pre-wrap; word-break: break-all; }
  .json-key { color: #4EC3F7; }
  .json-string { color: #2ecc71; }
  .json-number { color: #e67e22; }
  .json-boolean { color: #e74c3c; }
  .json-null { color: #95a5a6; }
  .empty-state { text-align: center; padding: 80px 20px; color: #bdc3c7; }
  .empty-state .icon { font-size: 48px; margin-bottom: 12px; }
  .loading { text-align: center; padding: 40px; color: #4EC3F7; }
</style>
</head>
<body>
<div class="header">
  <h1>寻旅记 API 预览</h1>
  <p>v2.1.7 · Mock 模式（荆州古城测试数据）</p>
</div>
<div class="container">
  <div class="grid">
    <div class="sidebar">
      <h2>接口列表</h2>
      <div class="api-group">
        <div class="api-group-title">内容浏览</div>
        <button class="api-item" data-api="recommend"><span class="method GET">GET</span>推荐流</button>
        <button class="api-item" data-api="guide"><span class="method GET">GET</span>攻略流</button>
        <button class="api-item" data-api="place"><span class="method GET">GET</span>打卡地流</button>
        <button class="api-item" data-api="route"><span class="method GET">GET</span>路线流</button>
        <button class="api-item" data-api="nearby"><span class="method GET">GET</span>附近流</button>
      </div>
      <div class="api-group">
        <div class="api-group-title">内容详情</div>
        <button class="api-item" data-api="route-detail"><span class="method GET">GET</span>路线详情</button>
        <button class="api-item" data-api="place-detail"><span class="method GET">GET</span>地点详情</button>
        <button class="api-item" data-api="guide-detail"><span class="method GET">GET</span>攻略详情</button>
      </div>
      <div class="api-group">
        <div class="api-group-title">系统</div>
        <button class="api-item" data-api="health"><span class="method GET">GET</span>健康检查</button>
      </div>
    </div>
    <div class="main">
      <div id="content">
        <div class="empty-state">
          <div class="icon">👈</div>
          <p>从左侧选择一个接口开始测试</p>
        </div>
      </div>
    </div>
  </div>
</div>
<script>
const APIs = {
  recommend: {
    method: 'GET', path: '/api/v1/content/feed/recommend',
    desc: '首页推荐流，混合展示路线/地点/攻略',
    params: [
      {name: 'page', label: '页码', type: 'input', value: '1'},
      {name: 'page_size', label: '每页数量', type: 'input', value: '10'}
    ]
  },
  guide: {
    method: 'GET', path: '/api/v1/content/feed/guide',
    desc: '攻略内容流，支持分类筛选',
    params: [
      {name: 'category', label: '分类', type: 'select', value: '', options: [
        {v:'', l:'全部'}, {v:'nature', l:'自然生态'}, {v:'history', l:'历史文化'},
        {v:'entertainment', l:'人工娱乐'}, {v:'urban', l:'城市公共'},
        {v:'transport', l:'交通枢纽'}, {v:'red_tourism', l:'红色旅游'}, {v:'religion', l:'宗教场所'}
      ]},
      {name: 'city', label: '城市', type: 'input', value: ''},
      {name: 'sort', label: '排序', type: 'select', value: 'newest', options: [
        {v:'newest', l:'最新'}, {v:'popular', l:'最热'}
      ]},
      {name: 'page', label: '页码', type: 'input', value: '1'},
      {name: 'page_size', label: '每页数量', type: 'input', value: '20'}
    ]
  },
  place: {
    method: 'GET', path: '/api/v1/content/feed/place',
    desc: '打卡地内容流，支持分类筛选',
    params: [
      {name: 'category', label: '分类', type: 'select', value: '', options: [
        {v:'', l:'全部'}, {v:'food', l:'吃喝'}, {v:'fun', l:'玩乐'},
        {v:'sightseeing', l:'逛看'}, {v:'outdoor', l:'户外'},
        {v:'shopping', l:'逛街'}, {v:'accommodation', l:'住宿'}
      ]},
      {name: 'city', label: '城市', type: 'input', value: ''},
      {name: 'sort', label: '排序', type: 'select', value: 'newest', options: [
        {v:'newest', l:'最新'}, {v:'popular', l:'最热'}, {v:'rating', l:'评分最高'}
      ]},
      {name: 'page', label: '页码', type: 'input', value: '1'},
      {name: 'page_size', label: '每页数量', type: 'input', value: '20'}
    ]
  },
  route: {
    method: 'GET', path: '/api/v1/content/feed/route',
    desc: '路线内容流，沉浸式整屏滑动',
    params: [
      {name: 'category', label: '分类', type: 'select', value: '', options: [
        {v:'', l:'全部'}, {v:'half_day', l:'半日游'}, {v:'one_day', l:'一日游'},
        {v:'two_day', l:'两日游'}, {v:'three_day_plus', l:'三日及以上'}
      ]},
      {name: 'city', label: '城市', type: 'input', value: ''},
      {name: 'sort', label: '排序', type: 'select', value: 'hot', options: [
        {v:'hot', l:'最热'}, {v:'newest', l:'最新'}, {v:'rating', l:'评分最高'}
      ]},
      {name: 'page', label: '页码', type: 'input', value: '1'},
      {name: 'page_size', label: '每页数量', type: 'input', value: '10'}
    ]
  },
  nearby: {
    method: 'GET', path: '/api/v1/content/feed/nearby',
    desc: '附近流，按距离排序展示打卡地',
    params: [
      {name: 'location_mode', label: '定位模式', type: 'select', value: 'travel', options: [
        {v:'travel', l:'travel (GPS)'}, {v:'plan', l:'plan (规划城市)'}, {v:'manual', l:'manual (手动城市)'}
      ]},
      {name: 'city', label: '城市(plan/manual必填)', type: 'input', value: ''},
      {name: 'longitude', label: '经度(travel必填)', type: 'input', value: '112.2400'},
      {name: 'latitude', label: '纬度(travel必填)', type: 'input', value: '30.3300'},
      {name: 'radius', label: '搜索半径(米)', type: 'input', value: '10000'},
      {name: 'page', label: '页码', type: 'input', value: '1'},
      {name: 'page_size', label: '每页数量', type: 'input', value: '20'}
    ]
  },
  'route-detail': {
    method: 'GET', path: '/api/v1/content/route/{route_id}',
    desc: '路线详情，含时间轴点位列表',
    params: [
      {name: 'route_id', label: '路线ID', type: 'input', value: '2001'}
    ]
  },
  'place-detail': {
    method: 'GET', path: '/api/v1/content/place/{place_id}',
    desc: '地点详情，含营业状态、距离计算',
    params: [
      {name: 'place_id', label: '地点ID', type: 'input', value: '3001'},
      {name: 'longitude', label: '经度(计算距离)', type: 'input', value: '112.2400'},
      {name: 'latitude', label: '纬度(计算距离)', type: 'input', value: '30.3300'}
    ]
  },
  'guide-detail': {
    method: 'GET', path: '/api/v1/content/guide/{guide_id}',
    desc: '攻略详情，含 Markdown 正文',
    params: [
      {name: 'guide_id', label: '攻略ID', type: 'input', value: '4001'}
    ]
  },
  health: {
    method: 'GET', path: '/health',
    desc: '服务健康检查',
    params: []
  }
};

let currentApi = null;

document.querySelectorAll('.api-item').forEach(btn => {
  btn.addEventListener('click', () => {
    document.querySelectorAll('.api-item').forEach(b => b.classList.remove('active'));
    btn.classList.add('active');
    const key = btn.dataset.api;
    renderApi(key);
  });
});

function renderApi(key) {
  currentApi = key;
  const api = APIs[key];
  let pathDisplay = api.path;
  let paramsHtml = '';
  if (api.params.length > 0) {
    paramsHtml = '<div class="params"><div class="params-title">请求参数</div>';
    api.params.forEach(p => {
      paramsHtml += '<div class="param-row"><label>' + p.label + '</label>';
      if (p.type === 'select') {
        paramsHtml += '<select data-param="' + p.name + '">';
        p.options.forEach(o => {
          paramsHtml += '<option value="' + o.v + '"' + (o.v === p.value ? ' selected' : '') + '>' + o.l + '</option>';
        });
        paramsHtml += '</select>';
      } else {
        paramsHtml += '<input type="text" data-param="' + p.name + '" value="' + p.value + '">';
      }
      paramsHtml += '</div>';
    });
    paramsHtml += '</div>';
  }
  document.getElementById('content').innerHTML =
    '<div class="endpoint-info">' +
      '<div class="path"><span class="method ' + api.method + '">' + api.method + '</span>' + api.path + '</div>' +
      '<div class="desc">' + api.desc + '</div>' +
    '</div>' +
    paramsHtml +
    '<button class="btn-send" id="btnSend">发送请求</button>' +
    '<div class="response-section" id="responseSection" style="display:none;">' +
      '<div class="response-meta">' +
        '<span>状态码: <span class="status-badge" id="statusBadge"></span></span>' +
        '<span>耗时: <span id="duration"></span></span>' +
        '<span>Request-ID: <span id="requestId"></span></span>' +
      '</div>' +
      '<div class="json-viewer" id="jsonViewer"></div>' +
    '</div>';
  document.getElementById('btnSend').addEventListener('click', sendRequest);
}

function sendRequest() {
  const api = APIs[currentApi];
  const inputs = document.querySelectorAll('[data-param]');
  let url = api.path;
  let pathParams = {};
  let queryParams = {};
  inputs.forEach(inp => {
    const name = inp.dataset.param;
    const val = inp.value;
    if (api.path.includes('{' + name + '}')) {
      pathParams[name] = val;
    } else {
      if (val) queryParams[name] = val;
    }
  });
  Object.keys(pathParams).forEach(k => {
    url = url.replace('{' + k + '}', pathParams[k]);
  });
  const qs = Object.keys(queryParams).map(k => k + '=' + encodeURIComponent(queryParams[k])).join('&');
  if (qs) url += '?' + qs;

  const btn = document.getElementById('btnSend');
  btn.disabled = true;
  btn.textContent = '请求中...';
  document.getElementById('responseSection').style.display = 'none';

  const start = performance.now();
  fetch(url, { method: api.method })
    .then(r => {
      const duration = (performance.now() - start).toFixed(0);
      const statusBadge = document.getElementById('statusBadge');
      statusBadge.textContent = r.status;
      statusBadge.className = 'status-badge ' + (r.status < 300 ? 'status-2xx' : (r.status < 500 ? 'status-4xx' : 'status-5xx'));
      document.getElementById('duration').textContent = duration + 'ms';
      document.getElementById('requestId').textContent = r.headers.get('X-Request-Id') || '-';
      return r.text().then(t => ({text: t, status: r.status}));
    })
    .then(({text}) => {
      try {
        const json = JSON.parse(text);
        document.getElementById('jsonViewer').innerHTML = syntaxHighlight(JSON.stringify(json, null, 2));
      } catch (e) {
        document.getElementById('jsonViewer').textContent = text;
      }
      document.getElementById('responseSection').style.display = 'block';
    })
    .catch(err => {
      document.getElementById('jsonViewer').textContent = '请求失败: ' + err.message;
      document.getElementById('responseSection').style.display = 'block';
    })
    .finally(() => {
      btn.disabled = false;
      btn.textContent = '发送请求';
    });
}

function syntaxHighlight(json) {
  json = json.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
  return json.replace(/("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g, function (match) {
    let cls = 'json-number';
    if (/^"/.test(match)) {
      if (/:$/.test(match)) { cls = 'json-key'; }
      else { cls = 'json-string'; }
    } else if (/true|false/.test(match)) { cls = 'json-boolean'; }
    else if (/null/.test(match)) { cls = 'json-null'; }
    return '<span class="' + cls + '">' + match + '</span>';
  });
}
</script>
</body>
</html>`

// RegisterDebugRoutes 注册调试页面路由
func RegisterDebugRoutes(r *gin.Engine) {
	r.GET("/debug", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, debugPage)
	})
	// 兼容根路径跳转
	r.GET("/", func(c *gin.Context) {
		if strings.Contains(c.GetHeader("Accept"), "text/html") {
			c.Redirect(http.StatusFound, "/debug")
			return
		}
		c.JSON(http.StatusOK, gin.H{"app": "xunlvji", "version": "v2.1.7", "debug_page": "/debug"})
	})
}
