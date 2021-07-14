const USER_ADDED = 0;
const USER_UPDATED = 1;
const USER_SAME = 2;

class UserState {
    constructor() {
        this.users = [];
    }

    update(name, spec) {
        const user = this.users.find(u => u.name === name);
        if (!user) {
            console.debug(`[state] adding user to state: ${name}`);
            this.users.push({ name, spec });
            return USER_ADDED;
        }

        if (user.spec.name !== spec.name || user.spec.email !== spec.email) {
            console.debug(`[state] updating user ${name}: ${user.spec.name} -> ${spec.name} and ${user.spec.email} -> ${spec.email}`);
            user.spec = { ...user.spec, name: spec.name, email: spec.email };
            return USER_UPDATED;
        }

        console.debug(`[state] leaving user: ${name}`);
        return USER_SAME;
    }

    remove(name) {
        const idx = this.users.findIndex(u => u.name === name);
        if (idx < 0) {
            return false;
        }

        this.users = this.users.splice(idx);
    }
};

module.exports = {
    UserState,
    States: {
        USER_ADDED,
        USER_UPDATED,
        USER_SAME
    }
};
