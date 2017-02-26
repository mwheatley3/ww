import 'whatwg-fetch';
import Promise from 'bluebird';

export function urlEncode(params) {
    const p = [];

    for (const k in params) {
        p.push(encodeURIComponent(k) + '=' + encodeURIComponent(params[k]));
    }

    return p.join('&');
}

export function pathString(url, data) {
    const qs = urlEncode(data);
    if (!qs) {
        return url;
    }

    return url + '?' + qs;
}

export function handleResponse(resp, body, Type) {
    if (resp.status >= 400) {
        let msg = body;
        let data;

        if (body && body.error) {
            msg = body.error.message;
            data = body.error.data;
        }

        throw new FetchError(resp.status, msg, data);
    }

    return Type ? Type(body.data) : null;
}

function call({ url, options = {}, Type }) {
    if (!options.headers) {
        options.headers = {};
    }

    options.headers.Accept = 'application/json';
    options.credentials = 'same-origin';

    return Promise.resolve(fetch(url, options)
        .then(resp => (
            (resp.headers.get('Content-Type') === 'application/json' ? resp.json() : resp.text())
                .then(body => handleResponse(resp, body, Type))
        )));
}

export function get({ url, data = {}, Type }) {
    url = pathString(url, data);

    return call({ url, Type, options: {
        method: 'GET',
    } });
}

export function post({ url, data = {}, Type, jsonReqBody = true }) {
    return call({ url, Type, options: {
        method: 'POST',
        headers: headers(jsonReqBody),
        body: handleBody(data),
    } });
}

export function put({ url, data = {}, Type }) {
    return call({ url, Type, options: {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
        },
        body: handleBody(data),
    } });
}

export function del({ url, data = {}, Type, jsonReqBody = true }) {
    return call({ url, Type, options: {
        method: 'DELETE',
        headers: headers(jsonReqBody),
        body: handleBody(data),
    } });
}

function handleBody(data) {
    if (data instanceof Blob || data instanceof FormData) {
        return data;
    }

    return JSON.stringify(data);
}

export function jsonp({ url, data = {}, Type, param = 'jsonp' }) {
    return new Promise((resolve, reject) => {
        const cb = '__jsonp_cb__' + String(Math.random()).substr(2);
        data[param] = cb;
        url = pathString(url, data);

        window[cb] = d => {
            delete window[cb];
            if (!Type) {
                return resolve();
            }

            resolve(Type(d));
        };

        const sc = document.createElement('script');
        sc.async = true;
        sc.src = url;
        sc.onerror = ev => reject(new Error("JSONP loading error from " + url));

        document.getElementsByTagName('head')[0].appendChild(sc);
    });
}


export class FetchError extends Error {
    statusCode: number;
    data: any;

    constructor(code: number, message: string, data: ?any) {
        super(message);

        this.statusCode = code;
        this.data = data;
    }
}

function headers(jsonReqBody) {
    if (jsonReqBody) {
        return { 'Content-Type': 'application/json' };
    }

    return {};
}
