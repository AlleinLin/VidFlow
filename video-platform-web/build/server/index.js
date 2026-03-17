import { w as with_request_store, r as r$1, o as o$1, a as o$2, s as s$2, t as t$1, i as i$1, b as t$2, c as a$2, d as b$1, e as w$2, S as S$1, m as m$2, f as o$3, g as r$2, h as e$1, n as n$1 } from './chunks/ssr2-BIEsMiuD.js';
import { c as c$2, u as u$1 } from './chunks/index-BnB10Knz.js';
import { x, a as a$1, g as g$1, w as w$1 } from './chunks/ssr-D91narGw.js';

/** @import { StandardSchemaV1 } from '@standard-schema/spec' */

class HttpError {
	/**
	 * @param {number} status
	 * @param {{message: string} extends App.Error ? (App.Error | string | undefined) : App.Error} body
	 */
	constructor(status, body) {
		this.status = status;
		if (typeof body === 'string') {
			this.body = { message: body };
		} else if (body) {
			this.body = body;
		} else {
			this.body = { message: `Error: ${status}` };
		}
	}

	toString() {
		return JSON.stringify(this.body);
	}
}

class Redirect {
	/**
	 * @param {300 | 301 | 302 | 303 | 304 | 305 | 306 | 307 | 308} status
	 * @param {string} location
	 */
	constructor(status, location) {
		this.status = status;
		this.location = location;
	}
}

/**
 * An error that was thrown from within the SvelteKit runtime that is not fatal and doesn't result in a 500, such as a 404.
 * `SvelteKitError` goes through `handleError`.
 * @extends Error
 */
class SvelteKitError extends Error {
	/**
	 * @param {number} status
	 * @param {string} text
	 * @param {string} message
	 */
	constructor(status, text, message) {
		super(message);
		this.status = status;
		this.text = text;
	}
}

/**
 * @template [T=undefined]
 */
class ActionFailure {
	/**
	 * @param {number} status
	 * @param {T} data
	 */
	constructor(status, data) {
		this.status = status;
		this.data = data;
	}
}

const text_encoder = new TextEncoder();
new TextDecoder();

/** @import { StandardSchemaV1 } from '@standard-schema/spec' */


// TODO 3.0: remove these types as they are not used anymore (we can't remove them yet because that would be a breaking change)
/**
 * @template {number} TNumber
 * @template {any[]} [TArray=[]]
 * @typedef {TNumber extends TArray['length'] ? TArray[number] : LessThan<TNumber, [...TArray, TArray['length']]>} LessThan
 */

/**
 * @template {number} TStart
 * @template {number} TEnd
 * @typedef {Exclude<TEnd | LessThan<TEnd>, LessThan<TStart>>} NumericRange
 */

// Keep the status codes as `number` because restricting to certain numbers makes it unnecessarily hard to use compared to the benefits
// (we have runtime errors already to check for invalid codes). Also see https://github.com/sveltejs/kit/issues/11780

// we have to repeat the JSDoc because the display for function overloads is broken
// see https://github.com/microsoft/TypeScript/issues/55056

/**
 * Throws an error with a HTTP status code and an optional message.
 * When called during request handling, this will cause SvelteKit to
 * return an error response without invoking `handleError`.
 * Make sure you're not catching the thrown error, which would prevent SvelteKit from handling it.
 * @param {number} status The [HTTP status code](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status#client_error_responses). Must be in the range 400-599.
 * @param {App.Error} body An object that conforms to the App.Error type. If a string is passed, it will be used as the message property.
 * @overload
 * @param {number} status
 * @param {App.Error} body
 * @return {never}
 * @throws {HttpError} This error instructs SvelteKit to initiate HTTP error handling.
 * @throws {Error} If the provided status is invalid (not between 400 and 599).
 */
/**
 * Throws an error with a HTTP status code and an optional message.
 * When called during request handling, this will cause SvelteKit to
 * return an error response without invoking `handleError`.
 * Make sure you're not catching the thrown error, which would prevent SvelteKit from handling it.
 * @param {number} status The [HTTP status code](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status#client_error_responses). Must be in the range 400-599.
 * @param {{ message: string } extends App.Error ? App.Error | string | undefined : never} [body] An object that conforms to the App.Error type. If a string is passed, it will be used as the message property.
 * @overload
 * @param {number} status
 * @param {{ message: string } extends App.Error ? App.Error | string | undefined : never} [body]
 * @return {never}
 * @throws {HttpError} This error instructs SvelteKit to initiate HTTP error handling.
 * @throws {Error} If the provided status is invalid (not between 400 and 599).
 */
/**
 * Throws an error with a HTTP status code and an optional message.
 * When called during request handling, this will cause SvelteKit to
 * return an error response without invoking `handleError`.
 * Make sure you're not catching the thrown error, which would prevent SvelteKit from handling it.
 * @param {number} status The [HTTP status code](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status#client_error_responses). Must be in the range 400-599.
 * @param {{ message: string } extends App.Error ? App.Error | string | undefined : never} body An object that conforms to the App.Error type. If a string is passed, it will be used as the message property.
 * @return {never}
 * @throws {HttpError} This error instructs SvelteKit to initiate HTTP error handling.
 * @throws {Error} If the provided status is invalid (not between 400 and 599).
 */
function error(status, body) {
	if ((isNaN(status) || status < 400 || status > 599)) {
		throw new Error(`HTTP error status codes must be between 400 and 599 — ${status} is invalid`);
	}

	throw new HttpError(status, body);
}

/**
 * Create a JSON `Response` object from the supplied data.
 * @param {any} data The value that will be serialized as JSON.
 * @param {ResponseInit} [init] Options such as `status` and `headers` that will be added to the response. `Content-Type: application/json` and `Content-Length` headers will be added automatically.
 */
function json(data, init) {
	// TODO deprecate this in favour of `Response.json` when it's
	// more widely supported
	const body = JSON.stringify(data);

	// we can't just do `text(JSON.stringify(data), init)` because
	// it will set a default `content-type` header. duplicated code
	// means less duplicated work
	const headers = new Headers(init?.headers);
	if (!headers.has('content-length')) {
		headers.set('content-length', text_encoder.encode(body).byteLength.toString());
	}

	if (!headers.has('content-type')) {
		headers.set('content-type', 'application/json');
	}

	return new Response(body, {
		...init,
		headers
	});
}

/**
 * Create a `Response` object from the supplied body.
 * @param {string} body The value that will be used as-is.
 * @param {ResponseInit} [init] Options such as `status` and `headers` that will be added to the response. A `Content-Length` header will be added automatically.
 */
function text(body, init) {
	const headers = new Headers(init?.headers);
	if (!headers.has('content-length')) {
		const encoded = text_encoder.encode(body);
		headers.set('content-length', encoded.byteLength.toString());
		return new Response(encoded, {
			...init,
			headers
		});
	}

	return new Response(body, {
		...init,
		headers
	});
}

/**
 * @template {{ tracing: { enabled: boolean, root: import('@opentelemetry/api').Span, current: import('@opentelemetry/api').Span } }} T
 * @param {T} event_like
 * @param {import('@opentelemetry/api').Span} current
 * @returns {T}
 */
function merge_tracing(event_like, current) {
	return {
		...event_like,
		tracing: {
			...event_like.tracing,
			current
		}
	};
}

/** @type {Record<string, string>} */
const escaped = {
	'<': '\\u003C',
	'\\': '\\\\',
	'\b': '\\b',
	'\f': '\\f',
	'\n': '\\n',
	'\r': '\\r',
	'\t': '\\t',
	'\u2028': '\\u2028',
	'\u2029': '\\u2029'
};

class DevalueError extends Error {
	/**
	 * @param {string} message
	 * @param {string[]} keys
	 * @param {any} [value] - The value that failed to be serialized
	 * @param {any} [root] - The root value being serialized
	 */
	constructor(message, keys, value, root) {
		super(message);
		this.name = 'DevalueError';
		this.path = keys.join('');
		this.value = value;
		this.root = root;
	}
}

/** @param {any} thing */
function is_primitive(thing) {
	return Object(thing) !== thing;
}

const object_proto_names = /* @__PURE__ */ Object.getOwnPropertyNames(
	Object.prototype
)
	.sort()
	.join('\0');

/** @param {any} thing */
function is_plain_object(thing) {
	const proto = Object.getPrototypeOf(thing);

	return (
		proto === Object.prototype ||
		proto === null ||
		Object.getPrototypeOf(proto) === null ||
		Object.getOwnPropertyNames(proto).sort().join('\0') === object_proto_names
	);
}

/** @param {any} thing */
function get_type(thing) {
	return Object.prototype.toString.call(thing).slice(8, -1);
}

/** @param {string} char */
function get_escaped_char(char) {
	switch (char) {
		case '"':
			return '\\"';
		case '<':
			return '\\u003C';
		case '\\':
			return '\\\\';
		case '\n':
			return '\\n';
		case '\r':
			return '\\r';
		case '\t':
			return '\\t';
		case '\b':
			return '\\b';
		case '\f':
			return '\\f';
		case '\u2028':
			return '\\u2028';
		case '\u2029':
			return '\\u2029';
		default:
			return char < ' '
				? `\\u${char.charCodeAt(0).toString(16).padStart(4, '0')}`
				: '';
	}
}

/** @param {string} str */
function stringify_string(str) {
	let result = '';
	let last_pos = 0;
	const len = str.length;

	for (let i = 0; i < len; i += 1) {
		const char = str[i];
		const replacement = get_escaped_char(char);
		if (replacement) {
			result += str.slice(last_pos, i) + replacement;
			last_pos = i + 1;
		}
	}

	return `"${last_pos === 0 ? str : result + str.slice(last_pos)}"`;
}

/** @param {Record<string | symbol, any>} object */
function enumerable_symbols(object) {
	return Object.getOwnPropertySymbols(object).filter(
		(symbol) => Object.getOwnPropertyDescriptor(object, symbol).enumerable
	);
}

const is_identifier = /^[a-zA-Z_$][a-zA-Z_$0-9]*$/;

/** @param {string} key */
function stringify_key(key) {
	return is_identifier.test(key) ? '.' + key : '[' + JSON.stringify(key) + ']';
}

/** @param {string} s */
function is_valid_array_index(s) {
	if (s.length === 0) return false;
	if (s.length > 1 && s.charCodeAt(0) === 48) return false; // leading zero
	for (let i = 0; i < s.length; i++) {
		const c = s.charCodeAt(i);
		if (c < 48 || c > 57) return false;
	}
	// by this point we know it's a string of digits, but it has to be within the range of valid array indices
	const n = +s;
	if (n >= 2 ** 32 - 1) return false;
	if (n < 0) return false;
	return true;
}

/**
 * Finds the populated indices of an array.
 * @param {unknown[]} array
 */
function valid_array_indices(array) {
	const keys = Object.keys(array);
	for (var i = keys.length - 1; i >= 0; i--) {
		if (is_valid_array_index(keys[i])) {
			break;
		}
	}
	keys.length = i + 1;
	return keys;
}

const chars = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_$';
const unsafe_chars = /[<\b\f\n\r\t\0\u2028\u2029]/g;
const reserved =
	/^(?:do|if|in|for|int|let|new|try|var|byte|case|char|else|enum|goto|long|this|void|with|await|break|catch|class|const|final|float|short|super|throw|while|yield|delete|double|export|import|native|return|switch|throws|typeof|boolean|default|extends|finally|package|private|abstract|continue|debugger|function|volatile|interface|protected|transient|implements|instanceof|synchronized)$/;

/**
 * Turn a value into the JavaScript that creates an equivalent value
 * @param {any} value
 * @param {(value: any, uneval: (value: any) => string) => string | void} [replacer]
 */
function uneval(value, replacer) {
	const counts = new Map();

	/** @type {string[]} */
	const keys = [];

	const custom = new Map();

	/** @param {any} thing */
	function walk(thing) {
		if (!is_primitive(thing)) {
			if (counts.has(thing)) {
				counts.set(thing, counts.get(thing) + 1);
				return;
			}

			counts.set(thing, 1);

			if (replacer) {
				const str = replacer(thing, (value) => uneval(value, replacer));

				if (typeof str === 'string') {
					custom.set(thing, str);
					return;
				}
			}

			if (typeof thing === 'function') {
				throw new DevalueError(`Cannot stringify a function`, keys, thing, value);
			}

			const type = get_type(thing);

			switch (type) {
				case 'Number':
				case 'BigInt':
				case 'String':
				case 'Boolean':
				case 'Date':
				case 'RegExp':
				case 'URL':
				case 'URLSearchParams':
					return;

				case 'Array':
					/** @type {any[]} */ (thing).forEach((value, i) => {
						keys.push(`[${i}]`);
						walk(value);
						keys.pop();
					});
					break;

				case 'Set':
					Array.from(thing).forEach(walk);
					break;

				case 'Map':
					for (const [key, value] of thing) {
						keys.push(
							`.get(${is_primitive(key) ? stringify_primitive$1(key) : '...'})`
						);
						walk(value);
						keys.pop();
					}
					break;

				case 'Int8Array':
				case 'Uint8Array':
				case 'Uint8ClampedArray':
				case 'Int16Array':
				case 'Uint16Array':
				case 'Int32Array':
				case 'Uint32Array':
				case 'Float32Array':
				case 'Float64Array':
				case 'BigInt64Array':
				case 'BigUint64Array':
					walk(thing.buffer);
					return;

				case 'ArrayBuffer':
					return;

				case 'Temporal.Duration':
				case 'Temporal.Instant':
				case 'Temporal.PlainDate':
				case 'Temporal.PlainTime':
				case 'Temporal.PlainDateTime':
				case 'Temporal.PlainMonthDay':
				case 'Temporal.PlainYearMonth':
				case 'Temporal.ZonedDateTime':
					return;

				default:
					if (!is_plain_object(thing)) {
						throw new DevalueError(
							`Cannot stringify arbitrary non-POJOs`,
							keys,
							thing,
							value
						);
					}

					if (enumerable_symbols(thing).length > 0) {
						throw new DevalueError(
							`Cannot stringify POJOs with symbolic keys`,
							keys,
							thing,
							value
						);
					}

					for (const key of Object.keys(thing)) {
						if (key === '__proto__') {
							throw new DevalueError(
								`Cannot stringify objects with __proto__ keys`,
								keys,
								thing,
								value
							);
						}

						keys.push(stringify_key(key));
						walk(thing[key]);
						keys.pop();
					}
			}
		}
	}

	walk(value);

	const names = new Map();

	Array.from(counts)
		.filter((entry) => entry[1] > 1)
		.sort((a, b) => b[1] - a[1])
		.forEach((entry, i) => {
			names.set(entry[0], get_name(i));
		});

	/**
	 * @param {any} thing
	 * @returns {string}
	 */
	function stringify(thing) {
		if (names.has(thing)) {
			return names.get(thing);
		}

		if (is_primitive(thing)) {
			return stringify_primitive$1(thing);
		}

		if (custom.has(thing)) {
			return custom.get(thing);
		}

		const type = get_type(thing);

		switch (type) {
			case 'Number':
			case 'String':
			case 'Boolean':
				return `Object(${stringify(thing.valueOf())})`;

			case 'RegExp':
				return `new RegExp(${stringify_string(thing.source)}, "${
					thing.flags
				}")`;

			case 'Date':
				return `new Date(${thing.getTime()})`;

			case 'URL':
				return `new URL(${stringify_string(thing.toString())})`;

			case 'URLSearchParams':
				return `new URLSearchParams(${stringify_string(thing.toString())})`;

			case 'Array': {
				// For dense arrays (no holes), we iterate normally.
				// When we encounter the first hole, we call Object.keys
				// to determine the sparseness, then decide between:
				//   - Array literal with holes: [,"a",,] (default)
				//   - Object.assign: Object.assign(Array(n),{...}) (for very sparse arrays)
				// Only the Object.assign path avoids iterating every slot, which
				// is what protects against the DoS of e.g. `arr[1000000] = 1`.
				let has_holes = false;

				let result = '[';

				for (let i = 0; i < thing.length; i += 1) {
					if (i > 0) result += ',';

					if (Object.hasOwn(thing, i)) {
						result += stringify(thing[i]);
					} else if (!has_holes) {
						// Decide between array literal and Object.assign.
						//
						// Array literal: holes are consecutive commas.
						// For example, [, "a", ,] is written as [,"a",,].
						// Each hole costs 1 char (a comma).
						//
						// Object.assign: populated indices are listed explicitly.
						// For example, [, "a", ,] would be written as
						// Object.assign(Array(3),{1:"a"}). This avoids paying
						// per-hole, but has a large fixed overhead for the
						// "Object.assign(Array(n),{...})" wrapper, and each
						// element costs extra chars for its index and colon.
						//
						// The serialized values are the same size either way, so
						// the choice comes down to the structural overhead:
						//
						//   Array literal overhead:
						//     1 char per element or hole (comma separators)
						//     + 2 chars for "[" and "]"
						//     = L + 2
						//
						//   Object.assign overhead:
						//     "Object.assign(Array(" — 20 chars
						//     + length              — d chars
						//     + "),{"               — 3 chars
						//     + for each populated element:
						//       index + ":" + ","   — (d + 2) chars
						//     + "})"                — 2 chars
						//     = (25 + d) + P * (d + 2)
						//
						// where L is the array length, P is the number of
						// populated elements, and d is the number of digits
						// in L (an upper bound on the digits in any index).
						//
						// Object.assign is cheaper when:
						//   (25 + d) + P * (d + 2) < L + 2
						const populated_keys = valid_array_indices(/** @type {any[]} */ (thing));
						const population = populated_keys.length;
						const d = String(thing.length).length;

						const hole_cost = thing.length + 2;
						const sparse_cost = (25 + d) + population * (d + 2);

						if (hole_cost > sparse_cost) {
							const entries = populated_keys
								.map((k) => `${k}:${stringify(thing[k])}`)
								.join(',');
							return `Object.assign(Array(${thing.length}),{${entries}})`;
						}

						// Re-process this index as a hole in the array literal
						has_holes = true;
						i -= 1;
					}
					// else: already decided on array literal, hole is just an empty slot
					// (the comma separator is all we need — no content for this position)
				}

				const tail = thing.length === 0 || thing.length - 1 in thing ? '' : ',';
				return result + tail + ']';
			}

			case 'Set':
			case 'Map':
				return `new ${type}([${Array.from(thing).map(stringify).join(',')}])`;

			case 'Int8Array':
			case 'Uint8Array':
			case 'Uint8ClampedArray':
			case 'Int16Array':
			case 'Uint16Array':
			case 'Int32Array':
			case 'Uint32Array':
			case 'Float32Array':
			case 'Float64Array':
			case 'BigInt64Array':
			case 'BigUint64Array': {
				let str = `new ${type}`;

				if (counts.get(thing.buffer) === 1) {
					const array = new thing.constructor(thing.buffer);
					str += `([${array}])`;
				} else {
					str += `([${stringify(thing.buffer)}])`;
				}

				const a = thing.byteOffset;
				const b = a + thing.byteLength;

				// handle subarrays
				if (a > 0 || b !== thing.buffer.byteLength) {
					const m = +/(\d+)/.exec(type)[1] / 8;
					str += `.subarray(${a / m},${b / m})`;
				}

				return str;
			}

			case 'ArrayBuffer': {
				const ui8 = new Uint8Array(thing);
				return `new Uint8Array([${ui8.toString()}]).buffer`;
			}

			case 'Temporal.Duration':
			case 'Temporal.Instant':
			case 'Temporal.PlainDate':
			case 'Temporal.PlainTime':
			case 'Temporal.PlainDateTime':
			case 'Temporal.PlainMonthDay':
			case 'Temporal.PlainYearMonth':
			case 'Temporal.ZonedDateTime':
				return `${type}.from(${stringify_string(thing.toString())})`;

			default:
				const keys = Object.keys(thing);
				const obj = keys
					.map((key) => `${safe_key(key)}:${stringify(thing[key])}`)
					.join(',');
				const proto = Object.getPrototypeOf(thing);
				if (proto === null) {
					return keys.length > 0
						? `{${obj},__proto__:null}`
						: `{__proto__:null}`;
				}

				return `{${obj}}`;
		}
	}

	const str = stringify(value);

	if (names.size) {
		/** @type {string[]} */
		const params = [];

		/** @type {string[]} */
		const statements = [];

		/** @type {string[]} */
		const values = [];

		names.forEach((name, thing) => {
			params.push(name);

			if (custom.has(thing)) {
				values.push(/** @type {string} */ (custom.get(thing)));
				return;
			}

			if (is_primitive(thing)) {
				values.push(stringify_primitive$1(thing));
				return;
			}

			const type = get_type(thing);

			switch (type) {
				case 'Number':
				case 'String':
				case 'Boolean':
					values.push(`Object(${stringify(thing.valueOf())})`);
					break;

				case 'RegExp':
					values.push(thing.toString());
					break;

				case 'Date':
					values.push(`new Date(${thing.getTime()})`);
					break;

				case 'Array':
					values.push(`Array(${thing.length})`);
					/** @type {any[]} */ (thing).forEach((v, i) => {
						statements.push(`${name}[${i}]=${stringify(v)}`);
					});
					break;

				case 'Set':
					values.push(`new Set`);
					statements.push(
						`${name}.${Array.from(thing)
							.map((v) => `add(${stringify(v)})`)
							.join('.')}`
					);
					break;

				case 'Map':
					values.push(`new Map`);
					statements.push(
						`${name}.${Array.from(thing)
							.map(([k, v]) => `set(${stringify(k)}, ${stringify(v)})`)
							.join('.')}`
					);
					break;

				case 'ArrayBuffer':
					values.push(
						`new Uint8Array([${new Uint8Array(thing).join(',')}]).buffer`
					);
					break;

				default:
					values.push(
						Object.getPrototypeOf(thing) === null ? 'Object.create(null)' : '{}'
					);
					Object.keys(thing).forEach((key) => {
						statements.push(
							`${name}${safe_prop(key)}=${stringify(thing[key])}`
						);
					});
			}
		});

		statements.push(`return ${str}`);

		return `(function(${params.join(',')}){${statements.join(
			';'
		)}}(${values.join(',')}))`;
	} else {
		return str;
	}
}

