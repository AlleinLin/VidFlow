import { x, g, j, r, y, E, h as h$1 } from './ssr-D91narGw.js';
import { m as m$1, u as u$1 } from './stores-BgOinkfD.js';
import { t } from './utils2-C_vlxP_i.js';
import './index-BnB10Knz.js';
import 'clsx';
import 'tailwind-merge';

function n(e,t){const s=document.documentElement,o=localStorage.getItem("mode-watcher-mode")||e,a="light"===o||"system"===o&&window.matchMedia("(prefers-color-scheme: light)").matches;if(s.classList[a?"remove":"add"]("dark"),s.style.colorScheme=a?"light":"dark",t){const e=document.querySelector('meta[name="theme-color"]');e&&e.setAttribute("content","light"===o?t.light:t.dark);}localStorage.setItem("mode-watcher-mode",o);}const m=x((e,s,o,a)=>{let{track:r=true}=s,{defaultMode:c="system"}=s,{themeColors:l}=s,{disableTransitions:m=true}=s;m$1.set(l),u$1.set(m);const h=`"${c}"${l?`, ${JSON.stringify(l)}`:""}`;return void 0===s.track&&o.track&&void 0!==r&&o.track(r),void 0===s.defaultMode&&o.defaultMode&&void 0!==c&&o.defaultMode(c),void 0===s.themeColors&&o.themeColors&&void 0!==l&&o.themeColors(l),void 0===s.disableTransitions&&o.disableTransitions&&void 0!==m&&o.disableTransitions(m),""+(e.head+=`\x3c!-- HEAD_svelte-cpyj77_START --\x3e${l?`   <meta name="theme-color"${j("content",l.dark,0)}>`:""}\x3c!-- HTML_TAG_START --\x3e${'<script nonce="%sveltekit.nonce%">('+n.toString()+")("+h+");<\/script>"}\x3c!-- HTML_TAG_END --\x3e\x3c!-- HEAD_svelte-cpyj77_END --\x3e`,"")}),h=x((e,t$1,c,d)=>{let i,n=r(t$1,["className"]),{className:m=""}=t$1;return void 0===t$1.className&&c.className&&void 0!==m&&c.className(m),i=t("toast-container",m),`<div${y([{class:E(i)},h$1(n)],{})}>${d.default?d.default({}):""}</div>`}),u=x((e,t,s,o)=>`${g(m,"ModeWatcher").$$render(e,{},{},{})} ${g(h,"Toaster").$$render(e,{},{},{})} ${o.default?o.default({}):""}`);

export { u as default };
//# sourceMappingURL=_layout.svelte-CTWW-JTC.js.map
