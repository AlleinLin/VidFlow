import { t, F, u as u$1, e, o as o$1 } from './ssr-D91narGw.js';

const o=[];function c(t,n){return {subscribe:u(t,n).subscribe}}function u(e,s=t){let r;const c=new Set;function u(t){if(F(e,t)&&(e=t,r)){const t=!o.length;for(const n of c)n[1](),o.push(n,e);if(t){for(let t=0;t<o.length;t+=2)o[t][0](o[t+1]);o.length=0;}}}function i(t){u(t(e));}return {set:u,update:i,subscribe:function(n,o=t){const a=[n,o];return c.add(a),1===c.size&&(r=s(u,i)||t),n(e),()=>{c.delete(a),0===c.size&&r&&(r(),r=null);}}}}function i(n,o,u){const i=!Array.isArray(n),a=i?[n]:n;if(!a.every(Boolean))throw new Error("derived() expects stores as input, got a falsy value");const f=o.length<2;return c(u,(n,c)=>{let u=false;const l=[];let d=0,h=t;const p=()=>{if(d)return;h();const e=o(i?l[0]:l,n,c);f?n(e):h=o$1(e)?e:t;},b=a.map((t,n)=>u$1(t,t=>{l[n]=t,d&=~(1<<n),u&&p();},()=>{d|=1<<n;}));return u=true,p(),function(){e(b),h(),u=false;}})}

export { c, i, u };
//# sourceMappingURL=index-BnB10Knz.js.map
