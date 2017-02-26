import qs from 'query-string';
// import { omit } from 'lodash';

export function getQuery(loc, key) {
    return qs.parse(loc.search)[key];
}

export function removeQueryKeys(loc, ...keys) {
    let q = qs.parse(loc.search);
    // q = omit(q, keys);

    return { ...loc, search: qs.stringify(q) };
}

export function mergeQuery(loc, q) {
    q = {
        ...qs.parse(loc.search),
        ...q,
    };

    return { ...loc, search: qs.stringify(q) };
}