/** @param {number} num */
function get_name(num) {
	let name = '';

	do {
		name = chars[num % chars.length] + name;
		num = ~~(num / chars.length) - 1;
	} while (num >= 0);

	return reserved.test(name) ? `${name}0` : name;
}

/** @param {string} c */
function escape_unsafe_char(c) {
	return escaped[c] || c;
}

/** @param {string} str */
function escape_unsafe_chars(str) {
	return str.replace(unsafe_chars, escape_unsafe_char);
}

/** @param {string} key */
function safe_key(key) {
	return /^[_$a-zA-Z][_$a-zA-Z0-9]*$/.test(key)
		? key
		: escape_unsafe_chars(JSON.stringify(key));
}

/** @param {string} key */
function safe_prop(key) {
	return /^[_$a-zA-Z][_$a-zA-Z0-9]*$/.test(key)
		? `.${key}`
		: `[${escape_unsafe_chars(JSON.stringify(key))}]`;
}

/** @param {any} thing */
function stringify_primitive$1(thing) {
	if (typeof thing === 'string') return stringify_string(thing);
	if (thing === void 0) return 'void 0';
	if (thing === 0 && 1 / thing < 0) return '-0';
	const str = String(thing);
	if (typeof thing === 'number') return str.replace(/^(-)?0\./, '$1.');
	if (typeof thing === 'bigint') return thing + 'n';
	return str;
}

/**
 * Base64 Encodes an arraybuffer
 * @param {ArrayBuffer} arraybuffer
 * @returns {string}
 */
function encode64(arraybuffer) {
  const dv = new DataView(arraybuffer);
  let binaryString = "";

  for (let i = 0; i < arraybuffer.byteLength; i++) {
    binaryString += String.fromCharCode(dv.getUint8(i));
  }

  return binaryToAscii(binaryString);
}

/**
 * Decodes a base64 string into an arraybuffer
 * @param {string} string
 * @returns {ArrayBuffer}
 */
function decode64(string) {
  const binaryString = asciiToBinary(string);
  const arraybuffer = new ArrayBuffer(binaryString.length);
  const dv = new DataView(arraybuffer);

  for (let i = 0; i < arraybuffer.byteLength; i++) {
    dv.setUint8(i, binaryString.charCodeAt(i));
  }

  return arraybuffer;
}

const KEY_STRING =
  "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";

/**
 * Substitute for atob since it's deprecated in node.
 * Does not do any input validation.
 *
 * @see https://github.com/jsdom/abab/blob/master/lib/atob.js
 *
 * @param {string} data
 * @returns {string}
 */
function asciiToBinary(data) {
  if (data.length % 4 === 0) {
    data = data.replace(/==?$/, "");
  }

  let output = "";
  let buffer = 0;
  let accumulatedBits = 0;

  for (let i = 0; i < data.length; i++) {
    buffer <<= 6;
    buffer |= KEY_STRING.indexOf(data[i]);
    accumulatedBits += 6;
    if (accumulatedBits === 24) {
      output += String.fromCharCode((buffer & 0xff0000) >> 16);
      output += String.fromCharCode((buffer & 0xff00) >> 8);
      output += String.fromCharCode(buffer & 0xff);
      buffer = accumulatedBits = 0;
    }
  }
  if (accumulatedBits === 12) {
    buffer >>= 4;
    output += String.fromCharCode(buffer);
  } else if (accumulatedBits === 18) {
    buffer >>= 2;
    output += String.fromCharCode((buffer & 0xff00) >> 8);
    output += String.fromCharCode(buffer & 0xff);
  }
  return output;
}

/**
 * Substitute for btoa since it's deprecated in node.
 * Does not do any input validation.
 *
 * @see https://github.com/jsdom/abab/blob/master/lib/btoa.js
 *
 * @param {string} str
 * @returns {string}
 */
function binaryToAscii(str) {
  let out = "";
  for (let i = 0; i < str.length; i += 3) {
    /** @type {[number, number, number, number]} */
    const groupsOfSix = [undefined, undefined, undefined, undefined];
    groupsOfSix[0] = str.charCodeAt(i) >> 2;
    groupsOfSix[1] = (str.charCodeAt(i) & 0x03) << 4;
    if (str.length > i + 1) {
      groupsOfSix[1] |= str.charCodeAt(i + 1) >> 4;
      groupsOfSix[2] = (str.charCodeAt(i + 1) & 0x0f) << 2;
    }
    if (str.length > i + 2) {
      groupsOfSix[2] |= str.charCodeAt(i + 2) >> 6;
      groupsOfSix[3] = str.charCodeAt(i + 2) & 0x3f;
    }
    for (let j = 0; j < groupsOfSix.length; j++) {
      if (typeof groupsOfSix[j] === "undefined") {
        out += "=";
      } else {
        out += KEY_STRING[groupsOfSix[j]];
      }
    }
  }
  return out;
}

const UNDEFINED = -1;
const HOLE = -2;
const NAN = -3;
const POSITIVE_INFINITY = -4;
const NEGATIVE_INFINITY = -5;
const NEGATIVE_ZERO = -6;
const SPARSE = -7;

/**
 * Revive a value serialized with `devalue.stringify`
 * @param {string} serialized
 * @param {Record<string, (value: any) => any>} [revivers]
 */
function parse(serialized, revivers) {
	return unflatten(JSON.parse(serialized), revivers);
}

/**
 * Revive a value flattened with `devalue.stringify`
 * @param {number | any[]} parsed
 * @param {Record<string, (value: any) => any>} [revivers]
 */
function unflatten(parsed, revivers) {
	if (typeof parsed === 'number') return hydrate(parsed, true);

	if (!Array.isArray(parsed) || parsed.length === 0) {
		throw new Error('Invalid input');
	}

	const values = /** @type {any[]} */ (parsed);

	const hydrated = Array(values.length);

	/**
	 * A set of values currently being hydrated with custom revivers,
	 * used to detect invalid cyclical dependencies
	 * @type {Set<number> | null}
	 */
	let hydrating = null;

	/**
	 * @param {number} index
	 * @returns {any}
	 */
	function hydrate(index, standalone = false) {
		if (index === UNDEFINED) return undefined;
		if (index === NAN) return NaN;
		if (index === POSITIVE_INFINITY) return Infinity;
		if (index === NEGATIVE_INFINITY) return -Infinity;
		if (index === NEGATIVE_ZERO) return -0;

		if (standalone || typeof index !== 'number') {
			throw new Error(`Invalid input`);
		}

		if (index in hydrated) return hydrated[index];

		const value = values[index];

		if (!value || typeof value !== 'object') {
			hydrated[index] = value;
		} else if (Array.isArray(value)) {
			if (typeof value[0] === 'string') {
				const type = value[0];

				const reviver =
					revivers && Object.hasOwn(revivers, type)
						? revivers[type]
						: undefined;

				if (reviver) {
					let i = value[1];
					if (typeof i !== 'number') {
						// if it's not a number, it was serialized by a builtin reviver
						// so we need to munge it into the format expected by a custom reviver
						i = values.push(value[1]) - 1;
					}

					hydrating ??= new Set();

					if (hydrating.has(i)) {
						throw new Error('Invalid circular reference');
					}

					hydrating.add(i);
					hydrated[index] = reviver(hydrate(i));
					hydrating.delete(i);

					return hydrated[index];
				}

				switch (type) {
					case 'Date':
						hydrated[index] = new Date(value[1]);
						break;

					case 'Set':
						const set = new Set();
						hydrated[index] = set;
						for (let i = 1; i < value.length; i += 1) {
							set.add(hydrate(value[i]));
						}
						break;

					case 'Map':
						const map = new Map();
						hydrated[index] = map;
						for (let i = 1; i < value.length; i += 2) {
							map.set(hydrate(value[i]), hydrate(value[i + 1]));
						}
						break;

					case 'RegExp':
						hydrated[index] = new RegExp(value[1], value[2]);
						break;

					case 'Object':
						const object = Object(value[1]);

						if (Object.hasOwn(object, '__proto__')) {
							throw new Error('Cannot parse an object with a `__proto__` property');
						}

						hydrated[index] = object;
						break;

					case 'BigInt':
						hydrated[index] = BigInt(value[1]);
						break;

					case 'null':
						const obj = Object.create(null);
						hydrated[index] = obj;
						for (let i = 1; i < value.length; i += 2) {
							if (value[i] === '__proto__') {
								throw new Error('Cannot parse an object with a `__proto__` property');
							}

							obj[value[i]] = hydrate(value[i + 1]);
						}
						break;

					case 'Int8Array':
					case 'Uint8Array':
					case 'Uint8ClampedArray':
					case 'Int16Array':
					case 'Uint16Array':
					case 'Int32Array':
					case 'Uint32Array':
					case 'Float32Array':
					case 'Float64Array':
					case 'BigInt64Array':
					case 'BigUint64Array': {
						if (values[value[1]][0] !== 'ArrayBuffer') {
							// without this, if we receive malformed input we could
							// end up trying to hydrate in a circle or allocate
							// huge amounts of memory when we call `new TypedArrayConstructor(buffer)`
							throw new Error('Invalid data');
						}

						const TypedArrayConstructor = globalThis[type];
						const buffer = hydrate(value[1]);
						const typedArray = new TypedArrayConstructor(buffer);

						hydrated[index] =
							value[2] !== undefined
								? typedArray.subarray(value[2], value[3])
								: typedArray;

						break;
					}

					case 'ArrayBuffer': {
						const base64 = value[1];
						if (typeof base64 !== 'string') {
							throw new Error('Invalid ArrayBuffer encoding');
						}
						const arraybuffer = decode64(base64);
						hydrated[index] = arraybuffer;
						break;
					}

					case 'Temporal.Duration':
					case 'Temporal.Instant':
					case 'Temporal.PlainDate':
					case 'Temporal.PlainTime':
					case 'Temporal.PlainDateTime':
					case 'Temporal.PlainMonthDay':
					case 'Temporal.PlainYearMonth':
					case 'Temporal.ZonedDateTime': {
						const temporalName = type.slice(9);
						// @ts-expect-error TS doesn't know about Temporal yet
						hydrated[index] = Temporal[temporalName].from(value[1]);
						break;
					}

					case 'URL': {
						const url = new URL(value[1]);
						hydrated[index] = url;
						break;
					}

					case 'URLSearchParams': {
						const url = new URLSearchParams(value[1]);
						hydrated[index] = url;
						break;
					}

					default:
						throw new Error(`Unknown type ${type}`);
				}
			} else if (value[0] === SPARSE) {
				// Sparse array encoding: [SPARSE, length, idx, val, idx, val, ...]
				const len = value[1];

				if (!Number.isInteger(len) || len < 0) {
					throw new Error('Invalid input');
				}

				const array = new Array(len);
				hydrated[index] = array;

				for (let i = 2; i < value.length; i += 2) {
					const idx = value[i];

					if (!Number.isInteger(idx) || idx < 0 || idx >= len) {
						throw new Error('Invalid input');
					}

					array[idx] = hydrate(value[i + 1]);
				}
			} else {
				const array = new Array(value.length);
				hydrated[index] = array;

				for (let i = 0; i < value.length; i += 1) {
					const n = value[i];
					if (n === HOLE) continue;

					array[i] = hydrate(n);
				}
			}
		} else {
			/** @type {Record<string, any>} */
			const object = {};
			hydrated[index] = object;

			for (const key of Object.keys(value)) {
				if (key === '__proto__') {
					throw new Error('Cannot parse an object with a `__proto__` property');
				}

				const n = value[key];
				object[key] = hydrate(n);
			}
		}

		return hydrated[index];
	}

	return hydrate(0);
}

/**
 * Turn a value into a JSON string that can be parsed with `devalue.parse`
 * @param {any} value
 * @param {Record<string, (value: any) => any>} [reducers]
 */
