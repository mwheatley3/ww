// import uuid from 'uuid/v4';
export class User {
    id: string
    email: string

    constructor(id, email) {
        this.id = id;
        this.email = email;
    }

    static fromJSON(obj) {
        return new User(obj.id, obj.email);
    }
}
