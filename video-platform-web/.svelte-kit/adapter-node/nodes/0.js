import * as server from '../entries/pages/_layout.server.ts.js';

export const index = 0;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/_layout.svelte.js')).default;
export { server };
export const server_id = "src/routes/+layout.server.ts";
export const imports = ["_app/immutable/nodes/0.C_IbWXFP.js","_app/immutable/chunks/S1E9SwF_.js","_app/immutable/chunks/XFbntdAW.js","_app/immutable/chunks/Cw9zqSKr.js","_app/immutable/chunks/CzxeiFMY.js","_app/immutable/chunks/BdJdH1g_.js"];
export const stylesheets = ["_app/immutable/assets/0.BMr8-D9T.css"];
export const fonts = [];
