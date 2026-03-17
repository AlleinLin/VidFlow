const manifest = (() => {
function __memo(fn) {
	let value;
	return () => value ??= (value = fn());
}

return {
	appDir: "_app",
	appPath: "_app",
	assets: new Set([]),
	mimeTypes: {},
	_: {
		client: {start:"_app/immutable/entry/start.D8ALiy_m.js",app:"_app/immutable/entry/app.DF2JtMHR.js",imports:["_app/immutable/entry/start.D8ALiy_m.js","_app/immutable/chunks/taLHk_yT.js","_app/immutable/chunks/S1E9SwF_.js","_app/immutable/chunks/CzxeiFMY.js","_app/immutable/entry/app.DF2JtMHR.js","_app/immutable/chunks/S1E9SwF_.js","_app/immutable/chunks/XFbntdAW.js"],stylesheets:[],fonts:[],uses_env_dynamic_public:false},
		nodes: [
			__memo(() => import('./chunks/0-COEbYt4r.js')),
			__memo(() => import('./chunks/1-1RgFhSJU.js')),
			__memo(() => import('./chunks/2-DkhPe2oV.js')),
			__memo(() => import('./chunks/3-CTHVZN85.js')),
			__memo(() => import('./chunks/4-DaoaOCPR.js'))
		],
		remotes: {
			
		},
		routes: [
			{
				id: "/",
				pattern: /^\/$/,
				params: [],
				page: { layouts: [0,], errors: [1,], leaf: 2 },
				endpoint: null
			},
			{
				id: "/login",
				pattern: /^\/login\/?$/,
				params: [],
				page: { layouts: [0,], errors: [1,], leaf: 3 },
				endpoint: null
			},
			{
				id: "/register",
				pattern: /^\/register\/?$/,
				params: [],
				page: { layouts: [0,], errors: [1,], leaf: 4 },
				endpoint: null
			}
		],
		prerendered_routes: new Set([]),
		matchers: async () => {
			
			return {  };
		},
		server_assets: {}
	}
}
})();

const prerendered = new Set([]);

const base = "";

export { base, manifest, prerendered };
//# sourceMappingURL=manifest.js.map
