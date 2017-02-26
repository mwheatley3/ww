import Value from 'js/common/value';
import { computed } from 'mobx';
import Api from 'js/personal/api';

export default class AuthStore {
    user = new Value();
    api: Api;

    constructor(api: Api) {
        this.api = api;
    }

    init() {
        // see if we get a user, but if we get an error don't worry about it
        this.user.trackPromise(this.api.getUser().catch(err => void 0));
    }

    @computed get loggedIn(): boolean {
        return !!this.user.value;
    }

    login(username: string, password: string) {
        this.user.trackPromise(this.api.login(username, password));
    }

    logout() {
        this.user.setValue(null);
        this.api.logout().catch(err => console.error("logout error", err));
    }
}