function stringify(value, reducers) {
	/** @type {any[]} */
	const stringified = [];

	/** @type {Map<any, number>} */
	const indexes = new Map();

	/** @type {Array<{ key: string, fn: (value: any) => any }>} */
	const custom = [];
	if (reducers) {
		for (const key of Object.getOwnPropertyNames(reducers)) {
			custom.push({ key, fn: reducers[key] });
		}
	}

	/** @type {string[]} */
	const keys = [];

	let p = 0;

	/** @param {any} thing */
	function flatten(thing) {
		if (thing === undefined) return UNDEFINED;
		if (Number.isNaN(thing)) return NAN;
		if (thing === Infinity) return POSITIVE_INFINITY;
		if (thing === -Infinity) return NEGATIVE_INFINITY;
		if (thing === 0 && 1 / thing < 0) return NEGATIVE_ZERO;

		if (indexes.has(thing)) return indexes.get(thing);

		const index = p++;
		indexes.set(thing, index);

		for (const { key, fn } of custom) {
			const value = fn(thing);
			if (value) {
				stringified[index] = `["${key}",${flatten(value)}]`;
				return index;
			}
		}

		if (typeof thing === 'function') {
			throw new DevalueError(`Cannot stringify a function`, keys, thing, value);
		}

		let str = '';

		if (is_primitive(thing)) {
			str = stringify_primitive(thing);
		} else {
			const type = get_type(thing);

			switch (type) {
				case 'Number':
				case 'String':
				case 'Boolean':
					str = `["Object",${stringify_primitive(thing)}]`;
					break;

				case 'BigInt':
					str = `["BigInt",${thing}]`;
					break;

				case 'Date':
					const valid = !isNaN(thing.getDate());
					str = `["Date","${valid ? thing.toISOString() : ''}"]`;
					break;

				case 'URL':
					str = `["URL",${stringify_string(thing.toString())}]`;
					break;

				case 'URLSearchParams':
					str = `["URLSearchParams",${stringify_string(thing.toString())}]`;
					break;

				case 'RegExp':
					const { source, flags } = thing;
					str = flags
						? `["RegExp",${stringify_string(source)},"${flags}"]`
						: `["RegExp",${stringify_string(source)}]`;
					break;

				case 'Array': {
					// For dense arrays (no holes), we iterate normally.
					// When we encounter the first hole, we call Object.keys
					// to determine the sparseness, then decide between:
					//   - HOLE encoding: [-2, val, -2, ...] (default)
					//   - Sparse encoding: [-7, length, idx, val, ...] (for very sparse arrays)
					// Only the sparse path avoids iterating every slot, which
					// is what protects against the DoS of e.g. `arr[1000000] = 1`.
					let mostly_dense = false;

					str = '[';

					for (let i = 0; i < thing.length; i += 1) {
						if (i > 0) str += ',';

						if (Object.hasOwn(thing, i)) {
							keys.push(`[${i}]`);
							str += flatten(thing[i]);
							keys.pop();
						} else if (mostly_dense) {
							// Use dense encoding. The heuristic guarantees the
							// array is only mildly sparse, so iterating over every
							// slot is fine.
							str += HOLE;
						} else {
							// Decide between HOLE encoding and sparse encoding.
							//
							// HOLE encoding: each hole is serialized as the HOLE
							// sentinel (-2). For example, [, "a", ,] becomes
							// [-2, 0, -2]. Each hole costs 3 chars ("-2" + comma).
							//
							// Sparse encoding: lists only populated indices.
							// For example, [, "a", ,] becomes [-7, 3, 1, 0] — the
							// -7 sentinel, the array length (3), then index-value
							// pairs. This avoids paying per-hole, but each element
							// costs extra chars to write its index.
							//
							// The values are the same size either way, so the
							// choice comes down to structural overhead:
							//
							//   HOLE overhead:
							//     3 chars per hole ("-2" + comma)
							//     = (L - P) * 3
							//
							//   Sparse overhead:
							//     "-7,"          — 3 chars (sparse sentinel + comma)
							//     + length + "," — (d + 1) chars (array length + comma)
							//     + per element: index + "," — (d + 1) chars
							//     = (4 + d) + P * (d + 1)
							//
							// where L is the array length, P is the number of
							// populated elements, and d is the number of digits
							// in L (an upper bound on the digits in any index).
							//
							// Sparse encoding is cheaper when:
							//   (4 + d) + P * (d + 1) < (L - P) * 3
							const populated_keys = valid_array_indices(/** @type {any[]} */ (thing));
							const population = populated_keys.length;
							const d = String(thing.length).length;

							const hole_cost = (thing.length - population) * 3;
							const sparse_cost = 4 + d + population * (d + 1);

							if (hole_cost > sparse_cost) {
								str = '[' + SPARSE + ',' + thing.length;
								for (let j = 0; j < populated_keys.length; j++) {
									const key = populated_keys[j];
									keys.push(`[${key}]`);
									str += ',' + key + ',' + flatten(thing[key]);
									keys.pop();
								}
								break;
							} else {
								mostly_dense = true;
								str += HOLE;
							}
						}
					}

					str += ']';

					break;
				}

				case 'Set':
					str = '["Set"';

					for (const value of thing) {
						str += `,${flatten(value)}`;
					}

					str += ']';
					break;

				case 'Map':
					str = '["Map"';

					for (const [key, value] of thing) {
						keys.push(
							`.get(${is_primitive(key) ? stringify_primitive(key) : '...'})`
						);
						str += `,${flatten(key)},${flatten(value)}`;
						keys.pop();
					}

					str += ']';
					break;

				case 'Int8Array':
				case 'Uint8Array':
				case 'Uint8ClampedArray':
				case 'Int16Array':
				case 'Uint16Array':
				case 'Int32Array':
				case 'Uint32Array':
				case 'Float32Array':
				case 'Float64Array':
				case 'BigInt64Array':
				case 'BigUint64Array': {
					/** @type {import("./types.js").TypedArray} */
					const typedArray = thing;
					str = '["' + type + '",' + flatten(typedArray.buffer);

					const a = thing.byteOffset;
					const b = a + thing.byteLength;

					// handle subarrays
					if (a > 0 || b !== typedArray.buffer.byteLength) {
						const m = +/(\d+)/.exec(type)[1] / 8;
						str += `,${a / m},${b / m}`;
					}

					str += ']';
					break;
				}

				case 'ArrayBuffer': {
					/** @type {ArrayBuffer} */
					const arraybuffer = thing;
					const base64 = encode64(arraybuffer);

					str = `["ArrayBuffer","${base64}"]`;
					break;
				}

				case 'Temporal.Duration':
				case 'Temporal.Instant':
				case 'Temporal.PlainDate':
				case 'Temporal.PlainTime':
				case 'Temporal.PlainDateTime':
				case 'Temporal.PlainMonthDay':
				case 'Temporal.PlainYearMonth':
				case 'Temporal.ZonedDateTime':
					str = `["${type}",${stringify_string(thing.toString())}]`;
					break;

				default:
					if (!is_plain_object(thing)) {
						throw new DevalueError(
							`Cannot stringify arbitrary non-POJOs`,
							keys,
							thing,
							value
						);
					}

					if (enumerable_symbols(thing).length > 0) {
						throw new DevalueError(
							`Cannot stringify POJOs with symbolic keys`,
							keys,
							thing,
							value
						);
					}

					if (Object.getPrototypeOf(thing) === null) {
						str = '["null"';
						for (const key of Object.keys(thing)) {
							if (key === '__proto__') {
								throw new DevalueError(
									`Cannot stringify objects with __proto__ keys`,
									keys,
									thing,
									value
								);
							}

							keys.push(stringify_key(key));
							str += `,${stringify_string(key)},${flatten(thing[key])}`;
							keys.pop();
						}
						str += ']';
					} else {
						str = '{';
						let started = false;
						for (const key of Object.keys(thing)) {
							if (key === '__proto__') {
								throw new DevalueError(
									`Cannot stringify objects with __proto__ keys`,
									keys,
									thing,
									value
								);
							}

							if (started) str += ',';
							started = true;
							keys.push(stringify_key(key));
							str += `${stringify_string(key)}:${flatten(thing[key])}`;
							keys.pop();
						}
						str += '}';
					}
			}
		}

		stringified[index] = str;
		return index;
	}

	const index = flatten(value);

	// special case — value is represented as a negative index
	if (index < 0) return `${index}`;

	return `[${stringified.join(',')}]`;
}

/**
 * @param {any} thing
 * @returns {string}
 */
function stringify_primitive(thing) {
	const type = typeof thing;
	if (type === 'string') return stringify_string(thing);
	if (thing instanceof String) return stringify_string(thing.toString());
	if (thing === void 0) return UNDEFINED.toString();
	if (thing === 0 && 1 / thing < 0) return NEGATIVE_ZERO.toString();
	if (type === 'bigint') return `["BigInt","${thing}"]`;
	return String(thing);
}

