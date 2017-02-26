import { observable, computed, action, autorun } from 'mobx';
import Promise from 'bluebird';

export const EMPTY = 'empty';
export const LOADING = 'loading';
export const SUCCESS = 'success';
export const ERROR = 'error';

export default class Value {
    @observable state = EMPTY;
    @observable value = void 0;
    @observable error = void 0;

    constructor({ value, error } = {}) {
        if (value) {
            this.setValue(value);
        } else if (error) {
            this.setError(error);
        }
    }

    @computed get empty() {
        return this.state === EMPTY;
    }

    @computed get loaded() {
        return this.state !== EMPTY && this.state !== LOADING;
    }

    @computed get loading() {
        return this.state === LOADING;
    }

    @action setValue(v) {
        this.value = v;
        this.error = void 0;
        this.state = SUCCESS;
    }

    @action setError(err) {
        this.value = void 0;
        this.error = err;
        this.state = ERROR;
    }

    @action trackPromise(p) {
        this.state = LOADING;

        p.then(
            val => this.setValue(val),
            err => this.setError(err)
        );
    }

    update(fn) {
        this.setValue(fn(this.value));
    }

    onValue(fn) {
        const dispose = autorun(() => {
            if (this.state === SUCCESS) {
                fn(this.value);
                dispose();
                return;
            }
        });

        return this;
    }

    asPromise() {
        return new Promise((resolve, reject) => {
            const dispose = autorun(() => {
                if (this.state === SUCCESS) {
                    resolve(this.value);
                    dispose();
                    return;
                }

                if (this.state === ERROR) {
                    reject(this.error);
                    dispose();
                    return;
                }
            });
        });
    }

    static trackPromise(p) {
        const v = new Value();
        v.trackPromise(p);

        return v;
    }
}
