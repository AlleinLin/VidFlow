// eslint-disable-next-line n/prefer-global/process
const IN_WEBCONTAINER = !!globalThis.process?.versions?.webcontainer;

/** @import { RequestEvent } from '@sveltejs/kit' */
/** @import { RequestStore } from 'types' */
/** @import { AsyncLocalStorage } from 'node:async_hooks' */


/** @type {RequestStore | null} */
let sync_store = null;

/** @type {AsyncLocalStorage<RequestStore | null> | null} */
let als;

import('node:async_hooks')
	.then((hooks) => (als = new hooks.AsyncLocalStorage()))
	.catch(() => {
		// can't use AsyncLocalStorage, but can still call getRequestEvent synchronously.
		// this isn't behind `supports` because it's basically just StackBlitz (i.e.
		// in-browser usage) that doesn't support it AFAICT
	});

/**
 * @template T
 * @param {RequestStore | null} store
 * @param {() => T} fn
 */
function with_request_store(store, fn) {
	try {
		sync_store = store;
		return als ? als.run(store, fn) : fn();
	} finally {
		// Since AsyncLocalStorage is not working in webcontainers, we don't reset `sync_store`
		// and handle only one request at a time in `src/runtime/server/index.js`.
		if (!IN_WEBCONTAINER) {
			sync_store = null;
		}
	}
}

const t$1=new TextEncoder,r$1=new TextDecoder;function e(t,r){const e=t.split(/[/\\]/),n=r.split(/[/\\]/);for(e.pop();e[0]===n[0];)e.shift(),n.shift();let o=e.length;for(;o--;)e[o]="..";return e.concat(n).join("/")}function n$2(t){if(globalThis.Buffer)return globalThis.Buffer.from(t).toString("base64");let r="";for(let e=0;e<t.length;e++)r+=String.fromCharCode(t[e]);return btoa(r)}function o$2(t){if(globalThis.Buffer){const r=globalThis.Buffer.from(t,"base64");return new Uint8Array(r)}const r=atob(t),e=new Uint8Array(r.length);for(let n=0;n<r.length;n++)e[n]=r.charCodeAt(n);return e}

const n$1=new URL("sveltekit-internal://");function r(e,r){if("/"===r[0]&&"/"===r[1])return r;let t=new URL(e,n$1);return t=new URL(r,t),t.protocol===n$1.protocol?t.pathname+t.search+t.hash:t.href}function t(e,n){return "/"===e||"ignore"===n?e:"never"===n?e.endsWith("/")?e.slice(0,-1):e:"always"!==n||e.endsWith("/")?e:e+"/"}function s(e){return e.split("%25").map(decodeURI).join("%25")}function a(e){for(const n in e)e[n]=decodeURIComponent(e[n]);return e}function o$1(e,n,r,t=false){const s=new URL(e);Object.defineProperty(s,"searchParams",{value:new Proxy(s.searchParams,{get(e,t){if("get"===t||"getAll"===t||"has"===t)return (n,...s)=>(r(n),e[t](n,...s));n();const s=Reflect.get(e,t);return "function"==typeof s?s.bind(e):s}}),enumerable:true,configurable:true});const a=["href","pathname","search","toString","toJSON"];t&&a.push("hash");for(const o of a)Object.defineProperty(s,o,{get:()=>(n(),e[o]),enumerable:true,configurable:true});return s[Symbol.for("nodejs.util.inspect.custom")]=(n,r,t)=>t(e,r),s.searchParams[Symbol.for("nodejs.util.inspect.custom")]=(n,r,t)=>t(e.searchParams,r),t||function(e){c(e),Object.defineProperty(e,"hash",{get(){throw new Error("Cannot access event.url.hash. Consider using `page.url.hash` inside a component instead")}});}(s),s}function i(e){c(e);for(const n of ["search","searchParams"])Object.defineProperty(e,n,{get(){throw new Error(`Cannot access url.${n} on a page with prerendering enabled`)}});}function c(e){e[Symbol.for("nodejs.util.inspect.custom")]=(n,r,t)=>t(new URL(e),r);}function h(e){return function(n,r){if(n)for(const t in n){if("_"===t[0]||e.has(t))continue;const n=[...e.values()],s=u(t,r?.slice(r.lastIndexOf(".")))??`valid exports are ${n.join(", ")}, or anything with a '_' prefix`;throw new Error(`Invalid export '${t}'${r?` in ${r}`:""} (${s})`)}}}function u(e,n=".js"){const r=[];if(l.has(e)&&r.push(`+layout${n}`),f.has(e)&&r.push(`+page${n}`),p.has(e)&&r.push(`+layout.server${n}`),d.has(e)&&r.push(`+page.server${n}`),g.has(e)&&r.push(`+server${n}`),r.length>0)return `'${e}' is a valid export in ${r.slice(0,-1).join(", ")}${r.length>1?" or ":""}${r.at(-1)}`}const l=new Set(["load","prerender","csr","ssr","trailingSlash","config"]),f=new Set([...l,"entries"]),p=new Set([...l]),d=new Set([...p,"actions","entries"]),g=new Set(["GET","POST","PATCH","PUT","DELETE","OPTIONS","HEAD","fallback","prerender","trailingSlash","config","entries"]),w=h(l),m=h(f),b=h(p),S=h(d);

function n(){}function o(){}

export { S, o as a, t$1 as b, a as c, b as d, w as e, o$1 as f, r as g, e as h, i, n as j, m, n$2 as n, o$2 as o, r$1 as r, s, t, with_request_store as w };
//# sourceMappingURL=ssr2-BIEsMiuD.js.map