const u=false,l$1="/_svelte_kit_assets",c$1=["GET","POST","PUT","PATCH","DELETE","OPTIONS","HEAD"],h=["POST","PUT","PATCH","DELETE"],p$1=["GET","POST","HEAD"];function d$1(e,t,s){t.startsWith("n:")?(t=t.slice(2),s=""===s?void 0:parseFloat(s)):t.startsWith("b:")&&(t=t.slice(2),s="on"===s),E(e,function(e){if(!v.test(e))throw new Error(`Invalid path ${e}`);return e.split(/\.|\[|\]/).filter(Boolean)}(t),s);}function m$1(e){const t={};for(let s of e.keys()){const n=s.endsWith("[]");let r=e.getAll(s);if(n&&(s=s.slice(0,-2)),r.length>1&&!n)throw new Error(`Form cannot contain duplicated keys — "${s}" has ${r.length} values`);r=r.filter(e=>"string"==typeof e||""!==e.name||e.size>0),s.startsWith("n:")?(s=s.slice(2),r=r.map(e=>""===e?void 0:parseFloat(e))):s.startsWith("b:")&&(s=s.slice(2),r=r.map(e=>"on"===e)),d$1(t,s,n?r:r[0]);}return t}const y="application/x-sveltekit-formdata";async function b(e){if(e.headers.get("content-type")!==y){const t=await e.formData();return {data:m$1(t),meta:{},form_data:t}}if(!e.body)throw g("no body");const t=parseInt(e.headers.get("content-length")??"");if(Number.isNaN(t))throw g("invalid Content-Length header");const s=e.body.getReader(),n=[];function r(e){if(e in n)return n[e];let t=n.length;for(;t<=e;)n[t]=s.read().then(e=>e.value),t++;return n[e]}async function o(e,t){let s,n,a=0;for(n=0;;n++){const t=await r(n);if(!t)return null;const i=a+t.byteLength;if(e>=a&&e<i){s=t;break}a=i;}if(e+t<=a+s.byteLength)return s.subarray(e-a,e+t-a);const i=[s.subarray(e-a)];let o=s.byteLength-e+a;for(;o<t;){n++;let e=await r(n);if(!e)return null;e.byteLength>t-o&&(e=e.subarray(0,t-o)),i.push(e),o+=e.byteLength;}const f=new Uint8Array(t);o=0;for(const r of i)f.set(r,o),o+=r.byteLength;return f}const f=await o(0,7);if(!f)throw g("too short");if(0!==f[0])throw g(`got version ${f[0]}, expected version 0`);const u=new DataView(f.buffer,f.byteOffset,f.byteLength),l=u.getUint32(1,true);if(7+l>t)throw g("data overflow");const c=u.getUint16(5,true);if(7+l+c>t)throw g("file offset table overflow");const h=await o(7,l);if(!h)throw g("data too short");let p,d;if(c>0){const e=await o(7+l,c);if(!e)throw g("file offset table too short");const t=JSON.parse(r$1.decode(e));if(!Array.isArray(t)||t.some(e=>"number"!=typeof e||!Number.isInteger(e)||e<0))throw g("invalid file offset table");p=t,d=7+l+c;}const b=[],[v,k]=parse(r$1.decode(h),{File:([e,s,n,a,i])=>{if("string"!=typeof e||"string"!=typeof s||"number"!=typeof n||"number"!=typeof a||"number"!=typeof i)throw g("invalid file metadata");let o=p[i];if(void 0===o)throw g("duplicate file offset table index");if(p[i]=void 0,o+=d,o+n>t)throw g("file data overflow");return b.push({offset:o,size:n}),new Proxy(new w(e,s,n,a,r,o),{getPrototypeOf:()=>File.prototype})}});b.sort((e,t)=>e.offset-t.offset||e.size-t.size);for(let a=1;a<b.length;a++){const e=b[a-1],t=b[a],s=e.offset+e.size;if(s<t.offset)throw g("gaps in file data");if(s>t.offset)throw g("overlapping file data")}return (async()=>{let e=true;for(;e;){e=!!(await r(n.length));}})(),{data:v,meta:k,form_data:null}}function g(e){return new SvelteKitError(400,"Bad Request",`Could not deserialize binary form: ${e}`)}class w{#e;#t;constructor(e,t,s,n,r,a){this.name=e,this.type=t,this.size=s,this.lastModified=n,this.webkitRelativePath="",this.#e=r,this.#t=a,this.arrayBuffer=this.arrayBuffer.bind(this),this.bytes=this.bytes.bind(this),this.slice=this.slice.bind(this),this.stream=this.stream.bind(this),this.text=this.text.bind(this);}#s;async arrayBuffer(){return this.#s??=await new Response(this.stream()).arrayBuffer(),this.#s}async bytes(){return new Uint8Array(await this.arrayBuffer())}slice(e=0,t=this.size,s=this.type){e=e<0?Math.max(this.size+e,0):Math.min(e,this.size),t=t<0?Math.max(this.size+t,0):Math.min(t,this.size);const n=Math.max(t-e,0);return new w(this.name,s,n,this.lastModified,this.#e,this.#t+e)}stream(){let e=0,t=0;return new ReadableStream({start:async s=>{let n,r=0;for(t=0;;t++){const e=await this.#e(t);if(!e)return null;const s=r+e.byteLength;if(this.#t>=r&&this.#t<s){n=e;break}r=s;}this.#t+this.size<=r+n.byteLength?(s.enqueue(n.subarray(this.#t-r,this.#t+this.size-r)),s.close()):(s.enqueue(n.subarray(this.#t-r)),e=n.byteLength-this.#t+r);},pull:async s=>{t++;let n=await this.#e(t);if(!n)return s.error("incomplete file data"),void s.close();n.byteLength>this.size-e&&(n=n.subarray(0,this.size-e)),s.enqueue(n),e+=n.byteLength,e>=this.size&&s.close();}})}async text(){return r$1.decode(await this.arrayBuffer())}}const v=/^[a-zA-Z_$]\w*(\.[a-zA-Z_$]\w*|\[\d+\])*$/;function k(e){if("__proto__"===e||"constructor"===e||"prototype"===e)throw new Error(`Invalid key "${e}"`)}function E(e,t,s){let n=e;for(let a=0;a<t.length-1;a+=1){const e=t[a];k(e);const s=/^\d+$/.test(t[a+1]),r=Object.hasOwn(n,e),i=n[e];if(r&&s!==Array.isArray(i))throw new Error(`Invalid array key ${t[a+1]}`);r||(n[e]=s?[]:{}),n=n[e];}const r=t[t.length-1];k(r),n[r]=s;}function j(e,t){const s=[];let n;e.split(",").forEach((e,t)=>{const n=/([^/ \t]+)\/([^; \t]+)[ \t]*(?:;[ \t]*q=([0-9.]+))?/.exec(e);if(n){const[,e,r,a="1"]=n;s.push({type:e,subtype:r,q:+a,i:t});}}),s.sort((e,t)=>e.q!==t.q?t.q-e.q:"*"===e.subtype!=("*"===t.subtype)?"*"===e.subtype?1:-1:"*"===e.type!=("*"===t.type)?"*"===e.type?1:-1:e.i-t.i);let r=1/0;for(const a of t){const[e,t]=a.split("/"),i=s.findIndex(s=>!(s.type!==e&&"*"!==s.type||s.subtype!==t&&"*"!==s.subtype));-1!==i&&i<r&&(n=a,r=i);}return n}function z(e){return function(e,...t){const s=e.headers.get("content-type")?.split(";",1)[0].trim()??"";return t.includes(s.toLowerCase())}(e,"application/x-www-form-urlencoded","multipart/form-data","text/plain",y)}function A(e){return e}function O(e){return e instanceof HttpError||e instanceof SvelteKitError?e.status:500}const T={"&":"&amp;",'"':"&quot;"},P={"&":"&amp;","<":"&lt;"},L="[\\ud800-\\udbff](?![\\udc00-\\udfff])|[\\ud800-\\udbff][\\udc00-\\udfff]|[\\udc00-\\udfff]",q=new RegExp(`[${Object.keys(T).join("")}]|`+L,"g"),D=new RegExp(`[${Object.keys(P).join("")}]|`+L,"g");function I(e,t){const s=t?T:P;return e.replace(t?q:D,e=>2===e.length?e:s[e]??`&#${e.charCodeAt(0)};`)}function B(e,s){return text(`${s} method not allowed`,{status:405,headers:{allow:R(e).join(", ")}})}function R(e){const t=c$1.filter(t=>t in e);return "GET"in e&&!("HEAD"in e)&&t.push("HEAD"),t}function S(e){return `__sveltekit_${e.version_hash}`}function M(e,s,n){let r=e.templates.error({status:s,message:I(n)});return text(r,{headers:{"content-type":"text/html; charset=utf-8"},status:s})}async function U(t,s,r,a){var i;const o=O(a=a instanceof HttpError?a:(i=a)instanceof Error||i&&i.name&&i.message?i:new Error(JSON.stringify(i))),f=await F(t,s,r,a),u=j(t.request.headers.get("accept")||"text/html",["application/json","text/html"]);return t.isDataRequest||"application/json"===u?json(f,{status:o}):M(r,o,f.message)}async function F(e,t,a,i){if(i instanceof HttpError)return {message:"Unknown Error",...i.body};const o=O(i),f=function(e){return e instanceof SvelteKitError?e.text:"Internal Error"}(i);return await with_request_store({event:e,state:t},()=>a.hooks.handleError({error:i,event:e,status:o,message:f}))??{message:f}}function N(e,t){return new Response(void 0,{status:e,headers:{location:t}})}function C(e,t){return t.path?`Data returned from \`load\` while rendering ${e.route.id} is not serializable: ${t.message} (${t.path}). If you need to serialize/deserialize custom types, use transport hooks: https://svelte.dev/docs/kit/hooks#Universal-hooks-transport.`:""===t.path?`Data returned from \`load\` while rendering ${e.route.id} is not a plain object`:t.message}function H(e){const t={};return e.uses&&e.uses.dependencies.size>0&&(t.dependencies=Array.from(e.uses.dependencies)),e.uses&&e.uses.search_params.size>0&&(t.search_params=Array.from(e.uses.search_params)),e.uses&&e.uses.params.size>0&&(t.params=Array.from(e.uses.params)),e.uses?.parent&&(t.parent=1),e.uses?.route&&(t.route=1),e.uses?.url&&(t.url=1),t}function W(e,t){return e._.prerendered_routes.has(t)||"/"===t.at(-1)&&e._.prerendered_routes.has(t.slice(0,-1))}function G(e,t,s){const n=`\n[1;31m[${e}] ${s.request.method} ${s.url.pathname}[0m`;return 404===e?n:`${n}\n${t.stack}`}function J(e){const t=e?.split("/"),s=t?.at(-1);if(!s)return "unknown";return s.split(".").slice(0,-1).join(".")}const Z="x-sveltekit-invalidated",V="x-sveltekit-trailing-slash";function K(e,t){const s=Object.fromEntries(Object.entries(t).map(([e,t])=>[e,t.encode]));return stringify(e,s)}function X(e,t){if(!e)return;const s=r$1.decode(o$1(e.replaceAll("-","+").replaceAll("_","/"))),n=Object.fromEntries(Object.entries(t).map(([e,t])=>[e,t.decode]));return parse(s,n)}function Y(e,t){return e+"/"+t}

let s$1="",a=s$1;const t="_app",e=true,n={base:s$1,assets:a};function o(t){s$1=t.base,a=t.assets;}function c(){s$1=n.base,a=n.assets;}

let s={};function r(t){}function i(t){s=t;}let d=null;function l(t){d=t;}const m={app_template_contains_nonce:false,async:false,csp:{mode:"auto",directives:{"upgrade-insecure-requests":false,"block-all-mixed-content":false},reportOnly:{"upgrade-insecure-requests":false,"block-all-mixed-content":false}},csrf_check_origin:true,csrf_trusted_origins:[],embedded:false,env_public_prefix:"PUBLIC_",env_private_prefix:"",hash_routing:false,hooks:null,preload_strategy:"modulepreload",root:x((t,s,r,i)=>{let d,l,{stores:c}=s,{page:m}=s,{constructors:p}=s,{components:v=[]}=s,{form:h}=s,{data_0:g=null}=s,{data_1:u=null}=s;a$1("__svelte__",c),o$2(c.page.notify),void 0===s.stores&&r.stores&&void 0!==c&&r.stores(c),void 0===s.page&&r.page&&void 0!==m&&r.page(m),void 0===s.constructors&&r.constructors&&void 0!==p&&r.constructors(p),void 0===s.components&&r.components&&void 0!==v&&r.components(v),void 0===s.form&&r.form&&void 0!==h&&r.form(h),void 0===s.data_0&&r.data_0&&void 0!==g&&r.data_0(g),void 0===s.data_1&&r.data_1&&void 0!==u&&r.data_1(u);let f=t.head;do{d=true,t.head=f,c.page.set(m),l=`  ${p[1]?`${g$1(p[0]||w$1,"svelte:component").$$render(t,{data:g,params:m.params,this:v[0]},{this:t=>{v[0]=t,d=false;}},{default:()=>`${g$1(p[1]||w$1,"svelte:component").$$render(t,{data:u,form:h,params:m.params,this:v[1]},{this:t=>{v[1]=t,d=false;}},{})}`})}`:`${g$1(p[0]||w$1,"svelte:component").$$render(t,{data:g,form:h,params:m.params,this:v[0]},{this:t=>{v[0]=t,d=false;}},{})}`} `;}while(!d);return l}),service_worker:false,service_worker_options:void 0,server_error_boundaries:false,templates:{app:({head:t,body:n,assets:e,nonce:a,env:o})=>'<!doctype html>\n<html lang="en">\n  <head>\n    <meta charset="utf-8" />\n    <link rel="icon" href="'+e+'/favicon.png" />\n    <meta name="viewport" content="width=device-width, initial-scale=1" />\n    '+t+'\n  </head>\n  <body data-sveltekit-preload-data="hover">\n    <div style="display: contents">'+n+"</div>\n  </body>\n</html>\n",error:({status:t,message:n})=>'<!doctype html>\n<html lang="en">\n\t<head>\n\t\t<meta charset="utf-8" />\n\t\t<title>'+n+"</title>\n\n\t\t<style>\n\t\t\tbody {\n\t\t\t\t--bg: white;\n\t\t\t\t--fg: #222;\n\t\t\t\t--divider: #ccc;\n\t\t\t\tbackground: var(--bg);\n\t\t\t\tcolor: var(--fg);\n\t\t\t\tfont-family:\n\t\t\t\t\tsystem-ui,\n\t\t\t\t\t-apple-system,\n\t\t\t\t\tBlinkMacSystemFont,\n\t\t\t\t\t'Segoe UI',\n\t\t\t\t\tRoboto,\n\t\t\t\t\tOxygen,\n\t\t\t\t\tUbuntu,\n\t\t\t\t\tCantarell,\n\t\t\t\t\t'Open Sans',\n\t\t\t\t\t'Helvetica Neue',\n\t\t\t\t\tsans-serif;\n\t\t\t\tdisplay: flex;\n\t\t\t\talign-items: center;\n\t\t\t\tjustify-content: center;\n\t\t\t\theight: 100vh;\n\t\t\t\tmargin: 0;\n\t\t\t}\n\n\t\t\t.error {\n\t\t\t\tdisplay: flex;\n\t\t\t\talign-items: center;\n\t\t\t\tmax-width: 32rem;\n\t\t\t\tmargin: 0 1rem;\n\t\t\t}\n\n\t\t\t.status {\n\t\t\t\tfont-weight: 200;\n\t\t\t\tfont-size: 3rem;\n\t\t\t\tline-height: 1;\n\t\t\t\tposition: relative;\n\t\t\t\ttop: -0.05rem;\n\t\t\t}\n\n\t\t\t.message {\n\t\t\t\tborder-left: 1px solid var(--divider);\n\t\t\t\tpadding: 0 0 0 1rem;\n\t\t\t\tmargin: 0 0 0 1rem;\n\t\t\t\tmin-height: 2.5rem;\n\t\t\t\tdisplay: flex;\n\t\t\t\talign-items: center;\n\t\t\t}\n\n\t\t\t.message h1 {\n\t\t\t\tfont-weight: 400;\n\t\t\t\tfont-size: 1em;\n\t\t\t\tmargin: 0;\n\t\t\t}\n\n\t\t\t@media (prefers-color-scheme: dark) {\n\t\t\t\tbody {\n\t\t\t\t\t--bg: #222;\n\t\t\t\t\t--fg: #ddd;\n\t\t\t\t\t--divider: #666;\n\t\t\t\t}\n\t\t\t}\n\t\t</style>\n\t</head>\n\t<body>\n\t\t<div class=\"error\">\n\t\t\t<span class=\"status\">"+t+'</span>\n\t\t\t<div class="message">\n\t\t\t\t<h1>'+n+"</h1>\n\t\t\t</div>\n\t\t</div>\n\t</body>\n</html>\n"},version_hash:"o7lzq5"};async function p(){return {handle:void 0,handleFetch:void 0,handleError:void 0,handleValidationError:void 0,init:void 0,reroute:void 0,transport:void 0}}

var cookie = {};

/*!
 * cookie
 * Copyright(c) 2012-2014 Roman Shtylman
 * Copyright(c) 2015 Douglas Christopher Wilson
 * MIT Licensed
 */

var hasRequiredCookie;

function requireCookie () {
	if (hasRequiredCookie) return cookie;
	hasRequiredCookie = 1;

	/**
	 * Module exports.
	 * @public
	 */

	cookie.parse = parse;
	cookie.serialize = serialize;

	/**
	 * Module variables.
	 * @private
	 */

	var __toString = Object.prototype.toString;

	/**
	 * RegExp to match field-content in RFC 7230 sec 3.2
	 *
	 * field-content = field-vchar [ 1*( SP / HTAB ) field-vchar ]
	 * field-vchar   = VCHAR / obs-text
	 * obs-text      = %x80-FF
	 */

	var fieldContentRegExp = /^[\u0009\u0020-\u007e\u0080-\u00ff]+$/;

	/**
	 * Parse a cookie header.
	 *
	 * Parse the given cookie header string into an object
	 * The object has the various cookies as keys(names) => values
	 *
	 * @param {string} str
	 * @param {object} [options]
	 * @return {object}
	 * @public
	 */

	function parse(str, options) {
	  if (typeof str !== 'string') {
	    throw new TypeError('argument str must be a string');
	  }

	  var obj = {};
	  var opt = options || {};
	  var dec = opt.decode || decode;

	  var index = 0;
	  while (index < str.length) {
	    var eqIdx = str.indexOf('=', index);

	    // no more cookie pairs
	    if (eqIdx === -1) {
	      break
	    }

	    var endIdx = str.indexOf(';', index);

	    if (endIdx === -1) {
	      endIdx = str.length;
	    } else if (endIdx < eqIdx) {
	      // backtrack on prior semicolon
	      index = str.lastIndexOf(';', eqIdx - 1) + 1;
	      continue
	    }

	    var key = str.slice(index, eqIdx).trim();

	    // only assign once
	    if (undefined === obj[key]) {
	      var val = str.slice(eqIdx + 1, endIdx).trim();

	      // quoted values
	      if (val.charCodeAt(0) === 0x22) {
	        val = val.slice(1, -1);
	      }

	      obj[key] = tryDecode(val, dec);
	    }

	    index = endIdx + 1;
	  }

	  return obj;
	}

	/**
	 * Serialize data into a cookie header.
	 *
	 * Serialize the a name value pair into a cookie string suitable for
	 * http headers. An optional options object specified cookie parameters.
	 *
	 * serialize('foo', 'bar', { httpOnly: true })
	 *   => "foo=bar; httpOnly"
	 *
	 * @param {string} name
	 * @param {string} val
	 * @param {object} [options]
	 * @return {string}
	 * @public
	 */

	function serialize(name, val, options) {
	  var opt = options || {};
	  var enc = opt.encode || encode;

	  if (typeof enc !== 'function') {
	    throw new TypeError('option encode is invalid');
	  }

	  if (!fieldContentRegExp.test(name)) {
	    throw new TypeError('argument name is invalid');
	  }

	  var value = enc(val);

	  if (value && !fieldContentRegExp.test(value)) {
	    throw new TypeError('argument val is invalid');
	  }

	  var str = name + '=' + value;

	  if (null != opt.maxAge) {
	    var maxAge = opt.maxAge - 0;

	    if (isNaN(maxAge) || !isFinite(maxAge)) {
	      throw new TypeError('option maxAge is invalid')
	    }

	    str += '; Max-Age=' + Math.floor(maxAge);
	  }

	  if (opt.domain) {
	    if (!fieldContentRegExp.test(opt.domain)) {
	      throw new TypeError('option domain is invalid');
	    }

	    str += '; Domain=' + opt.domain;
	  }

	  if (opt.path) {
	    if (!fieldContentRegExp.test(opt.path)) {
	      throw new TypeError('option path is invalid');
	    }

	    str += '; Path=' + opt.path;
	  }

	  if (opt.expires) {
	    var expires = opt.expires;

	    if (!isDate(expires) || isNaN(expires.valueOf())) {
	      throw new TypeError('option expires is invalid');
	    }

	    str += '; Expires=' + expires.toUTCString();
	  }

	  if (opt.httpOnly) {
	    str += '; HttpOnly';
	  }

	  if (opt.secure) {
	    str += '; Secure';
	  }

	  if (opt.partitioned) {
	    str += '; Partitioned';
	  }

	  if (opt.priority) {
	    var priority = typeof opt.priority === 'string'
	      ? opt.priority.toLowerCase()
	      : opt.priority;

	    switch (priority) {
	      case 'low':
	        str += '; Priority=Low';
	        break
	      case 'medium':
	        str += '; Priority=Medium';
	        break
	      case 'high':
	        str += '; Priority=High';
	        break
	      default:
	        throw new TypeError('option priority is invalid')
	    }
	  }

	  if (opt.sameSite) {
	    var sameSite = typeof opt.sameSite === 'string'
	      ? opt.sameSite.toLowerCase() : opt.sameSite;

	    switch (sameSite) {
	      case true:
	        str += '; SameSite=Strict';
	        break;
	      case 'lax':
	        str += '; SameSite=Lax';
	        break;
	      case 'strict':
	        str += '; SameSite=Strict';
	        break;
	      case 'none':
	        str += '; SameSite=None';
	        break;
	      default:
	        throw new TypeError('option sameSite is invalid');
	    }
	  }

	  return str;
	}

	/**
	 * URL-decode string value. Optimized to skip native call when no %.
	 *
	 * @param {string} str
	 * @returns {string}
	 */

	function decode (str) {
	  return str.indexOf('%') !== -1
	    ? decodeURIComponent(str)
	    : str
	}

	/**
	 * URL-encode value.
	 *
	 * @param {string} val
	 * @returns {string}
	 */

	function encode (val) {
	  return encodeURIComponent(val)
	}

	/**
	 * Determine if value is a Date.
	 *
	 * @param {*} val
	 * @private
	 */

	function isDate (val) {
	  return __toString.call(val) === '[object Date]' ||
	    val instanceof Date
	}

	/**
	 * Try decoding a string using a decoding function.
	 *
	 * @param {string} str
	 * @param {function} decode
	 * @private
	 */

	function tryDecode(str, decode) {
	  try {
	    return decode(str);
	  } catch (e) {
	    return str;
	  }
	}
	return cookie;
}

var cookieExports = requireCookie();

var defaultParseOptions = {
  decodeValues: true,
  map: false,
  silent: false,
  split: "auto", // auto = split strings but not arrays
};

function isForbiddenKey(key) {
  return typeof key !== "string" || key in {};
}

function createNullObj() {
  return Object.create(null);
}

function isNonEmptyString(str) {
  return typeof str === "string" && !!str.trim();
}

function parseString(setCookieValue, options) {
  var parts = setCookieValue.split(";").filter(isNonEmptyString);

  var nameValuePairStr = parts.shift();
  var parsed = parseNameValuePair(nameValuePairStr);
  var name = parsed.name;
  var value = parsed.value;

  options = options
    ? Object.assign({}, defaultParseOptions, options)
    : defaultParseOptions;

  if (isForbiddenKey(name)) {
    return null;
  }

  try {
    value = options.decodeValues ? decodeURIComponent(value) : value; // decode cookie value
  } catch (e) {
    console.error(
      "set-cookie-parser: failed to decode cookie value. Set options.decodeValues=false to disable decoding.",
      e
    );
  }

  var cookie = createNullObj();
  cookie.name = name;
  cookie.value = value;

  parts.forEach(function (part) {
    var sides = part.split("=");
    var key = sides.shift().trimLeft().toLowerCase();
    if (isForbiddenKey(key)) {
      return;
    }
    var value = sides.join("=");
    if (key === "expires") {
      cookie.expires = new Date(value);
    } else if (key === "max-age") {
      var n = parseInt(value, 10);
      if (!Number.isNaN(n)) cookie.maxAge = n;
    } else if (key === "secure") {
      cookie.secure = true;
    } else if (key === "httponly") {
      cookie.httpOnly = true;
    } else if (key === "samesite") {
      cookie.sameSite = value;
    } else if (key === "partitioned") {
      cookie.partitioned = true;
    } else if (key) {
      cookie[key] = value;
    }
  });

  return cookie;
}

function parseNameValuePair(nameValuePairStr) {
  // Parses name-value-pair according to rfc6265bis draft

  var name = "";
  var value = "";
  var nameValueArr = nameValuePairStr.split("=");
  if (nameValueArr.length > 1) {
    name = nameValueArr.shift();
    value = nameValueArr.join("="); // everything after the first =, joined by a "=" if there was more than one part
  } else {
    value = nameValuePairStr;
  }

  return { name: name, value: value };
}

function parseSetCookie(input, options) {
  options = options
    ? Object.assign({}, defaultParseOptions, options)
    : defaultParseOptions;

  if (!input) {
    if (!options.map) {
      return [];
    } else {
      return createNullObj();
    }
  }

  if (input.headers) {
    if (typeof input.headers.getSetCookie === "function") {
      // for fetch responses - they combine headers of the same type in the headers array,
      // but getSetCookie returns an uncombined array
      input = input.headers.getSetCookie();
    } else if (input.headers["set-cookie"]) {
      // fast-path for node.js (which automatically normalizes header names to lower-case)
      input = input.headers["set-cookie"];
    } else {
      // slow-path for other environments - see #25
      var sch =
        input.headers[
          Object.keys(input.headers).find(function (key) {
            return key.toLowerCase() === "set-cookie";
          })
        ];
      // warn if called on a request-like object with a cookie header rather than a set-cookie header - see #34, 36
      if (!sch && input.headers.cookie && !options.silent) {
        console.warn(
          "Warning: set-cookie-parser appears to have been called on a request object. It is designed to parse Set-Cookie headers from responses, not Cookie headers from requests. Set the option {silent: true} to suppress this warning."
        );
      }
      input = sch;
    }
  }

  var split = options.split;
  var isArray = Array.isArray(input);

  if (split === "auto") {
    split = !isArray;
  }

  if (!isArray) {
    input = [input];
  }

  input = input.filter(isNonEmptyString);

  if (split) {
    input = input.map(splitCookiesString).flat();
  }

  if (!options.map) {
    return input
      .map(function (str) {
        return parseString(str, options);
      })
      .filter(Boolean);
  } else {
    var cookies = createNullObj();
    return input.reduce(function (cookies, str) {
      var cookie = parseString(str, options);
      if (cookie && !isForbiddenKey(cookie.name)) {
        cookies[cookie.name] = cookie;
      }
      return cookies;
    }, cookies);
  }
}

/*
  Set-Cookie header field-values are sometimes comma joined in one string. This splits them without choking on commas
  that are within a single set-cookie field-value, such as in the Expires portion.

  This is uncommon, but explicitly allowed - see https://tools.ietf.org/html/rfc2616#section-4.2
  Node.js does this for every header *except* set-cookie - see https://github.com/nodejs/node/blob/d5e363b77ebaf1caf67cd7528224b651c86815c1/lib/_http_incoming.js#L128
  React Native's fetch does this for *every* header, including set-cookie.

  Based on: https://github.com/google/j2objc/commit/16820fdbc8f76ca0c33472810ce0cb03d20efe25
  Credits to: https://github.com/tomball for original and https://github.com/chrusart for JavaScript implementation
*/
function splitCookiesString(cookiesString) {
  if (Array.isArray(cookiesString)) {
    return cookiesString;
  }
  if (typeof cookiesString !== "string") {
    return [];
  }

  var cookiesStrings = [];
  var pos = 0;
  var start;
  var ch;
  var lastComma;
  var nextStart;
  var cookiesSeparatorFound;

  function skipWhitespace() {
    while (pos < cookiesString.length && /\s/.test(cookiesString.charAt(pos))) {
      pos += 1;
    }
    return pos < cookiesString.length;
  }

  function notSpecialChar() {
    ch = cookiesString.charAt(pos);

    return ch !== "=" && ch !== ";" && ch !== ",";
  }

  while (pos < cookiesString.length) {
    start = pos;
    cookiesSeparatorFound = false;

    while (skipWhitespace()) {
      ch = cookiesString.charAt(pos);
      if (ch === ",") {
        // ',' is a cookie separator if we have later first '=', not ';' or ','
        lastComma = pos;
        pos += 1;

        skipWhitespace();
        nextStart = pos;

        while (pos < cookiesString.length && notSpecialChar()) {
          pos += 1;
        }

        // currently special character
        if (pos < cookiesString.length && cookiesString.charAt(pos) === "=") {
          // we found cookies separator
          cookiesSeparatorFound = true;
          // pos is inside the next cookie, so back up and return it.
          pos = nextStart;
          cookiesStrings.push(cookiesString.substring(start, lastComma));
          start = pos;
        } else {
          // in param ',' or param separator ';',
          // we continue from that comma
          pos = lastComma + 1;
        }
      } else {
        pos += 1;
      }
    }

    if (!cookiesSeparatorFound || pos >= cookiesString.length) {
      cookiesStrings.push(cookiesString.substring(start, cookiesString.length));
    }
  }

  return cookiesStrings;
}

// named export for CJS
parseSetCookie.parseSetCookie = parseSetCookie;
// for backwards compatibility
parseSetCookie.parse = parseSetCookie;
parseSetCookie.parseString = parseString;
parseSetCookie.splitCookiesString = splitCookiesString;

function ge(){let e,t;return {promise:new Promise((s,r)=>{e=s,t=r;}),resolve:e,reject:t}}const we=[101,103,204,205,304],ve=!!globalThis.process?.versions?.webcontainer;function $e(e){return e.filter(e=>null!=e)}const be="/__data.json",ke=".html__data.json";function je(e){return e.endsWith(".html")?e.replace(/\.html$/,ke):e.replace(/\/$/,"")+be}const Re="/__route.js";function Ee(e){return e.replace(/\/$/,"")+Re}const Se={spanContext:()=>qe,setAttribute(){return this},setAttributes(){return this},addEvent(){return this},setStatus(){return this},updateName(){return this},end(){return this},isRecording:()=>false,recordException(){return this},addLink(){return this},addLinks(){return this}},qe={traceId:"",spanId:"",traceFlags:0};async function xe({name:e,attributes:t,fn:s}){return s(Se)}function Oe(e){return "application/json"===j(e.request.headers.get("accept")??"*/*",["application/json","text/html"])&&"POST"===e.request.method}function Te(e){return e instanceof ActionFailure?new Error('Cannot "throw fail()". Use "return fail()"'):e}function Pe(e){return Ae({type:"redirect",status:e.status,location:e.location})}function Ae(e,t){return json(e,t)}function Ue(e){if(e.default&&Object.keys(e).length>1)throw new Error("When using named actions, the default action cannot be used. See the docs for more info: https://svelte.dev/docs/kit/form-actions#named-actions")}async function Ce(e,t,s){const r=new URL(e.request.url);let n="default";for(const o of r.searchParams)if(o[0].startsWith("/")){if(n=o[0].slice(1),"default"===n)throw new Error('Cannot use reserved action name "default"');break}const a=s[n];if(!a)throw new SvelteKitError(404,"Not Found",`No action with name '${n}' found`);if(!z(e.request))throw new SvelteKitError(415,"Unsupported Media Type",`Form actions expect form-encoded data — received ${e.request.headers.get("content-type")}`);return xe({name:"sveltekit.form_action",attributes:{"http.route":e.route.id||"unknown"},fn:async s=>{const r=merge_tracing(e,s);t.allows_commands=true;const n=await with_request_store({event:r,state:t},()=>a(r));return n instanceof ActionFailure&&s.setAttributes({"sveltekit.form_action.result.type":"failure","sveltekit.form_action.result.status":n.status}),n}})}function Ne(e,t,s){const r=Object.fromEntries(Object.entries(s).map(([e,t])=>[e,t.encode]));return He(e,e=>stringify(e,r),t)}function He(e,t,s){try{return t(e)}catch(r){const t=r;if(e instanceof Response)throw new Error(`Data returned from action inside ${s} is not serializable. Form actions need to return plain objects or fail(). E.g. return { success: true } or return fail(400, { message: "invalid" });`,{cause:r});if("path"in t){let e=`Data returned from action inside ${s} is not serializable: ${t.message}`;throw ""!==t.path&&(e+=` (data.${t.path})`),new Error(e,{cause:r})}throw t}}function ze(){let e=-1,t=-1;const s=[];return {iterate:(e=e=>e)=>({[Symbol.asyncIterator]:()=>({next:async()=>{const r=s[++t];if(!r)return {value:null,done:true};const n=await r.promise;return {value:e(n),done:false}}})}),add:t=>{s.push(ge()),t.then(t=>{s[++e].resolve(t);});}}}function Le(e,t,s){let r=1,a=-1;const o=ze(),i=S(s);const c=[];return {set_max_nodes(e){a=e;},add_node(a,d){try{if(!d)return void(c[a]="null");const u={type:"data",data:d.data,uses:H(d)};d.slash&&(u.slash=d.slash),c[a]=uneval(u,(p=a,function a(c){if("function"==typeof c?.then){const d=r++,l=c.then(e=>({data:e})).catch(async r=>({error:await F(e,t,s,r)})).then(async({data:r,error:o})=>{let c;try{c=uneval(o?[,o]:[r],a);}catch{o=await F(e,t,s,new Error(`Failed to serialize promise while rendering ${e.route.id}`)),c=uneval([,o],a);}return {index:p,str:`${i}.resolve(${d}, ${c.includes("app.decode")?`(app) => ${c}`:`() => ${c}`})`}});return o.add(l),`${i}.defer(${d})`}for(const e in s.hooks.transport){const t=s.hooks.transport[e].encode(c);if(t)return `app.decode('${e}', ${uneval(t,a)})`}}));}catch(h){throw h.path=h.path.slice(1),new Error(C(e,h),{cause:h})}var p;},get_data(e){const t=`<script${e.script_needs_nonce?` nonce="${e.nonce}"`:""}>`;return {data:`[${$e(a>-1?c.slice(0,a):c).join(",")}]`,chunks:r>1?o.iterate(({index:e,str:s})=>a>-1&&e>=a?"":t+s+"<\/script>\n"):null}}}}function We(e,t,s){let r=1;const a=ze(),o={...Object.fromEntries(Object.entries(s.hooks.transport).map(([e,t])=>[e,t.encode])),Promise:i=>{if("function"!=typeof i?.then)return;const c=r++;let d="data";const l=i.catch(async r=>(d="error",F(e,t,s,r))).then(async r=>{let a;try{a=stringify(r,o);}catch{const r=await F(e,t,s,new Error(`Failed to serialize promise while rendering ${e.route.id}`));d="error",a=stringify(r,o);}return `{"type":"chunk","id":${c},"${d}":${a}}\n`});return a.add(l),c}},i=[];return {add_node(t,s){try{if(!s)return void(i[t]="null");if("error"===s.type||"skip"===s.type)return void(i[t]=JSON.stringify(s));i[t]=`{"type":"data","data":${stringify(s.data,o)},"uses":${JSON.stringify(H(s))}${s.slash?`,"slash":${JSON.stringify(s.slash)}`:""}}`;}catch(r){throw r.path="data"+r.path,new Error(C(e,r),{cause:r})}},get_data:()=>({data:`{"type":"data","nodes":[${i.join(",")}]}\n`,chunks:r>1?a.iterate():null})}}async function Ie({event:e,event_state:t,state:s,node:r,parent:n}){if(!r?.server)return null;let a=true;const o={dependencies:new Set,params:new Set,parent:false,route:false,url:false,search_params:new Set},i=r.server.load,c=r.server.trailingSlash;if(!i)return {type:"data",data:null,uses:o,slash:c};const d=o$3(e.url,()=>{a&&(o.url=true);},e=>{a&&o.search_params.add(e);});s.prerendering&&i$1(d);return {type:"data",data:await xe({name:"sveltekit.load",attributes:{"sveltekit.load.node_id":r.server_id||"unknown","sveltekit.load.node_type":J(r.server_id),"http.route":e.route.id||"unknown"},fn:async s=>{const r=merge_tracing(e,s);return await with_request_store({event:r,state:t},()=>i.call(null,{...r,fetch:(t,s)=>(new URL(t instanceof Request?t.url:t,e.url),e.fetch(t,s)),depends:(...t)=>{for(const s of t){const{href:t}=new URL(s,e.url);o.dependencies.add(t);}},params:new Proxy(e.params,{get:(e,t)=>(a&&o.params.add(t),e[t])}),parent:async()=>(a&&(o.parent=!0),n()),route:new Proxy(e.route,{get:(e,t)=>(a&&(o.route=!0),e[t])}),url:d,untrack(e){a=!1;try{return e()}finally{a=!0;}}}))}})??null,uses:o,slash:c}}async function Me({event:e,event_state:t,fetched:s,node:r,parent:n,server_data_promise:a,state:o,resolve_opts:i,csr:c}){const d=await a,l=r?.universal?.load;if(!l)return d?.data??null;return await xe({name:"sveltekit.load",attributes:{"sveltekit.load.node_id":r.universal_id||"unknown","sveltekit.load.node_type":J(r.universal_id),"http.route":e.route.id||"unknown"},fn:async r=>{const a=merge_tracing(e,r);return await with_request_store({event:a,state:t},()=>l.call(null,{url:e.url,params:e.params,data:d?.data??null,route:e.route,fetch:Fe(e,o,s,c,i),setHeaders:e.setHeaders,depends:()=>{},parent:n,untrack:e=>e(),tracing:a.tracing}))}})??null}function Fe(e,t,s,r,n){const a=async(a,o)=>{const i=a instanceof Request&&a.body?a.clone().body:null,c=a instanceof Request&&[...a.headers].length?new Headers(a.headers):o?.headers;let d=await e.fetch(a,o);const l=new URL(a instanceof Request?a.url:a,e.url),u=l.origin===e.url.origin;let p,h;if(u)t.prerendering&&(p={response:d,body:null},t.prerendering.dependencies.set(l.pathname,p));else if("https:"===l.protocol||"http:"===l.protocol){if("no-cors"===(a instanceof Request?a.mode:o?.mode??"cors"))d=new Response("",{status:d.status,statusText:d.statusText,headers:d.headers});else {const t=d.headers.get("access-control-allow-origin");if(!t||t!==e.url.origin&&"*"!==t)throw new Error(`CORS error: ${t?"Incorrect":"No"} 'Access-Control-Allow-Origin' header is present on the requested resource`)}}const f=new Proxy(d,{get(t,r,n){async function d(r,n){const d=Number(t.status);if(isNaN(d))throw new Error(`response.status is not a number. value: "${t.status}" type: ${typeof t.status}`);s.push({url:u?l.href.slice(e.url.origin.length):l.href,method:e.request.method,request_body:a instanceof Request&&i?await De(i):o?.body,request_headers:c,response_body:r,response:t,is_b64:n});}if("body"===r){if(null===t.body)return null;if(h)return h;const[e,s]=t.body.tee();return (async()=>{let t=new Uint8Array;for await(const s of e){const e=new Uint8Array(t.length+s.length);e.set(t,0),e.set(s,t.length),t=e;}p&&(p.body=new Uint8Array(t)),d(n$1(t),true);})(),h=s}if("arrayBuffer"===r)return async()=>{const e=await t.arrayBuffer(),s=new Uint8Array(e);return p&&(p.body=s),e instanceof ArrayBuffer&&await d(n$1(s),true),e};async function f(){const e=await t.text();if(""!==e||!we.includes(t.status))return e&&"string"!=typeof e||await d(e,false),p&&(p.body=e),e;await d(void 0,false);}if("text"===r)return f;if("json"===r)return async()=>{const e=await f();return e?JSON.parse(e):void 0};const _=Reflect.get(t,r,t);return _ instanceof Function?Object.defineProperties(function(){return Reflect.apply(_,this===n?t:this,arguments)},{name:{value:_.name},length:{value:_.length}}):_}});if(r){const t=d.headers.get;d.headers.get=s=>{const r=s.toLowerCase(),a=t.call(d.headers,r);if(a&&!r.startsWith("x-sveltekit-")){if(!n.filterSerializedResponseHeaders(r,a))throw new Error(`Failed to get response header "${r}" — it must be included by the \`filterSerializedResponseHeaders\` option: https://svelte.dev/docs/kit/hooks#Server-hooks-handle (at ${e.route.id})`)}return a};}return f};return (e,t)=>{const s=a(e,t);return s.catch(()=>{}),s}}async function De(e){let t="";const s=e.getReader();for(;;){const{done:e,value:r}=await s.read();if(e)break;t+=r$1.decode(r);}return t}function Je(...e){let t=5381;for(const s of e)if("string"==typeof s){let e=s.length;for(;e;)t=33*t^s.charCodeAt(--e);}else {if(!ArrayBuffer.isView(s))throw new TypeError("value must be a string or TypedArray");{const e=new Uint8Array(s.buffer,s.byteOffset,s.byteLength);let r=e.length;for(;r;)t=33*t^e[--r];}}return (t>>>0).toString(36)}const Ge={"<":"\\u003C","\u2028":"\\u2028","\u2029":"\\u2029"},Be=new RegExp(`[${Object.keys(Ge).join("")}]`,"g");const Ve=JSON.stringify;function Ke(e){Xe[0]||function(){function e(e){return 4294967296*(e-Math.floor(e))}let t=2;for(let s=0;s<64;t++){let r=true;for(let e=2;e*e<=t;e++)if(t%e===0){r=false;break}r&&(s<8&&(Qe[s]=e(t**.5)),Xe[s]=e(t**(1/3)),s++);}}();const t=Qe.slice(0),s=function(e){const t=t$2.encode(e),s=8*t.length,r=512*Math.ceil((s+65)/512),n=new Uint8Array(r/8);n.set(t),n[t.length]=128,Ye(n);const a=new Uint32Array(n.buffer);return a[a.length-2]=Math.floor(s/4294967296),a[a.length-1]=s,a}(e);for(let n=0;n<s.length;n+=16){const e=s.subarray(n,n+16);let r,a,o,i=t[0],c=t[1],d=t[2],l=t[3],u=t[4],p=t[5],h=t[6],f=t[7];for(let t=0;t<64;t++)t<16?r=e[t]:(a=e[t+1&15],o=e[t+14&15],r=e[15&t]=(a>>>7^a>>>18^a>>>3^a<<25^a<<14)+(o>>>17^o>>>19^o>>>10^o<<15^o<<13)+e[15&t]+e[t+9&15]|0),r=r+f+(u>>>6^u>>>11^u>>>25^u<<26^u<<21^u<<7)+(h^u&(p^h))+Xe[t],f=h,h=p,p=u,u=l+r|0,l=d,d=c,c=i,i=r+(c&d^l&(c^d))+(c>>>2^c>>>13^c>>>22^c<<30^c<<19^c<<10)|0;t[0]=t[0]+i|0,t[1]=t[1]+c|0,t[2]=t[2]+d|0,t[3]=t[3]+l|0,t[4]=t[4]+u|0,t[5]=t[5]+p|0,t[6]=t[6]+h|0,t[7]=t[7]+f|0;}const r=new Uint8Array(t.buffer);return Ye(r),btoa(String.fromCharCode(...r))}const Qe=new Uint32Array(8),Xe=new Uint32Array(64);function Ye(e){for(let t=0;t<e.length;t+=4){const s=e[t+0],r=e[t+1],n=e[t+2],a=e[t+3];e[t+0]=a,e[t+1]=n,e[t+2]=r,e[t+3]=s;}}const Ze=new Uint8Array(16);const et=new Set(["self","unsafe-eval","unsafe-hashes","unsafe-inline","none","strict-dynamic","report-sample","wasm-unsafe-eval","script"]),tt=/^(nonce|sha\d\d\d)-/;class st{#e;#t;#s;#r;#n;#a;#o;#i;#c;#d;#l;#u;#p;#h;script_needs_nonce;style_needs_nonce;script_needs_hash;#f;constructor(e,t,s){this.#e=e,this.#c=t;const r=this.#c;this.#d=new Set,this.#l=new Set,this.#u=new Set,this.#p=new Set,this.#h=new Set;const n=r["script-src"]||r["default-src"],a=r["script-src-elem"],o=r["style-src"]||r["default-src"],i=r["style-src-attr"],c=r["style-src-elem"],d=e=>!!e&&!e.some(e=>"unsafe-inline"===e),l=e=>!!e&&(!e.some(e=>"unsafe-inline"===e)||e.some(e=>"strict-dynamic"===e));this.#s=l(n),this.#r=l(a),this.#a=d(o),this.#o=d(i),this.#i=d(c),this.#t=this.#s||this.#r,this.#n=this.#a||this.#o||this.#i,this.script_needs_nonce=this.#t&&!this.#e,this.style_needs_nonce=this.#n&&!this.#e,this.script_needs_hash=this.#t&&this.#e,this.#f=s;}add_script(e){if(!this.#t)return;const t=this.#e?`sha256-${Ke(e)}`:`nonce-${this.#f}`;this.#s&&this.#d.add(t),this.#r&&this.#l.add(t);}add_script_hashes(e){for(const t of e)this.#s&&this.#d.add(t),this.#r&&this.#l.add(t);}add_style(e){if(!this.#n)return;const t=this.#e?`sha256-${Ke(e)}`:`nonce-${this.#f}`;if(this.#a&&this.#u.add(t),this.#o&&this.#p.add(t),this.#i){const e="sha256-9OlNO0DNEeaVzHL4RZwCLsBHA8WBQ8toBp/4F5XV2nc=",s=this.#c;!s["style-src-elem"]||s["style-src-elem"].includes(e)||this.#h.has(e)||this.#h.add(e),t!==e&&this.#h.add(t);}}get_header(e=false){const t=[],s={...this.#c};this.#u.size>0&&(s["style-src"]=[...s["style-src"]||s["default-src"]||[],...this.#u]),this.#p.size>0&&(s["style-src-attr"]=[...s["style-src-attr"]||[],...this.#p]),this.#h.size>0&&(s["style-src-elem"]=[...s["style-src-elem"]||[],...this.#h]),this.#d.size>0&&(s["script-src"]=[...s["script-src"]||s["default-src"]||[],...this.#d]),this.#l.size>0&&(s["script-src-elem"]=[...s["script-src-elem"]||[],...this.#l]);for(const r in s){if(e&&("frame-ancestors"===r||"report-uri"===r||"sandbox"===r))continue;const n=s[r];if(!n)continue;const a=[r];Array.isArray(n)&&n.forEach(e=>{et.has(e)||tt.test(e)?a.push(`'${e}'`):a.push(e);}),t.push(a.join(" "));}return t.join("; ")}}class rt extends st{get_meta(){const e=this.get_header(true);if(e)return `<meta http-equiv="content-security-policy" content="${I(e,true)}">`}}class nt extends st{constructor(e,t,s){if(super(e,t,s),Object.values(t).filter(e=>!!e).length>0){const e=t["report-to"]?.length??false,s=t["report-uri"]?.length??false;if(!e&&!s)throw Error("`content-security-policy-report-only` must be specified with either the `report-to` or `report-uri` directives, or both")}}}class at{nonce=function(){return crypto.getRandomValues(Ze),btoa(String.fromCharCode(...Ze))}();csp_provider;report_only_provider;constructor({mode:e,directives:t,reportOnly:s},{prerender:r}){const n="hash"===e||"auto"===e&&r;this.csp_provider=new rt(n,t,this.nonce),this.report_only_provider=new nt(n,s,this.nonce);}get script_needs_hash(){return this.csp_provider.script_needs_hash||this.report_only_provider.script_needs_hash}get script_needs_nonce(){return this.csp_provider.script_needs_nonce||this.report_only_provider.script_needs_nonce}get style_needs_nonce(){return this.csp_provider.style_needs_nonce||this.report_only_provider.style_needs_nonce}add_script(e){this.csp_provider.add_script(e),this.report_only_provider.add_script(e);}add_script_hashes(e){this.csp_provider.add_script_hashes(e),this.report_only_provider.add_script_hashes(e);}add_style(e){this.csp_provider.add_style(e),this.report_only_provider.add_style(e);}}function ot(e,t,s){const r={},n=e.slice(1),a=n.filter(e=>void 0!==e);let o=0;for(let i=0;i<t.length;i+=1){const e=t[i];let c=n[i-o];if(e.chained&&e.rest&&o&&(c=n.slice(i-o,i+1).filter(e=>e).join("/"),o=0),void 0===c){if(!e.rest)continue;c="";}if(!e.matcher||s[e.matcher](c)){r[e.name]=c;const s=t[i+1],d=n[i+1];s&&!s.rest&&s.optional&&d&&e.chained&&(o=0),s||d||Object.keys(r).length!==a.length||(o=0);continue}if(!e.optional||!e.chained)return;o++;}if(!o)return r}function it(e,t,s){for(const r of t){const t=r.pattern.exec(e);if(!t)continue;const n=ot(t,r.params,s);if(n)return {route:r,params:a$2(n)}}return null}function ct(e,t,s){const{errors:r,layouts:n,leaf:a}=e,o=[...r,...n.map(e=>e?.[1]),a[1]].filter(e=>"number"==typeof e).map(e=>`'${e}': () => ${dt(s._.client.nodes?.[e],t)}`).join(",\n\t\t");return [`{\n\tid: ${Ve(e.id)}`,`errors: ${Ve(e.errors)}`,`layouts: ${Ve(e.layouts)}`,`leaf: ${Ve(e.leaf)}`,`nodes: {\n\t\t${o}\n\t}\n}`].join(",\n\t")}function dt(e,t){if(!e)return "Promise.resolve({})";if("/"===e[0])return `import('${e}')`;if(""!==a)return `import('${a}/${e}')`;let s=e$1(t.pathname,`${s$1}/${e}`);return "."!==s[0]&&(s=`./${s}`),`import('${s}')`}function lt(e,t,s,r){const n=new Headers({"content-type":"application/javascript; charset=utf-8"});if(e){const a$1=ct(e,s,r),o=`${function(e,t,s){const{errors:r,layouts:n,leaf:a$1}=e;let o="";for(const i of [...r,...n.map(e=>e?.[1]),a$1[1]]){if("number"!=typeof i)continue;const e=s._.client.css?.[i];for(const t of e??[])o+=`'${a||s$1}/${t}',`;}return o?`${dt(s._.client.start,t)}.then(x => x.load_css([${o}]));`:""}(e,s,r)}\nexport const route = ${a$1}; export const params = ${JSON.stringify(t)};`;return {response:text(o,{headers:n}),body:o}}return {response:text("",{headers:n}),body:""}}const ut={...c$2(false),check:()=>false};async function pt({branch:e$1,fetched:t$1,options:s$2,manifest:r,state:a$1,page_config:i,status:c$1,error:l=null,event:u,event_state:p,resolve_opts:m,action_result:y,data_serializer:g,error_components:w}){if(a$1.prerendering){if("nonce"===s$2.csp.mode)throw new Error('Cannot use prerendering if config.kit.csp.mode === "nonce"');if(s$2.app_template_contains_nonce)throw new Error("Cannot use prerendering if page template contains %sveltekit.nonce%")}const{client:v}=r._,$=new Set(v.imports),b=new Set(v.stylesheets),k=new Set(v.fonts),j=new Set,R=new Map;let E;const S$1="success"===y?.type||"failure"===y?.type?y.data??null:null;let x=s$1,O$1=a,T=Ve(s$1);const P=new at(s$2.csp,{prerender:!!a$1.prerendering});if(a$1.prerendering?.fallback)s$2.hash_routing&&(T="new URL('.', location).pathname.slice(0, -1)");else {const e=u.url.pathname.slice(s$1.length).split("/").slice(2);x=e.map(()=>"..").join("/")||".",T=`new URL(${Ve(x)}, location).pathname.slice(0, -1)`,(!a||"/"===a[0]&&a!==l$1)&&(O$1=x);}if(i.ssr){const t$1={stores:{page:u$1(null),navigating:u$1(null),updated:ut},constructors:await Promise.all(e$1.map(({node:e})=>{if(!e.component)throw new Error(`Missing +page.svelte component for route ${u.route.id}`);return e.component()})),form:S$1};w&&(l&&(t$1.error=l),t$1.errors=w);let r={};for(let s=0;s<e$1.length;s+=1)r={...r,...e$1[s].data},t$1[`data_${s}`]=r;t$1.page={error:l,params:u.params,route:u.route,status:c$1,url:u.url,data:r,form:S$1,state:{}};const a={context:new Map([["__request__",{page:t$1.page}]]),csp:P.script_needs_nonce?{nonce:P.nonce}:{hash:P.script_needs_hash},transformError:w?async e=>{const r=await F(u,p,s$2,e);return t$1.page.error=t$1.error=l=r,t$1.page.status=c$1=O(e),r}:void 0};try{p.allows_commands=!1,E=await with_request_store({event:u,state:p},async()=>{e&&o({base:x,assets:O$1});const e$1=s$2.root.render(t$1,a),r=s$2.async&&"then"in e$1?e$1.then(e=>e):e$1;s$2.async&&c();const{head:n,html:o$1,css:i,hashes:c$1}=s$2.async?await r:r;return c$1&&P.add_script_hashes(c$1.script),{head:n,html:o$1,css:i,hashes:c$1}});}finally{c();}for(const{node:s}of e$1){for(const e of s.imports)$.add(e);for(const e of s.stylesheets)b.add(e);for(const e of s.fonts)k.add(e);s.inline_styles&&!v.inline&&Object.entries(await s.inline_styles()).forEach(([e,t$1])=>{"string"!=typeof t$1?R.set(e,t$1(`${O$1}/${t}/immutable/assets`,O$1)):R.set(e,t$1);});}}else E={head:"",html:"",css:{code:"",map:null},hashes:{script:[]}};const A=new ht(E.head,!!a$1.prerendering);let C=E.html;const N=e=>e.startsWith("/")?s$1+e:`${O$1}/${e}`,D=v.inline?v.inline?.style:Array.from(R.values()).join("\n");if(D){const e=[];P.style_needs_nonce&&e.push(`nonce="${P.nonce}"`),P.add_style(D),A.add_style(D,e);}for(const n of b){const e=N(n),t=['rel="stylesheet"'];R.has(n)?t.push("disabled",'media="(max-width: 0)"'):m.preload({type:"css",path:e})&&j.add(`<${encodeURI(e)}>; rel="preload"; as="style"; nopush`),A.add_stylesheet(e,t);}for(const n of k){const e=N(n);if(m.preload({type:"font",path:e})){const t=n.slice(n.lastIndexOf(".")+1);A.add_link_tag(e,['rel="preload"','as="font"',`type="font/${t}"`,"crossorigin"]),j.add(`<${encodeURI(e)}>; rel="preload"; as="font"; type="font/${t}"; crossorigin; nopush`);}}const J=S(s$2),{data:G,chunks:B}=g.get_data(P);if(i.ssr&&i.csr&&(C+=`\n\t\t\t${t$1.map(e=>function(e,t,s=false){const r={};let n=null,a=null,o=false;for(const[l,u]of e.response.headers)t(l,u)&&(r[l]=u),"cache-control"===l?n=u:"age"===l?a=u:"vary"===l&&"*"===u.trim()&&(o=true);const i={status:e.response.status,statusText:e.response.statusText,headers:r,body:e.response_body},c=JSON.stringify(i).replace(Be,e=>Ge[e]),d=['type="application/json"',"data-sveltekit-fetched",`data-url="${I(e.url,true)}"`];if(e.is_b64&&d.push("data-b64"),e.request_headers||e.request_body){const t=[];e.request_headers&&t.push([...new Headers(e.request_headers)].join(",")),e.request_body&&t.push(e.request_body),d.push(`data-hash="${Je(...t)}"`);}if(!s&&"GET"===e.method&&n&&!o){const e=/s-maxage=(\d+)/g.exec(n)??/max-age=(\d+)/g.exec(n);if(e){const t=+e[1]-+(a??"0");d.push(`data-ttl="${t}"`);}}return `<script ${d.join(" ")}>${c}<\/script>`}(e,m.filterSerializedResponseHeaders,!!a$1.prerendering)).join("\n\t\t\t")}`),i.csr){const t$1=r._.client.routes?.find(e=>e.id===u.route.id)??null;if(v.uses_env_dynamic_public&&a$1.prerendering&&$.add(`${t}/env.js`),!v.inline){const e=Array.from($,e=>N(e)).filter(e=>m.preload({type:"js",path:e}));for(const t of e)j.add(`<${encodeURI(t)}>; rel="modulepreload"; nopush`),"modulepreload"!==s$2.preload_strategy?A.add_script_preload(t):A.add_link_tag(t,['rel="modulepreload"']);}if(r._.client.routes&&a$1.prerendering&&!a$1.prerendering.fallback){const e=Ee(u.url.pathname);a$1.prerendering.dependencies.set(e,lt(t$1,u.params,new URL(e,u.url),r));}const n=[],o=v.uses_env_dynamic_public&&a$1.prerendering,d=[`base: ${T}`];if(a&&d.push(`assets: ${Ve(a)}`),v.uses_env_dynamic_public&&d.push(`env: ${o?"null":Ve(s)}`),B){n.push("const deferred = new Map();"),d.push("defer: (id) => new Promise((fulfil, reject) => {\n\t\t\t\t\t\t\tdeferred.set(id, { fulfil, reject });\n\t\t\t\t\t\t})");let e="";Object.keys(s$2.hooks.transport).length>0&&(e=v.inline?`const app = __sveltekit_${s$2.version_hash}.app.app;`:v.app?`const app = await import(${Ve(N(v.app))});`:`const { app } = await import(${Ve(N(v.start))});`);const t=e?`${e}\n\t\t\t\t\t\t\tconst [data, error] = fn(app);`:"const [data, error] = fn();";d.push(`resolve: async (id, fn) => {\n\t\t\t\t\t\t\t${t}\n\n\t\t\t\t\t\t\tconst try_to_resolve = () => {\n\t\t\t\t\t\t\t\tif (!deferred.has(id)) {\n\t\t\t\t\t\t\t\t\tsetTimeout(try_to_resolve, 0);\n\t\t\t\t\t\t\t\t\treturn;\n\t\t\t\t\t\t\t\t}\n\t\t\t\t\t\t\t\tconst { fulfil, reject } = deferred.get(id);\n\t\t\t\t\t\t\t\tdeferred.delete(id);\n\t\t\t\t\t\t\t\tif (error) reject(error);\n\t\t\t\t\t\t\t\telse fulfil(data);\n\t\t\t\t\t\t\t}\n\t\t\t\t\t\t\ttry_to_resolve();\n\t\t\t\t\t\t}`);}n.push(`${J} = {\n\t\t\t\t\t\t${d.join(",\n\t\t\t\t\t\t")}\n\t\t\t\t\t};`);const h=["element"];if(n.push("const element = document.currentScript.parentElement;"),i.ssr){const n={form:"null",error:"null"};S$1&&(n.form=function(e,t,s){const r=e=>{for(const t in s){const n=s[t].encode(e);if(n)return `app.decode('${t}', ${uneval(n,r)})`}};return He(e,e=>uneval(e,r),t)}(S$1,u.route.id,s$2.hooks.transport)),l&&(n.error=uneval(l));const a=[`node_ids: [${e$1.map(({node:e})=>e.index).join(", ")}]`,`data: ${G}`,`form: ${n.form}`,`error: ${n.error}`];if(200!==c$1&&a.push(`status: ${c$1}`),r._.client.routes){if(t$1){const e=ct(t$1,u.url,r).replaceAll("\n","\n\t\t\t\t\t\t\t");a.push(`params: ${uneval(u.params)}`,`server_route: ${e}`);}}else s$2.embedded&&a.push(`params: ${uneval(u.params)}`,`route: ${Ve(u.route)}`);const i="\t".repeat(o?7:6);h.push(`{\n${i}\t${a.join(`,\n${i}\t`)}\n${i}}`);}const{remote_data:f}=p;let y="";if(f){const e={};for(const[s,r]of f)if(s.id)for(const t in r){const n=Y(s.id,t);if(void 0!==p.refreshes?.[n])e[n]=await r[t];else {const s=await Promise.race([Promise.resolve(r[t]).then(e=>({settled:true,value:e}),e=>({settled:true,error:e})),new Promise(e=>{queueMicrotask(()=>e({settled:false}));})]);if(s.settled){if("error"in s)throw s.error;e[n]=s.value;}}}const t=e=>{for(const r in s$2.hooks.transport){const n=s$2.hooks.transport[r].encode(e);if(n)return `app.decode('${r}', ${uneval(n,t)})`}};y=`${J}.data = ${uneval(e,t)};\n\n\t\t\t\t\t\t`;}const g=v.inline?`${v.inline.script}\n\n\t\t\t\t\t${y}${J}.app.start(${h.join(", ")});`:v.app?`Promise.all([\n\t\t\t\t\t\timport(${Ve(N(v.start))}),\n\t\t\t\t\t\timport(${Ve(N(v.app))})\n\t\t\t\t\t]).then(([kit, app]) => {\n\t\t\t\t\t\t${y}kit.start(app, ${h.join(", ")});\n\t\t\t\t\t});`:`import(${Ve(N(v.start))}).then((app) => {\n\t\t\t\t\t\t${y}app.start(${h.join(", ")})\n\t\t\t\t\t});`;if(o?n.push(`import(${Ve(`${x}/${t}/env.js`)}).then(({ env }) => {\n\t\t\t\t\t\t${J}.env = env;\n\n\t\t\t\t\t\t${g.replace(/\n/g,"\n\t")}\n\t\t\t\t\t});`):n.push(g),s$2.service_worker){let e="";if(null!=s$2.service_worker_options){const t={...s$2.service_worker_options};e=`, ${Ve(t)}`;}n.push(`if ('serviceWorker' in navigator) {\n\t\t\t\t\t\taddEventListener('load', function () {\n\t\t\t\t\t\t\tnavigator.serviceWorker.register('${N("service-worker.js")}'${e});\n\t\t\t\t\t\t});\n\t\t\t\t\t}`);}const w=`\n\t\t\t\t{\n\t\t\t\t\t${n.join("\n\n\t\t\t\t\t")}\n\t\t\t\t}\n\t\t\t`;P.add_script(w),C+=`\n\t\t\t<script${P.script_needs_nonce?` nonce="${P.nonce}"`:""}>${w}<\/script>\n\t\t`;}const V=new Headers({"x-sveltekit-page":"true","content-type":"text/html"});if(a$1.prerendering){const e=P.csp_provider.get_meta();e&&A.add_http_equiv(e),a$1.prerendering.cache&&A.add_http_equiv(`<meta http-equiv="cache-control" content="${a$1.prerendering.cache}">`);}else {const e=P.csp_provider.get_header();e&&V.set("content-security-policy",e);const t=P.report_only_provider.get_header();t&&V.set("content-security-policy-report-only",t),j.size&&V.set("link",Array.from(j).join(", "));}const K=s$2.templates.app({head:A.build(),body:C,assets:O$1,nonce:P.nonce,env:s}),Q=await m.transformPageChunk({html:K,done:true})||"";return B||V.set("etag",`"${Je(Q)}"`),B?new Response(new ReadableStream({async start(e){e.enqueue(t$2.encode(Q+"\n"));for await(const t of B)t.length&&e.enqueue(t$2.encode(t));e.close();},type:"bytes"}),{headers:V}):text(Q,{status:c$1,headers:V})}class ht{#_;#m;#y=[];#g=[];#w=[];#v=[];#$=[];constructor(e,t){this.#_=e,this.#m=t;}build(){return [...this.#y,...this.#g,...this.#w,this.#_,...this.#v,...this.#$].join("\n\t\t")}add_style(e,t){this.#v.push(`<style${t.length?" "+t.join(" "):""}>${e}</style>`);}add_stylesheet(e,t){this.#$.push(`<link href="${e}" ${t.join(" ")}>`);}add_script_preload(e){this.#w.push(`<link rel="preload" as="script" crossorigin="anonymous" href="${e}">`);}add_link_tag(e,t){this.#m&&this.#g.push(`<link href="${e}" ${t.join(" ")}>`);}add_http_equiv(e){this.#m&&this.#y.push(e);}}class ft{data;constructor(e){this.data=e;}layouts(){return this.data.slice(0,-1)}page(){return this.data.at(-1)}validate(){for(const t of this.layouts())t&&(b$1(t.server,t.server_id),w$2(t.universal,t.universal_id));const e=this.page();e&&(S$1(e.server,e.server_id),m$2(e.universal,e.universal_id));}#b(e){return this.data.reduce((t,s)=>s?.universal?.[e]??s?.server?.[e]??t,void 0)}csr(){return this.#b("csr")??true}ssr(){return this.#b("ssr")??true}prerender(){return this.#b("prerender")??false}trailing_slash(){return this.#b("trailingSlash")??"never"}get_config(){let e={};for(const t of this.data)(t?.universal?.config||t?.server?.config)&&(e={...e,...t?.universal?.config,...t?.server?.config});return Object.keys(e).length?e:void 0}should_prerender_data(){return this.data.some(e=>e?.server?.load||void 0!==e?.server?.trailingSlash)}}async function _t({event:e,event_state:t,options:s,manifest:r,state:a,status:i,error:c,resolve_opts:d}){if(e.request.headers.get("x-sveltekit-error"))return M(s,i,c.message);const l=[];try{const o=[],u=await r._.nodes[0](),p=new ft([u]),h=p.ssr(),f=p.csr(),_=Le(e,t,s);if(h){a.error=!0;const s=Ie({event:e,event_state:t,state:a,node:u,parent:async()=>({})}),n=await s;_.add_node(0,n);const i=await Me({event:e,event_state:t,fetched:l,node:u,parent:async()=>({}),resolve_opts:d,server_data_promise:s,state:a,csr:f});o.push({node:u,server_data:n,data:i},{node:await r._.nodes[1](),data:null,server_data:null});}return await pt({options:s,manifest:r,state:a,page_config:{ssr:h,csr:f},status:i,error:await F(e,t,s,c),branch:o,error_components:[],fetched:l,event:e,event_state:t,resolve_opts:d,data_serializer:_})}catch(u){return u instanceof Redirect?N(u.status,u.location):M(s,O(u),(await F(e,t,s,u)).message)}}async function mt(e,t,s,r,a){return xe({name:"sveltekit.remote.call",attributes:{},fn:o=>{const c=merge_tracing(e,o);return with_request_store({event:c,state:t},()=>async function(e,t,s,r,a){const[o,c,d]=a.split("/"),l=r._.remotes;l[o]||error(404);const u=(await l[o]()).default[c];u||error(404);const p=u.__,h=s.hooks.transport;let f;e.tracing.current.setAttributes({"sveltekit.remote.call.type":p.type,"sveltekit.remote.call.name":p.name});try{if("query_batch"===p.type){if("POST"!==e.request.method)throw new SvelteKitError(405,"Method Not Allowed",`\`query.batch\` functions must be invoked via POST request, not ${e.request.method}`);const{payloads:r}=await e.request.json(),n=await Promise.all(r.map(e=>X(e,h))),a=await with_request_store({event:e,state:t},()=>p.run(n,s));return json({type:"result",result:K(a,h)})}if("form"===p.type){if("POST"!==e.request.method)throw new SvelteKitError(405,"Method Not Allowed",`\`form\` functions must be invoked via POST request, not ${e.request.method}`);if(!z(e.request))throw new SvelteKitError(415,"Unsupported Media Type",`\`form\` functions expect form-encoded data — received ${e.request.headers.get("content-type")}`);const{data:s,meta:r,form_data:n}=await b(e.request);f=r.remote_refreshes,d&&!("id"in s)&&(s.id=JSON.parse(decodeURIComponent(d)));const a=p.fn,o=await with_request_store({event:e,state:t},()=>a(s,r,n));return json({type:"result",result:K(o,h),refreshes:o.issues?void 0:await _(r.remote_refreshes)})}if("command"===p.type){const{payload:s,refreshes:r}=await e.request.json(),n=X(s,h),a=await with_request_store({event:e,state:t},()=>u(n));return json({type:"result",result:K(a,h),refreshes:await _(r)})}const r="prerender"===p.type?d:new URL(e.request.url).searchParams.get("payload"),n=await with_request_store({event:e,state:t},()=>u(X(r,h)));return json({type:"result",result:K(n,h)})}catch(m){if(m instanceof Redirect)return json({type:"redirect",location:m.location,refreshes:await _(f)});const r=m instanceof HttpError||m instanceof SvelteKitError?m.status:500;return json({type:"error",error:await F(e,t,s,m),status:r},{status:t.prerendering?r:void 0,headers:{"cache-control":"private, no-store"}})}async function _(s){const n=t.refreshes??{};if(s)for(const a of s){if(void 0!==n[a])continue;const[s,o,i]=a.split("/"),c=r._.remotes[s],d=(await(c?.()))?.default?.[o];d||error(400,"Bad Request"),n[a]=with_request_store({event:e,state:t},()=>d(X(i,h)));}if(0!==Object.keys(n).length)return K(Object.fromEntries(await Promise.all(Object.entries(n).map(async([e,t])=>[e,await t]))),h)}}(c,t,s,r,a))}})}async function yt(e,t,s,r){return xe({name:"sveltekit.remote.form.post",attributes:{},fn:n=>{const a=merge_tracing(e,n);return with_request_store({event:a,state:t},()=>async function(e,t,s,r){const[n,a,o]=r.split("/"),i=s._.remotes,d=await(i[n]?.());let l=d?.default[a];if(!l)return e.setHeaders({allow:"GET"}),{type:"error",error:new SvelteKitError(405,"Method Not Allowed","POST method not allowed. No form actions exist for this page")};o&&(l=with_request_store({event:e,state:t},()=>l.for(JSON.parse(o))));try{const s=l.__.fn,{data:r,meta:n,form_data:a}=await b(e.request);return o&&!("id"in r)&&(r.id=JSON.parse(decodeURIComponent(o))),await with_request_store({event:e,state:t},()=>s(r,n,a)),{type:"success",status:200}}catch(u){const e=A(u);return e instanceof Redirect?{type:"redirect",status:e.status,location:e.location}:{type:"error",error:Te(e)}}}(a,t,s,r))}})}async function gt(e,t,s,r,i,d,l,u$1){if(d.depth>10)return text(`Not found: ${e.url.pathname}`,{status:404});if(Oe(e)){const a=await i._.nodes[s.leaf]();return async function(e,t,s,r){const a=r?.actions;if(!a){const r=new SvelteKitError(405,"Method Not Allowed","POST method not allowed. No form actions exist for this page");return Ae({type:"error",error:await F(e,t,s,r)},{status:r.status,headers:{allow:"GET"}})}Ue(a);try{const r=await Ce(e,t,a);return Ae(r instanceof ActionFailure?{type:"failure",status:r.status,data:Ne(r.data,e.route.id,s.hooks.transport)}:{type:"success",status:r?200:204,data:Ne(r,e.route.id,s.hooks.transport)})}catch(i){const r=A(i);return r instanceof Redirect?Pe(r):Ae({type:"error",error:await F(e,t,s,Te(r))},{status:O(r)})}}(e,t,r,a?.server)}try{const h=l.page();let f,_=200;if(function(e){return "POST"===e.request.method}(e)){const s=e.url.searchParams.get("/remote");if(f=s?await yt(e,t,i,s):await async function(e,t,s){const r=s?.actions;if(!r)return e.setHeaders({allow:"GET"}),{type:"error",error:new SvelteKitError(405,"Method Not Allowed","POST method not allowed. No form actions exist for this page")};Ue(r);try{const s=await Ce(e,t,r);return s instanceof ActionFailure?{type:"failure",status:s.status,data:s.data}:{type:"success",status:200,data:s}}catch(n){const e=A(n);return e instanceof Redirect?{type:"redirect",status:e.status,location:e.location}:{type:"error",error:Te(e)}}}(e,t,h.server),"redirect"===f?.type)return N(f.status,f.location);"error"===f?.type&&(_=O(f.error)),"failure"===f?.type&&(_=f.status);}const g=l.prerender();if(g){const e=h.server;if(e?.actions)throw new Error("Cannot prerender pages with actions")}else if(d.prerendering)return new Response(void 0,{status:204});d.prerender_default=g;const w=l.should_prerender_data(),v=je(e.url.pathname),$=[],b=l.ssr(),k=l.csr();if(!(!1!==b||d.prerendering&&w))return u&&f&&e.request.headers.has("x-sveltekit-action"),await pt({branch:[],fetched:$,page_config:{ssr:!1,csr:k},status:_,error:null,event:e,event_state:t,options:r,manifest:i,state:d,resolve_opts:u$1,data_serializer:Le(e,t,r)});const j=[];let R=null;const E=Le(e,t,r),S=d.prerendering&&w?We(e,t,r):null,x=l.data.map((s,r)=>{if(R)throw R;return Promise.resolve().then(async()=>{try{if(s===h&&"error"===f?.type)throw f.error;const n=await Ie({event:e,event_state:t,state:d,node:s,parent:async()=>{const e={};for(let t=0;t<r;t+=1){const s=await x[t];s&&Object.assign(e,s.data);}return e}});return s&&E.add_node(r,n),S?.add_node(r,n),n}catch(n){throw R=n,R}})}),A$1=l.data.map((s,r)=>{if(R)throw R;return Promise.resolve().then(async()=>{try{return await Me({event:e,event_state:t,fetched:$,node:s,parent:async()=>{const e={};for(let t=0;t<r;t+=1)Object.assign(e,await A$1[t]);return e},resolve_opts:u$1,server_data_promise:x[r],state:d,csr:k})}catch(n){throw R=n,R}})});for(const e of x)e.catch(()=>{});for(const e of A$1)e.catch(()=>{});for(let a=0;a<l.data.length;a+=1){const h=l.data[a];if(h)try{const e=await x[a],t=await A$1[a];j.push({node:h,server_data:e,data:t});}catch(p){const l=A(p);if(l instanceof Redirect){if(d.prerendering&&w){const e=JSON.stringify({type:"redirect",location:l.location});d.prerendering.dependencies.set(v,{response:text(e),body:e});}return N(l.status,l.location)}const h=O(l),f=await F(e,t,r,l);for(;a--;)if(s.errors[a]){const n=s.errors[a],o=await i._.nodes[n]();let c=a;for(;!j[c];)c-=1;E.set_max_nodes(c+1);const l=$e(j.slice(0,c+1)),p=new ft(l.map(e=>e.node)),_=l.concat({node:o,data:null,server_data:null});return await pt({event:e,event_state:t,options:r,manifest:i,state:d,resolve_opts:u$1,page_config:{ssr:p.ssr(),csr:p.csr()},status:h,error:f,error_components:await wt(r,b,_,s,i),branch:_,fetched:$,data_serializer:E})}return M(r,h,f.message)}else j.push(null);}if(d.prerendering&&S){let{data:e,chunks:t}=S.get_data();if(t)for await(const s of t)e+=s;d.prerendering.dependencies.set(v,{response:text(e),body:e});}return await pt({event:e,event_state:t,options:r,manifest:i,state:d,resolve_opts:u$1,page_config:{csr:k,ssr:b},status:_,error:null,branch:b?$e(j):[],action_result:f,fetched:$,data_serializer:b?E:Le(e,t,r),error_components:await wt(r,b,j,s,i)})}catch(p){return p instanceof Redirect?N(p.status,p.location):await _t({event:e,event_state:t,options:r,manifest:i,state:d,status:p instanceof HttpError?p.status:500,error:p,resolve_opts:u$1})}}async function wt(e,t,s,r,n){let a;if(e.server_error_boundaries&&t){let e=-1;a=await Promise.all(s.map((t,s)=>{if(0===s)return;if(!t)return null;for(s--;s>e+1&&void 0===r.errors[s];)s-=1;e=s;const a=r.errors[s];return null!=a?n._.nodes[a]?.().then(e=>e.component?.()).catch(()=>{}):void 0}).filter(e=>null!==e));}return a}function vt(e,t=200){return text("string"==typeof e?e:JSON.stringify(e),{status:t,headers:{"content-type":"application/json","cache-control":"private, no-store"}})}function $t(e){return vt({type:"redirect",location:e.location})}const bt=/[\x00-\x1F\x7F()<>@,;:"/[\]?={} \t]/;function kt(e){if(void 0===e?.path)throw new Error("You must specify a `path` when setting, deleting or serializing cookies")}function jt(e,t){const s=e.headers.get("cookie")??"",r=cookieExports.parse(s,{decode:e=>e});let n;const a=new Map,o={httpOnly:true,sameSite:"lax",secure:"localhost"!==t.hostname||"http:"!==t.protocol},i={get(e,r){const n=Array.from(a.values()).filter(s=>s.name===e&&Rt(t.hostname,s.options.domain)&&Et(t.pathname,s.options.path)).sort((e,t)=>t.options.path.length-e.options.path.length)[0];if(n)return 0===n.options.maxAge?void 0:n.value;return cookieExports.parse(s,{decode:r?.decode})[e]},getAll(e){const r=cookieExports.parse(s,{decode:e?.decode}),n=new Map;for(const s of a.values())if(Rt(t.hostname,s.options.domain)&&Et(t.pathname,s.options.path)){const e=n.get(s.name);(!e||s.options.path.length>e.options.path.length)&&n.set(s.name,s);}for(const t of n.values())r[t.name]=t.value;return Object.entries(r).map(([e,t])=>({name:e,value:t}))},set(e,t,s){const r=e.match(bt);r&&console.warn(`The cookie name "${e}" will be invalid in SvelteKit 3.0 as it contains ${r.join(" and ")}. See RFC 2616 for more details https://datatracker.ietf.org/doc/html/rfc2616#section-2.2`),kt(s),d(e,t,{...o,...s});},delete(e,t){kt(t),i.set(e,"",{...t,maxAge:0});},serialize(e,s,r){kt(r);let a=r.path;if(!r.domain||r.domain===t.hostname){if(!n)throw new Error("Cannot serialize cookies until after the route is determined");a=r$2(n,a);}return cookieExports.serialize(e,s,{...o,...r,path:a})}};const c=[];function d(e,s,r){if(!n)return void c.push(()=>d(e,s,r));let o=r.path;r.domain&&r.domain!==t.hostname||(o=r$2(n,o));const i=function(e,t,s){return `${e||""}${t}?${encodeURIComponent(s)}`}(r.domain,o,e),l={name:e,value:s,options:{...r,path:o}};a.set(i,l);}return {cookies:i,new_cookies:a,get_cookie_header:function(e,t){const s={...r};for(const r of a.values()){if(!Rt(e.hostname,r.options.domain))continue;if(!Et(e.pathname,r.options.path))continue;const t=r.options.encode||encodeURIComponent;s[r.name]=t(r.value);}if(t){const e=cookieExports.parse(t,{decode:e=>e});for(const t in e)s[t]=e[t];}return Object.entries(s).map(([e,t])=>`${e}=${t}`).join("; ")},set_internal:d,set_trailing_slash:function(e){n=t$1(t.pathname,e),c.forEach(e=>e());}}}function Rt(e,t){if(!t)return  true;const s="."===t[0]?t.slice(1):t;return e===s||e.endsWith("."+s)}function Et(e,t){if(!t)return  true;const s=t.endsWith("/")?t.slice(0,-1):t;return e===s||e.startsWith(s+"/")}function St(e,t){for(const s of t){const{name:t,value:r,options:n}=s;if(e.append("set-cookie",cookieExports.serialize(t,r,n)),n.path.endsWith(".html")){const s=je(n.path);e.append("set-cookie",cookieExports.serialize(t,r,{...n,path:s}));}}}function qt({event:e,options:t,manifest:s,state:r,get_cookie_header:n,set_internal:a$1}){const o=async(o,i)=>{const c=xt(o,i,e.url);let d$1=(o instanceof Request?o.mode:i?.mode)??"cors",l=(o instanceof Request?o.credentials:i?.credentials)??"same-origin";return t.hooks.handleFetch({event:e,request:c,fetch:async(o,i)=>{const u=xt(o,i,e.url),p=new URL(u.url);u.headers.has("origin")||u.headers.set("origin",e.url.origin),o!==c&&(d$1=(o instanceof Request?o.mode:i?.mode)??"cors",l=(o instanceof Request?o.credentials:i?.credentials)??"same-origin"),"GET"!==u.method&&"HEAD"!==u.method||("no-cors"!==d$1||p.origin===e.url.origin)&&p.origin!==e.url.origin||u.headers.delete("origin");const h=decodeURIComponent(p.pathname);if(p.origin!==e.url.origin||s$1&&h!==s$1&&!h.startsWith(`${s$1}/`)){if(`.${p.hostname}`.endsWith(`.${e.url.hostname}`)&&"omit"!==l){const e=n(p,u.headers.get("cookie"));e&&u.headers.set("cookie",e);}return fetch(u)}const f=a||s$1,_=(h.startsWith(f)?h.slice(f.length):h).slice(1),m=`${_}/index.html`,y=s.assets.has(_)||_ in s._.server_assets,g=s.assets.has(m)||m in s._.server_assets;if(y||g){const e=y?_:m;if(r.read){const t=y?s.mimeTypes[_.slice(_.lastIndexOf("."))]:"text/html";return new Response(r.read(e),{headers:t?{"content-type":t}:{}})}if(d&&e in s._.server_assets){const t=s._.server_assets[e],r=s.mimeTypes[e.slice(e.lastIndexOf("."))];return new Response(d(e),{headers:{"Content-Length":""+t,"Content-Type":r}})}return await fetch(u)}if(W(s,s$1+h))return await fetch(u);if("omit"!==l){const t=n(p,u.headers.get("cookie"));t&&u.headers.set("cookie",t);const s=e.request.headers.get("authorization");s&&!u.headers.has("authorization")&&u.headers.set("authorization",s);}u.headers.has("accept")||u.headers.set("accept","*/*"),u.headers.has("accept-language")||u.headers.set("accept-language",e.request.headers.get("accept-language"));const w=await async function(e,t,s,r){if(e.signal){if(e.signal.aborted)throw new DOMException("The operation was aborted.","AbortError");let n=()=>{};const a=new Promise((t,s)=>{const r=()=>{s(new DOMException("The operation was aborted.","AbortError"));};e.signal.addEventListener("abort",r,{once:true}),n=()=>e.signal.removeEventListener("abort",r);}),o=await Promise.race([zt(e,t,s,{...r,depth:r.depth+1}),a]);return n(),o}return await zt(e,t,s,{...r,depth:r.depth+1})}(u,t,s,r),v=w.headers.get("set-cookie");if(v)for(const e of splitCookiesString(v)){const{name:t,value:s,...r}=parseString(e,{decodeValues:false}),n=r.path??(p.pathname.split("/").slice(0,-1).join("/")||"/");a$1(t,s,{path:n,encode:e=>e,...r});}return w}})};return (e,t)=>{const s=o(e,t);return s.catch(()=>{}),s}}function xt(e,t,s){return e instanceof Request?e:new Request("string"==typeof e?new URL(e,s):e,t)}let Ot,Tt,Pt;const At=({html:e})=>e,Ut=()=>false,Ct=({type:e})=>"js"===e||"css"===e,Nt=new Set(["GET","HEAD","POST"]),Ht=new Set(["GET","HEAD","OPTIONS"]);const zt=(Lt=async function(o,d,l,u$1){const p=new URL(o.url),h$1=p.pathname.endsWith(Re),f=function(e){return e.endsWith(be)||e.endsWith(ke)}(p.pathname),_=function(e){return e.pathname.startsWith(`${s$1}/${t}/remote/`)&&e.pathname.replace(`${s$1}/${t}/remote/`,"")}(p);{const e=o.headers.get("origin");if(_){if("GET"!==o.method&&e!==p.origin)return json({message:"Cross-site remote requests are forbidden"},{status:403})}else if(d.csrf_check_origin&&z(o)&&("POST"===o.method||"PUT"===o.method||"PATCH"===o.method||"DELETE"===o.method)&&e!==p.origin&&(!e||!d.csrf_trusted_origins.includes(e))){const e=`Cross-site ${o.method} form submissions are forbidden`,t={status:403};return "application/json"===o.headers.get("accept")?json({message:e},t):text(e,t)}}if(d.hash_routing&&p.pathname!==s$1+"/"&&"/[fallback]"!==p.pathname)return text("Not found",{status:404});let m;h$1?p.pathname=function(e){return e.slice(0,-11)}(p.pathname):f?(p.pathname=function(e){return e.endsWith(ke)?e.slice(0,-16)+".html":e.slice(0,-12)}(p.pathname)+("1"===p.searchParams.get(V)?"/":"")||"/",p.searchParams.delete(V),m=p.searchParams.get(Z)?.split("").map(e=>"1"===e),p.searchParams.delete(Z)):_&&(p.pathname=o.headers.get("x-sveltekit-pathname")??s$1,p.search=o.headers.get("x-sveltekit-search")??"");const g={},{cookies:w,new_cookies:v,get_cookie_header:E,set_internal:x,set_trailing_slash:P}=jt(o,p),N$1={prerendering:u$1.prerendering,transport:d.hooks.transport,handleValidationError:d.hooks.handleValidationError,tracing:{record_span:xe},is_in_remote_function:false},H={cookies:w,fetch:null,getClientAddress:u$1.getClientAddress||(()=>{throw new Error("@sveltejs/adapter-node does not specify getClientAddress. Please raise an issue")}),locals:{},params:{},platform:u$1.platform,request:o,route:{id:null},setHeaders:e=>{for(const t in e){const s=t.toLowerCase(),r=e[t];if("set-cookie"===s)throw new Error("Use `event.cookies.set(name, value, options)` instead of `event.setHeaders` to set cookies");if(s in g){if("server-timing"!==s)throw new Error(`"${t}" header is already set`);g[s]+=", "+r;}else g[s]=r,u$1.prerendering&&"cache-control"===s&&(u$1.prerendering.cache=r);}},url:p,isDataRequest:f,isSubRequest:u$1.depth>0,isRemoteRequest:!!_};H.fetch=qt({event:H,options:d,manifest:l,state:u$1,get_cookie_header:E,set_internal:x}),u$1.emulator?.platform&&(H.platform=await u$1.emulator.platform({config:{},prerender:!!u$1.prerendering?.fallback}));let W$1=p.pathname;if(!_){const e=u$1.prerendering?.inside_reroute;try{u$1.prerendering&&(u$1.prerendering.inside_reroute=!0),W$1=await d.hooks.reroute({url:new URL(p),fetch:H.fetch})??p.pathname;}catch{return text("Internal Server Error",{status:500})}finally{u$1.prerendering&&(u$1.prerendering.inside_reroute=e);}}try{W$1=s$2(W$1);}catch{return text("Malformed URI",{status:400})}if(W$1!==s$2(p.pathname)&&!u$1.prerendering?.fallback&&W(l,W$1)){const e=new URL(o.url);e.pathname=f?je(W$1):h$1?Ee(W$1):W$1;try{const t=await fetch(e,o),s=new Headers(t.headers);return s.has("content-encoding")&&(s.delete("content-encoding"),s.delete("content-length")),new Response(t.body,{headers:s,status:t.status,statusText:t.statusText})}catch(D){return await U(H,N$1,d,D)}}let I=null;if(s$1&&!u$1.prerendering?.fallback){if(!W$1.startsWith(s$1))return text("Not found",{status:404});W$1=W$1.slice(s$1.length)||"/";}if(h$1)return async function(e,t,s){if(!s._.client.routes)return text("Server-side route resolution disabled",{status:400});const r=await s._.matchers(),n=it(e,s._.client.routes,r);return lt(n?.route??null,n?.params??{},t,s).response}(W$1,new URL(o.url),l);if(W$1===`/${t}/env.js`)return function(e){return Ot??=`export const env=${JSON.stringify(s)}`,Tt??=`W/${Date.now()}`,Pt??=new Headers({"content-type":"application/javascript; charset=utf-8",etag:Tt}),e.headers.get("if-none-match")===Tt?new Response(void 0,{status:304,headers:Pt}):new Response(Ot,{headers:Pt})}(o);if(!_&&W$1.startsWith(`/${t}`)){const e=new Headers;return e.set("cache-control","public, max-age=0, must-revalidate"),text("Not found",{status:404,headers:e})}if(!u$1.prerendering?.fallback){const e=await l._.matchers(),t=it(W$1,l._.routes,e);t&&(I=t.route,H.route={id:I.id},H.params=t.params);}let M={transformPageChunk:At,filterSerializedResponseHeaders:Ut,preload:Ct},F$1="never";try{const i=I?.page?new ft(await function(e,t){return Promise.all([...e.layouts.map(e=>null==e?e:t._.nodes[e]()),t._.nodes[e.leaf]()])}(I.page,l)):void 0;if(I&&!_){if(p.pathname===s$1||p.pathname===s$1+"/")F$1="always";else if(i)F$1=i.trailing_slash();else if(I.endpoint){const e=await I.endpoint();F$1=e.trailingSlash??"never";}if(!f){const e=t$1(p.pathname,F$1);if(e!==p.pathname&&!u$1.prerendering?.fallback)return new Response(void 0,{status:308,headers:{"x-sveltekit-normalize":"1",location:(e.startsWith("//")?p.origin+e:e)+("?"===p.search?"":p.search)}})}if(u$1.before_handle||u$1.emulator?.platform){let e={},t=!1;if(I.endpoint){const s=await I.endpoint();e=s.config??e,t=s.prerender??t;}else i&&(e=i.get_config()??e,t=i.prerender());u$1.before_handle&&u$1.before_handle(H,e,t),u$1.emulator?.platform&&(H.platform=await u$1.emulator.platform({config:e,prerender:t}));}}P(F$1),!u$1.prerendering||u$1.prerendering.fallback||u$1.prerendering.inside_reroute||i$1(p);const h$1=await xe({name:"sveltekit.handle.root",attributes:{"http.route":H.route.id||"unknown","http.method":H.request.method,"http.url":H.url.href,"sveltekit.is_data_request":f,"sveltekit.is_sub_request":H.isSubRequest},fn:async p=>{const h$1={...H,tracing:{enabled:!1,root:p,current:p}};return N$1.allows_commands=h.includes(o.method),await with_request_store({event:h$1,state:N$1},()=>d.hooks.handle({event:h$1,resolve:(p,h)=>xe({name:"sveltekit.resolve",attributes:{"http.route":p.route.id||"unknown"},fn:y=>with_request_store(null,()=>async function(i,p,h){try{if(h&&(M={transformPageChunk:h.transformPageChunk||At,filterSerializedResponseHeaders:h.filterSerializedResponseHeaders||Ut,preload:h.preload||Ct}),d.hash_routing||u$1.prerendering?.fallback)return await pt({event:i,event_state:N$1,options:d,manifest:l,state:u$1,page_config:{ssr:!1,csr:!0},status:200,error:null,branch:[],fetched:[],resolve_opts:M,data_serializer:Le(i,N$1,d)});if(_)return await mt(i,N$1,d,l,_);if(I){const a=i.request.method;let h;if(f)h=await async function(e,t,s,r,a,o,i,d){if(!s.page)return new Response(void 0,{status:404});try{const c=[...s.page.layouts,s.page.leaf],l=i??c.map(()=>!0);let u=!1;const p=new URL(e.url);p.pathname=t$1(p.pathname,d);const h={...e,url:p},f=c.map((e,s)=>function(e){let t,s=!1;return ()=>s?t:(s=!0,t=e())}(async()=>{try{if(u)return {type:"skip"};const r=null==e?e:await a._.nodes[e]();return Ie({event:h,event_state:t,state:o,node:r,parent:async()=>{const e={};for(let t=0;t<s;t+=1){const s=await f[t]();s&&Object.assign(e,s.data);}return e}})}catch(r){throw u=!0,r}})),_=f.map(async(e,t)=>l[t]?e():{type:"skip"});let m=_.length;const y=await Promise.all(_.map((s,a)=>s.catch(async s=>{if(s instanceof Redirect)throw s;return m=Math.min(m,a+1),{type:"error",error:await F(e,t,r,s),status:s instanceof HttpError||s instanceof SvelteKitError?s.status:void 0}}))),g=We(e,t,r);for(let e=0;e<y.length;e++)g.add_node(e,y[e]);const{data:w,chunks:v}=g.get_data();return v?new Response(new ReadableStream({async start(e){e.enqueue(t$2.encode(w));for await(const t of v)e.enqueue(t$2.encode(t));e.close();},type:"bytes"}),{headers:{"content-type":"text/sveltekit-data","cache-control":"private, no-store"}}):vt(w)}catch(l){const s=A(l);return s instanceof Redirect?$t(s):vt(await F(e,t,r,s),500)}}(i,N$1,I,d,l,u$1,m,F$1);else if(!I.endpoint||I.page&&!function(r){const{method:n,headers:a}=r.request;if(c$1.includes(n)&&!p$1.includes(n))return !0;if("POST"===n&&"true"===a.get("x-sveltekit-action"))return !1;const o=r.request.headers.get("accept")??"*/*";return "text/html"!==j(o,["*","text/html"])}(i)){if(!I.page)throw new Error("Route is neither page nor endpoint. This should never happen");if(!p)throw new Error("page_nodes not found. This should never happen");if(Nt.has(a))h=await gt(i,N$1,I.page,d,l,u$1,p,M);else {const e=new Set(Ht),t=await l._.nodes[I.page.leaf]();if(t?.server?.actions&&e.add("POST"),"OPTIONS"===a)h=new Response(null,{status:204,headers:{allow:Array.from(e.values()).join(", ")}});else {const t=[...e].reduce((e,t)=>(e[t]=!0,e),{});h=B(t,a);}}}else h=await async function(e,t,s,n){const a=e.request.method;let o=s[a]||s.fallback;if("HEAD"===a&&!s.HEAD&&s.GET&&(o=s.GET),!o)return B(s,a);const i=s.prerender??n.prerender_default;if(i&&(s.POST||s.PATCH||s.PUT||s.DELETE))throw new Error("Cannot prerender endpoints that have mutative methods");if(n.prerendering&&!n.prerendering.inside_reroute&&!i){if(n.depth>0)throw new Error(`${e.route.id} is not prerenderable`);return new Response(void 0,{status:204})}try{t.allows_commands=!0;const s=await with_request_store({event:e,state:t},()=>o(e));if(!(s instanceof Response))throw new Error(`Invalid response from route ${e.url.pathname}: handler should return a Response object`);if(n.prerendering&&(!n.prerendering.inside_reroute||i)){const t=new Response(s.clone().body,{status:s.status,statusText:s.statusText,headers:new Headers(s.headers)});if(t.headers.set("x-sveltekit-prerender",String(i)),!n.prerendering.inside_reroute||!i)return t;t.headers.set("x-sveltekit-routeid",encodeURI(e.route.id)),n.prerendering.dependencies.set(e.url.pathname,{response:t,body:null});}return s}catch(c){if(c instanceof Redirect)return new Response(void 0,{status:c.status,headers:{location:c.location}});throw c}}(i,N$1,await I.endpoint(),u$1);if("GET"===o.method&&I.page&&I.endpoint){const e=h.headers.get("vary")?.split(",")?.map(e=>e.trim().toLowerCase());e?.includes("accept")||e?.includes("*")||(h=new Response(h.body,{status:h.status,statusText:h.statusText,headers:new Headers(h.headers)}),h.headers.append("Vary","Accept"));}return h}if(u$1.error&&i.isSubRequest){const e=new Headers(o.headers);return e.set("x-sveltekit-error","true"),await fetch(o,{headers:e})}if(u$1.error)return text("Internal Server Error",{status:500});if(0===u$1.depth)return u&&i.url.pathname,await _t({event:i,event_state:N$1,options:d,manifest:l,state:u$1,status:404,error:new SvelteKitError(404,"Not Found",`Not found: ${i.url.pathname}`),resolve_opts:M});if(u$1.prerendering)return text("not found",{status:404});const y=await fetch(o);return new Response(y.body,y)}catch(y){return await U(i,N$1,d,y)}finally{i.cookies.set=()=>{throw new Error("Cannot use `cookies.set(...)` after the response has been generated")},i.setHeaders=()=>{throw new Error("Cannot use `setHeaders(...)` after the response has been generated")};}}(merge_tracing(p,y),i,h).then(e=>{for(const t in g){const s=g[t];e.headers.set(t,s);}return St(e.headers,v.values()),u$1.prerendering&&null!==p.route.id&&e.headers.set("x-sveltekit-routeid",encodeURI(p.route.id)),y.setAttributes({"http.response.status_code":e.status,"http.response.body.size":e.headers.get("content-length")||"unknown"}),e}))})}))}});if(200===h$1.status&&h$1.headers.has("etag")){let e=o.headers.get("if-none-match");e?.startsWith('W/"')&&(e=e.substring(2));const t=h$1.headers.get("etag");if(e===t){const e=new Headers({etag:t});for(const t of ["cache-control","content-location","date","expires","vary","set-cookie"]){const s=h$1.headers.get(t);s&&e.set(t,s);}return new Response(void 0,{status:304,headers:e})}}if(f&&h$1.status>=300&&h$1.status<=308){const e=h$1.headers.get("location");if(e)return $t(new Redirect(h$1.status,e))}return h$1}catch(G){if(G instanceof Redirect){const e=f||_?$t(G):I?.page&&Oe(H)?Pe(G):N(G.status,G.location);return St(e.headers,v.values()),e}return await U(H,N$1,d,G)}},async(e,...t)=>Lt(e,...t));var Lt;function Wt(e,t,s){return Object.fromEntries(Object.entries(e).filter(([e])=>e.startsWith(t)&&(""===s||!e.startsWith(s))))}let It,Mt=null;class Ft{#k;#j;constructor(e){if(this.#k=m,this.#j=e,ve){const e=this.respond.bind(this);this.respond=async(...t)=>{const{promise:s,resolve:r}=ge(),n=Mt;return Mt=s,await n,e(...t).finally(r)};}}async init({env:e,read:t}){const{env_public_prefix:s,env_private_prefix:r$1}=this.#k;if(r(Wt(e,r$1,s)),i(Wt(e,s,r$1)),t){l(e=>{const s=t(e);return s instanceof ReadableStream?s:new ReadableStream({async start(e){try{const t=await Promise.resolve(s);if(!t)return void e.close();const r=t.getReader();for(;;){const{done:t,value:s}=await r.read();if(t)break;e.enqueue(s);}e.close();}catch(t){e.error(t);}}})});}await(It??=(async()=>{try{const e=await p();this.#k.hooks={handle:e.handle||(({event:e,resolve:t})=>t(e)),handleError:e.handleError||(({status:e,error:t,event:s})=>{const r=G(e,t,s);console.error(r);}),handleFetch:e.handleFetch||(({request:e,fetch:t})=>t(e)),handleValidationError:e.handleValidationError||(({issues:e})=>(console.error("Remote function schema validation failed:",e),{message:"Bad Request"})),reroute:e.reroute||(()=>{}),transport:e.transport||{}},e.transport&&Object.fromEntries(Object.entries(e.transport).map(([e,t])=>[e,t.decode])),e.init&&await e.init();}catch(e){throw e}})());}async respond(e,t){return zt(e,this.#k,this.#j,{...t,error:false,depth:0})}}

export { Ft as Server };
//# sourceMappingURL=index.js.map
