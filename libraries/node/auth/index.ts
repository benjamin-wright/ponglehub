import axios from "axios";

const authUrl = "http://localhost:4000";

export interface UserData {
  name: string
}

function parseUserData(data: string): UserData {
  const parsed = JSON.parse(data);
  if (typeof parsed.name !== "string") {
    throw new Error(`parsing error, name property not found on "${data}"`);
  }

  return {
    name: parsed.name
  };
}

export class Auth {
  private storage: Storage;

  constructor(storage: Storage) {
    this.storage = storage;
  }

  loggedIn(): boolean {
    return !!this.storage.getItem('userData');
  }

  loading(): boolean {
    return this.storage.getItem('loading') == 'true';
  }

  async restore(): Promise<UserData> {
    const userString = this.storage.getItem('userData');
    if (userString == null) {
      throw new Error('userData not found');
    }

    return parseUserData(userString);
  }

  async load(): Promise<UserData> {
    if (!this.loading()) {
      throw new Error('tried to load user data without logging in');
    }

    const response = await axios.get(`${authUrl}/auth/user`, {
      withCredentials: true,
    });

    if (response.status != 200) {
      throw new Error(`failed to get user data: status code ${response.status}`);
    }

    if (typeof response.data.name !== "string") {
      throw new Error(`parsing error, name property not found on "${JSON.stringify(response.data)}"`);
    }

    const data: UserData = {
      name: response.data.name
    };

    this.storage.removeItem('loading');
    this.storage.setItem('userData', JSON.stringify(data));

    return data;
  }

  async logOut(): Promise<any> {
    this.storage.clear();

    const response = await axios.post(
      `${authUrl}/auth/logout`,
      {},
      { withCredentials: true }
    );

    if (response.status != 204) {
      throw new Error(`failed signing out: status code ${response.status}`);
    };
  }

  /** Redirect to the login page. */
  logIn() {
    this.storage.setItem('loading', 'true');
    window.location.href = `${authUrl}/auth/login?redirect=${window.location.toString()}`;
  }
}
