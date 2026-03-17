const e="http://localhost:8080",t={async list(t,o=1,i=20){const a={"Content-Type":"application/json"};t&&(a.Authorization=`Bearer ${t}`);const r=await fetch(`${e}/api/v1/videos?page=${o}&page_size=${i}`,{headers:a});if(!r.ok)throw new Error("Failed to fetch videos");return r.json()},async get(t,o){const i={"Content-Type":"application/json"};o&&(i.Authorization=`Bearer ${o}`);const a=await fetch(`${e}/api/v1/videos/${t}`,{headers:i});if(!a.ok)throw new Error("Failed to fetch video");return a.json()},async search(t,o=1,i=20){const a=await fetch(`${e}/api/v1/videos?keyword=${encodeURIComponent(t)}&page=${o}&page_size=${i}`);if(!a.ok)throw new Error("Search failed");return a.json()},async create(t,o){const i=await fetch(`${e}/api/v1/videos`,{method:"POST",headers:{"Content-Type":"application/json",Authorization:`Bearer ${t}`},body:JSON.stringify(o)});if(!i.ok)throw new Error("Failed to create video");return i.json()},async delete(t,o){if(!(await fetch(`${e}/api/v1/videos/${o}`,{method:"DELETE",headers:{Authorization:`Bearer ${t}`}})).ok)throw new Error("Failed to delete video")},async publish(t,o){if(!(await fetch(`${e}/api/v1/videos/${o}/publish`,{method:"POST",headers:{Authorization:`Bearer ${t}`}})).ok)throw new Error("Failed to publish video")}},o=async({cookies:e})=>{const o=e.get("access_token");try{return {videos:(await t.list(o)).videos,isAuthenticated:!!o}}catch(i){return {videos:[],isAuthenticated:!!o}}};

var _page_server_ts = /*#__PURE__*/Object.freeze({
	__proto__: null,
	load: o
});

const index = 2;
let component_cache;
const component = async () => component_cache ??= (await import('./_page.svelte-D_umLjY0.js')).default;
const server_id = "src/routes/+page.server.ts";
const imports = ["_app/immutable/nodes/2.De7aQ-Di.js","_app/immutable/chunks/S1E9SwF_.js","_app/immutable/chunks/XFbntdAW.js","_app/immutable/chunks/n1todNSJ.js","_app/immutable/chunks/BdJdH1g_.js","_app/immutable/chunks/Cw9zqSKr.js","_app/immutable/chunks/CzxeiFMY.js","_app/immutable/chunks/taLHk_yT.js"];
const stylesheets = [];
const fonts = [];

export { component, fonts, imports, index, _page_server_ts as server, server_id, stylesheets };
//# sourceMappingURL=2-DkhPe2oV.js.map
