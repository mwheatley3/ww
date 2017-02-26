import { provider } from 'js/common/context';
import AuthStore from 'js/personal/store/auth';
import Api from 'js/personal/api';

export const withStore = provider('store');

export default class Store {
    _inited = false;
    auth: AuthStore;

    constructor(api: Api) {
        this.auth = new AuthStore(api);
    }

    init() {
        if (this._inited) {
            return;
        }

        this._inited = true;
        this.auth.init();
    }
}
