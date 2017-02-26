import { get, post, del } from 'js/common/fetch';
import { User } from 'js/personal/models';

const identity = x => x;

export default class API {
    baseURL: string;

    constructor(baseURL: string) {
        this.baseURL = baseURL;
    }

    url(...parts: string[]) {
        return [this.baseURL, ...parts].join('/');
    }

    projectURL(projectID: string, ...parts: string[]) {
        return this.url('projects', projectID, ...parts);
    }

    login(email: string, password: string) {
        return post({
            url: this.url('auth'),
            data: { email, password },
            Type: User.fromJSON,
        });
    }

    logout() {
        return del({
            url: this.url('auth'),
            Type: identity,
        });
    }

    getUser(userID: string = 'me') {
        return get({
            url: this.url('users', userID),
            Type: User.fromJSON,
        });
    }
}
